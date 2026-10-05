package infraestructure

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"quinielas/src/premios/application"
	"quinielas/src/premios/domain"
	"quinielas/src/shared/money"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakePremioRepo struct {
	guardado *domain.Premio

	// Estado que usan las pruebas de reparto.
	premio      *domain.Premio
	reparto     *domain.RepartoPremio
	guardados   int
	yaRepartido bool
}

func (f *fakePremioRepo) SavePremio(idQuiniela int32, montoBruto money.Money) (*domain.Premio, error) {
	p := domain.NewPremio(idQuiniela, montoBruto)
	p.IDPremio = 1
	f.guardado = p
	return p, nil
}

func (f *fakePremioRepo) DeletePremio(idPremio int32) error { return nil }

func (f *fakePremioRepo) GetPremioByQuinielaID(idQuiniela int32) (*domain.Premio, error) {
	return f.premio, nil
}

func (f *fakePremioRepo) GetAllPremios() ([]domain.Premio, error) { return nil, nil }

func (f *fakePremioRepo) SaveReparto(reparto *domain.RepartoPremio) error {
	f.reparto = reparto
	f.guardados++
	return nil
}

func (f *fakePremioRepo) GetRepartoByQuiniela(idQuiniela int32) (*domain.RepartoPremio, error) {
	return f.reparto, nil
}

func (f *fakePremioRepo) ExistsReparto(idQuiniela int32) (bool, error) {
	return f.yaRepartido, nil
}

func newController(repo domain.IPremio) *CreatePremioController {
	return NewCreatePremioController(application.NewCreatePremio(repo))
}

func TestPostPremioDevuelveMontoNetoExacto(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakePremioRepo{}
	router := gin.New()
	router.POST("/premios", newController(repo).Execute)

	// Cuerpo identico al que reportaste
	cuerpo := `{"id_quiniela": 1, "monto_bruto": 1000.00}`

	req := httptest.NewRequest(http.MethodPost, "/premios", strings.NewReader(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("codigo esperado 201, obtenido %d. Cuerpo: %s", rec.Code, rec.Body.String())
	}

	var respuesta struct {
		Message string `json:"message"`
		Premio  struct {
			IDPremio            int32       `json:"id_premio"`
			IDQuiniela          int32       `json:"id_quiniela"`
			MontoBruto          json.Number `json:"monto_bruto"`
			PorcentajeRetencion json.Number `json:"porcentaje_retencion"`
			MontoNeto           json.Number `json:"monto_neto"`
		} `json:"premio"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &respuesta); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}

	if respuesta.Premio.MontoNeto.String() != "930.00" {
		t.Errorf("monto_neto esperado 930.00, obtenido %q", respuesta.Premio.MontoNeto.String())
	}
	if respuesta.Premio.MontoBruto.String() != "1000.00" {
		t.Errorf("monto_bruto esperado 1000.00, obtenido %q", respuesta.Premio.MontoBruto.String())
	}
	if respuesta.Premio.PorcentajeRetencion.String() != "7.0000" {
		t.Errorf("porcentaje_retencion esperado 7.0000, obtenido %q", respuesta.Premio.PorcentajeRetencion.String())
	}

	if strings.Contains(rec.Body.String(), "9999") || strings.Contains(rec.Body.String(), "929.99") {
		t.Errorf("la respuesta aun contiene ruido de punto flotante: %s", rec.Body.String())
	}

	t.Logf("respuesta: %s", rec.Body.String())
}

func TestPostPremioRechazaMontoCeroYNegativo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	casos := []struct {
		nombre string
		cuerpo string
		codigo int
	}{
		{"cero", `{"id_quiniela": 1, "monto_bruto": 0}`, http.StatusBadRequest},
		{"negativo", `{"id_quiniela": 1, "monto_bruto": -50}`, http.StatusBadRequest},
		{"ausente", `{"id_quiniela": 1}`, http.StatusBadRequest},
	}

	for _, c := range casos {
		router := gin.New()
		router.POST("/premios", newController(&fakePremioRepo{}).Execute)

		req := httptest.NewRequest(http.MethodPost, "/premios", strings.NewReader(c.cuerpo))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != c.codigo {
			t.Errorf("%s: codigo esperado %d, obtenido %d", c.nombre, c.codigo, rec.Code)
		}
	}
}
