package application

import (
	aportesDomain "quinielas/src/aportes/domain"
	quinielasDomain "quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
)

// UpdateAporte reemplaza el monto de un aporte existente. El nuevo monto no
// puede hacer que el total recaudado de la quiniela supere su meta.
type UpdateAporte struct {
	repo         aportesDomain.IAporte
	quinielaRepo quinielasDomain.IQuiniela
	cobro        aportesDomain.CobroAporte
}

func NewUpdateAporte(repo aportesDomain.IAporte, quinielaRepo quinielasDomain.IQuiniela, cobro aportesDomain.CobroAporte) *UpdateAporte {
	return &UpdateAporte{repo: repo, quinielaRepo: quinielaRepo, cobro: cobro}
}

func (ua *UpdateAporte) Execute(idAporte int32, montoNuevo money.Money) (*aportesDomain.Aporte, error) {
	aporte, err := ua.repo.GetAporteByID(idAporte)
	if err != nil {
		return nil, err
	}

	quiniela, err := ua.quinielaRepo.GetQuinielaByID(aporte.IDQuiniela)
	if err != nil {
		return nil, aportesDomain.ErrQuinielaNoEncontrada
	}

	// Total recaudado actual (incluye este aporte). Al editarlo se descuenta el
	// monto viejo y se suma el nuevo para validar contra la meta.
	totalActual, err := ua.totalRecaudado(aporte.IDQuiniela)
	if err != nil {
		return nil, err
	}
	totalProyectado := totalActual.Sub(aporte.MontoTotalAcumulado).Add(montoNuevo)

	if totalProyectado.GreaterThan(quiniela.MontoMeta) {
		return nil, aportesDomain.ErrMetaExcedida
	}

	resultado, err := ua.repo.UpdateAporte(idAporte, montoNuevo, quiniela.MontoMeta)
	if err != nil {
		return nil, err
	}

	// El movimiento del aporte se ajusta al nuevo monto acumulado.
	if ua.cobro != nil {
		if err := ua.cobro.SincronizarAporte(resultado.IDCliente, resultado.IDQuiniela, resultado.MontoTotalAcumulado, quiniela.CanalVenta); err != nil {
			return nil, err
		}
	}

	return resultado, nil
}

func (ua *UpdateAporte) totalRecaudado(idQuiniela int32) (money.Money, error) {
	aportes, err := ua.repo.GetAportesByQuinielaID(idQuiniela)
	if err != nil {
		return money.Zero(), err
	}
	total := money.Zero()
	for _, a := range aportes {
		total = total.Add(a.MontoTotalAcumulado)
	}
	return total, nil
}
