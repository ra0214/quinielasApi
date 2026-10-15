package application

import "quinielas/src/quinielas/domain"

type DeleteQuiniela struct {
	repo    domain.IQuiniela
	cobro   domain.CobroDirecta
	aportes domain.AportesDeQuiniela
	premios domain.PremiosDeQuiniela
}

func NewDeleteQuiniela(repo domain.IQuiniela, cobro domain.CobroDirecta, aportes domain.AportesDeQuiniela, premios domain.PremiosDeQuiniela) *DeleteQuiniela {
	return &DeleteQuiniela{repo: repo, cobro: cobro, aportes: aportes, premios: premios}
}

func (dq *DeleteQuiniela) Execute(id int32) error {
	// Reembolsa el cobro: lo pagado queda como saldo a favor y lo fiado se
	// descarta.
	if dq.cobro != nil {
		if err := dq.cobro.ReembolsarCobro(id); err != nil {
			return err
		}
	}
	// Los aportes y premios tienen claves foráneas hacia la quiniela: hay que
	// eliminarlos antes para poder borrarla.
	if dq.aportes != nil {
		if err := dq.aportes.DeleteAportesByQuinielaID(id); err != nil {
			return err
		}
	}
	if dq.premios != nil {
		if err := dq.premios.DeletePremioByQuinielaID(id); err != nil {
			return err
		}
	}
	return dq.repo.DeleteQuiniela(id)
}
