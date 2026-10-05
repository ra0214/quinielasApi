package application

import (
	"quinielas/src/clientes/domain"
)

type DeleteCliente struct {
	repo domain.ICliente
}

func NewDeleteCliente(repo domain.ICliente) *DeleteCliente {
	return &DeleteCliente{repo: repo}
}

func (dc *DeleteCliente) Execute(id int32) error {
	return dc.repo.DeleteCliente(id)
}