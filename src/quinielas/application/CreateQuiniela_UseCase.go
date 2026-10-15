package application

import (
	"quinielas/src/quinielas/domain"
)

type CreateQuiniela struct {
	repo  domain.IQuiniela
	cobro domain.CobroDirecta
}

func NewCreateQuiniela(repo domain.IQuiniela, cobro domain.CobroDirecta) *CreateQuiniela {
	return &CreateQuiniela{repo: repo, cobro: cobro}
}

func (cq *CreateQuiniela) Execute(q *domain.Quiniela) (*domain.Quiniela, error) {
	// Validar que la edición no tenga ya 5 quinielas
	count, err := cq.repo.CountByEdicionID(q.IDEdicion)
	if err != nil {
		return nil, err
	}

	if count >= 5 {
		return nil, domain.ErrLimitQuinielasExceeded
	}

	// Un mismo analista (variante) no puede repetirse dentro del mismo número.
	existentes, err := cq.repo.GetByEdicionID(q.IDEdicion)
	if err != nil {
		return nil, err
	}
	for _, x := range existentes {
		if x.NombreVariante == q.NombreVariante {
			return nil, domain.ErrAnalistaDuplicado
		}
	}

	// El estado de fallo se deriva de si la quiniela está sujeta a cancelación:
	// sujeto -> 1 (cancelar), no sujeto -> 2 (el creador absorbe el faltante).
	if q.TipoQuiniela == "DESARROLLO" && q.EstadoFallo == 0 {
		if q.SujetoCancelacion {
			q.EstadoFallo = 1
		} else {
			q.EstadoFallo = 2
		}
	}

	saved, err := cq.repo.SaveQuiniela(q)
	if err != nil {
		return nil, err
	}

	// Las quinielas DESARROLLO no generan un cobro directo al cliente
	// responsable: se fondean con los aportes de los clientes.
	if cq.cobro != nil && saved.IDCliente != nil && saved.TipoQuiniela != "DESARROLLO" {
		if err := cq.cobro.SincronizarCobro(*saved.IDCliente, saved.ID, saved.TotalACobrar(), saved.MontoPagado, saved.CanalVenta); err != nil {
			return nil, err
		}
	}

	return saved, nil
}
