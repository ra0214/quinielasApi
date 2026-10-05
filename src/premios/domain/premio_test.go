package domain

import (
	"bytes"
	"encoding/json"
	"quinielas/src/shared/money"
	"testing"
)

func TestNewPremioAplicaRetencionDelSietePorCiento(t *testing.T) {
	bruto := money.NewFromInt64(1000)
	premio := NewPremio(1, bruto)

	if premio.MontoNeto.String() != "930.00" {
		t.Errorf("esperado monto_neto 930.00, obtenido %s", premio.MontoNeto.String())
	}
	if premio.PorcentajeRetencion.String() != "7.0000" {
		t.Errorf("esperado porcentaje_retencion 7.0000, obtenido %s", premio.PorcentajeRetencion.String())
	}
	if premio.MontoBruto.String() != "1000.00" {
		t.Errorf("esperado monto_bruto 1000.00, obtenido %s", premio.MontoBruto.String())
	}
}

func TestNewPremioNetoMenosRetencionRecuperaElBruto(t *testing.T) {
	casos := []struct {
		bruto     string
		retencion string
		neto      string
	}{
		{"0.01", "0.00", "0.01"},
		{"1.00", "0.07", "0.93"},
		{"19.99", "1.40", "18.59"},
		{"100.00", "7.00", "93.00"},
		{"1000.00", "70.00", "930.00"},
		{"55555.55", "3888.89", "51666.66"},
		{"999999.99", "70000.00", "929999.99"},
	}

	for _, c := range casos {
		bruto, err := money.Parse(c.bruto)
		if err != nil {
			t.Fatalf("Parse(%q) error inesperado: %v", c.bruto, err)
		}

		premio := NewPremio(1, bruto)

		if premio.MontoNeto.String() != c.neto {
			t.Errorf("bruto %s: esperado neto %s, obtenido %s", c.bruto, c.neto, premio.MontoNeto.String())
		}

		retenido := premio.MontoBruto.Sub(premio.MontoNeto)
		if retenido.String() != c.retencion {
			t.Errorf("bruto %s: esperado retencion %s, obtenido %s", c.bruto, c.retencion, retenido.String())
		}
	}
}

func TestPremioSeSerializaComoNumeroJSON(t *testing.T) {
	premio := NewPremio(1, money.NewFromInt64(1000))

	data, err := json.Marshal(premio)
	if err != nil {
		t.Fatalf("Marshal error inesperado: %v", err)
	}

	var respuesta struct {
		IDPremio            int32       `json:"id_premio"`
		IDQuiniela          int32       `json:"id_quiniela"`
		MontoBruto          json.Number `json:"monto_bruto"`
		PorcentajeRetencion json.Number `json:"porcentaje_retencion"`
		MontoNeto           json.Number `json:"monto_neto"`
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&respuesta); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}

	if respuesta.MontoNeto.String() != "930.00" {
		t.Errorf("esperado monto_neto 930.00, obtenido %q", respuesta.MontoNeto.String())
	}
	if respuesta.MontoBruto.String() != "1000.00" {
		t.Errorf("esperado monto_bruto 1000.00, obtenido %q", respuesta.MontoBruto.String())
	}
	if respuesta.PorcentajeRetencion.String() != "7.0000" {
		t.Errorf("esperado porcentaje_retencion 7.0000, obtenido %q", respuesta.PorcentajeRetencion.String())
	}
}
