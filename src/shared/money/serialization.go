package money

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

func normalizeJSONNumber(raw []byte, scale int32, target *decimal.Decimal) error {
	text := strings.TrimSpace(string(raw))
	if text == "null" {
		*target = decimal.Zero
		return nil
	}

	if len(text) >= 2 && text[0] == '"' && text[len(text)-1] == '"' {
		text = text[1 : len(text)-1]
	}

	parsed, err := decimal.NewFromString(text)
	if err != nil {
		return fmt.Errorf("valor monetario inválido %q", string(raw))
	}

	*target = parsed.Round(scale)
	return nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(m.Decimal.StringFixed(MontoScale)), nil
}

func (m *Money) UnmarshalJSON(data []byte) error {
	return normalizeJSONNumber(data, MontoScale, &m.Decimal)
}

func (p Percent) MarshalJSON() ([]byte, error) {
	return []byte(p.Decimal.StringFixed(PercentScale)), nil
}

func (p *Percent) UnmarshalJSON(data []byte) error {
	return normalizeJSONNumber(data, PercentScale, &p.Decimal)
}

func (m *Money) Scan(src any) error {
	return scanIntoDecimal(src, MontoScale, &m.Decimal)
}

func (p *Percent) Scan(src any) error {
	return scanIntoDecimal(src, PercentScale, &p.Decimal)
}

func scanIntoDecimal(src any, scale int32, target *decimal.Decimal) error {
	switch v := src.(type) {
	case nil:
		*target = decimal.Zero
		return nil
	case []byte:
		return normalizeJSONNumber(v, scale, target)
	case string:
		return normalizeJSONNumber([]byte(v), scale, target)
	case float64:
		*target = decimal.NewFromFloat(v).Round(scale)
		return nil
	case int64:
		*target = decimal.New(v, 0)
		return nil
	default:
		return fmt.Errorf("no se puede convertir %T a decimal", src)
	}
}

func (m Money) Value() (driver.Value, error) {
	return m.Decimal.String(), nil
}

func (p Percent) Value() (driver.Value, error) {
	return p.Decimal.String(), nil
}

var (
	_ json.Marshaler   = Money{}
	_ json.Unmarshaler = (*Money)(nil)
	_ driver.Valuer    = Money{}

	_ json.Marshaler   = Percent{}
	_ json.Unmarshaler = (*Percent)(nil)
	_ driver.Valuer    = Percent{}
)
