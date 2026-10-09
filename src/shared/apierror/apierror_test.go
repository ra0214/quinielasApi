package apierror

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func TestForeignKeyViolation(t *testing.T) {
	casos := []struct {
		nombre string
		err    error
		espera bool
	}{
		{"error 1451 detectado", &mysql.MySQLError{Number: 1451, Message: "Cannot delete or update a parent row"}, true},
		{"error 1452 detectado", &mysql.MySQLError{Number: 1452, Message: "Cannot add or update a child row"}, true},
		{"otro error mysql", &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}, false},
		{"error envuelto", errors.Join(&mysql.MySQLError{Number: 1451}, errors.New("contexto")), true},
		{"error generico", errors.New("algo salió mal"), false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := ForeignKeyViolation(c.err); got != c.espera {
				t.Errorf("ForeignKeyViolation(%v) = %v, se esperaba %v", c.err, got, c.espera)
			}
		})
	}
}

func TestEliminarRespondeMensajeLimpio(t *testing.T) {
	gin.SetMode(gin.TestMode)
	detalle := errors.Join(&mysql.MySQLError{Number: 1451, Message: "FK fails"}, errors.New("cadena técnica"))

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)

	Eliminar(c, detalle, "(cliente demo)")

	if rec.Code != http.StatusConflict {
		t.Errorf("se esperaba 409, se obtuvo %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("respuesta no es JSON válido: %v", err)
	}

	if strings.Contains(body["error"], "FK fails") || strings.Contains(body["error"], "cadena técnica") {
		t.Errorf("el mensaje expone el detalle técnico: %q", body["error"])
	}
	if !strings.Contains(body["error"], "No se puede eliminar") {
		t.Errorf("mensaje no es amigable: %q", body["error"])
	}
}

func TestEliminarResponde500SinDetalle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	detalle := errors.New("error interno de la base de datos")

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)

	Eliminar(c, detalle, "(cliente demo)")

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("se esperaba 500, se obtuvo %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("respuesta no es JSON válido: %v", err)
	}

	if strings.Contains(body["error"], "error interno") {
		t.Errorf("el mensaje expone el detalle técnico: %q", body["error"])
	}
	if body["error"] != "Error al eliminar (cliente demo)" {
		t.Errorf("mensaje inesperado: %q", body["error"])
	}
}