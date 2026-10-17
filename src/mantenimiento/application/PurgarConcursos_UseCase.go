package application

import (
	"log"
	"time"

	"quinielas/src/mantenimiento/domain"
)

// PurgarConcursos elimina la informacion de los concursos que ya cerraron hace
// mas de N meses (retencion), para que la base de datos no crezca sin limite.
// No toca los saldos de los clientes: el saldo es el balance actual y borrar
// historial viejo no debe cambiarlo.
type PurgarConcursos struct {
	repo  domain.IPurga
	meses int
}

func NewPurgarConcursos(repo domain.IPurga, meses int) *PurgarConcursos {
	return &PurgarConcursos{repo: repo, meses: meses}
}

// Procesar elimina los concursos cerrados antes del limite (ahora - N meses).
// Es idempotente: al eliminar los concursos dejan de aparecer en la busqueda.
func (p *PurgarConcursos) Procesar(ahora time.Time) error {
	if p.meses <= 0 {
		return nil
	}

	limite := ahora.AddDate(0, -p.meses, 0)
	eliminados, err := p.repo.PurgarConcursosAnteriores(limite)
	if err != nil {
		return err
	}
	if eliminados > 0 {
		log.Printf("[purga] concursos eliminados: %d (cerrados antes de %s)", eliminados, limite.Format("2006-01-02"))
	}
	return nil
}
