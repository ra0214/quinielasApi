package application

import "quinielas/src/premios/domain"

type DeletePremio struct {
	repo domain.IPremio
}

func NewDeletePremio(repo domain.IPremio) *DeletePremio {
	return &DeletePremio{repo: repo}
}

func (dp *DeletePremio) Execute(idPremio int32) error {
	return dp.repo.DeletePremio(idPremio)
}