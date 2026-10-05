package application

import (
	"quinielas/src/aportes/domain"
	"quinielas/src/shared/money"
)

type CreateOrUpdateAporte struct {
	repo domain.IAporte
}

func NewCreateOrUpdateAporte(repo domain.IAporte) *CreateOrUpdateAporte {
	return &CreateOrUpdateAporte{repo: repo}
}

func (cua *CreateOrUpdateAporte) Execute(idCliente int32, idQuiniela int32, montoSumar money.Money, montoMetaQuiniela money.Money) (*domain.Aporte, error) {
	return cua.repo.SaveOrUpdateAporte(idCliente, idQuiniela, montoSumar, montoMetaQuiniela)
}