package application

import "quinielas/src/ediciones/domain"

type ViewEdicion struct {
	repo domain.IEdicion
}

func NewViewEdicion(repo domain.IEdicion) *ViewEdicion {
	return &ViewEdicion{repo: repo}
}

func (ve *ViewEdicion) Execute() ([]domain.Edicion, error) {
	return ve.repo.GetAll()
}

func (ve *ViewEdicion) ExecuteByID(id int32) (*domain.Edicion, error) {
	return ve.repo.GetEdicionByID(id)
}