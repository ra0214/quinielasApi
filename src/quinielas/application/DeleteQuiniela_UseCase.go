package application

import "quinielas/src/quinielas/domain"

type DeleteQuiniela struct {
	repo domain.IQuiniela
}

func NewDeleteQuiniela(repo domain.IQuiniela) *DeleteQuiniela {
	return &DeleteQuiniela{repo: repo}
}

func (dq *DeleteQuiniela) Execute(id int32) error {
	return dq.repo.DeleteQuiniela(id)
}