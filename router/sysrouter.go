package sysrouter

import (
	"context"
	"go-gin-restful-service/config"
	"go-gin-restful-service/controller"
	"go-gin-restful-service/database"
	"go-gin-restful-service/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(ctx context.Context, cfg *config.Config, router *gin.Engine) *gin.Engine {
	n4j := database.NewNeo4jDriver(cfg)
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
	apiRouter.POST("/create", personController.CreatePerson)
	apiRouter.POST("/relation", personController.CreateRelationShip)
	apiRouter.GET("/person/:pid", personController.GetPersonBy)
	apiRouter.GET("/search", personController.SearchPerson)
	return router
}
