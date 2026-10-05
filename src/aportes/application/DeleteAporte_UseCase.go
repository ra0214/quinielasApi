package application

import "quinielas/src/aportes/domain"

type DeleteAporte struct {
	repo domain.IAporte
}

func NewDeleteAporte(repo domain.IAporte) *DeleteAporte {
	return &DeleteAporte{repo: repo}
}

func (da *DeleteAporte) Execute(idAporte int32) error {
	return da.repo.DeleteAporte(idAporte)
}