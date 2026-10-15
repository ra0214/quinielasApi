package application

import (
	"fmt"
	"strings"

	movimientos "quinielas/src/movimientos/domain"
	"quinielas/src/shared/money"
)

// CobradorAporte implementa el puerto domain.CobroAporte del módulo de
// aportes: cada aporte de un cliente a una quiniela queda registrado como un
// movimiento PAGO_EFECTIVO (uno por cliente y quiniela) cuyo monto es el total
// acumulado del aporte. Al editar o eliminar el aporte, el movimiento se
// resincroniza o se revierte, de modo que el saldo del cliente refleje siempre
// lo aportado.
type CobradorAporte struct {
	repo   movimientos.IMovimiento
	create *CreateMovimiento
	delete *DeleteMovimiento
}

func NewCobradorAporte(repo movimientos.IMovimiento, create *CreateMovimiento, delete *DeleteMovimiento) *CobradorAporte {
	return &CobradorAporte{repo: repo, create: create, delete: delete}
}

// AplicarAporte registra incrementalmente un aporte: el monto aportado como
// FIADO (deuda) y lo pagado como PAGO_EFECTIVO que la reduce. No revierte los
// movimientos anteriores, de modo que cada aporte se acumula sobre el saldo del
// cliente. Un aporte de contado (monto == pagado) no altera el saldo.
func (c *CobradorAporte) AplicarAporte(idCliente int32, idQuiniela int32, monto money.Money, pagado money.Money, canalVenta string) error {
	if monto.IsZero() || monto.IsNegative() {
		return nil
	}
	if pagado.IsNegative() {
		pagado = money.Zero()
	}
	if pagado.GreaterThan(monto) {
		pagado = monto
	}

	canal := "VENTANILLA"
	if canalVenta != "" {
		canal = canalVenta
	}

	if _, err := c.create.Execute(idCliente, &idQuiniela, "FIADO", monto, fmt.Sprintf("Aporte fiado quiniela #%d (%s)", idQuiniela, canal)); err != nil {
		return err
	}
	if pagado.GreaterThan(money.Zero()) {
		if _, err := c.create.Execute(idCliente, &idQuiniela, "PAGO_EFECTIVO", pagado, fmt.Sprintf("Aporte a quiniela #%d (%s)", idQuiniela, canal)); err != nil {
			return err
		}
	}
	return nil
}

// SincronizarAporte deja los movimientos del aporte acordes al total y a lo
// pagado: PAGO_EFECTIVO por la parte pagada y FIADO por el resto. Es
// idempotente: revierte los movimientos anteriores (si existen) y registra los
// nuevos. canalVenta ('EN_LINEA'/'VENTANILLA') queda en la descripción.
func (c *CobradorAporte) SincronizarAporte(idCliente int32, idQuiniela int32, total money.Money, pagado money.Money, canalVenta string) error {
	if err := c.RevertirAporte(idCliente, idQuiniela); err != nil {
		return err
	}

	if total.IsZero() {
		return nil
	}

	if pagado.IsNegative() {
		pagado = money.Zero()
	}
	if pagado.GreaterThan(total) {
		pagado = total
	}

	canal := "VENTANILLA"
	if canalVenta != "" {
		canal = canalVenta
	}

	// El total del aporte se carga como FIADO (deuda) y lo pagado se aplica como
	// PAGO_EFECTIVO que la reduce. Así un aporte de contado deja el saldo en cero
	// y solo lo realmente fiado queda como deuda.
	if _, err := c.create.Execute(idCliente, &idQuiniela, "FIADO", total, fmt.Sprintf("Aporte fiado quiniela #%d (%s)", idQuiniela, canal)); err != nil {
		return err
	}
	if pagado.GreaterThan(money.Zero()) {
		if _, err := c.create.Execute(idCliente, &idQuiniela, "PAGO_EFECTIVO", pagado, fmt.Sprintf("Aporte a quiniela #%d (%s)", idQuiniela, canal)); err != nil {
			return err
		}
	}
	return nil
}

// RevertirAporte elimina (y revierte en el saldo) el movimiento de aporte del
// cliente en la quiniela.
func (c *CobradorAporte) RevertirAporte(idCliente int32, idQuiniela int32) error {
	movs, err := c.repo.GetMovimientosByQuinielaID(idQuiniela)
	if err != nil {
		return err
	}
	for _, m := range movs {
		if m.IDCliente == idCliente && esAporte(m) {
			if err := c.delete.Execute(m.IDMovimiento); err != nil {
				return err
			}
		}
	}
	return nil
}

func esAporte(m movimientos.Movimiento) bool {
	return strings.HasPrefix(m.Descripcion, "Aporte ")
}