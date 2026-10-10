package application

import "quinielas/src/quinielas/domain"

type DeleteQuiniela struct {
	repo  domain.IQuiniela
	cobro domain.CobroDirecta
}

func NewDeleteQuiniela(repo domain.IQuiniela, cobro domain.CobroDirecta) *DeleteQuiniela {
	return &DeleteQuiniela{repo: repo, cobro: cobro}
}

func (dq *DeleteQuiniela) Execute(id int32) error {
	if dq.cobro != nil {
		if err := dq.cobro.RevertirCobro(id); err != nil {
			return err
		}
	}
	return dq.repo.DeleteQuiniela(id)
}
