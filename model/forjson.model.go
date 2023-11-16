package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// 需要把对象保存入 postgresql json字段，需要添加下面 hookup function,保证 struct->json, json->struct
type mapObject map[string]interface{}

func (c mapObject) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *mapObject) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

type arrayString []string

func (c arrayString) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *arrayString) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}
