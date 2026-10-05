package application

import "quinielas/src/aportes/domain"

type ViewAporte struct {
	repo domain.IAporte
}

func NewViewAporte(repo domain.IAporte) *ViewAporte {
	return &ViewAporte{repo: repo}
}

func (va *ViewAporte) ExecuteByQuiniela(idQuiniela int32) ([]domain.AporteDetalle, error) {
	return va.repo.GetAportesByQuinielaID(idQuiniela)
}

func (va *ViewAporte) ExecuteByCliente(idCliente int32) ([]domain.Aporte, error) {
	return va.repo.GetAportesByClienteID(idCliente)
}