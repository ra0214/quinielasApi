package application

import (
	"quinielas/src/movimientos/domain"
	"quinielas/src/shared/money"
)

type CreateMovimiento struct {
	repo      domain.IMovimiento
	ajustador SaldoAjustador
}

func NewCreateMovimiento(repo domain.IMovimiento, ajustador SaldoAjustador) *CreateMovimiento {
	return &CreateMovimiento{repo: repo, ajustador: ajustador}
}

func (cm *CreateMovimiento) Execute(idCliente int32, idQuiniela *int32, tipo string, monto money.Money, descripcion string) (*domain.Movimiento, error) {
	// Aplicamos el efecto del movimiento sobre el saldo. Si el saldo es
	// insuficiente (RETIRO/PAGO_SALDO) no se registra el movimiento.
	if err := cm.ajustador.Ajustar(idCliente, tipo, monto); err != nil {
		return nil, err
	}

	movimiento, err := cm.repo.SaveMovimiento(idCliente, idQuiniela, tipo, monto, descripcion)
	if err != nil {
		// Revertimos el ajuste para no dejar el saldo alterado si el guardado falla.
		_ = cm.ajustador.Revertir(idCliente, tipo, monto)
		return nil, err
	}

	return movimiento, nil
}