package application

import "quinielas/src/ediciones/domain"

type CreateEdicion struct {
	repo domain.IEdicion
}

func NewCreateEdicion(repo domain.IEdicion) *CreateEdicion {
	return &CreateEdicion{repo: repo}
}

func (ce *CreateEdicion) Execute(tipoEdicion string, nombreEdicion string) (*domain.Edicion, error) {
	return ce.repo.SaveEdicion(tipoEdicion, nombreEdicion)
}