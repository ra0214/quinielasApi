package application

import (
	"quinielas/src/aportes/domain"
)

type DeleteAporte struct {
	repo  domain.IAporte
	cobro domain.CobroAporte
}

func NewDeleteAporte(repo domain.IAporte, cobro domain.CobroAporte) *DeleteAporte {
	return &DeleteAporte{repo: repo, cobro: cobro}
}

func (da *DeleteAporte) Execute(idAporte int32) error {
	aporte, err := da.repo.GetAporteByID(idAporte)
	if err != nil {
		return err
	}

	if da.cobro != nil {
		if err := da.cobro.RevertirAporte(aporte.IDCliente, aporte.IDQuiniela); err != nil {
			return err
		}
	}

	return da.repo.DeleteAporte(idAporte)
}