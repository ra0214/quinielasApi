package infraestructure

import (
	"net/http"
	"quinielas/src/clientes/application"
	"quinielas/src/shared/apierror"
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
		apierror.Eliminar(c, err, "el cliente")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cliente eliminado correctamente"})
}