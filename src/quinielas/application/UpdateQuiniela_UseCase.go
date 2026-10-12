package application

import (
	"quinielas/src/quinielas/domain"
)

type UpdateQuiniela struct {
	repo  domain.IQuiniela
	cobro domain.CobroDirecta
}

func NewUpdateQuiniela(repo domain.IQuiniela, cobro domain.CobroDirecta) *UpdateQuiniela {
	return &UpdateQuiniela{repo: repo, cobro: cobro}
}

func (uq *UpdateQuiniela) Execute(q *domain.Quiniela) error {
	if err := uq.repo.UpdateQuiniela(q); err != nil {
		return err
	}

	if uq.cobro != nil {
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
