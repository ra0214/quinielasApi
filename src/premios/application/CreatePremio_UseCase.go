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

// Execute registra el premio de una quiniela. El monto bruto es el que se
// captura en pantalla y sobre él se calcula la retención estándar (7%). Si en
// su lugar se captura el monto neto, este se usa tal cual (sin retención).
// Ambos montos son opcionales, pero al menos uno debe ser mayor a 0.
func (cp *CreatePremio) Execute(idQuiniela int32, montoBruto money.Money, montoNeto money.Money) (*domain.Premio, error) {
	return cp.repo.SavePremio(idQuiniela, montoBruto, montoNeto)
}
