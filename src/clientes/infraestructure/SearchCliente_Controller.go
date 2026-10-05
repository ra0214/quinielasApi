package infraestructure

import (
	"net/http"
	"quinielas/src/clientes/application"

	"github.com/gin-gonic/gin"
)

type SearchClienteController struct {
	useCase *application.SearchCliente
}

func NewSearchClienteController(useCase *application.SearchCliente) *SearchClienteController {
	return &SearchClienteController{useCase: useCase}
}

func (sc_c *SearchClienteController) Execute(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro de búsqueda 'q' es requerido"})
		return
	}

	clientes, err := sc_c.useCase.Execute(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clientes)
}