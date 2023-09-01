package controller

import (
	"go-gin-restful-service/database/neo4jdb"
	"go-gin-restful-service/dto"
	"go-gin-restful-service/model"
	"go-gin-restful-service/response"

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

func (ctr *PersonController) CreateNewNode(ctx *gin.Context) {
	var persion dto.CreateNodeDTO
	if err := ctx.ShouldBindJSON(&persion); err != nil {
		response.Fail(ctx, response.ParamsValidError)
		return
	}

	nj := *ctr.Neo4j
	res, err := nj.CreateNode(persion.Label, persion.NodeAttr)
	if err != nil {
		response.FailWithMsg(ctx, response.Failed, err.Error())
		return
	}
	response.OkWithData(ctx, res)
}

func (ctr *PersonController) UpdateNodeBy(ctx *gin.Context) {
	name := ctx.Param("name")
	label := ctx.Param("label")
	var payload map[string]interface{}
	if len(name) == 0 || len(label) == 0 {
		response.Fail(ctx, response.ParamsValidError)
		return
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response.Fail(ctx, response.ParamsValidError)
		return
	}

	nj := *ctr.Neo4j

	res, err := nj.UpdateNodeAttr(model.NodeFrom{Label: label, Name: name}, payload)
	if err != nil {
		response.FailWithMsg(ctx, response.Failed, err.Error())
		return
	}
	response.OkWithData(ctx, res)
}

func (ctr *PersonController) GetNodeBy(ctx *gin.Context) {
	name := ctx.Param("name")
	label := ctx.Param("label")
	if len(name) == 0 || len(label) == 0 {
		response.Fail(ctx, response.ParamsValidError)
		return
	}
	nj := *ctr.Neo4j
	res, err := nj.GetNodeBy(model.NodeFrom{Label: label, Name: name})
	if err != nil {
		response.FailWithMsg(ctx, response.Failed, err.Error())
		return
	}
	response.OkWithData(ctx, res)
}

func (ctr *PersonController) CreateRelationShip(ctx *gin.Context) {
	var relation dto.RelationShipDTO
	if err := ctx.ShouldBindJSON(&relation); err != nil {
		response.Fail(ctx, response.ParamsValidError)
		return
	}

	nj := *ctr.Neo4j
	err := nj.CreateRelationship(relation.Node1, relation.Node2, relation.RelationShip)
	if err != nil {
		response.FailWithMsg(ctx, response.Failed, err.Error())
		return
	}
	response.OkWithMsg(ctx, "ok")
}

func (ctr *PersonController) DelteNodeBy(ctx *gin.Context) {
	name := ctx.Param("name")
	label := ctx.Param("label")
	if len(name) == 0 || len(label) == 0 {
		response.Fail(ctx, response.ParamsValidError)
		return
	}
	nj := *ctr.Neo4j
	res, err := nj.DeleteNodeByLabel(label, name)
	if err != nil {
		response.FailWithMsg(ctx, response.Failed, err.Error())
		return
	}
	response.OkWithData(ctx, res)
}

func (ctr *PersonController) CountNodeBy(ctx *gin.Context) {
	label := ctx.Query("label")
	nj := *ctr.Neo4j
	res, err := nj.CountNodes(label)
	if err != nil {
		response.FailWithMsg(ctx, response.Failed, err.Error())
		return
	}
	response.OkWithData(ctx, res)
}

func (ctr *PersonController) SearchNodes(ctx *gin.Context) {
	var params dto.QueryParamsDTO
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.Fail(ctx, response.ParamsValidError)
		return
	}
	nj := *ctr.Neo4j
	res, err := nj.SearchNodes(params)
	if err != nil {
		response.FailWithMsg(ctx, response.Failed, err.Error())
		return
	}
	response.OkWithData(ctx, res)
}
