package sioc

type InjectorInsertion interface {
	InjectorInsertion() any
}

type Injector[T any] interface {
	GetInstance() T
	GetNewInstance() T
	MatchWithName(name string) bool
	AddInstance(instance T) Injector[T]
}
