package application

import (
	"quinielas/src/ediciones/domain"
	"time"
)

type UpdateEdicion struct {
	repo domain.IEdicion
}

func NewUpdateEdicion(repo domain.IEdicion) *UpdateEdicion {
	return &UpdateEdicion{repo: repo}
}

func (ue *UpdateEdicion) Execute(id int32, tipoEdicion string, nombreEdicion string, fechaLimite *time.Time) error {
	return ue.repo.UpdateEdicion(id, tipoEdicion, nombreEdicion, fechaLimite)
}
