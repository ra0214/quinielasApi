package infraestructure

import (
	"errors"
	"net/http"
	"quinielas/src/premios/application"
	"quinielas/src/premios/domain"
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
	IDQuiniela      int32       `json:"id_quiniela" binding:"required"`
	MontoBruto      money.Money `json:"monto_bruto"`
	PremioPorBoleto money.Money `json:"premio_por_boleto"`
}

func (cp *CreatePremioController) Execute(c *gin.Context) {
	var body CreatePremioRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obligatorios faltantes o JSON inválido"})
		return
	}

	if body.MontoBruto.IsZero() && body.PremioPorBoleto.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monto_bruto (o premio_por_boleto en coperacha) es obligatorio y debe ser mayor a 0"})
		return
	}

	if body.MontoBruto.IsNegative() || body.PremioPorBoleto.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "los montos del premio no pueden ser negativos"})
		return
	}

	premio, err := cp.useCase.Execute(body.IDQuiniela, body.MontoBruto, body.PremioPorBoleto)
	if err != nil {
		if errors.Is(err, domain.ErrPremioSinMontoPorBoleto) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Premio registrado y calculado correctamente",
		"premio":  premio,
	})
}