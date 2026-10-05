package domain

import (
	"errors"
	"quinielas/src/shared/money"
	"time"
)

var ErrLimitQuinielasExceeded = errors.New("límite alcanzado: no se pueden crear más de 5 quinielas en la misma edición")

type IQuiniela interface {
	SaveQuiniela(idEdicion int32, nombreVariante string, precio money.Money, montoMeta money.Money, fechaLimite time.Time) (*Quiniela, error)
	DeleteQuiniela(id int32) error
	UpdateQuiniela(id int32, nombreVariante string, precio money.Money, montoMeta money.Money, fechaLimite time.Time, estado string) error
	GetAll() ([]Quiniela, error)
	GetQuinielaByID(id int32) (*Quiniela, error)
	GetByEdicionID(idEdicion int32) ([]Quiniela, error)
	CountByEdicionID(idEdicion int32) (int, error)
}

type Quiniela struct {
	ID             int32       `json:"id_quiniela"`
	IDEdicion      int32       `json:"id_edicion"`
	NombreVariante string      `json:"nombre_variante"` // 'Neto', 'Otros', 'Tlacuache', 'Salomón', 'Extra'
	Precio         money.Money `json:"precio"`
	MontoMeta      money.Money `json:"monto_meta"`
	FechaLimite    time.Time   `json:"fecha_limite"`
	Estado         string      `json:"estado"` // 'EN_JUEGO', 'COMPLETADA', 'CANCELADA'
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
