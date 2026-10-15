package domain

import (
	"errors"

	"quinielas/src/shared/money"
)

var (
	ErrQuinielaNoEncontrada = errors.New("la quiniela no existe")
	ErrMetaExcedida         = errors.New("el aporte excede el monto total solicitado de la quiniela")
)

// CobroAporte refleja un aporte en movimientos y saldo. Lo implementa el
// módulo de movimientos; los aportes solo lo consumen. La parte pagada del
// aporte queda como PAGO_EFECTIVO y lo pendiente como FIADO, con su canal de
// venta.
type CobroAporte interface {
	SincronizarAporte(idCliente int32, idQuiniela int32, total money.Money, pagado money.Money, canalVenta string) error
	RevertirAporte(idCliente int32, idQuiniela int32) error
}

type IAporte interface {
	SaveOrUpdateAporte(idCliente int32, idQuiniela int32, montoSumar money.Money, montoMetaQuiniela money.Money) (*Aporte, error)
	// UpdateAporte reemplaza el monto acumulado de un aporte existente por un
	// monto nuevo, recalcando el porcentaje de participación.
	UpdateAporte(idAporte int32, montoNuevo money.Money, montoMetaQuiniela money.Money) (*Aporte, error)
	DeleteAporte(idAporte int32) error
	DeleteAportesByQuinielaID(idQuiniela int32) error
	GetAporteByID(idAporte int32) (*Aporte, error)
	GetAportesByQuinielaID(idQuiniela int32) ([]AporteDetalle, error)
	GetAportesByClienteID(idCliente int32) ([]Aporte, error)
	GetAllAportes() ([]AporteDetalle, error)
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
