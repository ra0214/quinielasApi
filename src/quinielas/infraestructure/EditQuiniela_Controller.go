package infraestructure

import (
	"net/http"
	"quinielas/src/quinielas/application"
	"quinielas/src/shared/money"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EditQuinielaController struct {
	useCase *application.UpdateQuiniela
}

func NewEditQuinielaController(useCase *application.UpdateQuiniela) *EditQuinielaController {
	return &EditQuinielaController{useCase: useCase}
}

type EditQuinielaRequestBody struct {
	NombreVariante string      `json:"nombre_variante" binding:"required"`
	Precio         money.Money `json:"precio"`
	MontoMeta      money.Money `json:"monto_meta"`
	FechaLimite    time.Time   `json:"fecha_limite" binding:"required"`
	Estado         string      `json:"estado" binding:"required"`
}

func (eq *EditQuinielaController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de quiniela inválido"})
		return
	}

	var body EditQuinielaRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if body.Precio.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "precio es obligatorio y debe ser mayor a 0"})
		return
	}

	if body.Precio.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "precio no puede ser negativo"})
		return
	}

	if body.MontoMeta.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monto_meta no puede ser negativo"})
		return
	}

	err = eq.useCase.Execute(int32(id), body.NombreVariante, body.Precio, body.MontoMeta, body.FechaLimite, body.Estado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la quiniela", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quiniela actualizada correctamente"})
}