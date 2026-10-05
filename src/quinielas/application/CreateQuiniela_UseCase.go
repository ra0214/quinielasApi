package application

import (
	"quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
	"time"
)

type CreateQuiniela struct {
	repo domain.IQuiniela
}

func NewCreateQuiniela(repo domain.IQuiniela) *CreateQuiniela {
	return &CreateQuiniela{repo: repo}
}

func (cq *CreateQuiniela) Execute(idEdicion int32, nombreVariante string, precio money.Money, montoMeta money.Money, fechaLimite time.Time) (*domain.Quiniela, error) {
	// Validar que la edición no tenga ya 5 quinielas
	count, err := cq.repo.CountByEdicionID(idEdicion)
	if err != nil {
		return nil, err
	}

	if count >= 5 {
		return nil, domain.ErrLimitQuinielasExceeded
	}

	return cq.repo.SaveQuiniela(idEdicion, nombreVariante, precio, montoMeta, fechaLimite)
}