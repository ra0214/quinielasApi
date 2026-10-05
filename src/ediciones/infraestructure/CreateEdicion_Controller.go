package infraestructure

import (
	"net/http"
	"quinielas/src/ediciones/application"

	"github.com/gin-gonic/gin"
)

type CreateEdicionController struct {
	useCase *application.CreateEdicion
}

func NewCreateEdicionController(useCase *application.CreateEdicion) *CreateEdicionController {
	return &CreateEdicionController{useCase: useCase}
}

type CreateEdicionRequestBody struct {
	TipoEdicion   string `json:"tipo_edicion" binding:"required"`
	NombreEdicion string `json:"nombre_edicion" binding:"required"`
}

func (cc *CreateEdicionController) Execute(c *gin.Context) {
	var body CreateEdicionRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obligatorios requeridos o JSON inválido"})
		return
	}

	edicion, err := cc.useCase.Execute(body.TipoEdicion, body.NombreEdicion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Edición creada correctamente",
		"edicion": edicion,
	})
}