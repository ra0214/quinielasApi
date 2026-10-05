package application

import "quinielas/src/ediciones/domain"

type DeleteEdicion struct {
	repo domain.IEdicion
}

func NewDeleteEdicion(repo domain.IEdicion) *DeleteEdicion {
	return &DeleteEdicion{repo: repo}
}

func (de *DeleteEdicion) Execute(id int32) error {
	return de.repo.DeleteEdicion(id)
}