package v1_interfaces

type InjectorInsertion interface {
	InjectorInsertion() any
}

type Injector[T any] interface {
	GetInstance() T
	GetNewInstance() T
	MatchWithName(name string) bool
	AddInstance(instance T) Injector[T]
}
