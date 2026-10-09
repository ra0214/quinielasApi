package infraestructure

import (
	"net/http"
	"quinielas/src/ediciones/application"
	"quinielas/src/shared/apierror"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteEdicionController struct {
	useCase *application.DeleteEdicion
}

func NewDeleteEdicionController(useCase *application.DeleteEdicion) *DeleteEdicionController {
	return &DeleteEdicionController{useCase: useCase}
}

func (dc *DeleteEdicionController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edición inválido"})
		return
	}

	err = dc.useCase.Execute(int32(id))
	if err != nil {
		apierror.Eliminar(c, err, "la edición")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Edición eliminada correctamente"})
}