package application

import (
	"database/sql"
	"errors"

	"quinielas/src/premios/domain"
	"quinielas/src/shared/money"
)

type UpdatePremio struct {
	repo         domain.IPremio
	quinielaRepo QuinielaDePremio
}

func NewUpdatePremio(repo domain.IPremio, quinielaRepo QuinielaDePremio) *UpdatePremio {
	return &UpdatePremio{repo: repo, quinielaRepo: quinielaRepo}
}

// Execute modifica el monto de un premio existente y recalcula la retención y
// el monto neto. En DIRECTA_EN_COPERACHA el monto capturado es por boleto.
func (up *UpdatePremio) Execute(idPremio int32, montoBruto money.Money, premioPorBoleto money.Money) (*domain.Premio, error) {
	premio, err := up.repo.GetPremioByID(idPremio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPremioNoEncontrado
		}
		return nil, err
	}

	monto := montoBruto
	fijo := premioPorBoleto

	quiniela, err := up.quinielaRepo.GetQuinielaByID(premio.IDQuiniela)
	if err == nil && quiniela.TipoQuiniela == "DIRECTA_EN_COPERACHA" {
		if premioPorBoleto.IsZero() {
			return nil, domain.ErrPremioSinMontoPorBoleto
		}
		cantidad := quiniela.CantidadBoletos
		if cantidad <= 0 {
			cantidad = 1
		}
		monto = premioPorBoleto.MulInt(int64(cantidad))
		fijo = premioPorBoleto
	}

	premio, err = up.repo.UpdatePremio(idPremio, monto, fijo)
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