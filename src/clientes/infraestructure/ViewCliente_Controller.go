package infraestructure

import (
	"net/http"
	"quinielas/src/clientes/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ViewClienteController struct {
	useCase *application.ViewCliente
}

func NewViewClienteController(useCase *application.ViewCliente) *ViewClienteController {
	return &ViewClienteController{useCase: useCase}
}

func (vc_c *ViewClienteController) GetAll(c *gin.Context) {
	clientes, err := vc_c.useCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clientes)
}

func (vc_c *ViewClienteController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente inválido"})
		return
	}

	cliente, err := vc_c.useCase.ExecuteByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}

	c.JSON(http.StatusOK, cliente)
}