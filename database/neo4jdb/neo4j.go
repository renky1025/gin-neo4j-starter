package neo4jdb

import (
	"context"
	"fmt"
	"go-gin-restful-service/config"
	"go-gin-restful-service/dto"
	"go-gin-restful-service/log"
	"go-gin-restful-service/model"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jDriver struct {
	DBCONN *neo4j.DriverWithContext
}

var ctx = context.Background()

func NewNeo4jDriver(cfg *config.Config) *Neo4jDriver {
	neo := InitNeo4j(cfg)
	return &Neo4jDriver{
		DBCONN: &neo,
	}
}

func InitNeo4j(cfg *config.Config) neo4j.DriverWithContext {
	dbUri := cfg.Neo4J.URI // scheme://host(:port) (default port is 7687)
	neodriver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(cfg.Neo4J.UserName, cfg.Neo4J.Password, ""))
	if err != nil {
		log.Logger.Panic(err)
	}
	//session := neodriver.NewSession(context.Background(), neo4j.SessionConfig{})
	log.Logger.Info("Connection with neo4j successed!")
	return neodriver
}

func (n *Neo4jDriver) CreateNode(label string, nodeMap map[string]interface{}) (int64, error) {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	cypher := `CREATE(node:` + label + `) SET node = $prop RETURN id(node)`
	result, err := session.Run(ctx, cypher, map[string]interface{}{
		"prop": nodeMap,
	})

	if err != nil {
		return 0, err
	}

	record, err := result.Single(ctx)
	if err != nil {
		return 0, err
	}

	id, ok := record.Values[0].(int64)
	if !ok {
		return 0, fmt.Errorf("invalid ID type")
	}

	return id, nil
}

func (n *Neo4jDriver) GetNodeBy(node model.NodeFrom) (interface{}, error) {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	//match (n) where id(n)=1 return n
	result, err := session.Run(ctx,
		"MATCH (p:"+node.Label+" {name: $name}) RETURN p LIMIT 1",
		map[string]interface{}{"name": node.Name},
	)
	if err != nil {
		return nil, err
	}

	record, err := result.Single(ctx)
	if err != nil {
		return nil, err
	}

	return record.Values[0], nil
}

func buildValues(updateData map[string]interface{}) string {
	temp := make([]string, 0)
	for k, v := range updateData {
		fmt.Println(reflect.TypeOf(v))
		switch v.(type) {
		case float64:
			temp = append(temp, fmt.Sprintf(" p.%s = %v ", k, v))
		case string:
			temp = append(temp, fmt.Sprintf(" p.%s = '%v' ", k, v))
		default:
			temp = append(temp, fmt.Sprintf(" p.%s = '%v' ", k, v))
		}
	}
	return strings.Join(temp, ",")
}

func (n *Neo4jDriver) UpdateNodeAttr(node model.NodeFrom, updateData map[string]interface{}) (bool, error) {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	// 增加节点属性 MATCH (a:Person {name:'Shawn'}) SET a.city="上海"
	temp := buildValues(updateData)
	result, err := session.Run(ctx,
		"MATCH (p:"+node.Label+") WHERE p.name = $name SET "+temp+" RETURN p.name",
		map[string]interface{}{"name": node.Name},
	)
	if err != nil {
		return false, err
	}

	record, err := result.Single(ctx)
	if err != nil {
		return false, err
	}

	_, ok := record.Values[0].(string)
	if !ok {
		return false, fmt.Errorf("invalid name type")
	}

	return true, nil
}

func (n *Neo4jDriver) DeleteNodeByLabel(label string, name string) (bool, error) {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	// 删除有关系的节点 MATCH (a:Person {name:'Todd'})-[rel]-(b:Person) DELETE a,b,rel
	// 删除节点 MATCH (a:Location {city:'Portland'}) DELETE a
	// 删除节点的属性 MATCH (a:Person {name:'Mike'}) REMOVE a.test;
	cypher := "MATCH (n:" + label + " {name: " + name + "}) detach delete n"
	if len(name) == 0 {
		cypher = "MATCH (n:" + label + ") detach delete n"
	}

	_, err := session.Run(ctx,
		cypher, map[string]interface{}{},
		//"MATCH (p:Person) WHERE id(p) = $id DELETE p",map[string]interface{}{"id": id},
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (n *Neo4jDriver) CleanAllDB(id int64) error {
	// match (e)-[l]->(x) delete e,l,x;
	// match (n) delete n;
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.Run(ctx,
		"match (e)-[l]->(x) delete e,l,x;match (n) delete n;",
		map[string]interface{}{},
	)
	if err != nil {
		return err
	}

	return nil
}

func buildMergeValues(inputData map[string]interface{}) string {
	temp := make([]string, 0)
	for k, v := range inputData {
		fmt.Println(reflect.TypeOf(v))
		switch v.(type) {
		case float64:
			temp = append(temp, fmt.Sprintf(" %s: %v ", k, v))
		case string:
			temp = append(temp, fmt.Sprintf(" %s: '%v' ", k, v))
		default:
			temp = append(temp, fmt.Sprintf(" %s: '%v' ", k, v))
		}
	}

	return "{" + strings.Join(temp, ",") + "}"
}

func (n *Neo4jDriver) CreateRelationship(node1 model.NodeFrom, node2 model.NodeFrom, relation model.RelationMap) error {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	// MATCH (a:Person {name:'Shawn'}),
	// (b:Person {name:'Sally'})
	// MERGE (a)-[:FRIENDS {since:2001}]->(b)
	cypher := "MATCH (a:" + node1.Label + " {name: $node1}), (b:" + node2.Label + " {name: $node2}) MERGE (a)-[:" + relation.RelationShip + "]->(b)"
	if relation.Attr != nil {
		cypher = "MATCH (a:" + node1.Label + " {name: $node1}), (b:" + node2.Label + " {name: $node2}) MERGE (a)-[:" + relation.RelationShip + " " + buildMergeValues(relation.Attr) + "]->(b)"
	}
	_, err := session.Run(ctx,
		cypher,
		map[string]interface{}{"node1": node1.Name, "node2": node2.Name},
	)
	if err != nil {
		return err
	}

	return nil
}

// query condition: name, relationship, times, relation properties
func (n *Neo4jDriver) SearchNodes(params dto.QueryParamsDTO) ([]map[string]interface{}, error) {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	queryMap := map[string]interface{}{
		"offset": (params.PageNo - 1) * params.PageSize,
		"limit":  params.PageSize,
	}
	where := `WHERE `
	cypher := `MATCH (f)`
	wherefrom := ""
	whereto := ""
	if params.NodeFrom != nil {
		cypher = `MATCH (f:` + params.NodeFrom.Label + `)`
		if len(params.NodeFrom.Name) > 0 {
			wherefrom = ` f.name =~ '.*` + params.NodeFrom.Name + `.*' `
		}
	}

	if params.RelationShip != nil {
		cypher += `-[rels:` + params.RelationShip.Name + `*` + strconv.Itoa(params.RelationShip.Number) + `]`
	}
	if params.NodeTo != nil {
		if len(params.NodeTo.Name) > 0 {
			whereto += ` t.name =~ '.*` + params.NodeTo.Name + `.*'  `
		}
		cypher += `-(t:` + params.NodeTo.Label + `) `
		if len(wherefrom) > 0 {
			cypher += where + wherefrom
		}
		if len(whereto) > 0 {
			cypher += ` AND ` + whereto
		}
		cypher += ` RETURN t `
	} else {
		if params.RelationShip != nil {
			cypher += `-() `
		}
		if len(wherefrom) > 0 {
			cypher += where + wherefrom
		}
		cypher += ` RETURN f `
	}
	if params.Sort != nil && len(*params.Sort) > 0 && params.Order != nil && len(*params.Order) > 0 {
		cypher += ` ORDER BY f.` + *params.Order + ` ` + *params.Sort
	}
	cypher += ` SKIP $offset LIMIT $limit `

	//WHERE p.name =~ '.*'+$name+'.*' RETURN p ORDER BY n.name DESC SKIP $offset LIMIT $limit`
	//MATCH (n:Person  {person_id:'180'})-[rels:FRIEND*2]-(m:Person)
	// 查询有关系的节点
	// cypher1 := `MATCH (n)-[:MARRIED]-() RETURN n SKIP $offset LIMIT $limit`
	// 查找 节点 关系的关系的 信息 MATCH (a:Person {name:'Mike'})-[r1:FRIENDS]-()-[r2:FRIENDS]-(friend_of_a_friend) RETURN friend_of_a_friend.name AS fofName
	result, err := session.Run(ctx, cypher, queryMap)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	res := make([]map[string]interface{}, 0)
	for result.Next(ctx) {
		record := result.Record()
		if value, ok := record.Get("f"); ok {
			node := value.(neo4j.Node)
			props := node.Props
			person := map[string]interface{}{}
			mapstructure.Decode(props, &person)
			res = append(res, person)
		}
	}
	return res, nil
}
