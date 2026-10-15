package application

import (
	"time"

	"quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
)

type UpdateQuiniela struct {
	repo    domain.IQuiniela
	cobro   domain.CobroDirecta
	aportes domain.AportesDeQuiniela
}

func NewUpdateQuiniela(repo domain.IQuiniela, cobro domain.CobroDirecta, aportes domain.AportesDeQuiniela) *UpdateQuiniela {
	return &UpdateQuiniela{repo: repo, cobro: cobro, aportes: aportes}
}

func (uq *UpdateQuiniela) totalAportado(idQuiniela int32) (money.Money, error) {
	if uq.aportes == nil {
		return money.Zero(), nil
	}
	lista, err := uq.aportes.GetAportesByQuinielaID(idQuiniela)
	if err != nil {
		return money.Zero(), err
	}
	total := money.Zero()
	for _, a := range lista {
		total = total.Add(a.MontoTotalAcumulado)
	}
	return total, nil
}

func (uq *UpdateQuiniela) Execute(q *domain.Quiniela, puntoRetorno *bool) error {
	actual, err := uq.repo.GetQuinielaByID(q.ID)
	if err != nil {
		return err
	}

	// Los campos del ciclo de vida no vienen en la edición normal: se preservan.
	q.SujetoCancelacion = actual.SujetoCancelacion
	q.EstadoFallo = actual.EstadoFallo
	q.MontoMetaAnterior = actual.MontoMetaAnterior
	q.FechaMetaAnterior = actual.FechaMetaAnterior

	// Si se sube la meta de una quiniela con aportes, se exige indicar si marca
	// punto de retorno y se actualiza el estado de fallo (3 o 4).
	if actual.TipoQuiniela == "DESARROLLO" && q.MontoMeta.GreaterThan(actual.MontoMeta) {
		recaudado, err := uq.totalAportado(q.ID)
		if err != nil {
			return err
		}
		if recaudado.GreaterThan(money.Zero()) {
			if puntoRetorno == nil {
				return domain.ErrPuntoRetornoRequerido
			}
			metaAnterior := actual.MontoMeta
			ahora := time.Now()
			q.MontoMetaAnterior = &metaAnterior
			q.FechaMetaAnterior = &ahora
			if *puntoRetorno {
				q.EstadoFallo = 3
			} else {
				q.EstadoFallo = 4
			}
		}
	}

	if err := uq.repo.UpdateQuiniela(q); err != nil {
		return err
	}

	// Las quinielas DESARROLLO no generan cobro directo (se fondean por aportes).
	if uq.cobro != nil && q.TipoQuiniela != "DESARROLLO" {
		if q.IDCliente == nil {
			if err := uq.cobro.RevertirCobro(q.ID); err != nil {
				return err
			}
		} else if err := uq.cobro.SincronizarCobro(*q.IDCliente, q.ID, q.TotalACobrar(), q.MontoPagado, q.CanalVenta); err != nil {
			return err
		}
	}

	return nil
}
