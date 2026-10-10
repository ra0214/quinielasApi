package infraestructure

import (
	"errors"
	"net/http"
	movimientosDomain "quinielas/src/movimientos/domain"
	"quinielas/src/movimientos/application"
	saldosDomain "quinielas/src/saldos/domain"
	"quinielas/src/shared/money"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateMovimientoController struct {
	useCase *application.UpdateMovimiento
}

func NewUpdateMovimientoController(useCase *application.UpdateMovimiento) *UpdateMovimientoController {
	return &UpdateMovimientoController{useCase: useCase}
}

type UpdateMovimientoRequestBody struct {
	IDCliente   int32       `json:"id_cliente" binding:"required"`
	IDQuiniela  *int32      `json:"id_quiniela"`
	Tipo        string      `json:"tipo" binding:"required"`
	Monto       money.Money `json:"monto"`
	Descripcion string      `json:"descripcion"`
}

func (um *UpdateMovimientoController) Execute(c *gin.Context) {
	idMovimiento, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de movimiento inválido"})
		return
	}

	var body UpdateMovimientoRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obligatorios faltantes o JSON inválido"})
		return
	}

	if body.Monto.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monto es obligatorio y debe ser mayor a 0"})
		return
	}

	if body.Monto.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monto no puede ser negativo"})
		return
	}

	nuevo := &movimientosDomain.Movimiento{
		IDMovimiento: int32(idMovimiento),
		IDCliente:    body.IDCliente,
		IDQuiniela:   body.IDQuiniela,
		Tipo:         body.Tipo,
		Monto:        body.Monto,
		Descripcion:  body.Descripcion,
	}

	movimiento, err := um.useCase.Execute(int32(idMovimiento), nuevo)
	if err != nil {
		if errors.Is(err, saldosDomain.ErrSaldoInsuficiente) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Movimiento actualizado correctamente",
		"movimiento": movimiento,
	})
}