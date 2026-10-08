package application

import (
	"quinielas/src/movimientos/domain"
)

type DeleteMovimiento struct {
	repo      domain.IMovimiento
	ajustador SaldoAjustador
}

func NewDeleteMovimiento(repo domain.IMovimiento, ajustador SaldoAjustador) *DeleteMovimiento {
	return &DeleteMovimiento{repo: repo, ajustador: ajustador}
}

func (dm *DeleteMovimiento) Execute(idMovimiento int32) error {
	movimiento, err := dm.repo.GetMovimientoByID(idMovimiento)
	if err != nil {
		return err
	}

	// Revertimos primero el impacto sobre el saldo y eliminamos el movimiento.
	if err := dm.ajustador.Revertir(movimiento.IDCliente, movimiento.Tipo, movimiento.Monto); err != nil {
		return err
	}

	if err := dm.repo.DeleteMovimiento(idMovimiento); err != nil {
		// Reaplicamos el ajuste para dejar el saldo consistente si el borrado falla.
		_ = dm.ajustador.Ajustar(movimiento.IDCliente, movimiento.Tipo, movimiento.Monto)
		return err
	}

	return nil
}