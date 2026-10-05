package domain

import (
	"errors"
	"quinielas/src/shared/money"
	"sort"
	"time"
)

var (
	ErrPremioNoRegistrado = errors.New("no hay un premio registrado para esta quiniela")
	ErrSinAportes         = errors.New("la quiniela no tiene aportes registrados")
	ErrAportesSinMonto    = errors.New("la quiniela tiene aportes pero el total recaudado es cero")
	ErrMetaInvalida       = errors.New("la quiniela no tiene una meta valida, no se puede repartir")
	ErrMetaNoAlcanzada    = errors.New("la meta de la quiniela no fue alcanzada, el premio no se reparte")
	ErrRepartoYaGenerado  = errors.New("esta quiniela ya tiene un reparto generado")
)

// precisionReparto es la escala con la que se calcula la porcion exacta de cada
// participante antes de truncar a centavos. Ocho decimales dejan margen de sobra
// para decidir a que participantes van los centavos sobrantes.
const precisionReparto int32 = 8

// AporteParticipante es el insumo minimo que el reparto necesita de cada cliente
// que participo en la quiniela.
type AporteParticipante struct {
	IDCliente     int32
	NombreCliente string
	Monto         money.Money
}

type ParticipanteReparto struct {
	IDCliente            int32         `json:"id_cliente"`
	NombreCliente        string        `json:"nombre_cliente"`
	MontoAporte          money.Money   `json:"monto_aporte"`
	PorcentajeAporte     money.Percent `json:"porcentaje_aporte"`     // sobre el total recaudado
	PorcentajeMeta       money.Percent `json:"porcentaje_sobre_meta"` // contra la meta de la quiniela
	PorcentajeRecompensa money.Percent `json:"porcentaje_recompensa"` // sobre el premio neto ya repartido
	MontoBrutoAsignado   money.Money   `json:"monto_bruto_asignado"`  // parte del premio antes de retencion
	MontoRetencion       money.Money   `json:"monto_retencion"`       // 7% aplicado a esa parte
	MontoNetoAsignado    money.Money   `json:"monto_neto_asignado"`   // parte del premio ya retenido
}

type RepartoPremio struct {
	IDPremio                int32                 `json:"id_premio"`
	IDQuiniela              int32                 `json:"id_quiniela"`
	MontoBruto              money.Money           `json:"monto_bruto"`
	PorcentajeRetencion     money.Percent         `json:"porcentaje_retencion"`
	MontoNeto               money.Money           `json:"monto_neto"`
	TotalRecaudado          money.Money           `json:"total_recaudado"`
	MontoMeta               money.Money           `json:"monto_meta"`
	PorcentajeMetaAlcanzado money.Percent         `json:"porcentaje_meta_alcanzado"`
	TotalParticipantes      int                   `json:"total_participantes"`
	Participantes           []ParticipanteReparto `json:"participantes"`
	FechaReparto            time.Time             `json:"fecha_reparto"`
}

// CalcularReparto reparte el premio de una quiniela entre quienes aportaron
// dinero, en proporcion a su aporte. El 7% ya viene descontado en el premio,
// asi que cada participante recibe su parte del monto neto.
//
// El reparto trunca cada parte a centavos y los centavos sobrantes se asignan a
// las fracciones mas altas (metodo del residuo mayor). Gracias a eso la suma de
// las partes es exactamente igual al monto a repartir: no se pierde ni se
// inventa ni un centavo.
func CalcularReparto(premio *Premio, aportes []AporteParticipante, montoMeta money.Money) (*RepartoPremio, error) {
	if premio == nil {
		return nil, ErrPremioNoRegistrado
	}
	if len(aportes) == 0 {
		return nil, ErrSinAportes
	}

	total := money.Zero()
	for _, a := range aportes {
		total = total.Add(a.Monto)
	}
	if !total.GreaterThan(money.Zero()) {
		return nil, ErrAportesSinMonto
	}
	if !montoMeta.GreaterThan(money.Zero()) {
		return nil, ErrMetaInvalida
	}
	if total.Cmp(montoMeta) < 0 {
		return nil, ErrMetaNoAlcanzada
	}

	brutoPorCliente := repartirExacto(premio.MontoBruto, aportes)
	netoPorCliente := repartirExacto(premio.MontoNeto, aportes)

	participantes := make([]ParticipanteReparto, 0, len(aportes))
	for i, a := range aportes {
		participantes = append(participantes, ParticipanteReparto{
			IDCliente:            a.IDCliente,
			NombreCliente:        a.NombreCliente,
			MontoAporte:          a.Monto,
			PorcentajeAporte:     money.PercentOf(a.Monto, total),
			PorcentajeMeta:       money.PercentOf(a.Monto, montoMeta),
			PorcentajeRecompensa: money.PercentOf(netoPorCliente[i], premio.MontoNeto),
			MontoBrutoAsignado:   brutoPorCliente[i],
			MontoRetencion:       brutoPorCliente[i].Sub(netoPorCliente[i]),
			MontoNetoAsignado:    netoPorCliente[i],
		})
	}

	// Mayor aporte primero: es el orden en que se quiere mostrar al usuario.
	sort.SliceStable(participantes, func(i, j int) bool {
		return participantes[i].MontoAporte.Cmp(participantes[j].MontoAporte) > 0
	})

	return &RepartoPremio{
		IDPremio:                premio.IDPremio,
		IDQuiniela:              premio.IDQuiniela,
		MontoBruto:              premio.MontoBruto,
		PorcentajeRetencion:     premio.PorcentajeRetencion,
		MontoNeto:               premio.MontoNeto,
		TotalRecaudado:          total,
		MontoMeta:               montoMeta,
		PorcentajeMetaAlcanzado: money.PercentOf(total, montoMeta),
		TotalParticipantes:      len(participantes),
		Participantes:           participantes,
		FechaReparto:            time.Now(),
	}, nil
}

// repartirExacto reparte un monto entre los aportes de forma que la suma de las
// partes sea identica al monto original.
func repartirExacto(monto money.Money, aportes []AporteParticipante) []money.Money {
	total := money.Zero()
	for _, a := range aportes {
		total = total.Add(a.Monto)
	}

	partes := make([]money.Money, len(aportes))
	fracciones := make([]money.Money, len(aportes))
	asignado := money.Zero()

	for i, a := range aportes {
		exacta := monto.MulDiv(a.Monto, total, precisionReparto)
		truncada := exacta.RoundDown(money.MontoScale)

		partes[i] = truncada
		fracciones[i] = exacta.Sub(truncada)
		asignado = asignado.Add(truncada)
	}

	// Centavos que se perdieron al truncar: se reparten de a uno.
	sobrantes := monto.Sub(asignado)
	centimos := int(sobrantes.Div(money.Centimos()).RoundDown(0).Decimal.IntPart())
	if centimos <= 0 {
		return partes
	}

	orden := make([]int, len(aportes))
	for i := range orden {
		orden[i] = i
	}
	sort.SliceStable(orden, func(a, b int) bool {
		return fracciones[orden[a]].Cmp(fracciones[orden[b]]) > 0
	})

	for k := 0; k < centimos; k++ {
		idx := orden[k%len(orden)]
		partes[idx] = partes[idx].Add(money.Centimos())
	}

	return partes
}
