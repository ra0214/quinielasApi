package application

import (
	"quinielas/src/movimientos/domain"
	"quinielas/src/shared/money"
)

type CreateMovimiento struct {
	repo domain.IMovimiento
}

func NewCreateMovimiento(repo domain.IMovimiento) *CreateMovimiento {
	return &CreateMovimiento{repo: repo}
}

func (cm *CreateMovimiento) Execute(idCliente int32, idQuiniela *int32, tipo string, monto money.Money, descripcion string) (*domain.Movimiento, error) {
	return cm.repo.SaveMovimiento(idCliente, idQuiniela, tipo, monto, descripcion)
}