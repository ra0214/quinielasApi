package apierror

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// ForeignKeyViolation indica si el error proviene de una restricción de clave
// foránea (MySQL 1451/1452), p. ej. al borrar un registro referenciado.
func ForeignKeyViolation(err error) bool {
	var myErr *mysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number == 1451 || myErr.Number == 1452
	}
	return false
}

// Eliminar responde 409 (Conflicto) si el registro está referenciado por otra
// tabla, o 500 con un mensaje genérico en cualquier otro caso. El detalle
// técnico del error solo se registra en el log del servidor, nunca se expone
// al usuario.
func Eliminar(c *gin.Context, err error, entidad string) {
	log.Printf("[API ERROR] detalle al eliminar %s: %v", entidad, err)

	if ForeignKeyViolation(err) {
		c.JSON(http.StatusConflict, gin.H{
			"error": "No se puede eliminar " + entidad + ": tiene registros asociados.",
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar " + entidad})
}