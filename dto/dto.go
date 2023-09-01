package dto

import "go-gin-restful-service/model"

type RelationShipDTO struct {
	Node1        model.NodeFrom    `json:"node1" from:"node1" binding:"required"`
	Node2        model.NodeFrom    `json:"node2" from:"node2" binding:"required"`
	RelationShip model.RelationMap `json:"relationShip" from:"relationShip" binding:"required"`
}

type CreateNodeDTO struct {
	Label    string                 `json:"label" from:"label" binding:"required"`
	NodeAttr map[string]interface{} `json:"nodeAttr" from:"nodeAttr" binding:"required"`
}

type RelationShipQuery struct {
	Name   string                 `json:"name" from:"name"`
	Number int                    `json:"number" from:"number"`
	Prop   map[string]interface{} `json:"prop" from:"prop"`
}

type QueryParamsDTO struct {
	NodeFrom     *model.NodeFrom    `json:"nodeFrom" from:"nodeFrom"`
	NodeTo       *model.NodeFrom    `json:"nodeTo" from:"nodeTo"`
	RelationShip *RelationShipQuery `json:"relationShip" from:"relationShip"`
	PageSize     int64              `json:"pageSize" from:"pageSize" binding:"required"`
	PageNo       int64              `json:"pageNo" from:"pageNo" binding:"required"`
	Order        *string            `json:"order" from:"order"`
	Sort         *string            `json:"sort" from:"sort"`
}
