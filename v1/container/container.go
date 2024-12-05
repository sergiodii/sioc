package v1_container

import (
	"sync"

	v1_utils "github.com/sergiodii/sIOC/v1/utils"
)

type Container interface {
	Set(instanceType string, instance any)
	Get(interfaceType string) (any, bool)
	GetAll() []any
	Len() int
}

type container struct {
	mapList sync.Map
}

func New() Container {
	return &container{}
}

func (c *container) Set(instanceType string, instance any) {
	c.mapList.Store(v1_utils.SanitizeName(instanceType), instance)
}

func (c *container) Get(instanceType string) (any, bool) {
	return c.mapList.Load(v1_utils.SanitizeName(instanceType))
}

func (c *container) GetAll() []any {
	var list []any
	c.mapList.Range(func(key, value any) bool {
		list = append(list, value)
		return true
	})
	return list
}

func (c *container) Len() int {
	var count int
	c.mapList.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}
