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

type abonoRegistro struct {
	idCliente  int32
	idQuiniela int32
	monto      money.Money
}

type fakeAbonos struct {
	abonos    []abonoRegistro
	deshechos []int32
}

func (f *fakeAbonos) Abonar(idCliente int32, idQuiniela int32, monto money.Money) (int32, error) {
	f.abonos = append(f.abonos, abonoRegistro{idCliente, idQuiniela, monto})
	return int32(len(f.abonos)), nil
}

func (f *fakeAbonos) DeshacerAbono(idMovimiento int32) error {
	f.deshechos = append(f.deshechos, idMovimiento)
	return nil
}

func montarRouterConAbonos(p *fakePremioRepo, a *fakeAporteRepo, q *fakeQuinielaRepo, ab *fakeAbonos) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	SetupRouterPremios(r, p, a, q, ab)
	return r
}

func montarRouter(p *fakePremioRepo, a *fakeAporteRepo, q *fakeQuinielaRepo) *gin.Engine {
	return montarRouterConAbonos(p, a, q, &fakeAbonos{})
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

// TestGenerarRepartoAbonaSaldos verifica que cada ganador reciba un movimiento
// de abono por su monto neto asignado.
func TestGenerarRepartoAbonaSaldos(t *testing.T) {
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
	ab := &fakeAbonos{}

	router := montarRouterConAbonos(repo, aportes, quiniela, ab)

	req := httptest.NewRequest(http.MethodPost, "/premios/quiniela/1/reparto", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("codigo esperado 201, obtenido %d. Cuerpo: %s", rec.Code, rec.Body.String())
	}
	if len(ab.abonos) != 3 {
		t.Fatalf("esperaba 3 abonos, hubo %d", len(ab.abonos))
	}

	esperados := []struct {
		idCliente int32
		monto     string
	}{
		{1, "465.00"},
		{2, "279.00"},
		{3, "186.00"},
	}
	for i, e := range esperados {
		if ab.abonos[i].idCliente != e.idCliente {
			t.Errorf("abono %d: cliente esperado %d, obtenido %d", i, e.idCliente, ab.abonos[i].idCliente)
		}
		if ab.abonos[i].monto.String() != e.monto {
			t.Errorf("abono %d: monto esperado %s, obtenido %s", i, e.monto, ab.abonos[i].monto.String())
		}
		if ab.abonos[i].idQuiniela != 1 {
			t.Errorf("abono %d: quiniela esperada 1, obtenida %d", i, ab.abonos[i].idQuiniela)
		}
	}
}

// TestGenerarRepartoDeshaceAbonosSiFallaElGuardado verifica el rollback: si el
// snapshot del reparto no se puede guardar, los abonos ya creados se deshacen.
func TestGenerarRepartoDeshaceAbonosSiFallaElGuardado(t *testing.T) {
	premio := domain.NewPremio(1, money.NewFromInt64(1000))
	premio.IDPremio = 7

	repo := &fakePremioRepo{premio: premio, fallaGuardarReparto: true}
	aportes := &fakeAporteRepo{aportes: []aportesDomain.AporteDetalle{
		aporte(1, "Ana", "5000.00"),
		aporte(2, "Luis", "3000.00"),
		aporte(3, "Sofia", "2000.00"),
	}}
	quiniela := &fakeQuinielaRepo{quiniela: &quinielasDomain.Quiniela{
		ID:        1,
		MontoMeta: money.MustParse("10000.00"),
	}}
	ab := &fakeAbonos{}

	router := montarRouterConAbonos(repo, aportes, quiniela, ab)

	req := httptest.NewRequest(http.MethodPost, "/premios/quiniela/1/reparto", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("codigo esperado 500, obtenido %d. Cuerpo: %s", rec.Code, rec.Body.String())
	}
	if len(ab.abonos) != 3 {
		t.Fatalf("esperaba 3 abonos creados antes de fallar, hubo %d", len(ab.abonos))
	}
	if len(ab.deshechos) != 3 {
		t.Errorf("esperaba deshacer los 3 abonos, se deshicieron %d", len(ab.deshechos))
	}
}
