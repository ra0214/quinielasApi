package domain

import (
	"quinielas/src/shared/money"
	"time"
)

type IPremio interface {
	SavePremio(idQuiniela int32, montoBruto money.Money) (*Premio, error)
	// UpdatePremio modifica el monto bruto de un premio ya registrado y
	// recalcula la retención y el monto neto.
	UpdatePremio(idPremio int32, montoBruto money.Money) (*Premio, error)
	DeletePremio(idPremio int32) error
	GetPremioByQuinielaID(idQuiniela int32) (*Premio, error)
	GetAllPremios() ([]Premio, error)

	SaveReparto(reparto *RepartoPremio) error
	GetRepartoByQuiniela(idQuiniela int32) (*RepartoPremio, error)
	ExistsReparto(idQuiniela int32) (bool, error)
}

type Premio struct {
	IDPremio            int32         `json:"id_premio"`
	IDQuiniela          int32         `json:"id_quiniela"`
	MontoBruto          money.Money   `json:"monto_bruto"`
	PorcentajeRetencion money.Percent `json:"porcentaje_retencion"` // Por defecto 7.0%
	MontoNeto           money.Money   `json:"monto_neto"`            // Monto Bruto - 7%
	FechaRegistro       time.Time     `json:"fecha_registro"`
}

var RetencionDefault = money.NewPercentFromInt64(7)

func NewPremio(idQuiniela int32, montoBruto money.Money) *Premio {
	retencion := RetencionDefault
	montoNeto := montoBruto.AfterRetention(retencion)

	return &Premio{
		IDQuiniela:          idQuiniela,
		MontoBruto:          montoBruto,
		PorcentajeRetencion: retencion,
		MontoNeto:           montoNeto,
		FechaRegistro:       time.Now(),
	}
}
