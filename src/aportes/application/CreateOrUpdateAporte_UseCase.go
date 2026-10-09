package application

import (
	aportesDomain "quinielas/src/aportes/domain"
	quinielasDomain "quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
)

type CreateOrUpdateAporte struct {
	repo         aportesDomain.IAporte
	quinielaRepo quinielasDomain.IQuiniela
}

func NewCreateOrUpdateAporte(repo aportesDomain.IAporte, quinielaRepo quinielasDomain.IQuiniela) *CreateOrUpdateAporte {
	return &CreateOrUpdateAporte{repo: repo, quinielaRepo: quinielaRepo}
}

func (cua *CreateOrUpdateAporte) Execute(idCliente int32, idQuiniela int32, montoSumar money.Money) (*aportesDomain.Aporte, error) {
	// La meta (tope) sale de la quiniela real en BD, no de lo que mande el cliente.
	quiniela, err := cua.quinielaRepo.GetQuinielaByID(idQuiniela)
	if err != nil {
		return nil, aportesDomain.ErrQuinielaNoEncontrada
	}

	// Suma de lo ya recaudado en esta quiniela (no superar el total solicitado).
	totalActual, err := cua.totalRecaudado(idQuiniela)
	if err != nil {
		return nil, err
	}

	if montoSumar.GreaterThan(quiniela.MontoMeta.Sub(totalActual)) {
		return nil, aportesDomain.ErrMetaExcedida
	}

	return cua.repo.SaveOrUpdateAporte(idCliente, idQuiniela, montoSumar, quiniela.MontoMeta)
}

func (cua *CreateOrUpdateAporte) totalRecaudado(idQuiniela int32) (money.Money, error) {
	aportes, err := cua.repo.GetAportesByQuinielaID(idQuiniela)
	if err != nil {
		return money.Zero(), err
	}
	total := money.Zero()
	for _, a := range aportes {
		total = total.Add(a.MontoTotalAcumulado)
	}
	return total, nil
}