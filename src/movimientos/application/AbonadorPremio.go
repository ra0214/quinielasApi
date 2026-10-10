package application

import (
	"fmt"

	"quinielas/src/shared/money"
)

// AbonadorPremio registra los abonos del reparto de un premio como
// movimientos de tipo PREMIO_ABONO, de modo que el saldo del ganador suba
// automáticamente. Cumple el puerto AbonosDePremio del módulo de premios.
type AbonadorPremio struct {
	create *CreateMovimiento
	delete *DeleteMovimiento
}

func NewAbonadorPremio(create *CreateMovimiento, delete *DeleteMovimiento) *AbonadorPremio {
	return &AbonadorPremio{create: create, delete: delete}
}

// Abonar acredita el monto neto del premio al saldo del cliente. Devuelve el
// ID del movimiento creado para poder deshacerlo si el reparto no se congela.
func (a *AbonadorPremio) Abonar(idCliente int32, idQuiniela int32, monto money.Money) (int32, error) {
	descripcion := fmt.Sprintf("Premio ganado en la quiniela #%d", idQuiniela)
	mov, err := a.create.Execute(idCliente, &idQuiniela, "PREMIO_ABONO", monto, descripcion)
	if err != nil {
		return 0, err
	}
	return mov.IDMovimiento, nil
}

// DeshacerAbono elimina el movimiento de abono y revierte el saldo (rollback).
func (a *AbonadorPremio) DeshacerAbono(idMovimiento int32) error {
	return a.delete.Execute(idMovimiento)
}