package application

import (
	"fmt"

	movimientos "quinielas/src/movimientos/domain"
	"quinielas/src/shared/money"
)

// CobradorQuiniela implementa el puerto domain.CobroDirecta del módulo de
// quinielas: al crear o editar una quiniela registra los movimientos de cobro
// (PAGO_EFECTIVO por lo pagado y FIADO por el saldo pendiente) para que el
// saldo del cliente quede actualizado; al eliminarla los revierte.
type CobradorQuiniela struct {
	repo   movimientos.IMovimiento
	create *CreateMovimiento
	delete *DeleteMovimiento
}

func NewCobradorQuiniela(repo movimientos.IMovimiento, create *CreateMovimiento, delete *DeleteMovimiento) *CobradorQuiniela {
	return &CobradorQuiniela{repo: repo, create: create, delete: delete}
}

// SincronizarCobro deja los movimientos de la quiniela acordes a total/pagado.
// Es idempotente: primero revierte los movimientos existentes de esa quiniela
// y luego aplica los nuevos. canalVenta ('EN_LINEA'/'VENTANILLA') se refleja en
// la descripción para que quede registrado cómo se realizó la compra.
func (c *CobradorQuiniela) SincronizarCobro(idCliente int32, idQuiniela int32, total money.Money, pagado money.Money, canalVenta string) error {
	if err := c.RevertirCobro(idQuiniela); err != nil {
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

	if pagado.IsZero() {
		_, err := c.create.Execute(idCliente, &idQuiniela, "FIADO", total, fmt.Sprintf("Fiado por compra %s quiniela #%d", canal, idQuiniela))
		return err
	}

	if pagado.Cmp(total) == 0 {
		_, err := c.create.Execute(idCliente, &idQuiniela, "PAGO_EFECTIVO", pagado, fmt.Sprintf("Compra %s quiniela #%d", canal, idQuiniela))
		return err
	}

	// Pago parcial: pago efectivo + resto fiado.
	if _, err := c.create.Execute(idCliente, &idQuiniela, "PAGO_EFECTIVO", pagado, fmt.Sprintf("Pago parcial %s quiniela #%d", canal, idQuiniela)); err != nil {
		return err
	}
	faltante := total.Sub(pagado)
	if faltante.GreaterThan(money.Zero()) {
		if _, err := c.create.Execute(idCliente, &idQuiniela, "FIADO", faltante, fmt.Sprintf("Saldo fiado %s quiniela #%d", canal, idQuiniela)); err != nil {
			return err
		}
	}
	return nil
}

// RevertirCobro elimina (y revierte en el saldo) todos los movimientos de la
// quiniela.
func (c *CobradorQuiniela) RevertirCobro(idQuiniela int32) error {
	movs, err := c.repo.GetMovimientosByQuinielaID(idQuiniela)
	if err != nil {
		return err
	}
	for _, m := range movs {
		if err := c.delete.Execute(m.IDMovimiento); err != nil {
			return err
		}
	}
	return nil
}
