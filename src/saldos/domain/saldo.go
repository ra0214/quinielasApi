package domain

import (
	"errors"
	"quinielas/src/shared/money"
)

var ErrSaldoInsuficiente = errors.New("saldo insuficiente a favor para realizar el retiro u operación")

type ISaldo interface {
	GetSaldoByClienteID(idCliente int32) (*Saldo, error)
	GetAllSaldos() ([]Saldo, error)
	GetSaldosFiltrados(filtro string) ([]Saldo, error) // 'deben', 'favor', 'ceros'
	InitSaldoCliente(idCliente int32) error
	UpdateSaldoCliente(idCliente int32, saldoFavor money.Money, saldoDeuda money.Money) error
}

type Saldo struct {
	IDCliente  int32       `json:"id_cliente"`
	SaldoFavor money.Money `json:"saldo_favor"`
	SaldoDeuda money.Money `json:"saldo_deuda"`
}

func NewSaldo(idCliente int32) *Saldo {
	return &Saldo{
		IDCliente:  idCliente,
		SaldoFavor: money.Zero(),
		SaldoDeuda: money.Zero(),
	}
}
