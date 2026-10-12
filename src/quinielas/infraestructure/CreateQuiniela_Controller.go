package infraestructure

import (
	"net/http"
	"quinielas/src/quinielas/application"
	"quinielas/src/quinielas/domain"
	"quinielas/src/shared/money"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateQuinielaController struct {
	useCase *application.CreateQuiniela
}

func NewCreateQuinielaController(useCase *application.CreateQuiniela) *CreateQuinielaController {
	return &CreateQuinielaController{useCase: useCase}
}

type CreateQuinielaRequestBody struct {
	IDEdicion       int32       `json:"id_edicion" binding:"required"`
	NombreVariante  string      `json:"nombre_variante"`
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

func validarCanalVenta(canal string) bool {
	return canal == "" || canal == "EN_LINEA" || canal == "VENTANILLA"
}

func (cq *CreateQuinielaController) Execute(c *gin.Context) {
	var body CreateQuinielaRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obligatorios faltantes o formato de fecha/número inválido"})
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
		IDEdicion:      body.IDEdicion,
		NombreVariante: body.NombreVariante,
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
	quiniela, err := cq.useCase.Execute(q)
	if err != nil {
		if err == domain.ErrLimitQuinielasExceeded {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Quiniela registrada correctamente",
		"quiniela": quiniela,
	})
}