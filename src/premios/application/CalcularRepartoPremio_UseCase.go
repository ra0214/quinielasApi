package application

import (
	"database/sql"
	"errors"

	aportesDomain "quinielas/src/aportes/domain"
	"quinielas/src/premios/domain"
	quinielasDomain "quinielas/src/quinielas/domain"
)

// AportesDeQuiniela es todo lo que el reparto necesita saber de los aportes.
type AportesDeQuiniela interface {
	GetAportesByQuinielaID(idQuiniela int32) ([]aportesDomain.AporteDetalle, error)
}

// MetaDeQuiniela expone la meta de la quiniela, que es contra la que se valida
// que el premio sea repartible.
type MetaDeQuiniela interface {
	GetQuinielaByID(id int32) (*quinielasDomain.Quiniela, error)
}

type CalcularRepartoPremio struct {
	premioRepo   domain.IPremio
	aporteRepo   AportesDeQuiniela
	quinielaRepo MetaDeQuiniela
}

func NewCalcularRepartoPremio(premioRepo domain.IPremio, aporteRepo AportesDeQuiniela, quinielaRepo MetaDeQuiniela) *CalcularRepartoPremio {
	return &CalcularRepartoPremio{
		premioRepo:   premioRepo,
		aporteRepo:   aporteRepo,
		quinielaRepo: quinielaRepo,
	}
}

// Execute calcula el reparto del premio de una quiniela y lo guarda como
// snapshot. A partir de ese momento el monto a pagar queda congelado: aunque
// alguien aporte mas a la misma quiniela, el reparto ya no cambia.
//
// No deja repartir dos veces la misma quiniela, que es lo que evita un doble
// pago del premio.
func (c *CalcularRepartoPremio) Execute(idQuiniela int32) (*domain.RepartoPremio, error) {
	yaRepartido, err := c.premioRepo.ExistsReparto(idQuiniela)
	if err != nil {
		return nil, err
	}
	if yaRepartido {
		return nil, domain.ErrRepartoYaGenerado
	}

	premio, err := c.premioRepo.GetPremioByQuinielaID(idQuiniela)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPremioNoRegistrado
		}
		return nil, err
	}

	quiniela, err := c.quinielaRepo.GetQuinielaByID(idQuiniela)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMetaInvalida
		}
		return nil, err
	}

	aportes, err := c.aporteRepo.GetAportesByQuinielaID(idQuiniela)
	if err != nil {
		return nil, err
	}

	participantes := make([]domain.AporteParticipante, 0, len(aportes))
	for _, a := range aportes {
		participantes = append(participantes, domain.AporteParticipante{
			IDCliente:     a.IDCliente,
			NombreCliente: a.NombreCliente,
			Monto:         a.MontoTotalAcumulado,
		})
	}

	reparto, err := domain.CalcularReparto(premio, participantes, quiniela.MontoMeta)
	if err != nil {
		return nil, err
	}

	if err := c.premioRepo.SaveReparto(reparto); err != nil {
		return nil, err
	}

	return reparto, nil
}
