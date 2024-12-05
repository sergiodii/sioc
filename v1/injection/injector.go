package v1_injection

import (
	"reflect"

	v1_interfaces "github.com/sergiodii/sIOC/v1/interfaces"
	v1_utils "github.com/sergiodii/sIOC/v1/utils"
)

type Injector[T any] struct {
	instance      T
	incjetionName string
}

func NewInjector[T any]() v1_interfaces.Injector[T] {
	return &Injector[T]{}
}

func (i Injector[T]) transformName(name string) string {
	return v1_utils.SanitizeName(name)
}

func (i *Injector[T]) AddInstance(instance T) v1_interfaces.Injector[T] {
	i.incjetionName = i.transformName(reflect.ValueOf(instance).Type().String())
	i.instance = instance
	return i
}

func (i Injector[T]) GetInstance() T {
	return i.instance
}

func (i Injector[T]) GetNewInstance() T {
	n := new(T)
	*n = i.instance
	return *n
}

func (i Injector[T]) MatchWithName(name string) bool {
	return i.incjetionName == i.transformName(name)
}
