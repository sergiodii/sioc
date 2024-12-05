package v0_injection

import (
	"reflect"
	"strings"
)

type Injector[T any] struct {
	instance      T
	incjetionName string
	initialized   bool
}

func NewInjector[T any]() Injector[T] {
	return Injector[T]{}
}

func (i Injector[T]) transformName(name string) string {

	return strings.ReplaceAll(name, "*", "")
}

func (i *Injector[T]) AddInstance(instance T) {
	i.incjetionName = i.transformName(reflect.ValueOf(instance).Type().String())
	i.instance = instance
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

func (i Injector[T]) Initialized() bool {
	return i.initialized
}

func (i *Injector[T]) SetInitialization() {
	i.initialized = true
}
