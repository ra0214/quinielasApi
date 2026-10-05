package infraestructure

import (
	"net/http"
	"quinielas/src/premios/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ViewPremioController struct {
	useCase *application.ViewPremio
}

func NewViewPremioController(useCase *application.ViewPremio) *ViewPremioController {
	return &ViewPremioController{useCase: useCase}
}

func (vp *ViewPremioController) GetAll(c *gin.Context) {
	premios, err := vp.useCase.ExecuteGetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, premios)
}

func (vp *ViewPremioController) GetByQuinielaID(c *gin.Context) {
	idStr := c.Param("idQuiniela")
	idQuiniela, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de quiniela inválido"})
		return
	}

	premio, err := vp.useCase.ExecuteByQuinielaID(int32(idQuiniela))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No hay un premio registrado para esta quiniela"})
		return
	}

	c.JSON(http.StatusOK, premio)
}