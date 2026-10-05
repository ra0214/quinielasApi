package infraestructure

import (
	"net/http"
	"quinielas/src/clientes/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EditClienteController struct {
	useCase *application.UpdateCliente
}

func NewEditClienteController(useCase *application.UpdateCliente) *EditClienteController {
	return &EditClienteController{useCase: useCase}
}

type EditClienteRequestBody struct {
	Nombre   string `json:"nombre" binding:"required"`
	Telefono string `json:"telefono"`
}

func (ec_c *EditClienteController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente inválido"})
		return
	}

	var body EditClienteRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	err = ec_c.useCase.Execute(int32(id), body.Nombre, body.Telefono)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el cliente", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cliente actualizado correctamente"})
}