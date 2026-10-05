package infraestructure

import (
	"net/http"
	"quinielas/src/movimientos/application"
	"quinielas/src/shared/money"

	"github.com/gin-gonic/gin"
)

type CreateMovimientoController struct {
	useCase *application.CreateMovimiento
}

func NewCreateMovimientoController(useCase *application.CreateMovimiento) *CreateMovimientoController {
	return &CreateMovimientoController{useCase: useCase}
}

type CreateMovimientoRequestBody struct {
	IDCliente   int32       `json:"id_cliente" binding:"required"`
	IDQuiniela  *int32      `json:"id_quiniela"`
	Tipo        string      `json:"tipo" binding:"required"` // 'PAGO_EFECTIVO', 'PAGO_SALDO', 'FIADO', 'RETIRO', 'PREMIO_ABONO', 'DEVOLUCION'
	Monto       money.Money `json:"monto"`
	Descripcion string      `json:"descripcion"`
}

func (cm *CreateMovimientoController) Execute(c *gin.Context) {
	var body CreateMovimientoRequestBody
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

	movimiento, err := cm.useCase.Execute(body.IDCliente, body.IDQuiniela, body.Tipo, body.Monto, body.Descripcion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Movimiento registrado correctamente",
		"movimiento": movimiento,
	})
}