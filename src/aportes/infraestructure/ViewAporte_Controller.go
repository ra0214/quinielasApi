package infraestructure

import (
	"net/http"
	"quinielas/src/aportes/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ViewAporteController struct {
	useCase *application.ViewAporte
}

func NewViewAporteController(useCase *application.ViewAporte) *ViewAporteController {
	return &ViewAporteController{useCase: useCase}
}

func (va *ViewAporteController) GetByQuiniela(c *gin.Context) {
	idStr := c.Param("idQuiniela")
	idQuiniela, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de quiniela inválido"})
		return
	}

	aportes, err := va.useCase.ExecuteByQuiniela(int32(idQuiniela))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aportes)
}

func (va *ViewAporteController) GetByCliente(c *gin.Context) {
	idStr := c.Param("idCliente")
	idCliente, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente inválido"})
		return
	}

	aportes, err := va.useCase.ExecuteByCliente(int32(idCliente))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aportes)
}