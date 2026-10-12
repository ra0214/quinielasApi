package application

import (
	"database/sql"
	"errors"
	"fmt"

	aportesDomain "quinielas/src/aportes/domain"
	clientesDomain "quinielas/src/clientes/domain"
	"quinielas/src/premios/domain"
	quinielasDomain "quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
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

// NombreDeCliente resuelve el nombre del comprador en una quiniela DIRECTA,
// que es quien recibe el 100% del premio.
type NombreDeCliente interface {
	GetClienteByID(id int32) (*clientesDomain.Cliente, error)
}

// AbonosDePremio abona a cada ganador el monto neto que le tocó. Se registra
// como un movimiento PREMIO_ABONO para que el saldo del cliente suba.
// DeshacerAbono elimina el movimiento y revierte el saldo (rollback).
type AbonosDePremio interface {
	Abonar(idCliente int32, idQuiniela int32, monto money.Money) (int32, error)
	DeshacerAbono(idMovimiento int32) error
}

type CalcularRepartoPremio struct {
	premioRepo   domain.IPremio
	aporteRepo   AportesDeQuiniela
	quinielaRepo MetaDeQuiniela
	abonos       AbonosDePremio
	clientes     NombreDeCliente
}

func NewCalcularRepartoPremio(premioRepo domain.IPremio, aporteRepo AportesDeQuiniela, quinielaRepo MetaDeQuiniela, abonos AbonosDePremio, clientes NombreDeCliente) *CalcularRepartoPremio {
	return &CalcularRepartoPremio{
		premioRepo:   premioRepo,
		aporteRepo:   aporteRepo,
		quinielaRepo: quinielaRepo,
		abonos:       abonos,
		clientes:     clientes,
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

	if quiniela.TipoQuiniela == "DIRECTA" {
		return c.repartirDirecta(premio, quiniela)
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

	return c.abonarYCongelar(reparto)
}

// repartirDirecta abona el premio completo al único comprador.
func (c *CalcularRepartoPremio) repartirDirecta(premio *domain.Premio, quiniela *quinielasDomain.Quiniela) (*domain.RepartoPremio, error) {
	if quiniela.IDCliente == nil {
		return nil, domain.ErrDirectaSinCliente
	}

	cliente, err := c.clientes.GetClienteByID(*quiniela.IDCliente)
	if err != nil {
		return nil, fmt.Errorf("no se pudo identificar al comprador de la quiniela: %w", err)
	}

	reparto, err := domain.CalcularRepartoDirecta(premio, domain.RepartoCliente{
		IDCliente:     cliente.ID,
		NombreCliente: cliente.Nombre,
	}, quiniela.MontoMeta)
	if err != nil {
		return nil, err
	}

	return c.abonarYCongelar(reparto)
}

// abonarYCongelar registra los abonos de cada ganador y luego congela el
// snapshot del reparto. Los abonos se escriben antes del snapshot para que, si
// algo falla a mitad, se puedan deshacer y el usuario reintente.
func (c *CalcularRepartoPremio) abonarYCongelar(reparto *domain.RepartoPremio) (*domain.RepartoPremio, error) {
	abonosCreados := make([]int32, 0, len(reparto.Participantes))
	for _, p := range reparto.Participantes {
		id, err := c.abonos.Abonar(p.IDCliente, reparto.IDQuiniela, p.MontoNetoAsignado)
		if err != nil {
			_ = c.deshacerAbonos(abonosCreados)
			return nil, fmt.Errorf("no se pudo abonar el premio a %s: %w", p.NombreCliente, err)
		}
		abonosCreados = append(abonosCreados, id)
	}

	if err := c.premioRepo.SaveReparto(reparto); err != nil {
		_ = c.deshacerAbonos(abonosCreados)
		return nil, err
	}

	return reparto, nil
}

func (c *CalcularRepartoPremio) deshacerAbonos(ids []int32) error {
	var primerError error
	for _, id := range ids {
		if err := c.abonos.DeshacerAbono(id); err != nil && primerError == nil {
			primerError = err
		}
	}
	return primerError
}
