package application

import "quinielas/src/movimientos/domain"

type DeleteMovimiento struct {
	repo domain.IMovimiento
}

func NewDeleteMovimiento(repo domain.IMovimiento) *DeleteMovimiento {
	return &DeleteMovimiento{repo: repo}
}

func (dm *DeleteMovimiento) Execute(idMovimiento int32) error {
	return dm.repo.DeleteMovimiento(idMovimiento)
}