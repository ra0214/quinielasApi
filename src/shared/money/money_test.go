package money

import (
	"encoding/json"
	"testing"
)

func TestAfterRetentionSinErrorDePuntoFlotante(t *testing.T) {
	casos := []struct {
		bruto string
		neto  string
	}{
		{"1000.00", "930.00"},
		{"0.01", "0.01"},
		{"10.00", "9.30"},
		{"123456.78", "114814.81"},
		{"1.00", "0.93"},
		{"333.33", "310.00"},
	}

	for _, c := range casos {
		bruto, err := Parse(c.bruto)
		if err != nil {
			t.Fatalf("Parse(%q) error inesperado: %v", c.bruto, err)
		}

		neto := bruto.AfterRetention(NewPercentFromInt64(7))
		if neto.String() != c.neto {
			t.Errorf("bruto %s: esperado neto %s, obtenido %s", c.bruto, c.neto, neto.String())
		}
	}
}

func TestSumaAcumuladaNoAcumulaError(t *testing.T) {
	acumulado := Zero()
	for i := 0; i < 1000; i++ {
		acumulado = acumulado.Add(NewFromInt64(1))
	}

	if acumulado.String() != "1000.00" {
		t.Errorf("esperado 1000.00, obtenido %s", acumulado.String())
	}
}

func TestSumarCentavosExactos(t *testing.T) {
	acumulado, _ := Parse("0.10")
	for i := 0; i < 3; i++ {
		acumulado = acumulado.Add(NewFromInt64(1))
	}

	if acumulado.String() != "3.10" {
		t.Errorf("esperado 3.10, obtenido %s", acumulado.String())
	}
}

func TestRestarEvitaPerdidaDePrecision(t *testing.T) {
	a, _ := Parse("1000.00")
	b, _ := Parse("0.10")

	if got := a.Sub(b).String(); got != "999.90" {
		t.Errorf("esperado 999.90, obtenido %s", got)
	}
}

func TestMarshalJSONEmiteNumeroNoString(t *testing.T) {
	monto, _ := Parse("930.00")

	data, err := json.Marshal(monto)
	if err != nil {
		t.Fatalf("Marshal error inesperado: %v", err)
	}

	if string(data) != "930.00" {
		t.Errorf("esperado 930.00 como numero JSON, obtenido %s", string(data))
	}
}

func TestMarshalJSONDePremioCompleto(t *testing.T) {
	type respuesta struct {
		MontoBruto Money   `json:"monto_bruto"`
		Retencion  Percent `json:"porcentaje_retencion"`
		MontoNeto  Money   `json:"monto_neto"`
	}

	bruto, _ := Parse("1000.00")
	neto := bruto.AfterRetention(NewPercentFromInt64(7))

	data, err := json.Marshal(respuesta{bruto, NewPercentFromInt64(7), neto})
	if err != nil {
		t.Fatalf("Marshal error inesperado: %v", err)
	}

	esperado := `{"monto_bruto":1000.00,"porcentaje_retencion":7.0000,"monto_neto":930.00}`
	if string(data) != esperado {
		t.Errorf("esperado %s, obtenido %s", esperado, string(data))
	}
}

func TestUnmarshalJSONAceptaNumeroYString(t *testing.T) {
	casos := []struct {
		entrada  string
		esperado string
	}{
		{`{"monto_bruto": 1000.00}`, "1000.00"},
		{`{"monto_bruto": "1000.00"}`, "1000.00"},
		{`{"monto_bruto": 1000}`, "1000.00"},
		{`{"monto_bruto": null}`, "0.00"},
	}

	for _, c := range casos {
		var wrapper struct {
			MontoBruto Money `json:"monto_bruto"`
		}
		if err := json.Unmarshal([]byte(c.entrada), &wrapper); err != nil {
			t.Fatalf("Unmarshal(%s) error inesperado: %v", c.entrada, err)
		}
		if wrapper.MontoBruto.String() != c.esperado {
			t.Errorf("entrada %s: esperado %s, obtenido %s", c.entrada, c.esperado, wrapper.MontoBruto.String())
		}
	}
}

func TestUnmarshalJSONRechazaNoNumero(t *testing.T) {
	var m Money
	if err := json.Unmarshal([]byte(`"abc"`), &m); err == nil {
		t.Error("se esperaba error al deserializar un valor no numérico")
	}
}

func TestRedondeoAEscenaDeDos(t *testing.T) {
	casos := []struct {
		entrada  string
		esperado string
	}{
		{"10.005", "10.01"},
		{"10.004", "10.00"},
		{"2.675", "2.68"},
		{"0.999", "1.00"},
	}

	for _, c := range casos {
		m, err := Parse(c.entrada)
		if err != nil {
			t.Fatalf("Parse(%q) error inesperado: %v", c.entrada, err)
		}
		if m.String() != c.esperado {
			t.Errorf("entrada %s: esperado %s, obtenido %s", c.entrada, c.esperado, m.String())
		}
	}
}

func TestPercentOf(t *testing.T) {
	casos := []struct {
		parte    string
		total    string
		esperado string
	}{
		{"1000.00", "5000.00", "20.0000"},
		{"333.33", "1000.00", "33.3330"},
		{"1.00", "3.00", "33.3333"},
		{"500.00", "0.00", "0.0000"},
		{"0.00", "0.00", "0.0000"},
	}

	for _, c := range casos {
		parte, _ := Parse(c.parte)
		total, _ := Parse(c.total)
		got := PercentOf(parte, total)
		if got.String() != c.esperado {
			t.Errorf("parte %s de %s: esperado %s, obtenido %s", c.parte, c.total, c.esperado, got.String())
		}
	}
}

func TestMoneyCeroNoPanica(t *testing.T) {
	var m Money

	if !m.IsZero() {
		t.Error("un Money sin inicializar debe ser cero")
	}
	if m.IsNegative() {
		t.Error("un Money sin inicializar no debe ser negativo")
	}
	if m.String() != "0.00" {
		t.Errorf("esperado 0.00, obtenido %s", m.String())
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Marshal error inesperado: %v", err)
	}
	if string(data) != "0.00" {
		t.Errorf("esperado 0.00, obtenido %s", string(data))
	}
}

func TestScanDesdeBytesDeMySQL(t *testing.T) {
	var m Money
	if err := m.Scan([]byte("929.9999999999999")); err != nil {
		t.Fatalf("Scan error inesperado: %v", err)
	}

	if m.String() != "930.00" {
		t.Errorf("esperado 930.00 tras escanear valor almacenado, obtenido %s", m.String())
	}
}

func TestScanDesdeVariosTipos(t *testing.T) {
	casos := []struct {
		entrada  any
		esperado string
	}{
		{nil, "0.00"},
		{[]byte("1234.56"), "1234.56"},
		{"1234.56", "1234.56"},
		{int64(1234), "1234.00"},
		{float64(1234.56), "1234.56"},
	}

	for _, c := range casos {
		var m Money
		if err := m.Scan(c.entrada); err != nil {
			t.Fatalf("Scan(%v) error inesperado: %v", c.entrada, err)
		}
		if m.String() != c.esperado {
			t.Errorf("entrada %v: esperado %s, obtenido %s", c.entrada, c.esperado, m.String())
		}
	}
}

func TestValueParaMySQL(t *testing.T) {
	m, _ := Parse("930.00")

	v, err := m.Value()
	if err != nil {
		t.Fatalf("Value error inesperado: %v", err)
	}
	if v != "930" {
		t.Errorf("esperado \"930\", obtenido %v", v)
	}
}

func TestMulPercent(t *testing.T) {
	bruto, _ := Parse("1000.00")
	if got := bruto.MulPercent(NewPercentFromInt64(7)).String(); got != "70.00" {
		t.Errorf("esperado 70.00, obtenido %s", got)
	}

	if got := bruto.MulPercent(NewPercentFromInt64(0)).String(); got != "0.00" {
		t.Errorf("esperado 0.00, obtenido %s", got)
	}
}

func TestComparaciones(t *testing.T) {
	a, _ := Parse("100.00")
	b, _ := Parse("50.00")

	if !a.GreaterThan(b) {
		t.Error("100.00 debe ser mayor que 50.00")
	}
	if b.Cmp(a) != -1 {
		t.Error("50.00 debe ser menor que 100.00")
	}
	if !moneyZeroIsZero() {
		t.Error("Zero() debe reportarse como cero")
	}
}

func moneyZeroIsZero() bool {
	return Zero().IsZero()
}
