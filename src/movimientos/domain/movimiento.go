package domain

import (
	"quinielas/src/shared/money"
	"time"
)

type IMovimiento interface {
	SaveMovimiento(idCliente int32, idQuiniela *int32, tipo string, monto money.Money, descripcion string) (*Movimiento, error)
	DeleteMovimiento(idMovimiento int32) error
	GetMovimientoByID(idMovimiento int32) (*Movimiento, error)
	GetMovimientosByClienteID(idCliente int32) ([]Movimiento, error)
	GetAllMovimientos() ([]Movimiento, error)
}

type Movimiento struct {
	IDMovimiento int32       `json:"id_movimiento"`
	IDCliente    int32       `json:"id_cliente"`
	IDQuiniela   *int32      `json:"id_quiniela,omitempty"` // Puntero porque puede ser nulo en retiros/pagos generales
	Tipo         string      `json:"tipo"`                 // 'PAGO_EFECTIVO', 'PAGO_SALDO', 'FIADO', 'RETIRO', 'PREMIO_ABONO', 'DEVOLUCION'
	Monto        money.Money `json:"monto"`
	Descripcion  string      `json:"descripcion"`
	Fecha        time.Time   `json:"fecha"`
}

func NewMovimiento(idCliente int32, idQuiniela *int32, tipo string, monto money.Money, descripcion string) *Movimiento {
	return &Movimiento{
		IDCliente:   idCliente,
		IDQuiniela:  idQuiniela,
		Tipo:        tipo,
		Monto:       monto,
		Descripcion: descripcion,
		Fecha:       time.Now(),
	}
}
