package application

import "quinielas/src/movimientos/domain"

type ViewMovimiento struct {
	repo domain.IMovimiento
}

func NewViewMovimiento(repo domain.IMovimiento) *ViewMovimiento {
	return &ViewMovimiento{repo: repo}
}

func (vm *ViewMovimiento) ExecuteByClienteID(idCliente int32) ([]domain.Movimiento, error) {
	return vm.repo.GetMovimientosByClienteID(idCliente)
}

func (vm *ViewMovimiento) ExecuteGetAll() ([]domain.Movimiento, error) {
	return vm.repo.GetAllMovimientos()
}