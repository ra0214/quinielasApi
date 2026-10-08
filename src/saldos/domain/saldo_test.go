package domain

import (
	"errors"
	"quinielas/src/shared/money"
	"testing"
)

func m(s string) money.Money {
	v, err := money.Parse(s)
	if err != nil {
		panic(err)
	}
	return v
}

func TestCompensar(t *testing.T) {
	casos := []struct {
		nombre        string
		favor, deuda  string
		espFavor      string
		espDeuda      string
	}{
		{"solo favor", "100.00", "0.00", "100.00", "0.00"},
		{"solo deuda", "0.00", "50.00", "0.00", "50.00"},
		{"netean a favor", "100.00", "20.00", "80.00", "0.00"},
		{"netean a deuda", "20.00", "100.00", "0.00", "80.00"},
		{"cero", "0.00", "0.00", "0.00", "0.00"},
		{"iguales quedan en cero", "50.00", "50.00", "0.00", "0.00"},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			fav, deu := Compensar(m(c.favor), m(c.deuda))
			if fav.Cmp(m(c.espFavor)) != 0 || deu.Cmp(m(c.espDeuda)) != 0 {
				t.Errorf("Compensar(%s, %s) = (%s, %s), se esperaba (%s, %s)",
					c.favor, c.deuda, fav, deu, c.espFavor, c.espDeuda)
			}
		})
	}
}

func TestAplicarMovimiento(t *testing.T) {
	casos := []struct {
		nombre   string
		favor    string
		deuda    string
		tipo     string
		monto    string
		espFavor string
		espDeuda string
		err      bool
	}{
		{"FIADO suma deuda", "0.00", "0.00", "FIADO", "20.00", "0.00", "20.00", false},
		{"FIADO compensa contra favor", "100.00", "0.00", "FIADO", "20.00", "80.00", "0.00", false},
		{"PAGO_EFECTIVO paga deuda", "0.00", "50.00", "PAGO_EFECTIVO", "80.00", "30.00", "0.00", false},
		{"PAGO_EFECTIVO sobra a favor", "0.00", "50.00", "PAGO_EFECTIVO", "30.00", "0.00", "20.00", false},
		{"PAGO_SALDO consume favor", "80.00", "0.00", "PAGO_SALDO", "50.00", "30.00", "0.00", false},
		{"PAGO_SALDO sin favor da error", "0.00", "0.00", "PAGO_SALDO", "10.00", "", "", true},
		{"RETIRO consume favor", "80.00", "0.00", "RETIRO", "30.00", "50.00", "0.00", false},
		{"RETIRO sin favor da error", "10.00", "0.00", "RETIRO", "50.00", "", "", true},
		{"PREMIO_ABONO suma favor", "10.00", "0.00", "PREMIO_ABONO", "90.00", "100.00", "0.00", false},
		{"DEVOLUCION suma favor", "10.00", "20.00", "DEVOLUCION", "50.00", "40.00", "0.00", false},
		{"tipo desconocido da error", "0.00", "0.00", "OTRO", "10.00", "", "", true},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			fav, deu, err := AplicarMovimiento(m(c.favor), m(c.deuda), c.tipo, m(c.monto))
			if c.err {
				if err == nil {
					t.Fatalf("se esperaba error, no hubo")
				}
				if errors.Is(err, ErrSaldoInsuficiente) && (c.tipo != "PAGO_SALDO" && c.tipo != "RETIRO") {
					t.Errorf("ErrSaldoInsuficiente no aplica al tipo %s", c.tipo)
				}
				return
			}
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if fav.Cmp(m(c.espFavor)) != 0 || deu.Cmp(m(c.espDeuda)) != 0 {
				t.Errorf("AplicarMovimiento(%s, %s, %s, %s) = (%s, %s), se esperaba (%s, %s)",
					c.favor, c.deuda, c.tipo, c.monto, fav, deu, c.espFavor, c.espDeuda)
			}
		})
	}
}

func TestRevertirMovimientoEsInverso(t *testing.T) {
	casos := []struct {
		favor, deuda string
		tipo         string
		monto        string
	}{
		{"100.00", "0.00", "FIADO", "40.00"},
		{"0.00", "50.00", "PAGO_EFECTIVO", "80.00"},
		{"80.00", "0.00", "PAGO_SALDO", "50.00"},
		{"80.00", "0.00", "RETIRO", "30.00"},
		{"100.00", "0.00", "PREMIO_ABONO", "25.00"},
		{"40.00", "0.00", "DEVOLUCION", "15.00"},
	}

	for _, c := range casos {
		t.Run(c.tipo+"_"+c.monto, func(t *testing.T) {
			fav, deu, nilErr := AplicarMovimiento(m(c.favor), m(c.deuda), c.tipo, m(c.monto))
			if nilErr != nil {
				t.Fatalf("error al aplicar: %v", nilErr)
			}
			rfav, rdeu, revErr := RevertirMovimiento(fav, deu, c.tipo, m(c.monto))
			if revErr != nil {
				t.Fatalf("error al revertir: %v", revErr)
			}
			if rfav.Cmp(m(c.favor)) != 0 || rdeu.Cmp(m(c.deuda)) != 0 {
				t.Errorf("aplicar y revertir (%s) no regresa al (%s, %s): quedó (%s, %s)",
					c.tipo, c.favor, c.deuda, rfav, rdeu)
			}
		})
	}
}