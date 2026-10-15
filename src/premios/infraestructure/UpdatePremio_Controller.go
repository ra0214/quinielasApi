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
	MontoBruto money.Money `json:"monto_bruto"`
	MontoNeto  money.Money `json:"monto_neto"`
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

	if body.MontoBruto.IsNegative() || body.MontoNeto.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "los montos no pueden ser negativos"})
		return
	}

	if !body.MontoBruto.GreaterThan(money.Zero()) && !body.MontoNeto.GreaterThan(money.Zero()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "debes indicar monto_bruto o monto_neto, mayor a 0"})
		return
	}

	premio, err := up.useCase.Execute(int32(idPremio), body.MontoBruto, body.MontoNeto)
	if err != nil {
		if errors.Is(err, application.ErrPremioNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Premio actualizado correctamente",
		"premio":  premio,
	})
}
