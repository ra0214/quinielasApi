package application

import (
	quinielasDomain "quinielas/src/quinielas/domain"
	"quinielas/src/premios/domain"
	"quinielas/src/shared/money"
)

type CreatePremio struct {
	repo         domain.IPremio
	quinielaRepo QuinielaDePremio
}

func NewCreatePremio(repo domain.IPremio, quinielaRepo QuinielaDePremio) *CreatePremio {
	return &CreatePremio{repo: repo, quinielaRepo: quinielaRepo}
}

// QuinielaDePremio es lo único que el registro del premio necesita de una
// quiniela: saber su tipo y su cantidad de boletos.
type QuinielaDePremio interface {
	GetQuinielaByID(id int32) (*quinielasDomain.Quiniela, error)
}

// Execute registra el premio. En DIRECTA_EN_COPERACHA el monto que se captura
// es el premio por boleto (el mismo para todos los ganadores), así que se
// guarda ese valor y como bruto el total nominal (premio x cantidad de boletos).
// En DIRECTA y DESARROLLO se registra el monto bruto con la retención estándar.
func (cp *CreatePremio) Execute(idQuiniela int32, montoBruto money.Money, premioPorBoleto money.Money) (*domain.Premio, error) {
	quiniela, err := cp.quinielaRepo.GetQuinielaByID(idQuiniela)
	if err != nil {
		return nil, err
	}

	if quiniela.TipoQuiniela == "DIRECTA_EN_COPERACHA" {
		if premioPorBoleto.IsZero() {
			return nil, domain.ErrPremioSinMontoPorBoleto
		}
		cantidad := quiniela.CantidadBoletos
		if cantidad <= 0 {
			cantidad = 1
		}
		total := premioPorBoleto.MulInt(int64(cantidad))
		return cp.repo.SavePremio(idQuiniela, total, premioPorBoleto)
	}

	if montoBruto.IsZero() {
		return nil, domain.ErrMetaInvalida
	}
	return cp.repo.SavePremio(idQuiniela, montoBruto, money.Zero())
}