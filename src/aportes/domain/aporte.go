package domain

import (
	"errors"

	"quinielas/src/shared/money"
)

var (
	ErrQuinielaNoEncontrada = errors.New("la quiniela no existe")
	ErrMetaExcedida         = errors.New("el aporte excede el monto total solicitado de la quiniela")
)

type IAporte interface {
	SaveOrUpdateAporte(idCliente int32, idQuiniela int32, montoSumar money.Money, montoMetaQuiniela money.Money) (*Aporte, error)
	DeleteAporte(idAporte int32) error
	GetAporteByID(idAporte int32) (*Aporte, error)
	GetAportesByQuinielaID(idQuiniela int32) ([]AporteDetalle, error)
	GetAportesByClienteID(idCliente int32) ([]Aporte, error)
}

type Aporte struct {
	IDAporte                int32         `json:"id_aporte"`
	IDCliente               int32         `json:"id_cliente"`
	IDQuiniela              int32         `json:"id_quiniela"`
	MontoTotalAcumulado     money.Money   `json:"monto_total_acumulado"`
	PorcentajeParticipacion money.Percent `json:"porcentaje_participacion"`
}

// AporteDetalle incluye el nombre del cliente para vistas y monitoreo de quinielas
type AporteDetalle struct {
	IDAporte                int32         `json:"id_aporte"`
	IDCliente               int32         `json:"id_cliente"`
	NombreCliente           string        `json:"nombre_cliente"`
	IDQuiniela              int32         `json:"id_quiniela"`
	MontoTotalAcumulado     money.Money   `json:"monto_total_acumulado"`
	PorcentajeParticipacion money.Percent `json:"porcentaje_participacion"`
}

func NewAporte(idCliente int32, idQuiniela int32, montoAcumulado money.Money, porcentaje money.Percent) *Aporte {
	return &Aporte{
		IDCliente:               idCliente,
		IDQuiniela:              idQuiniela,
		MontoTotalAcumulado:     montoAcumulado,
		PorcentajeParticipacion: porcentaje,
	}
}
