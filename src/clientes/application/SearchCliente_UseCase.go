package application

import (
	"quinielas/src/clientes/domain"
)

type SearchCliente struct {
	repo domain.ICliente
}

func NewSearchCliente(repo domain.ICliente) *SearchCliente {
	return &SearchCliente{repo: repo}
}

func (sc *SearchCliente) Execute(query string) ([]domain.Cliente, error) {
	return sc.repo.SearchClientesByName(query)
}