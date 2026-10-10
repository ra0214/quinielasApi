package application

import (
	"database/sql"
	"errors"

	"quinielas/src/premios/domain"
	"quinielas/src/shared/money"
)

type UpdatePremio struct {
	repo domain.IPremio
}

func NewUpdatePremio(repo domain.IPremio) *UpdatePremio {
	return &UpdatePremio{repo: repo}
}

// Execute modifica el monto bruto de un premio existente y recalcula la
// retención de 7% y el monto neto que se reparte.
func (up *UpdatePremio) Execute(idPremio int32, montoBruto money.Money) (*domain.Premio, error) {
	premio, err := up.repo.UpdatePremio(idPremio, montoBruto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPremioNoEncontrado
		}
		return nil, err
	}
	return premio, nil
}

// ErrPremioNoEncontrado indica que el id de premio no existe.
var ErrPremioNoEncontrado = errors.New("el premio no existe")
