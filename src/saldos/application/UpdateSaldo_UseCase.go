package application

import (
	"quinielas/src/saldos/domain"
	"quinielas/src/shared/money"
)

type UpdateSaldo struct {
	repo domain.ISaldo
}

func NewUpdateSaldo(repo domain.ISaldo) *UpdateSaldo {
	return &UpdateSaldo{repo: repo}
}

func (us *UpdateSaldo) Execute(idCliente int32, saldoFavor money.Money, saldoDeuda money.Money) error {
	favor, deuda := domain.Compensar(saldoFavor, saldoDeuda)
	return us.repo.UpdateSaldoCliente(idCliente, favor, deuda)
}