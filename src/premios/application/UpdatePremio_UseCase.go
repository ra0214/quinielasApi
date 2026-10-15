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

// Execute modifica el monto de un premio existente y recalcula la retención
// (7% si se captura el bruto, 0% si se captura el neto) y el monto neto.
func (up *UpdatePremio) Execute(idPremio int32, montoBruto money.Money, montoNeto money.Money) (*domain.Premio, error) {
	premio, err := up.repo.UpdatePremio(idPremio, montoBruto, montoNeto)
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
