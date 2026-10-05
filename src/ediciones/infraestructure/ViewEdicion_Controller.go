package infraestructure

import (
	"net/http"
	"quinielas/src/ediciones/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ViewEdicionController struct {
	useCase *application.ViewEdicion
}

func NewViewEdicionController(useCase *application.ViewEdicion) *ViewEdicionController {
	return &ViewEdicionController{useCase: useCase}
}

func (vc *ViewEdicionController) GetAll(c *gin.Context) {
	ediciones, err := vc.useCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ediciones)
}

func (vc *ViewEdicionController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edición inválido"})
		return
	}

	edicion, err := vc.useCase.ExecuteByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Edición no encontrada"})
		return
	}

	c.JSON(http.StatusOK, edicion)
}