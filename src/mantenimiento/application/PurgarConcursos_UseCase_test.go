package application

import (
	"testing"
	"time"

	"quinielas/src/mantenimiento/domain"
)

type repoFalso struct {
	llamado  bool
	limite   time.Time
	devuelve int
	err      error
}

var _ domain.IPurga = (*repoFalso)(nil)

func (r *repoFalso) PurgarConcursosAnteriores(limite time.Time) (int, error) {
	r.llamado = true
	r.limite = limite
	return r.devuelve, r.err
}

func TestProcesarCalculaLimiteDeSeisMeses(t *testing.T) {
	repo := &repoFalso{devuelve: 2}
	uc := NewPurgarConcursos(repo, 6)

	ahora := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)
	if err := uc.Procesar(ahora); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if !repo.llamado {
		t.Fatal("se esperaba que el caso de uso llamara al repositorio")
	}
	esperado := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)
	if !repo.limite.Equal(esperado) {
		t.Errorf("limite = %v, se esperaba %v", repo.limite, esperado)
	}
}

func TestProcesarConRetencionInvalidaNoHaceNada(t *testing.T) {
	repo := &repoFalso{}
	uc := NewPurgarConcursos(repo, 0)

	if err := uc.Procesar(time.Now()); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repo.llamado {
		t.Error("no deberia llamar al repositorio cuando los meses no son positivos")
	}
}
