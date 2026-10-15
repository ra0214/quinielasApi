package domain

import (
	"errors"
	"quinielas/src/shared/money"
	"time"
)

var ErrLimitQuinielasExceeded = errors.New("límite alcanzado: no se pueden crear más de 5 quinielas en la misma edición")

// CobroDirecta refleja el cobro de una quiniela en movimientos y saldo. Lo
// implementa el módulo de movimientos; las quinielas solo lo consumen.
// canalVenta identifica cómo se vendió la quiniela ('EN_LINEA' o 'VENTANILLA')
// y se registra en la descripción de los movimientos.
type CobroDirecta interface {
	SincronizarCobro(idCliente int32, idQuiniela int32, total money.Money, pagado money.Money, canalVenta string) error
	RevertirCobro(idQuiniela int32) error
	// ReembolsarCobro elimina los movimientos de la quiniela dejando lo pagado
	// (PAGO_EFECTIVO) como saldo a favor del cliente y descartando lo fiado.
	ReembolsarCobro(idQuiniela int32) error
}

// AportesDeQuiniela borra los aportes registrados en una quiniela cuando esta
// se elimina. Lo implementa el módulo de aportes.
type AportesDeQuiniela interface {
	DeleteAportesByQuinielaID(idQuiniela int32) error
}

// PremiosDeQuiniela borra el premio y su reparto (si existen) de una quiniela
// cuando esta se elimina. Lo implementa el módulo de premios.
type PremiosDeQuiniela interface {
	DeletePremioByQuinielaID(idQuiniela int32) error
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
	ID             int32       `json:"id_quiniela"`
	IDEdicion      int32       `json:"id_edicion"`
	NombreVariante string      `json:"nombre_variante"` // 'Neto', 'Otros', 'Tlacuache', 'Salomón', 'Extra'
	MontoMeta      money.Money `json:"monto_meta"`      // costo total de la quiniela
	FechaLimite    time.Time   `json:"fecha_limite"`
	Estado         string      `json:"estado"` // 'EN_JUEGO', 'COMPLETADA', 'CANCELADA'
	TipoQuiniela   string      `json:"tipo_quiniela,omitempty"` // DIRECTA, DESARROLLO
	IDCliente      *int32      `json:"id_cliente,omitempty"`
	MontoPagado    money.Money `json:"monto_pagado"`
	Analista       string      `json:"analista,omitempty"`
	TipoPago       string      `json:"tipo_pago,omitempty"` // CONTADO, PAGO_PARCIAL, A_CREDITO
	CanalVenta     string      `json:"canal_venta,omitempty"` // EN_LINEA, VENTANILLA
}

func NewQuiniela(idEdicion int32, nombreVariante string, montoMeta money.Money, fechaLimite time.Time) *Quiniela {
	return &Quiniela{
		IDEdicion:      idEdicion,
		NombreVariante: nombreVariante,
		MontoMeta:      montoMeta,
		FechaLimite:    fechaLimite,
		Estado:         "EN_JUEGO",
	}
}

// TotalACobrar devuelve el monto total a cobrar: el costo total de la quiniela.
func (q *Quiniela) TotalACobrar() money.Money {
	return q.MontoMeta
}
