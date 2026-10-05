package infraestructure

import (
	"net/http"
	"quinielas/src/clientes/application"

	"github.com/gin-gonic/gin"
)

type CreateClienteController struct {
	useCase *application.CreateCliente
}

func NewCreateClienteController(useCase *application.CreateCliente) *CreateClienteController {
	return &CreateClienteController{useCase: useCase}
}

type CreateClienteRequestBody struct {
	Nombre   string `json:"nombre" binding:"required"`
	Telefono string `json:"telefono"`
}

func (cc_c *CreateClienteController) Execute(c *gin.Context) {
	var body CreateClienteRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre es requerido o formato JSON inválido"})
		return
	}

	cliente, err := cc_c.useCase.Execute(body.Nombre, body.Telefono)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cliente registrado correctamente",
		"cliente": cliente,
	})
}