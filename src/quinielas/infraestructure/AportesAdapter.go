package infraestructure

import (
	aportesDomain "quinielas/src/aportes/domain"
	"quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
)

// AportesAdapter adapta el repositorio de aportes al puerto que consume el
// ciclo de vida de la quiniela (domain.AportesDeQuiniela).
type AportesAdapter struct {
	repo aportesDomain.IAporte
}

var _ domain.AportesDeQuiniela = (*AportesAdapter)(nil)

func NewAportesAdapter(repo aportesDomain.IAporte) *AportesAdapter {
	return &AportesAdapter{repo: repo}
}

func (a *AportesAdapter) DeleteAportesByQuinielaID(idQuiniela int32) error {
	return a.repo.DeleteAportesByQuinielaID(idQuiniela)
}

func (a *AportesAdapter) GetAportesByQuinielaID(idQuiniela int32) ([]domain.AporteResumen, error) {
	lista, err := a.repo.GetAportesByQuinielaID(idQuiniela)
	if err != nil {
		return nil, err
	}
	resumen := make([]domain.AporteResumen, 0, len(lista))
	for _, ad := range lista {
		resumen = append(resumen, domain.AporteResumen{
			IDCliente:           ad.IDCliente,
			MontoTotalAcumulado: ad.MontoTotalAcumulado,
		})
	}
	return resumen, nil
}

func (a *AportesAdapter) ReducirAporte(idCliente int32, idQuiniela int32, monto money.Money, meta money.Money) error {
	return a.repo.ReducirAporte(idCliente, idQuiniela, monto, meta)
}

func (a *AportesAdapter) AgregarAporte(idCliente int32, idQuiniela int32, monto money.Money, meta money.Money) error {
	_, err := a.repo.SaveOrUpdateAporte(idCliente, idQuiniela, monto, meta)
	return err
}
