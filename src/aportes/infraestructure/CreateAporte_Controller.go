package infraestructure

import (
	"net/http"
	"quinielas/src/aportes/application"
	"quinielas/src/shared/money"

	"github.com/gin-gonic/gin"
)

type CreateAporteController struct {
	useCase *application.CreateOrUpdateAporte
}

func NewCreateAporteController(useCase *application.CreateOrUpdateAporte) *CreateAporteController {
	return &CreateAporteController{useCase: useCase}
}

type CreateAporteRequestBody struct {
	IDCliente         int32       `json:"id_cliente" binding:"required"`
	IDQuiniela        int32       `json:"id_quiniela" binding:"required"`
	Monto             money.Money `json:"monto"`
	MontoMetaQuiniela money.Money `json:"monto_meta"`
}

func (ca *CreateAporteController) Execute(c *gin.Context) {
	var body CreateAporteRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obligatorios faltantes o JSON inválido"})
		return
	}

	if body.Monto.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monto no puede ser negativo"})
		return
	}

	if body.MontoMetaQuiniela.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monto_meta no puede ser negativo"})
		return
	}

	aporte, err := ca.useCase.Execute(body.IDCliente, body.IDQuiniela, body.Monto, body.MontoMetaQuiniela)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Aporte registrado/actualizado correctamente",
		"aporte":  aporte,
	})
}