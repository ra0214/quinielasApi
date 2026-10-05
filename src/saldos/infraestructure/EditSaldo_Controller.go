package infraestructure

import (
	"net/http"
	"quinielas/src/saldos/application"
	"quinielas/src/shared/money"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EditSaldoController struct {
	useCase *application.UpdateSaldo
}

func NewEditSaldoController(useCase *application.UpdateSaldo) *EditSaldoController {
	return &EditSaldoController{useCase: useCase}
}

type EditSaldoRequestBody struct {
	SaldoFavor money.Money `json:"saldo_favor"`
	SaldoDeuda money.Money `json:"saldo_deuda"`
}

func (es *EditSaldoController) Execute(c *gin.Context) {
	idStr := c.Param("idCliente")
	idCliente, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente inválido"})
		return
	}

	var body EditSaldoRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if body.SaldoFavor.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "saldo_favor no puede ser negativo"})
		return
	}

	if body.SaldoDeuda.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "saldo_deuda no puede ser negativo"})
		return
	}

	err = es.useCase.Execute(int32(idCliente), body.SaldoFavor, body.SaldoDeuda)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el saldo", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Saldo del cliente actualizado correctamente"})
}