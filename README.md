# Starter for Gin + Neo4j

**Dependencies used:**
 * using [go 1.19](https://tip.golang.org/doc/go1.19)
 * using [gin-gonic](https://github.com/gin-gonic/gin#gin-web-framework) web framework
 * using [viper](https://github.com/spf13/viper) as a configuration solution
 * using [swag](https://github.com/swaggo/swag) as API Documents
 * using [neo4j](https://github.com/neo4j/neo4j-go-driver/v5) as a graph database


## How to start the application
### start neo4j at local
```shell
# or via docker-compose
docker-compose up -d
```
### start go 
```shell
go run main.go
```

## the structures of files
```shell
├─config     # project configs
├─controller # project controller
├─database   # database conn
├─dto        # tranlate data 
├─i18n
├─log
├─middleware
├─model
├─response
├─router
├─service
└─util
```
## update config in file ./config/config.yaml