package application

import "quinielas/src/ediciones/domain"

type UpdateEdicion struct {
	repo domain.IEdicion
}

func NewUpdateEdicion(repo domain.IEdicion) *UpdateEdicion {
	return &UpdateEdicion{repo: repo}
}

func (ue *UpdateEdicion) Execute(id int32, tipoEdicion string, nombreEdicion string) error {
	return ue.repo.UpdateEdicion(id, tipoEdicion, nombreEdicion)
}