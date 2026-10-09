package infraestructure

import (
	"net/http"
	"quinielas/src/premios/application"
	"quinielas/src/shared/apierror"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeletePremioController struct {
	useCase *application.DeletePremio
}

func NewDeletePremioController(useCase *application.DeletePremio) *DeletePremioController {
	return &DeletePremioController{useCase: useCase}
}

func (dp *DeletePremioController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	idPremio, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de premio inválido"})
		return
	}

	err = dp.useCase.Execute(int32(idPremio))
	if err != nil {
		apierror.Eliminar(c, err, "el premio")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Premio eliminado correctamente"})
}