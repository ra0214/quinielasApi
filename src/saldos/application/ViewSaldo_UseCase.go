package application

import "quinielas/src/saldos/domain"

type ViewSaldo struct {
	repo domain.ISaldo
}

func NewViewSaldo(repo domain.ISaldo) *ViewSaldo {
	return &ViewSaldo{repo: repo}
}

func (vs *ViewSaldo) ExecuteGetAll() ([]domain.Saldo, error) {
	return vs.repo.GetAllSaldos()
}

func (vs *ViewSaldo) ExecuteGetByClienteID(idCliente int32) (*domain.Saldo, error) {
	return vs.repo.GetSaldoByClienteID(idCliente)
}

func (vs *ViewSaldo) ExecuteFiltrados(filtro string) ([]domain.Saldo, error) {
	return vs.repo.GetSaldosFiltrados(filtro)
}