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
	IDEdicion      int32       `json:"id_edicion" binding:"required"`
	NombreVariante string      `json:"nombre_variante" binding:"required"`
	Precio         money.Money `json:"precio"`
	MontoMeta      money.Money `json:"monto_meta"`
	FechaLimite    time.Time   `json:"fecha_limite" binding:"required"`
}

func (cq *CreateQuinielaController) Execute(c *gin.Context) {
	var body CreateQuinielaRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obligatorios faltantes o formato de fecha/número inválido"})
		return
	}

	if body.Precio.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "precio es obligatorio y debe ser mayor a 0"})
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

	quiniela, err := cq.useCase.Execute(body.IDEdicion, body.NombreVariante, body.Precio, body.MontoMeta, body.FechaLimite)
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