package application

import (
	"quinielas/src/clientes/domain"
)

type UpdateCliente struct {
	repo domain.ICliente
}

func NewUpdateCliente(repo domain.ICliente) *UpdateCliente {
	return &UpdateCliente{repo: repo}
}

func (uc *UpdateCliente) Execute(id int32, nombre string, telefono string) error {
	return uc.repo.UpdateCliente(id, nombre, telefono)
}