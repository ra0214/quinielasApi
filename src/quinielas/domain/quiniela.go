package domain

import (
	"errors"
	"quinielas/src/shared/money"
	"time"
)

var ErrLimitQuinielasExceeded = errors.New("límite alcanzado: no se pueden crear más de 5 quinielas en la misma edición")

// ErrPuntoRetornoRequerido indica que se intentó subir la meta de una quiniela
// con aportes sin indicar si marca punto de retorno.
var ErrPuntoRetornoRequerido = errors.New("al aumentar la meta de una quiniela con aportes debes indicar si se marca punto de retorno")

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
	// ReembolsarCobroDesde hace lo mismo pero solo con los movimientos
	// posteriores a `desde`, y devuelve lo devuelto por cliente.
	ReembolsarCobroDesde(idQuiniela int32, desde time.Time) (map[int32]money.Money, error)
	// RegistrarFiado registra un FIADO a favor de la casa contra el cliente
	// indicado (se usa cuando el creador absorbe un faltante).
	RegistrarFiado(idCliente int32, idQuiniela int32, monto money.Money, descripcion string) error
}

// AportesDeQuiniela expone lo que el ciclo de vida necesita del módulo de
// aportes. Lo implementa el módulo de aportes.
type AportesDeQuiniela interface {
	DeleteAportesByQuinielaID(idQuiniela int32) error
	// GetAportesByQuinielaID devuelve el acumulado por cliente.
	GetAportesByQuinielaID(idQuiniela int32) ([]AporteResumen, error)
	// ReducirAporte descuenta `monto` del acumulado del cliente en la quiniela
	// y recalcula su porcentaje con `meta`.
	ReducirAporte(idCliente int32, idQuiniela int32, monto money.Money, meta money.Money) error
}

// AporteResumen es la vista mínima de un aporte para el ciclo de vida.
type AporteResumen struct {
	IDCliente           int32
	MontoTotalAcumulado money.Money
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
	GetVencidasDesarrollo(ahora time.Time) ([]Quiniela, error)
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

	// Ciclo de vida con fecha límite (solo DESARROLLO).
	SujetoCancelacion bool         `json:"sujeto_cancelacion"`
	EstadoFallo       int          `json:"estado_fallo"` // 1 cancelar, 2 creador absorbe, 3 devolver extra, 4 cliente absorbe aumento
	MontoMetaAnterior *money.Money `json:"monto_meta_anterior,omitempty"`
	FechaMetaAnterior *time.Time   `json:"fecha_meta_anterior,omitempty"`
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
