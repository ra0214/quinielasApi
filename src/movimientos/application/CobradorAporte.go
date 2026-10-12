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

// SincronizarAporte deja el movimiento de aporte del cliente acorde al monto
// total acumulado. Es idempotente: revierte el movimiento anterior (si existe)
// y registra el nuevo. canalVenta ('EN_LINEA'/'VENTANILLA') queda en la
// descripción para que el registro sea completo.
func (c *CobradorAporte) SincronizarAporte(idCliente int32, idQuiniela int32, monto money.Money, canalVenta string) error {
	if err := c.RevertirAporte(idCliente, idQuiniela); err != nil {
		return err
	}

	if monto.IsZero() {
		return nil
	}

	canal := "VENTANILLA"
	if canalVenta != "" {
		canal = canalVenta
	}

	_, err := c.create.Execute(idCliente, &idQuiniela, "PAGO_EFECTIVO", monto, fmt.Sprintf("Aporte a quiniela #%d (%s)", idQuiniela, canal))
	return err
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
	return strings.HasPrefix(m.Descripcion, "Aporte a quiniela ")
}