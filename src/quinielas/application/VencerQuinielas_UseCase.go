package application

import (
	"fmt"
	"log"
	"time"

	"quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
)

// VencerQuinielas aplica la consecuencia de la fecha límite a las quinielas
// DESARROLLO vencidas, según su estado_fallo:
//
//	1 -> se cancela (lo pagado vuelve a favor, lo fiado se descarta)
//	2 -> el creador absorbe el faltante (FIADO)
//	3 -> se devuelve el excedente sobre la meta previa y se vuelve al precio anterior
//	4 -> el cliente absorbe el faltante del aumento (FIADO)
type VencerQuinielas struct {
	repo    domain.IQuiniela
	cobro   domain.CobroDirecta
	aportes domain.AportesDeQuiniela
}

func NewVencerQuinielas(repo domain.IQuiniela, cobro domain.CobroDirecta, aportes domain.AportesDeQuiniela) *VencerQuinielas {
	return &VencerQuinielas{repo: repo, cobro: cobro, aportes: aportes}
}

// Procesar procesa las quinielas vencidas. Es idempotente: al cambiar el estado
// dejan de aparecer en GetVencidasDesarrollo.
func (v *VencerQuinielas) Procesar(ahora time.Time) error {
	vencidas, err := v.repo.GetVencidasDesarrollo(ahora)
	if err != nil {
		return err
	}

	for i := range vencidas {
		if err := v.procesar(vencidas[i]); err != nil {
			log.Printf("[vencimientos] error al procesar la quiniela #%d: %v", vencidas[i].ID, err)
		}
	}
	return nil
}

func (v *VencerQuinielas) totalAportado(idQuiniela int32) (money.Money, error) {
	aportes, err := v.aportes.GetAportesByQuinielaID(idQuiniela)
	if err != nil {
		return money.Zero(), err
	}
	total := money.Zero()
	for _, a := range aportes {
		total = total.Add(a.MontoTotalAcumulado)
	}
	return total, nil
}

func (v *VencerQuinielas) procesar(q domain.Quiniela) error {
	recaudado, err := v.totalAportado(q.ID)
	if err != nil {
		return err
	}

	// Si ya se alcanzó la meta, se marca como completada (no aplica el fallo).
	if !q.MontoMeta.GreaterThan(recaudado) {
		q.Estado = "COMPLETADA"
		return v.repo.UpdateQuiniela(&q)
	}

	fallo := q.EstadoFallo
	if fallo == 0 {
		if q.SujetoCancelacion {
			fallo = 1
		} else {
			fallo = 2
		}
	}

	switch fallo {
	case 1:
		return v.cancelar(&q)
	case 3:
		return v.devolverExcedente(&q)
	default: // 2, 4 y cualquier otro: el creador absorbe el faltante
		return v.absorberCreador(&q, recaudado)
	}
}

func (v *VencerQuinielas) cancelar(q *domain.Quiniela) error {
	if v.cobro != nil {
		if err := v.cobro.ReembolsarCobro(q.ID); err != nil {
			return err
		}
	}
	q.Estado = "CANCELADA"
	return v.repo.UpdateQuiniela(q)
}

func (v *VencerQuinielas) absorberCreador(q *domain.Quiniela, recaudado money.Money) error {
	faltante := q.MontoMeta.Sub(recaudado)
	if q.IDCliente == nil {
		log.Printf("Quiniela #%d sin creador (id_cliente nulo); no se puede absorber el faltante de %s", q.ID, faltante.String())
		q.Estado = "COMPLETADA"
		return v.repo.UpdateQuiniela(q)
	}
	if faltante.GreaterThan(money.Zero()) && v.cobro != nil {
		descripcion := fmt.Sprintf("Absorcion de faltante por cierre de la quiniela #%d", q.ID)
		if err := v.cobro.RegistrarFiado(*q.IDCliente, q.ID, faltante, descripcion); err != nil {
			return err
		}
	}
	q.Estado = "COMPLETADA"
	return v.repo.UpdateQuiniela(q)
}

func (v *VencerQuinielas) devolverExcedente(q *domain.Quiniela) error {
	if v.cobro == nil || q.MontoMetaAnterior == nil || q.FechaMetaAnterior == nil {
		// Sin aumento registrado: se trata como una cancelación normal.
		return v.cancelar(q)
	}

	excedentes, err := v.cobro.ReembolsarCobroDesde(q.ID, *q.FechaMetaAnterior)
	if err != nil {
		return err
	}

	metaAnterior := *q.MontoMetaAnterior
	for idCliente, monto := range excedentes {
		if err := v.aportes.ReducirAporte(idCliente, q.ID, monto, metaAnterior); err != nil {
			return err
		}
	}

	q.MontoMeta = metaAnterior
	q.Estado = "COMPLETADA"
	return v.repo.UpdateQuiniela(q)
}
