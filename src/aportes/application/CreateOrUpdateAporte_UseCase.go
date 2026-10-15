package application

import (
	aportesDomain "quinielas/src/aportes/domain"
	quinielasDomain "quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
)

type CreateOrUpdateAporte struct {
	repo         aportesDomain.IAporte
	quinielaRepo quinielasDomain.IQuiniela
	cobro        aportesDomain.CobroAporte
}

func NewCreateOrUpdateAporte(repo aportesDomain.IAporte, quinielaRepo quinielasDomain.IQuiniela, cobro aportesDomain.CobroAporte) *CreateOrUpdateAporte {
	return &CreateOrUpdateAporte{repo: repo, quinielaRepo: quinielaRepo, cobro: cobro}
}

func (cua *CreateOrUpdateAporte) Execute(idCliente int32, idQuiniela int32, montoSumar money.Money, montoPagado money.Money) (*aportesDomain.Aporte, error) {
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

	aporte, err := cua.repo.SaveOrUpdateAporte(idCliente, idQuiniela, montoSumar, quiniela.MontoMeta)
	if err != nil {
		return nil, err
	}

	// Cada aporte se registra en movimientos de forma incremental (con su canal
	// de venta): el monto aportado como FIADO y lo pagado como PAGO_EFECTIVO que
	// lo reduce, ajustando el saldo del cliente. Un aporte de contado
	// (montoPagado == montoSumar) no altera el saldo; uno a crédito (montoPagado
	// = 0) genera solo el FIADO.
	if cua.cobro != nil {
		if err := cua.cobro.AplicarAporte(aporte.IDCliente, aporte.IDQuiniela, montoSumar, montoPagado, quiniela.CanalVenta); err != nil {
			return nil, err
		}
	}

	return aporte, nil
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