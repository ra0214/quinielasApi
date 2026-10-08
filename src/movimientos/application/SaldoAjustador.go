package application

import "quinielas/src/shared/money"

// SaldoAjustador ajusta el saldo del cliente al registrar o eliminar
// movimientos. Lo implementa el módulo de saldos; los movimientos solo lo
// consumen para no acoplarse a su infraestructura.
type SaldoAjustador interface {
	// Ajustar aplica el efecto del movimiento sobre el saldo (FIADO suma deuda,
	// PAGO_EFECTIVO paga deuda, PAGO_SALDO/RETIRO consumen favor, etc.).
	Ajustar(idCliente int32, tipo string, monto money.Money) error
	// Revertir deshace el efecto del movimiento sobre el saldo.
	Revertir(idCliente int32, tipo string, monto money.Money) error
}