package application

import (
	"fmt"
	"time"

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

	// El total se carga como FIADO (deuda) y lo pagado se aplica como
	// PAGO_EFECTIVO que la reduce. Así un pago de contado deja el saldo en cero
	// y solo lo realmente fiado queda como deuda.
	if _, err := c.create.Execute(idCliente, &idQuiniela, "FIADO", total, fmt.Sprintf("Compra fiada %s quiniela #%d", canal, idQuiniela)); err != nil {
		return err
	}
	if pagado.GreaterThan(money.Zero()) {
		if _, err := c.create.Execute(idCliente, &idQuiniela, "PAGO_EFECTIVO", pagado, fmt.Sprintf("Compra %s quiniela #%d", canal, idQuiniela)); err != nil {
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

// ReembolsarCobro elimina los movimientos de la quiniela devolviendo a saldo a
// favor lo que el cliente pagó y descartando lo fiado: los PAGO_EFECTIVO se
// borran conservando su efecto en el saldo (queda como favor) y el resto de
// movimientos se revierten normalmente.
func (c *CobradorQuiniela) ReembolsarCobro(idQuiniela int32) error {
	movimientos, err := c.repo.GetMovimientosByQuinielaID(idQuiniela)
	if err != nil {
		return err
	}

	for _, m := range movimientos {
		if m.Tipo == "PAGO_EFECTIVO" {
			if err := c.repo.DeleteMovimiento(m.IDMovimiento); err != nil {
				return err
			}
			continue
		}
		if err := c.delete.Execute(m.IDMovimiento); err != nil {
			return err
		}
	}

	return nil
}

// ReembolsarCobroDesde hace lo mismo que ReembolsarCobro pero solo con los
// movimientos posteriores a `desde`, y devuelve por cliente el total aportado
// en ese tramo (el FIADO de cada compra), que es el excedente a descontar.
func (c *CobradorQuiniela) ReembolsarCobroDesde(idQuiniela int32, desde time.Time) (map[int32]money.Money, error) {
	movimientos, err := c.repo.GetMovimientosByQuinielaID(idQuiniela)
	if err != nil {
		return nil, err
	}

	excedentes := make(map[int32]money.Money)
	for _, m := range movimientos {
		if !m.Fecha.After(desde) {
			continue
		}
		if m.Tipo == "FIADO" {
			excedentes[m.IDCliente] = excedentes[m.IDCliente].Add(m.Monto)
		}
		if m.Tipo == "PAGO_EFECTIVO" {
			if err := c.repo.DeleteMovimiento(m.IDMovimiento); err != nil {
				return nil, err
			}
			continue
		}
		if err := c.delete.Execute(m.IDMovimiento); err != nil {
			return nil, err
		}
	}

	return excedentes, nil
}

// RegistrarFiado registra un FIADO contra el cliente indicado.
func (c *CobradorQuiniela) RegistrarFiado(idCliente int32, idQuiniela int32, monto money.Money, descripcion string) error {
	if monto.IsZero() || monto.IsNegative() {
		return nil
	}
	_, err := c.create.Execute(idCliente, &idQuiniela, "FIADO", monto, descripcion)
	return err
}
