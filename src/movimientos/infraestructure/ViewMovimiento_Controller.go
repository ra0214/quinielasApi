package infraestructure

import (
	"net/http"
	"quinielas/src/movimientos/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ViewMovimientoController struct {
	useCase *application.ViewMovimiento
}

func NewViewMovimientoController(useCase *application.ViewMovimiento) *ViewMovimientoController {
	return &ViewMovimientoController{useCase: useCase}
}

func (vm *ViewMovimientoController) GetAll(c *gin.Context) {
	movimientos, err := vm.useCase.ExecuteGetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movimientos)
}

func (vm *ViewMovimientoController) GetByClienteID(c *gin.Context) {
	idStr := c.Param("idCliente")
	idCliente, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente inválido"})
		return
	}

	movimientos, err := vm.useCase.ExecuteByClienteID(int32(idCliente))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movimientos)
}