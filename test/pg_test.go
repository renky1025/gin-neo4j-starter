package test

import (
	"fmt"
	"go-gin-restful-service/database/postgres"
	"testing"
	"time"
)

// user_info表
/*
    create table user_info
(
    user_id                 varchar(32) not null
        constraint user_info_pkey
            primary key,
    user_name               varchar(32) not null,
    age                     smallint not null,
    updated                 timestamp default now(),
    created                 timestamp default now()
);
*/

type UserInfoModel struct {
	tableName struct{}  `pg:"user_info"`
	UserID    string    `json:"user_id" pg:"user_id,pk"`  // 使用ORM建表后是text类型
	UserName  string    `json:"user_name" pg:"user_name"` // 使用ORM建表后是text类型
	Age       int8      `json:"age" pg:"age"`             // 使用ORM建表后是smallint类型
	Updated   time.Time `json:"updated" pg:"updated"`     // 使用ORM建表后不设置default！！！
	Created   time.Time `json:"created" pg:"created"`     // 使用ORM建表后不设置default！！！
}

func (u *UserInfoModel) TableName() string {
	return "user_info"
}

const (
	// 数据库链接及链接池相关配置
	pgURL              = "postgresql://postgres:'123456'@192.168.1.104:5432/postgres?sslmode=disable"
	connMaxLifetimeSec = 3600
	maxOpenConns       = 100
	minIdleConns       = 10
	maxRetries         = 5
)

// 测试数据库链接 &
func TestDBInit(t *testing.T) {
	// ==================================== 初始化数据库 ====================================================
	dbCfg := postgres.DatabaseConfig{
		URL:                    pgURL,
		ConnMaxLifetimeSeconds: connMaxLifetimeSec,
		MaxOpenConns:           maxOpenConns,
		MinIdleConns:           minIdleConns,
		MaxRetries:             maxRetries,
	}
	// init db
	db, dbVersion, err := postgres.InitDBService(&dbCfg)
	if err != nil {
		fmt.Println("init DB Error: ", err)
		panic(err)
	}
	// close db
	defer db.Close()
	fmt.Println("db: ", db, "dbVersion: ", dbVersion)

	dbService := postgres.DBService{
		DB: db,
	}

	dbService.CreateTable()
	dbService.InsertNodeInfo()
	dbService.SelectAllNodeInfo()
	dbService.SelectNodeInfo()
	dbService.UpdateNodeInfo()
	dbService.DeleteTable()
	// ==================================== 根据model创建表/删除表(实际建议使用原生SQL做) =========================
	// // =========== 创建表 ===========
	// // TODO 注意，使用ORM创建的表，字符串类型都是 text类型！如果是int的话默认是bigint类型！！！实际上还是建议根据实际使用原生SQL创建表！！
	// models := []interface{}{
	// 	// 声明类型～～
	// 	(*UserInfoModel)(nil),
	// }
	// // =========== 删除表 ===========
	// for _, model := range models {
	// 	err = db.Model(model).DropTable(&orm.DropTableOptions{
	// 		IfExists: true,
	// 		Cascade:  true,
	// 	})
	// 	if err != nil {
	// 		fmt.Println("删除表 ", model, "raise error: ", err)
	// 		panic(err)
	// 	}
	// }
	// // =========== 创建表 ===========
	// for _, model := range models {
	// 	// TODO 如果是一张表的话，里面的model可以用 &UserInfoModel{} 代替～
	// 	err = db.Model(model).CreateTable(&orm.CreateTableOptions{
	// 		IfNotExists: true,
	// 	})
	// 	if err != nil {
	// 		fmt.Println("创建表 ", model, "raise error: ", err)
	// 		panic(err)
	// 	}
	// }

	// // ==================================== insert 操作 ================================================
	// // =========== 插入单条数据 ===========
	// user1 := UserInfoModel{
	// 	UserID:   "xxx1232",
	// 	UserName: "whw1",
	// 	Age:      22,
	// }
	// result, err := db.Model(&user1).Insert()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("single insert rows affected:%d \n", result.RowsAffected()) // single insert rows affected:1

	// // =========== 插入多条数据 ===========
	// userList := []UserInfoModel{
	// 	{
	// 		UserID:   "xxdsx555",
	// 		UserName: "whw12",
	// 		Age:      22,
	// 	},
	// 	{
	// 		UserID:   "xxdsx5awd55",
	// 		UserName: "whw22",
	// 		Age:      23,
	// 	},
	// }
	// result, err = db.Model(&userList).Insert()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("batch insert rows affected:%d \n", result.RowsAffected()) // batch insert rows affected:2

	// =========== 执行原生SQL插入数据 避免主键冲突 ===========

}
