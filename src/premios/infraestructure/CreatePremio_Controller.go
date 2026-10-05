package infraestructure

import (
	"net/http"
	"quinielas/src/premios/application"
	"quinielas/src/shared/money"

	"github.com/gin-gonic/gin"
)

type CreatePremioController struct {
	useCase *application.CreatePremio
}

func NewCreatePremioController(useCase *application.CreatePremio) *CreatePremioController {
	return &CreatePremioController{useCase: useCase}
}

type CreatePremioRequestBody struct {
	IDQuiniela int32       `json:"id_quiniela" binding:"required"`
	MontoBruto money.Money `json:"monto_bruto"`
}

func (cp *CreatePremioController) Execute(c *gin.Context) {
	var body CreatePremioRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obligatorios faltantes o JSON inválido"})
		return
	}

	if body.MontoBruto.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monto_bruto es obligatorio y debe ser mayor a 0"})
		return
	}

	if body.MontoBruto.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monto_bruto no puede ser negativo"})
		return
	}

	premio, err := cp.useCase.Execute(body.IDQuiniela, body.MontoBruto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Premio registrado y calculado con 7% de retención correctamente",
		"premio":  premio,
	})
}