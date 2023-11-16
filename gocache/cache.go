package gocache

import (
	"errors"
)

type Cache interface {
	Get(key string) interface{}
	Set(key string, value interface{}, timeSec int) bool
	Delete(key string) error
	Clear() error
}

func NewCache(tp string) (Cache, error) {
	switch tp {
	case "redis":
		return NewRedisCache(), nil
	default:
		return NewMemoryCache(), nil
	}
	return nil, errors.New("can not found target cache")
}
