package domain

import (
	"quinielas/src/shared/money"
	"time"
)

type IPremio interface {
	SavePremio(idQuiniela int32, montoBruto money.Money, montoNeto money.Money) (*Premio, error)
	// UpdatePremio modifica el monto de un premio ya registrado y recalcula
	// la retención y el monto neto.
	UpdatePremio(idPremio int32, montoBruto money.Money, montoNeto money.Money) (*Premio, error)
	DeletePremio(idPremio int32) error
	DeletePremioByQuinielaID(idQuiniela int32) error
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

// SinRetencion se usa cuando el premio se captura directamente por su monto
// neto: no se descuenta nada.
var SinRetencion = money.NewPercentFromInt64(0)

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

// NewPremioDesdeNeto registra un premio capturado por su monto neto: no se
// aplica retención, por lo que el bruto coincide con el neto.
func NewPremioDesdeNeto(idQuiniela int32, montoNeto money.Money) *Premio {
	return &Premio{
		IDQuiniela:          idQuiniela,
		MontoBruto:          montoNeto,
		PorcentajeRetencion: SinRetencion,
		MontoNeto:           montoNeto,
		FechaRegistro:       time.Now(),
	}
}

// NewPremioDesdeMontos elige la vía según lo capturado: si hay monto bruto se
// le aplica la retención (7%); si no, se toma el monto neto tal cual, sin
// descuento. Ambos montos son opcionales, pero al menos uno debe ser mayor a 0.
func NewPremioDesdeMontos(idQuiniela int32, montoBruto, montoNeto money.Money) *Premio {
	if montoBruto.GreaterThan(money.Zero()) {
		return NewPremio(idQuiniela, montoBruto)
	}
	return NewPremioDesdeNeto(idQuiniela, montoNeto)
}
