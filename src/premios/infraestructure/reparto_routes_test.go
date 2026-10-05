package infraestructure

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	aportesDomain "quinielas/src/aportes/domain"
	"quinielas/src/premios/domain"
	quinielasDomain "quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"

	"github.com/gin-gonic/gin"
)

type fakeAporteRepo struct {
	aportes []aportesDomain.AporteDetalle
}

type fakeQuinielaRepo struct {
	quiniela *quinielasDomain.Quiniela
}

func (f *fakeAporteRepo) GetAportesByQuinielaID(id int32) ([]aportesDomain.AporteDetalle, error) {
	return f.aportes, nil
}

func (f *fakeQuinielaRepo) GetQuinielaByID(id int32) (*quinielasDomain.Quiniela, error) {
	return f.quiniela, nil
}

func montarRouter(p *fakePremioRepo, a *fakeAporteRepo, q *fakeQuinielaRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	SetupRouterPremios(r, p, a, q)
	return r
}

func aporte(id int32, nombre string, monto string) aportesDomain.AporteDetalle {
	return aportesDomain.AporteDetalle{
		IDCliente:           id,
		NombreCliente:       nombre,
		MontoTotalAcumulado: money.MustParse(monto),
	}
}

// TestRutasDeRepartoNoEntranEnPanic verifica que /premios/quiniela/:id y
// /premios/quiniela/:id/reparto conviven en el arbol de gin.
func TestRutasDeRepartoNoEntranEnPanic(t *testing.T) {
	router := montarRouter(&fakePremioRepo{}, &fakeAporteRepo{}, &fakeQuinielaRepo{})

	var encontradas int
	for _, r := range router.Routes() {
		if r.Path == "/premios/quiniela/:idQuiniela/reparto" {
			encontradas++
		}
	}

	if encontradas != 2 {
		t.Errorf("esperaba 2 rutas de reparto (GET y POST), encontré %d", encontradas)
	}
}

func TestGenerarRepartoPorHTTP(t *testing.T) {
	premio := domain.NewPremio(1, money.NewFromInt64(1000))
	premio.IDPremio = 7

	repo := &fakePremioRepo{premio: premio}
	aportes := &fakeAporteRepo{aportes: []aportesDomain.AporteDetalle{
		aporte(1, "Ana", "5000.00"),
		aporte(2, "Luis", "3000.00"),
		aporte(3, "Sofia", "2000.00"),
	}}
	quiniela := &fakeQuinielaRepo{quiniela: &quinielasDomain.Quiniela{
		ID:        1,
		MontoMeta: money.MustParse("10000.00"),
	}}

	router := montarRouter(repo, aportes, quiniela)

	req := httptest.NewRequest(http.MethodPost, "/premios/quiniela/1/reparto", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("codigo esperado 201, obtenido %d. Cuerpo: %s", rec.Code, rec.Body.String())
	}
	if repo.guardados != 1 {
		t.Errorf("se esperaba 1 guardado del snapshot, hubo %d", repo.guardados)
	}

	var respuesta struct {
		Reparto struct {
			MontoNeto          json.Number `json:"monto_neto"`
			TotalRecaudado     json.Number `json:"total_recaudado"`
			TotalParticipantes int         `json:"total_participantes"`
			Participantes      []struct {
				NombreCliente      string      `json:"nombre_cliente"`
				MontoAporte        json.Number `json:"monto_aporte"`
				PorcentajeAporte   json.Number `json:"porcentaje_aporte"`
				MontoBrutoAsignado json.Number `json:"monto_bruto_asignado"`
				MontoRetencion     json.Number `json:"monto_retencion"`
				MontoNetoAsignado  json.Number `json:"monto_neto_asignado"`
			} `json:"participantes"`
		} `json:"reparto"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &respuesta); err != nil {
		t.Fatalf("no se pudo decodificar: %v", err)
	}

	r := respuesta.Reparto
	if r.MontoNeto.String() != "930.00" {
		t.Errorf("monto_neto esperado 930.00, obtenido %s", r.MontoNeto.String())
	}
	if r.TotalRecaudado.String() != "10000.00" {
		t.Errorf("total recaudado esperado 10000.00, obtenido %s", r.TotalRecaudado.String())
	}
	if r.TotalParticipantes != 3 {
		t.Errorf("esperaba 3 participantes, obtenido %d", r.TotalParticipantes)
	}

	esperado := []struct {
		nombre     string
		porcentaje string
		bruto      string
		retencion  string
		neto       string
	}{
		{"Ana", "50.0000", "500.00", "35.00", "465.00"},
		{"Luis", "30.0000", "300.00", "21.00", "279.00"},
		{"Sofia", "20.0000", "200.00", "14.00", "186.00"},
	}

	if len(r.Participantes) != len(esperado) {
		t.Fatalf("esperaba %d participantes en la respuesta, llegaron %d", len(esperado), len(r.Participantes))
	}

	for i, e := range esperado {
		p := r.Participantes[i]
		if p.NombreCliente != e.nombre {
			t.Errorf("participante %d: esperaba %s, obtenido %s", i, e.nombre, p.NombreCliente)
			continue
		}
		if p.PorcentajeAporte.String() != e.porcentaje {
			t.Errorf("%s: porcentaje esperado %s, obtenido %s", e.nombre, e.porcentaje, p.PorcentajeAporte.String())
		}
		if p.MontoBrutoAsignado.String() != e.bruto {
			t.Errorf("%s: bruto esperado %s, obtenido %s", e.nombre, e.bruto, p.MontoBrutoAsignado.String())
		}
		if p.MontoRetencion.String() != e.retencion {
			t.Errorf("%s: retencion esperada %s, obtenido %s", e.nombre, e.retencion, p.MontoRetencion.String())
		}
		if p.MontoNetoAsignado.String() != e.neto {
			t.Errorf("%s: neto esperado %s, obtenido %s", e.nombre, e.neto, p.MontoNetoAsignado.String())
		}
	}

	t.Logf("respuesta: %s", rec.Body.String())
}

func TestGenerarRepartoRechazaMetaNoAlcanzada(t *testing.T) {
	premio := domain.NewPremio(1, money.NewFromInt64(1000))
	premio.IDPremio = 7

	repo := &fakePremioRepo{premio: premio}
	aportes := &fakeAporteRepo{aportes: []aportesDomain.AporteDetalle{
		aporte(1, "Ana", "1000.00"),
		aporte(2, "Luis", "2000.00"),
	}}
	quiniela := &fakeQuinielaRepo{quiniela: &quinielasDomain.Quiniela{
		ID:        1,
		MontoMeta: money.MustParse("5000.00"),
	}}

	router := montarRouter(repo, aportes, quiniela)

	req := httptest.NewRequest(http.MethodPost, "/premios/quiniela/1/reparto", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("codigo esperado 409, obtenido %d. Cuerpo: %s", rec.Code, rec.Body.String())
	}
	if repo.guardados != 0 {
		t.Error("no se debe guardar el snapshot si la meta no fue alcanzada")
	}
}

func TestGenerarRepartoNoReparteDosVeces(t *testing.T) {
	premio := domain.NewPremio(1, money.NewFromInt64(1000))
	premio.IDPremio = 7

	repo := &fakePremioRepo{premio: premio, yaRepartido: true}
	aportes := &fakeAporteRepo{aportes: []aportesDomain.AporteDetalle{aporte(1, "Ana", "10000.00")}}
	quiniela := &fakeQuinielaRepo{quiniela: &quinielasDomain.Quiniela{
		ID:        1,
		MontoMeta: money.MustParse("10000.00"),
	}}

	router := montarRouter(repo, aportes, quiniela)

	req := httptest.NewRequest(http.MethodPost, "/premios/quiniela/1/reparto", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("codigo esperado 409, obtenido %d. Cuerpo: %s", rec.Code, rec.Body.String())
	}
	if repo.guardados != 0 {
		t.Error("no se debe volver a guardar un reparto ya existente")
	}
}
