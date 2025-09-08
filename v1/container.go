package sioc

import (
	"sync"

	"github.com/sergiodii/sioc/extension/text"
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
	c.mapList.Store(text.Sanitize(instanceType), instance)
}

func (c *container) Get(instanceType string) (any, bool) {
	return c.mapList.Load(text.Sanitize(instanceType))
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
