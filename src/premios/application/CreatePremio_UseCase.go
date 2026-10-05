package application

import (
	"quinielas/src/premios/domain"
	"quinielas/src/shared/money"
)

type CreatePremio struct {
	repo domain.IPremio
}

func NewCreatePremio(repo domain.IPremio) *CreatePremio {
	return &CreatePremio{repo: repo}
}

func (cp *CreatePremio) Execute(idQuiniela int32, montoBruto money.Money) (*domain.Premio, error) {
	return cp.repo.SavePremio(idQuiniela, montoBruto)
}