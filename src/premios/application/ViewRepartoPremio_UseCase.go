package application

import (
	"database/sql"
	"errors"
	"quinielas/src/premios/domain"
)

var ErrRepartoNoGenerado = errors.New("esta quiniela todavia no tiene un reparto generado")

type ViewRepartoPremio struct {
	premioRepo domain.IPremio
}

func NewViewRepartoPremio(premioRepo domain.IPremio) *ViewRepartoPremio {
	return &ViewRepartoPremio{premioRepo: premioRepo}
}

// ExecuteByQuiniela devuelve el reparto ya generado. No recalcula nada: lee el
// snapshot, de modo que lo que se consulto es siempre lo que se agreed a pagar.
func (v *ViewRepartoPremio) ExecuteByQuiniela(idQuiniela int32) (*domain.RepartoPremio, error) {
	reparto, err := v.premioRepo.GetRepartoByQuiniela(idQuiniela)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRepartoNoGenerado
		}
		return nil, err
	}
	return reparto, nil
}
