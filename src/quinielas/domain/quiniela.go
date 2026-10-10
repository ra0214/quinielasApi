package domain

import (
	"errors"
	"quinielas/src/shared/money"
	"time"
)

var ErrLimitQuinielasExceeded = errors.New("límite alcanzado: no se pueden crear más de 5 quinielas en la misma edición")

// CobroDirecta refleja el cobro de una quiniela en movimientos y saldo. Lo
// implementa el módulo de movimientos; las quinielas solo lo consumen.
type CobroDirecta interface {
	SincronizarCobro(idCliente int32, idQuiniela int32, total money.Money, pagado money.Money) error
	RevertirCobro(idQuiniela int32) error
}

type IQuiniela interface {
	SaveQuiniela(q *Quiniela) (*Quiniela, error)
	DeleteQuiniela(id int32) error
	UpdateQuiniela(q *Quiniela) error
	GetAll() ([]Quiniela, error)
	GetQuinielaByID(id int32) (*Quiniela, error)
	GetByEdicionID(idEdicion int32) ([]Quiniela, error)
	CountByEdicionID(idEdicion int32) (int, error)
}

type Quiniela struct {
	ID                int32       `json:"id_quiniela"`
	IDEdicion         int32       `json:"id_edicion"`
	NombreVariante    string      `json:"nombre_variante"` // 'Neto', 'Otros', 'Tlacuache', 'Salomón', 'Extra'
	Precio            money.Money `json:"precio"`
	MontoMeta         money.Money `json:"monto_meta"`
	FechaLimite       time.Time   `json:"fecha_limite"`
	Estado            string      `json:"estado"` // 'EN_JUEGO', 'COMPLETADA', 'CANCELADA'
	TipoQuiniela      string      `json:"tipo_quiniela,omitempty"` // DIRECTA, DIRECTA_EN_COPERACHA, DESARROLLO
	IDCliente         *int32      `json:"id_cliente,omitempty"`
	MontoPagado       money.Money `json:"monto_pagado"`
	CantidadBoletos   int32       `json:"cantidad_boletos,omitempty"`
	Analista          string      `json:"analista,omitempty"`
	TipoPago          string      `json:"tipo_pago,omitempty"` // CONTADO, PAGO_PARCIAL, A_CREDITO
}

func NewQuiniela(idEdicion int32, nombreVariante string, precio money.Money, montoMeta money.Money, fechaLimite time.Time) *Quiniela {
	return &Quiniela{
		IDEdicion:      idEdicion,
		NombreVariante: nombreVariante,
		Precio:         precio,
		MontoMeta:      montoMeta,
		FechaLimite:    fechaLimite,
		Estado:         "EN_JUEGO",
	}
}

func (q *Quiniela) TotalACobrar() money.Money {
	switch q.TipoQuiniela {
	case "DIRECTA_EN_COPERACHA":
		if q.CantidadBoletos <= 0 {
			return q.Precio
		}
		return q.Precio.MulInt(int64(q.CantidadBoletos))
	case "DESARROLLO":
		return q.Precio // precio = costo total
	case "DIRECTA":
		fallthrough
	default:
		return q.Precio
	}
}
