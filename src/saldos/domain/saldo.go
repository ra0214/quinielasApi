package domain

import (
	"errors"
	"fmt"

	"quinielas/src/shared/money"
)

var ErrSaldoInsuficiente = errors.New("saldo insuficiente a favor para realizar el retiro u operación")

// SaldoPendienteError indica que el cliente no puede borrarse porque tiene
// saldo a favor o deuda pendiente.
type SaldoPendienteError struct {
	Favor money.Money
	Deuda money.Money
}

func (e *SaldoPendienteError) Error() string {
	switch {
	case !e.Favor.IsZero() && !e.Deuda.IsZero():
		return "tiene saldo a favor y deuda pendiente"
	case !e.Deuda.IsZero():
		return fmt.Sprintf("tiene una deuda de $%s", e.Deuda.String())
	default:
		return fmt.Sprintf("tiene saldo a favor de $%s", e.Favor.String())
	}
}

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

// Neto es el saldo neto del cliente: favor menos deuda (puede ser negativo).
func (s *Saldo) Neto() money.Money {
	return s.SaldoFavor.Sub(s.SaldoDeuda)
}

// Compensar reduce el par (favor, deuda) a su valor neto: nunca quedan ambos
// mayores a 0. Si el favor es mayor o igual a la deuda, la deuda se cancela.
func Compensar(favor, deuda money.Money) (money.Money, money.Money) {
	return Descomponer(favor.Sub(deuda))
}

// Descomponer convierte un neto en el par (favor, deuda): con neto >= 0 queda
// todo en favor; con neto < 0 queda todo en deuda.
func Descomponer(neto money.Money) (money.Money, money.Money) {
	if neto.IsNegative() {
		return money.Zero(), neto.Neg()
	}
	return neto, money.Zero()
}

// AplicarMovimiento calcula el nuevo saldo compensado tras registrar un
// movimiento. Reglas por tipo:
//
//	PAGO_EFECTIVO  -> acredita al neto (paga deuda; el sobrante pasa a favor)
//	PAGO_SALDO     -> consume favor (error si no alcanza)
//	FIADO          -> descuenta del neto (aumenta la deuda)
//	RETIRO         -> consume favor (error si no alcanza)
//	PREMIO_ABONO   -> acredita al neto (suma a favor)
//	DEVOLUCION     -> acredita al neto (suma a favor)
func AplicarMovimiento(favor, deuda money.Money, tipo string, monto money.Money) (money.Money, money.Money, error) {
	neto := favor.Sub(deuda)
	switch tipo {
	case "PAGO_EFECTIVO", "PREMIO_ABONO", "DEVOLUCION":
		neto = neto.Add(monto)
	case "FIADO":
		neto = neto.Sub(monto)
	case "PAGO_SALDO", "RETIRO":
		if neto.Cmp(monto) < 0 {
			return favor, deuda, ErrSaldoInsuficiente
		}
		neto = neto.Sub(monto)
	default:
		return favor, deuda, fmt.Errorf("tipo de movimiento desconocido: %s", tipo)
	}
	f, d := Descomponer(neto)
	return f, d, nil
}

// RevertirMovimiento revierte el efecto del movimiento sobre el saldo neto.
func RevertirMovimiento(favor, deuda money.Money, tipo string, monto money.Money) (money.Money, money.Money, error) {
	neto := favor.Sub(deuda)
	switch tipo {
	case "PAGO_EFECTIVO", "PREMIO_ABONO", "DEVOLUCION":
		neto = neto.Sub(monto)
	default:
		neto = neto.Add(monto)
	}
	f, d := Descomponer(neto)
	return f, d, nil
}
