package infraestructure

import (
	"net/http"
	"quinielas/src/movimientos/application"
	"quinielas/src/shared/apierror"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteMovimientoController struct {
	useCase *application.DeleteMovimiento
}

func NewDeleteMovimientoController(useCase *application.DeleteMovimiento) *DeleteMovimientoController {
	return &DeleteMovimientoController{useCase: useCase}
}

func (dm *DeleteMovimientoController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	idMovimiento, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de movimiento inválido"})
		return
	}

	err = dm.useCase.Execute(int32(idMovimiento))
	if err != nil {
		apierror.Eliminar(c, err, "el movimiento")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movimiento eliminado correctamente"})
}