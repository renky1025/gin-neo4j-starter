package collection

import (
	"context"
	"fmt"
	"go-gin-restful-service/config"
	"go-gin-restful-service/log"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	MongoClient *mongo.Client
)

const (
	defaultTimeout = 100 * time.Second
	maxPoolSize    = 10
)

// InitMongo
// @Description: 初始化mongo
// @return error
func InitMongo() error {
	var err error
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	credential := options.Credential{
		AuthSource: "admin",
		Username:   cfg.MongoDatabase.UserName,
		Password:   cfg.MongoDatabase.Password,
	}
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			fmt.Println(evt.Command)
		},
	}

	clientOptions := options.Client().ApplyURI(cfg.MongoDatabase.URI).SetAuth(credential).SetMonitor(cmdMonitor).SetMaxPoolSize(maxPoolSize)
	MongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Logger.Error(err)
		return err
	}

	if err := MongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	log.Logger.Infof("Connected to MongoDB!")
	return nil
}

// CreateMongoCollection
// @Description: 创建mongo集合的服务
// @return BaseCollection
func CreateMongoCollection(dbName string, colName string) BaseCollection {
	dataBase := MongoClient.Database(dbName)
	return &BaseCollectionImpl{
		DbName:     dbName,
		ColName:    colName,
		DataBase:   dataBase,
		Collection: dataBase.Collection(colName),
	}
}
