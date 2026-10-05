package application

import "quinielas/src/premios/domain"

type ViewPremio struct {
	repo domain.IPremio
}

func NewViewPremio(repo domain.IPremio) *ViewPremio {
	return &ViewPremio{repo: repo}
}

func (vp *ViewPremio) ExecuteByQuinielaID(idQuiniela int32) (*domain.Premio, error) {
	return vp.repo.GetPremioByQuinielaID(idQuiniela)
}

func (vp *ViewPremio) ExecuteGetAll() ([]domain.Premio, error) {
	return vp.repo.GetAllPremios()
}