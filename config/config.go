package config

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Neo4J         Neo4J         `json:"neo4j" yaml:"neo4j"`
	ServerPort    int           `json:"server_port" yaml:"server_port"`
	MongoDatabase MongoDatabase `json:"database" yaml:"database"`
	AwsConfig     AwsConfig     `json:"aws" yaml:"aws"`
	Redis         RedisConf     `json:"redis" yaml:"redis"`
}
type RedisConf struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Prefix   string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
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
type AwsConfig struct {
	Endpoint   string `json:"endpoint"`
	AccessKey  string `json:"access-key-id"`
	SecretKey  string `json:"access-key-secret"`
	BucketName string `json:"bucket-name"`
	BucketUrl  string `json:"bucket-url"`
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
