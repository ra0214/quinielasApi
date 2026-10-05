package money

import (
	"fmt"

	"github.com/shopspring/decimal"
)

const (
	MontoScale   int32 = 2
	PercentScale int32 = 4
)

var (
	OneHundred = decimal.New(100, 0)
)

type Money struct {
	decimal.Decimal
}

type Percent struct {
	decimal.Decimal
}

func Zero() Money {
	return Money{decimal.Zero}
}

func NewFromInt64(value int64) Money {
	return Money{decimal.NewFromInt(value)}
}

func Parse(value string) (Money, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Zero(), fmt.Errorf("monto inválido %q: %w", value, err)
	}
	return Money{d.Round(MontoScale)}, nil
}

func FromFloat64(value float64) Money {
	return Money{decimal.NewFromFloat(value).Round(MontoScale)}
}

func (m Money) Add(other Money) Money {
	return Money{m.Decimal.Add(other.Decimal).Round(MontoScale)}
}

func (m Money) Sub(other Money) Money {
	return Money{m.Decimal.Sub(other.Decimal).Round(MontoScale)}
}

func (m Money) MulPercent(p Percent) Money {
	return Money{m.Decimal.Mul(p.Decimal).Div(OneHundred).Round(MontoScale)}
}

func (m Money) AfterRetention(p Percent) Money {
	return m.Sub(m.MulPercent(p))
}

func (m Money) Cmp(other Money) int {
	return m.Decimal.Cmp(other.Decimal)
}

func (m Money) IsZero() bool {
	return m.Decimal.IsZero()
}

func (m Money) IsNegative() bool {
	return m.Decimal.IsNegative()
}

func (m Money) GreaterThan(other Money) bool {
	return m.Decimal.GreaterThan(other.Decimal)
}

func (m Money) String() string {
	return m.Decimal.StringFixed(MontoScale)
}

func NewPercentFromInt64(value int64) Percent {
	return Percent{decimal.NewFromInt(value)}
}

func PercentOf(parte, total Money) Percent {
	if !total.GreaterThan(Zero()) {
		return Percent{decimal.Zero}
	}
	return Percent{parte.Decimal.Div(total.Decimal).Mul(OneHundred).Round(PercentScale)}
}

func (p Percent) Cmp(other Percent) int {
	return p.Decimal.Cmp(other.Decimal)
}

func (p Percent) IsZero() bool {
	return p.Decimal.IsZero()
}

func (p Percent) String() string {
	return p.Decimal.StringFixed(PercentScale)
}
