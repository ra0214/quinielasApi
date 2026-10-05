package application

import (
	"quinielas/src/clientes/domain"
)

type CreateCliente struct {
	repo domain.ICliente
}

func NewCreateCliente(repo domain.ICliente) *CreateCliente {
	return &CreateCliente{repo: repo}
}

func (cc *CreateCliente) Execute(nombre string, telefono string) (*domain.Cliente, error) {
	return cc.repo.SaveCliente(nombre, telefono)
}