package controller

import (
	"go-gin-restful-service/database"
	"go-gin-restful-service/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PersonController struct {
	Neo4j *database.Neo4jDriver
}

func NewPersonController(n4j *database.Neo4jDriver) *PersonController {
	return &PersonController{
		Neo4j: n4j,
	}
}

func (ctr *PersonController) CreatePerson(ctx *gin.Context) {
	var persion database.Person
	if err := ctx.ShouldBindJSON(&persion); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, err.Error())
		return
	}

	nj := *ctr.Neo4j
	res, err := nj.CreatePerson(persion.Name, persion.Age)
	if err != nil {
		response.FailWithMsg(ctx, response.Failed, err.Error())
		return
	}
	response.OkWithData(ctx, res)
}

func (ctr *PersonController) GetPersonBy(ctx *gin.Context) {
	pid := ctx.Param("pid")
	if len(pid) == 0 {
		response.FailWithMsg(ctx, response.ParamsValidError, "pid not allow be nil")
		return
	}

	nj := *ctr.Neo4j
	id, _ := strconv.ParseInt(pid, 10, 64)
	res, err := nj.GetPersonById(id)
	if err != nil {
		response.FailWithMsg(ctx, response.Failed, err.Error())
		return
	}
	response.OkWithData(ctx, res)
}
