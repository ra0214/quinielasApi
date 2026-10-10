package infraestructure

import (
	"errors"
	"net/http"
	"quinielas/src/aportes/application"
	"quinielas/src/aportes/domain"
	"quinielas/src/shared/money"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateAporteController struct {
	useCase *application.UpdateAporte
}

func NewUpdateAporteController(useCase *application.UpdateAporte) *UpdateAporteController {
	return &UpdateAporteController{useCase: useCase}
}

type UpdateAporteRequestBody struct {
	Monto money.Money `json:"monto"`
}

func (ua *UpdateAporteController) Execute(c *gin.Context) {
	idAporte, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de aporte inválido"})
		return
	}

	var body UpdateAporteRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if body.Monto.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el monto debe ser mayor a 0"})
		return
	}

	if body.Monto.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monto no puede ser negativo"})
		return
	}

	aporte, err := ua.useCase.Execute(int32(idAporte), body.Monto)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMetaExcedida):
			c.JSON(http.StatusBadRequest, gin.H{"error": "No se puede exceder el monto total solicitado de la quiniela"})
		case errors.Is(err, domain.ErrQuinielaNoEncontrada):
			c.JSON(http.StatusNotFound, gin.H{"error": "La quiniela no existe"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Aporte actualizado correctamente",
		"aporte":  aporte,
	})
}
