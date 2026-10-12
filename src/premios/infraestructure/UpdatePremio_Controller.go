package infraestructure

import (
	"errors"
	"net/http"
	"quinielas/src/premios/application"
	"quinielas/src/shared/money"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdatePremioController struct {
	useCase *application.UpdatePremio
}

func NewUpdatePremioController(useCase *application.UpdatePremio) *UpdatePremioController {
	return &UpdatePremioController{useCase: useCase}
}

type UpdatePremioRequestBody struct {
	MontoBruto      money.Money `json:"monto_bruto"`
	PremioPorBoleto money.Money `json:"premio_por_boleto"`
}

func (up *UpdatePremioController) Execute(c *gin.Context) {
	idPremio, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de premio inválido"})
		return
	}

	var body UpdatePremioRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
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

	premio, err := up.useCase.Execute(int32(idPremio), body.MontoBruto, body.PremioPorBoleto)
	if err != nil {
		if errors.Is(err, application.ErrPremioNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Premio actualizado y recalculado correctamente",
		"premio":  premio,
	})
}