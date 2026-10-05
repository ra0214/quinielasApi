package domain

import (
	"quinielas/src/shared/money"
	"testing"
)

func aportes(pairs ...string) []AporteParticipante {
	var lista []AporteParticipante
	for i := 0; i < len(pairs); i += 2 {
		m, err := money.Parse(pairs[i+1])
		if err != nil {
			panic(err)
		}
		lista = append(lista, AporteParticipante{
			IDCliente:     int32(i/2 + 1),
			NombreCliente: pairs[i],
			Monto:         m,
		})
	}
	return lista
}

func sumaNeto(r *RepartoPremio) money.Money {
	t := money.Zero()
	for _, p := range r.Participantes {
		t = t.Add(p.MontoNetoAsignado)
	}
	return t
}

func sumaBruto(r *RepartoPremio) money.Money {
	t := money.Zero()
	for _, p := range r.Participantes {
		t = t.Add(p.MontoBrutoAsignado)
	}
	return t
}

func sumaRetencion(r *RepartoPremio) money.Money {
	t := money.Zero()
	for _, p := range r.Participantes {
		t = t.Add(p.MontoRetencion)
	}
	return t
}

func meta(v string) money.Money {
	m, err := money.Parse(v)
	if err != nil {
		panic(err)
	}
	return m
}

func TestRepartoSumaExactamenteElPremioNeto(t *testing.T) {
	casos := []struct {
		premio  string
		aportes []string
	}{
		{"1000.00", []string{"Ana", "1000.00", "Luis", "2000.00", "Sofia", "3000.00"}},
		{"1000.00", []string{"Ana", "0.01", "Luis", "0.01", "Sofia", "0.02"}},
		{"1.00", []string{"Ana", "0.01", "Luis", "0.01", "Sofia", "0.01"}},
		{"7.77", []string{"Ana", "3.33", "Luis", "3.33", "Sofia", "3.34"}},
		{"100.00", []string{"Ana", "33.33", "Luis", "33.33", "Sofia", "33.34"}},
		{"999999.99", []string{"Ana", "12345.67", "Luis", "23456.78", "Sofia", "34567.89", "Pedro", "1.00"}},
		{"50.00", []string{"Ana", "10000.00", "Luis", "10000.00", "Sofia", "10000.00"}},
		{"123.45", []string{"Ana", "500.00", "Luis", "1500.00"}},
	}

	for _, c := range casos {
		premio := NewPremio(1, money.MustParse(c.premio))
		total := money.Zero()
		for _, a := range aportes(c.aportes...) {
			total = total.Add(a.Monto)
		}

		reparto, err := CalcularReparto(premio, aportes(c.aportes...), total)
		if err != nil {
			t.Fatalf("premio %s: error inesperado: %v", c.premio, err)
		}

		if got := sumaNeto(reparto); got.Cmp(reparto.MontoNeto) != 0 {
			t.Errorf("premio %s: la suma de las partes es %s pero el neto es %s",
				c.premio, got.String(), reparto.MontoNeto.String())
		}
		if got := sumaBruto(reparto); got.Cmp(reparto.MontoBruto) != 0 {
			t.Errorf("premio %s: la suma de brutos es %s pero el bruto es %s",
				c.premio, got.String(), reparto.MontoBruto.String())
		}

		retencionEsperada := reparto.MontoBruto.Sub(reparto.MontoNeto)
		if got := sumaRetencion(reparto); got.Cmp(retencionEsperada) != 0 {
			t.Errorf("premio %s: la suma de retenciones es %s pero deberia ser %s",
				c.premio, got.String(), retencionEsperada.String())
		}
	}
}

func TestRepartoUnSoloParticipanteRecibeTodoElNeto(t *testing.T) {
	premio := NewPremio(1, money.NewFromInt64(1000))

	reparto, err := CalcularReparto(premio, aportes("Ana", "5000.00"), meta("5000.00"))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if len(reparto.Participantes) != 1 {
		t.Fatalf("esperaba 1 participante, obtenido %d", len(reparto.Participantes))
	}
	if got := reparto.Participantes[0].MontoNetoAsignado.String(); got != "930.00" {
		t.Errorf("esperaba 930.00 para el unico participante, obtenido %s", got)
	}
	if got := reparto.Participantes[0].MontoBrutoAsignado.String(); got != "1000.00" {
		t.Errorf("esperaba 1000.00 de bruto, obtenido %s", got)
	}
	if got := reparto.Participantes[0].MontoRetencion.String(); got != "70.00" {
		t.Errorf("esperaba 70.00 de retencion, obtenido %s", got)
	}
}

func TestRepartoProporcionalDeImagenConocida(t *testing.T) {
	premio := NewPremio(1, money.MustParse("1000.00"))

	// 50% / 30% / 20% de 1000 neto -> 465.00 / 279.00 / 186.00
	reparto, err := CalcularReparto(premio, aportes(
		"Ana", "5000.00",
		"Luis", "3000.00",
		"Sofia", "2000.00",
	), meta("10000.00"))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	esperado := map[string]struct {
		bruto string
		neto  string
	}{
		"Ana":   {"500.00", "465.00"},
		"Luis":  {"300.00", "279.00"},
		"Sofia": {"200.00", "186.00"},
	}

	for _, p := range reparto.Participantes {
		e, ok := esperado[p.NombreCliente]
		if !ok {
			t.Fatalf("participante inesperado: %s", p.NombreCliente)
		}
		if p.MontoBrutoAsignado.String() != e.bruto {
			t.Errorf("%s: bruto esperado %s, obtenido %s", p.NombreCliente, e.bruto, p.MontoBrutoAsignado.String())
		}
		if p.MontoNetoAsignado.String() != e.neto {
			t.Errorf("%s: neto esperado %s, obtenido %s", p.NombreCliente, e.neto, p.MontoNetoAsignado.String())
		}
	}
}

func TestRepartoOrdenaPorAporteDescendente(t *testing.T) {
	premio := NewPremio(1, money.NewFromInt64(1000))

	reparto, err := CalcularReparto(premio, aportes(
		"Chico", "100.00",
		"Grande", "900.00",
		"Mediano", "500.00",
	), meta("1500.00"))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	orden := []string{"Grande", "Mediano", "Chico"}
	for i, nombre := range orden {
		if reparto.Participantes[i].NombreCliente != nombre {
			t.Errorf("posicion %d: esperaba %s, obtenido %s", i, nombre, reparto.Participantes[i].NombreCliente)
		}
	}
}

func TestRepartoRechazaMetaNoAlcanzada(t *testing.T) {
	premio := NewPremio(1, money.NewFromInt64(1000))

	// Se recaudaron 3000 de una meta de 5000.
	_, err := CalcularReparto(premio, aportes(
		"Ana", "1000.00",
		"Luis", "2000.00",
	), meta("5000.00"))

	if err != ErrMetaNoAlcanzada {
		t.Errorf("esperaba ErrMetaNoAlcanzada, obtenido %v", err)
	}
}

func TestRepartoAceptaMetaExactamenteAlcanzada(t *testing.T) {
	premio := NewPremio(1, money.NewFromInt64(1000))

	reparto, err := CalcularReparto(premio, aportes(
		"Ana", "1000.00",
		"Luis", "4000.00",
	), meta("5000.00"))
	if err != nil {
		t.Fatalf("la meta alcanzada exactamente debe repartirse, obtenu error: %v", err)
	}
	if len(reparto.Participantes) != 2 {
		t.Errorf("esperaba 2 participantes, obtenido %d", len(reparto.Participantes))
	}
}

func TestRepartoRechazaSinAportes(t *testing.T) {
	premio := NewPremio(1, money.NewFromInt64(1000))

	if _, err := CalcularReparto(premio, nil, meta("5000.00")); err != ErrSinAportes {
		t.Errorf("esperaba ErrSinAportes, obtenido %v", err)
	}
}

func TestRepartoRechazaTotalEnCero(t *testing.T) {
	premio := NewPremio(1, money.NewFromInt64(1000))

	if _, err := CalcularReparto(premio, aportes("Ana", "0.00"), meta("5000.00")); err != ErrAportesSinMonto {
		t.Errorf("esperaba ErrAportesSinMonto, obtenido %v", err)
	}
}

func TestRepartoRechazaMetaCero(t *testing.T) {
	premio := NewPremio(1, money.NewFromInt64(1000))

	if _, err := CalcularReparto(premio, aportes("Ana", "100.00"), money.Zero()); err != ErrMetaInvalida {
		t.Errorf("esperaba ErrMetaInvalida, obtenido %v", err)
	}
}

func TestRepartoReportaPorcentajes(t *testing.T) {
	premio := NewPremio(1, money.NewFromInt64(1000))

	reparto, err := CalcularReparto(premio, aportes(
		"Ana", "5000.00",
		"Luis", "5000.00",
	), meta("10000.00"))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if got := reparto.TotalRecaudado.String(); got != "10000.00" {
		t.Errorf("total recaudado esperado 10000.00, obtenido %s", got)
	}
	if got := reparto.PorcentajeMetaAlcanzado.String(); got != "100.0000" {
		t.Errorf("porcentaje de meta esperado 100.0000, obtenido %s", got)
	}
	for _, p := range reparto.Participantes {
		if got := p.PorcentajeAporte.String(); got != "50.0000" {
			t.Errorf("%s: porcentaje de aporte esperado 50.0000, obtenido %s", p.NombreCliente, got)
		}
		if got := p.PorcentajeMeta.String(); got != "50.0000" {
			t.Errorf("%s: porcentaje sobre meta esperado 50.0000, obtenido %s", p.NombreCliente, got)
		}
		if got := p.PorcentajeRecompensa.String(); got != "50.0000" {
			t.Errorf("%s: porcentaje de recompensa esperado 50.0000, obtenido %s", p.NombreCliente, got)
		}
	}
}

func TestRepartoConMuchosParticipantesNoPierdeCentavos(t *testing.T) {
	// 7 participantes con aportes que no dividen exacto entre el premio:
	// es el caso que hace fallar un redondeo ingenuo.
	lista := []string{
		"A", "100.00",
		"B", "100.01",
		"C", "100.02",
		"D", "100.03",
		"E", "100.04",
		"F", "100.05",
		"G", "100.06",
	}

	premio := NewPremio(1, money.MustParse("1234.57"))
	reparto, err := CalcularReparto(premio, aportes(lista...), meta("700.21"))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if got := sumaNeto(reparto); got.Cmp(reparto.MontoNeto) != 0 {
		t.Errorf("suma %s, neto %s", got.String(), reparto.MontoNeto.String())
	}
	if got := sumaBruto(reparto); got.Cmp(reparto.MontoBruto) != 0 {
		t.Errorf("suma bruta %s, bruto %s", got.String(), reparto.MontoBruto.String())
	}
	if len(reparto.Participantes) != 7 {
		t.Errorf("esperaba 7 participantes, obtenido %d", len(reparto.Participantes))
	}
}
