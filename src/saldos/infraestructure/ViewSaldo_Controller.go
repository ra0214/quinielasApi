package infraestructure

import (
	"net/http"
	"quinielas/src/saldos/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ViewSaldoController struct {
	useCase *application.ViewSaldo
}

func NewViewSaldoController(useCase *application.ViewSaldo) *ViewSaldoController {
	return &ViewSaldoController{useCase: useCase}
}

func (vs *ViewSaldoController) GetAll(c *gin.Context) {
	filtro := c.Query("filtro") // 'deben', 'favor', 'ceros'

	if filtro != "" {
		saldos, err := vs.useCase.ExecuteFiltrados(filtro)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, saldos)
		return
	}

	saldos, err := vs.useCase.ExecuteGetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, saldos)
}

func (vs *ViewSaldoController) GetByClienteID(c *gin.Context) {
	idStr := c.Param("idCliente")
	idCliente, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente inválido"})
		return
	}

	saldo, err := vs.useCase.ExecuteGetByClienteID(int32(idCliente))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Saldo no encontrado para este cliente"})
		return
	}

	c.JSON(http.StatusOK, saldo)
}