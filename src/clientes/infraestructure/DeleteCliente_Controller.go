package infraestructure

import (
	"errors"
	"net/http"
	"quinielas/src/clientes/application"
	"quinielas/src/shared/apierror"
	saldosDomain "quinielas/src/saldos/domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteClienteController struct {
	useCase *application.DeleteCliente
}

func NewDeleteClienteController(useCase *application.DeleteCliente) *DeleteClienteController {
	return &DeleteClienteController{useCase: useCase}
}

func (dc_c *DeleteClienteController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente inválido"})
		return
	}

	err = dc_c.useCase.Execute(int32(id))
	if err != nil {
		var saldoErr *saldosDomain.SaldoPendienteError
		if errors.As(err, &saldoErr) {
			c.JSON(http.StatusConflict, gin.H{"error": "No se puede eliminar el cliente: " + saldoErr.Error()})
			return
		}
		apierror.Eliminar(c, err, "el cliente")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cliente eliminado correctamente"})
}