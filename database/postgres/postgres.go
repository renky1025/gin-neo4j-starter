package postgres

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// 定义表结构
type NodeInfoModel struct {
	tableName    struct{} `pg:"node_info"`
	NodeName     string   `json:"node_name" pg:"node_name,pk"`
	NodeIp       string   `json:"node_ip" pg:"node_ip"`
	NodePort     string   `json:"node_port" pg:"node_port"`
	NodeUsername string   `json:"node_username" pg:"node_username"`
	NodePassword string   `json:"node_password" pg:"node_password"`
}

func (u *NodeInfoModel) TableName() string {
	return "node_info"
}

// -----------------------------------
// DataBase的配置
type DatabaseConfig struct {
	URL                    string `yaml:"url"`
	ConnMaxLifetimeSeconds int    `yaml:"conn_max_lifetime_sec"`
	MaxOpenConns           int    `yaml:"max_open_conns"`
	MinIdleConns           int    `yaml:"min_idle_conns"`
	MaxRetries             int    `yaml:"max_retries"`
}

const defaultDBConnMaxRetries = 5

type DBService struct {
	DB *pg.DB
}

func InitDBService(dbCfg *DatabaseConfig) (*pg.DB, string, error) {
	// 1: 根据postgre的链接信息，生成一个 Options 对象
	opt, err := pg.ParseURL(dbCfg.URL)
	if err != nil {
		return nil, "", err
	}

	// 2：这里给 Options对象设置 链接、链接池 相关的参数～～～最长链接时间、最小空闲链接数、链接池最大链数（即链接池的大小）、最大重试次数～
	// 最大链接时间
	opt.MaxConnAge = time.Second * time.Duration(dbCfg.ConnMaxLifetimeSeconds)
	// 最小的空闲链接数
	opt.MinIdleConns = dbCfg.MinIdleConns
	// 链接池的最大连接数
	opt.PoolSize = dbCfg.MaxOpenConns
	// 最大retry次数
	if dbCfg.MaxRetries == 0 {
		opt.MaxRetries = defaultDBConnMaxRetries
	} else {
		opt.MaxRetries = dbCfg.MaxRetries
	}

	// 3: 通过 Options 对象生成DB对象～最终业务代码中用的是DB对象，里面有一个 newConnPool 方法用来设置链接池的～
	db := pg.Connect(opt)
	// 获取 version 信息
	var dbVersion string
	_, err = db.QueryOne(pg.Scan(&dbVersion), "select version()")
	if err != nil {
		return nil, "", err
	}
	return db, dbVersion, nil
}

func (pgDB *DBService) CreateTable() error {
	models := []interface{}{
		(*NodeInfoModel)(nil),
	}

	for _, model := range models {
		err := pgDB.DB.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (pgDB *DBService) DeleteTable() error {
	models := []interface{}{
		(*NodeInfoModel)(nil),
	}
	for _, model := range models {
		err := pgDB.DB.Model(model).DropTable(&orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			fmt.Println("删除表 ", model, "raise error: ", err)
			return err
		}
	}
	return nil
}

func (pgDB *DBService) InsertNodeInfo() error {
	// 插入数据方法一
	nodeinfodata := &NodeInfoModel{
		NodeName:     "123",
		NodeIp:       "10.0.0.5",
		NodePort:     "2222",
		NodeUsername: "123",
		NodePassword: "1234321",
	}
	result, err := pgDB.DB.Model(nodeinfodata).Insert()
	if err != nil {
		return err
	}
	fmt.Printf("single insert rows affected:%d", result.RowsAffected())
	return nil
}

func (pgDB *DBService) SelectAllNodeInfo() error {
	var nodeinfodata []NodeInfoModel
	err := pgDB.DB.Model(&nodeinfodata).Select()
	if err != nil {
		return err
	}
	fmt.Println(nodeinfodata)
	return nil
}

func (pgDB *DBService) SelectNodeInfo() (interface{}, error) {
	nodeinfodata := &NodeInfoModel{
		NodeName: "123",
	}
	err := pgDB.DB.Model(nodeinfodata).Select()
	if err != nil {
		return nil, err
	}
	fmt.Println(nodeinfodata)
	return nodeinfodata, err
}

func (pgDB *DBService) UpdateNodeInfo() error {
	var nodeinfodata []NodeInfoModel

	updatedata := &NodeInfoModel{
		NodeName:     "123",
		NodeIp:       "10.0.0.5",
		NodePort:     "3333",
		NodeUsername: "123",
		NodePassword: "1234321",
	}

	_, err := pgDB.DB.Model(&nodeinfodata).
		Set("node_port = ?", updatedata.NodePort).
		Where("node_name = ?", updatedata.NodeName).
		Returning("*").
		Update()
	if err != nil {
		return err
	}
	return nil
}
