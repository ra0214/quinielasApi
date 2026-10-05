package application

import (
	"quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
	"time"
)

type UpdateQuiniela struct {
	repo domain.IQuiniela
}

func NewUpdateQuiniela(repo domain.IQuiniela) *UpdateQuiniela {
	return &UpdateQuiniela{repo: repo}
}

func (uq *UpdateQuiniela) Execute(id int32, nombreVariante string, precio money.Money, montoMeta money.Money, fechaLimite time.Time, estado string) error {
	return uq.repo.UpdateQuiniela(id, nombreVariante, precio, montoMeta, fechaLimite, estado)
}