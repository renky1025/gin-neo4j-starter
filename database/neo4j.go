package database

import (
	"context"
	"fmt"
	"go-gin-restful-service/config"
	"go-gin-restful-service/log"

	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/pilinux/structs"
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

type Person struct {
	ID   int64  `json:"id" from:"id"`
	Name string `json:"name" from:"name"`
	Age  int    `json:"age" from:"age"`
	//Identifiers int64  `json:"identifiers,omitempty" from:"identifiers"`
}

func (n *Neo4jDriver) CreatePerson(name string, age int) (*Person, error) {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	cypher := `CREATE(user:Person) SET user = $prop RETURN user`
	result, err := session.Run(ctx, cypher, map[string]interface{}{
		"prop": structs.Map(Person{Name: name, Age: age}),
	})

	// result, err := session.Run(ctx,
	// 	"CREATE (p:Person {name: $name, age: $age}) RETURN id(p)",
	// 	map[string]interface{}{"name": name, "age": age, "id": util.GenerateSnowID()},
	// )
	if err != nil {
		return nil, err
	}

	record, err := result.Single(ctx)
	if err != nil {
		return nil, err
	}

	id, ok := record.Values[0].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid ID type")
	}

	return &Person{ID: id, Name: name, Age: age}, nil
}

func (n *Neo4jDriver) GetPersonByName(name string) (*Person, error) {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.Run(ctx,
		"MATCH (p:Person) WHERE p.name = $name RETURN id(p), p.age LIMIT 1",
		map[string]interface{}{"name": name},
	)
	if err != nil {
		return nil, err
	}

	record, err := result.Single(ctx)
	if err != nil {
		return nil, err
	}

	id, ok := record.Values[0].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid ID type")
	}

	age, ok := record.Values[1].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid age type")
	}

	return &Person{ID: id, Name: name, Age: int(age)}, nil
}

func (n *Neo4jDriver) GetPersonById(id int64) (*Person, error) {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	//match (n) where id(n)=1 return n
	result, err := session.Run(ctx,
		"MATCH (p:Person) WHERE id(p) = $id RETURN p.name, p.age LIMIT 1",
		map[string]interface{}{"id": id},
	)
	if err != nil {
		return nil, err
	}

	record, err := result.Single(ctx)
	if err != nil {
		return nil, err
	}

	name, ok := record.Values[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid name type")
	}

	age, ok := record.Values[1].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid age type")
	}

	return &Person{ID: id, Name: name, Age: int(age)}, nil
}

func (n *Neo4jDriver) UpdatePersonAge(id int64, age int) (*Person, error) {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	// 增加节点属性 MATCH (a:Person {name:'Shawn'}) SET a.city="上海"
	result, err := session.Run(ctx,
		"MATCH (p:Person) WHERE id(p) = $id SET p.age = $age RETURN p.name, p.age",
		map[string]interface{}{"id": id, "age": age},
	)
	if err != nil {
		return nil, err
	}

	record, err := result.Single(ctx)
	if err != nil {
		return nil, err
	}

	name, ok := record.Values[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid name type")
	}

	newAge, ok := record.Values[1].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid age type")
	}

	return &Person{ID: id, Name: name, Age: int(newAge)}, nil
}

func (n *Neo4jDriver) DeletePerson(id int64) error {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	// 删除有关系的节点 MATCH (a:Person {name:'Todd'})-[rel]-(b:Person) DELETE a,b,rel
	// 删除节点 MATCH (a:Location {city:'Portland'}) DELETE a
	// 删除节点的属性 MATCH (a:Person {name:'Mike'}) REMOVE a.test;
	_, err := session.Run(ctx,
		"MATCH (p:Person) WHERE id(p) = $id DELETE p",
		map[string]interface{}{"id": id},
	)
	if err != nil {
		return err
	}

	return nil
}

func (n *Neo4jDriver) CreateRelationship(node1 string, node2 string) error {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	// MATCH (a:Person {name:'Shawn'}),
	// (b:Person {name:'Sally'})
	// MERGE (a)-[:FRIENDS {since:2001}]->(b)
	_, err := session.Run(ctx,
		"MATCH (a:Person {name: $node1}), (b:Person {name: $node2}) MERGE (a)-[:FRIENDS]->(b)",
		map[string]interface{}{"node1": node1, "node2": node2},
	)
	if err != nil {
		return err
	}

	return nil
}

func (n *Neo4jDriver) SearchPerson(name string, offset int64, limit int64) ([]Person, error) {
	driver := *n.DBCONN
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	cypher := `MATCH (p:Person) WHERE p.name =~ '.*'+$name+'.*' RETURN p SKIP $offset LIMIT $limit`
	result, err := session.Run(ctx, cypher,
		map[string]interface{}{"name": name, "offset": offset, "limit": limit},
	)
	if err != nil {
		return nil, err
	}

	res := make([]Person, 0)
	for result.Next(ctx) {
		record := result.Record()
		if value, ok := record.Get("p"); ok {
			node := value.(neo4j.Node)
			props := node.Props
			person := Person{}
			mapstructure.Decode(props, &person)
			res = append(res, person)
		}
	}
	return res, nil
}
