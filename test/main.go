package main

import (
	"fmt"

	"github.com/sergiodii/sioc/v1"
)

type Sentimento interface {
	Falar()
}

type Amor struct {
	Alvo string
}

func (a *Amor) Shoot() string {
	return "Tiro de amor em " + a.Alvo
}

type Paixao struct {
	a *Amor
}

func (p *Paixao) Init(a *Amor) {
	println("Paixão inicializada")
	p.a = a
}

func (p *Paixao) Falar() {
	fmt.Println(p.a.Shoot())
}

func main() {
	println("Iniciando o contêiner de injeção de dependência...")

	container := sioc.NewContainer()

	// Registrando serviços
	sioc.Inject(&Paixao{}, container)
	sioc.Inject(&Amor{}, container)

	sioc.Init(container)

	p := sioc.Get[*Paixao](container)

	p.Falar()

	a := sioc.Get[*Amor](container)

	a.Alvo = "Coração"

	p.Falar()

	s := sioc.Get[Sentimento](container)

	s.Falar()

}
