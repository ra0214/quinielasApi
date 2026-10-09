package infraestructure

import (
	"net/http"
	"quinielas/src/aportes/application"
	"quinielas/src/shared/apierror"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteAporteController struct {
	useCase *application.DeleteAporte
}

func NewDeleteAporteController(useCase *application.DeleteAporte) *DeleteAporteController {
	return &DeleteAporteController{useCase: useCase}
}

func (da *DeleteAporteController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	idAporte, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de aporte inválido"})
		return
	}

	err = da.useCase.Execute(int32(idAporte))
	if err != nil {
		apierror.Eliminar(c, err, "el aporte")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Aporte eliminado correctamente"})
}