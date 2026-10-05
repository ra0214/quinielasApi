package application

import (
	"quinielas/src/clientes/domain"
)

type ViewCliente struct {
	repo domain.ICliente
}

func NewViewCliente(repo domain.ICliente) *ViewCliente {
	return &ViewCliente{repo: repo}
}

func (vc *ViewCliente) Execute() ([]domain.Cliente, error) {
	return vc.repo.GetAll()
}

func (vc *ViewCliente) ExecuteByID(id int32) (*domain.Cliente, error) {
	return vc.repo.GetClienteByID(id)
}