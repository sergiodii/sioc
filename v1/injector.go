package sioc

import (
	"reflect"

	"github.com/sergiodii/sioc/extension/text"
)

type injector[T any] struct {
	instance      T
	injectionName string
}

func NewInjector[T any]() Injector[T] {
	return &injector[T]{}
}

func (i injector[T]) transformName(name string) string {
	return text.Sanitize(name)
}

func (i *injector[T]) AddInstance(instance T) Injector[T] {
	i.injectionName = i.transformName(reflect.ValueOf(instance).Type().String())
	i.instance = instance
	return i
}

func (i injector[T]) GetInstance() T {
	return i.instance
}

func (i injector[T]) GetNewInstance() T {
	n := new(T)
	*n = i.instance
	return *n
}

func (i injector[T]) MatchWithName(name string) bool {
	return i.injectionName == i.transformName(name)
}
