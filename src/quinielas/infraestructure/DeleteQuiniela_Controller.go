package infraestructure

import (
	"fmt"
	"net/http"
	"quinielas/src/quinielas/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteQuinielaController struct {
	useCase *application.DeleteQuiniela
}

func NewDeleteQuinielaController(useCase *application.DeleteQuiniela) *DeleteQuinielaController {
	return &DeleteQuinielaController{useCase: useCase}
}

func (dq *DeleteQuinielaController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de quiniela inválido"})
		return
	}

	err = dq.useCase.Execute(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al eliminar la quiniela: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quiniela eliminada correctamente"})
}