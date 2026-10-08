package application

import (
	"quinielas/src/saldos/domain"
	"quinielas/src/shared/money"
)

// AjustarSaldo aplica (o revierte) el efecto de un movimiento sobre el saldo
// del cliente y persiste el resultado ya compensado. Lo usa el módulo de
// movimientos a traves del puerto SaldoAjustador.
type AjustarSaldo struct {
	repo domain.ISaldo
}

func NewAjustarSaldo(repo domain.ISaldo) *AjustarSaldo {
	return &AjustarSaldo{repo: repo}
}

// Ajustar aplica el efecto del movimiento al saldo actual y guarda el
// resultado compensado (nunca quedan favor y deuda positivos a la vez).
func (as *AjustarSaldo) Ajustar(idCliente int32, tipo string, monto money.Money) error {
	saldo, err := as.repo.GetSaldoByClienteID(idCliente)
	if err != nil {
		return err
	}

	favor, deuda, err := domain.AplicarMovimiento(saldo.SaldoFavor, saldo.SaldoDeuda, tipo, monto)
	if err != nil {
		return err
	}

	return as.repo.UpdateSaldoCliente(idCliente, favor, deuda)
}

// Revertir deshace el efecto del movimiento sobre el saldo actual y guarda el
// resultado compensado.
func (as *AjustarSaldo) Revertir(idCliente int32, tipo string, monto money.Money) error {
	saldo, err := as.repo.GetSaldoByClienteID(idCliente)
	if err != nil {
		return err
	}

	favor, deuda, err := domain.RevertirMovimiento(saldo.SaldoFavor, saldo.SaldoDeuda, tipo, monto)
	if err != nil {
		return err
	}

	return as.repo.UpdateSaldoCliente(idCliente, favor, deuda)
}