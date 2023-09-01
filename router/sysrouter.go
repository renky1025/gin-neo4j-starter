package sysrouter

import (
	"context"
	"go-gin-restful-service/config"
	"go-gin-restful-service/controller"
	"go-gin-restful-service/database/neo4jdb"
	"go-gin-restful-service/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(ctx context.Context, cfg *config.Config, router *gin.Engine) *gin.Engine {
	n4j := neo4jdb.NewNeo4jDriver(cfg)
	apiRouter := router.Group("/api/v1.0")
	apiRouter.Use(
		middleware.Cors(),
		middleware.TokenAuth(),
		middleware.GinI18nLocalize(),
		middleware.ErrorRecover())
	// middleware
	// routers
	// groups
	personController := controller.NewPersonController(n4j)
	apiRouter.POST("/create", personController.CreateNewNode)
	apiRouter.POST("/update/:label/:name", personController.UpdateNodeBy)
	apiRouter.POST("/relation", personController.CreateRelationShip)
	apiRouter.GET("/view/:label/:name", personController.GetNodeBy)
	apiRouter.GET("/count", personController.CountNodeBy)
	apiRouter.GET("/del/:label/:name", personController.DelteNodeBy)
	apiRouter.POST("/search", personController.SearchNodes)
	return router
}
