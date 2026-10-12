package infraestructure

import (
	"net/http"
	"quinielas/src/quinielas/application"
	"quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EditQuinielaController struct {
	useCase *application.UpdateQuiniela
}

func NewEditQuinielaController(useCase *application.UpdateQuiniela) *EditQuinielaController {
	return &EditQuinielaController{useCase: useCase}
}

type EditQuinielaRequestBody struct {
	NombreVariante  string      `json:"nombre_variante"`
	Precio          money.Money `json:"precio"`
	MontoMeta       money.Money `json:"monto_meta"`
	FechaLimite     time.Time   `json:"fecha_limite"`
	Estado          string      `json:"estado"`
	TipoQuiniela    string      `json:"tipo_quiniela"`
	IDCliente       *int32      `json:"id_cliente,omitempty"`
	MontoPagado     money.Money `json:"monto_pagado"`
	Analista        string      `json:"analista"`
	TipoPago        string      `json:"tipo_pago"`
	CanalVenta      string      `json:"canal_venta"`
}

func (eq *EditQuinielaController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de quiniela inválido"})
		return
	}

	var body EditQuinielaRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if body.Precio.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "precio no puede ser negativo"})
		return
	}

	if body.MontoMeta.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monto_meta no puede ser negativo"})
		return
	}

	if !validarCanalVenta(body.CanalVenta) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "canal_venta debe ser EN_LINEA o VENTANILLA"})
		return
	}

	q := &domain.Quiniela{
		ID:             int32(id),
		NombreVariante: body.NombreVariante,
		Precio:         body.Precio,
		MontoMeta:      body.MontoMeta,
		FechaLimite:    body.FechaLimite,
		Estado:         body.Estado,
		TipoQuiniela:   body.TipoQuiniela,
		IDCliente:      body.IDCliente,
		MontoPagado:    body.MontoPagado,
		Analista:       body.Analista,
		TipoPago:       body.TipoPago,
		CanalVenta:     body.CanalVenta,
	}
	if q.Estado == "" {
		q.Estado = "EN_JUEGO"
	}
	err = eq.useCase.Execute(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la quiniela", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quiniela actualizada correctamente"})
}