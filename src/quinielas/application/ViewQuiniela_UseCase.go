package application

import "quinielas/src/quinielas/domain"

type ViewQuiniela struct {
	repo domain.IQuiniela
}

func NewViewQuiniela(repo domain.IQuiniela) *ViewQuiniela {
	return &ViewQuiniela{repo: repo}
}

func (vq *ViewQuiniela) Execute() ([]domain.Quiniela, error) {
	return vq.repo.GetAll()
}

func (vq *ViewQuiniela) ExecuteByID(id int32) (*domain.Quiniela, error) {
	return vq.repo.GetQuinielaByID(id)
}

func (vq *ViewQuiniela) ExecuteByEdicionID(idEdicion int32) ([]domain.Quiniela, error) {
	return vq.repo.GetByEdicionID(idEdicion)
}