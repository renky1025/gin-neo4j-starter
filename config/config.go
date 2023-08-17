package config

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Neo4J         Neo4J         `json:"neo4j"`
	ServerPort    int           `json:"server_port"`
	MongoDatabase MongoDatabase `json:"database"`
}
type Neo4J struct {
	URI      string `json:"uri"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
type MongoDatabase struct {
	URI      string `json:"uri"`
	DbName   string `json:"dbname"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

var c *Config

// InitConfig ...
func InitConfig() {
	vp := viper.New()
	basePath, _ := os.Getwd()
	vp.SetConfigName("config")
	vp.AddConfigPath(basePath + "/config")
	err := vp.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = vp.Unmarshal(&c, func(config *mapstructure.DecoderConfig) {
		config.TagName = "json"
	})
	if err != nil {
		panic(err)
	}
}

// GetConfig ...
func GetConfig() *Config {
	if c == nil {
		InitConfig()
	}
	return c
}
