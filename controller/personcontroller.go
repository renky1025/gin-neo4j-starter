package controller

import (
	"go-gin-restful-service/database/neo4jdb"
	"go-gin-restful-service/dto"
	"go-gin-restful-service/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PersonController struct {
	Neo4j *neo4jdb.Neo4jDriver
}

func NewPersonController(n4j *neo4jdb.Neo4jDriver) *PersonController {
	return &PersonController{
		Neo4j: n4j,
	}
}

func (ctr *PersonController) CreatePerson(ctx *gin.Context) {
	var persion neo4jdb.Person
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

func (ctr *PersonController) CreateRelationShip(ctx *gin.Context) {
	var relation dto.RelationShipDTO
	if err := ctx.ShouldBindJSON(&relation); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, err.Error())
		return
	}

	nj := *ctr.Neo4j
	err := nj.CreateRelationship(relation.Node1, relation.Node2)
	if err != nil {
		response.FailWithMsg(ctx, response.Failed, err.Error())
		return
	}
	response.OkWithMsg(ctx, "ok")
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

func (ctr *PersonController) SearchPerson(ctx *gin.Context) {
	pageNo := strings.TrimSpace(ctx.Query("pageNo"))
	pageSize := strings.TrimSpace(ctx.Query("pageSize"))
	name := strings.TrimSpace(ctx.Query("name"))
	nj := *ctr.Neo4j
	pn, _ := strconv.ParseInt(pageNo, 10, 64)
	ps, _ := strconv.ParseInt(pageSize, 10, 64)
	res, err := nj.SearchPerson(name, (pn-1)*ps, ps)
	if err != nil {
		response.FailWithMsg(ctx, response.Failed, err.Error())
		return
	}
	response.OkWithData(ctx, res)
}
