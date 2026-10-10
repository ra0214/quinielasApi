package application

import (
	"quinielas/src/movimientos/domain"
)

// UpdateMovimiento reemplaza los datos de un movimiento ya registrado. Primero
// deshace el efecto del valor original sobre el saldo y luego aplica el del
// valor nuevo; si el nuevo no cabe en el saldo (RETIRO/PAGO_SALDO) se restaura
// el original y se aborta sin guardar nada.
type UpdateMovimiento struct {
	repo      domain.IMovimiento
	ajustador SaldoAjustador
}

func NewUpdateMovimiento(repo domain.IMovimiento, ajustador SaldoAjustador) *UpdateMovimiento {
	return &UpdateMovimiento{repo: repo, ajustador: ajustador}
}

func (um *UpdateMovimiento) Execute(idMovimiento int32, nuevo *domain.Movimiento) (*domain.Movimiento, error) {
	viejo, err := um.repo.GetMovimientoByID(idMovimiento)
	if err != nil {
		return nil, err
	}

	// 1. Deshacer el efecto del movimiento original sobre el saldo.
	if err := um.ajustador.Revertir(viejo.IDCliente, viejo.Tipo, viejo.Monto); err != nil {
		return nil, err
	}

	// 2. Aplicar el efecto del movimiento nuevo. Si el saldo no alcanza para
	//    un retiro o pago con saldo, se restaura el efecto original y se aborta.
	if err := um.ajustador.Ajustar(nuevo.IDCliente, nuevo.Tipo, nuevo.Monto); err != nil {
		_ = um.ajustador.Ajustar(viejo.IDCliente, viejo.Tipo, viejo.Monto)
		return nil, err
	}

	// 3. Persistir el movimiento editado.
	if err := um.repo.UpdateMovimiento(nuevo); err != nil {
		// Revertimos el nuevo efecto y restauramos el original.
		_ = um.ajustador.Revertir(nuevo.IDCliente, nuevo.Tipo, nuevo.Monto)
		_ = um.ajustador.Ajustar(viejo.IDCliente, viejo.Tipo, viejo.Monto)
		return nil, err
	}

	return nuevo, nil
}