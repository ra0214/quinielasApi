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

	saved, err := cq.repo.SaveQuiniela(q)
	if err != nil {
		return nil, err
	}

	if cq.cobro != nil && saved.IDCliente != nil {
		if err := cq.cobro.SincronizarCobro(*saved.IDCliente, saved.ID, saved.TotalACobrar(), saved.MontoPagado); err != nil {
			return nil, err
		}
	}

	return saved, nil
}
