package infraestructure

import (
	"net/http"
	"quinielas/src/ediciones/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EditEdicionController struct {
	useCase *application.UpdateEdicion
}

func NewEditEdicionController(useCase *application.UpdateEdicion) *EditEdicionController {
	return &EditEdicionController{useCase: useCase}
}

type EditEdicionRequestBody struct {
	TipoEdicion   string `json:"tipo_edicion" binding:"required"`
	NombreEdicion string `json:"nombre_edicion" binding:"required"`
}

func (ec *EditEdicionController) Execute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edición inválido"})
		return
	}

	var body EditEdicionRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	err = ec.useCase.Execute(int32(id), body.TipoEdicion, body.NombreEdicion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la edición", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Edición actualizada correctamente"})
}