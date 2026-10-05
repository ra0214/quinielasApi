package infraestructure

import (
	"net/http"
	"quinielas/src/quinielas/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ViewQuinielaController struct {
	useCase *application.ViewQuiniela
}

func NewViewQuinielaController(useCase *application.ViewQuiniela) *ViewQuinielaController {
	return &ViewQuinielaController{useCase: useCase}
}

func (vq *ViewQuinielaController) GetAll(c *gin.Context) {
	quinielas, err := vq.useCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quinielas)
}

func (vq *ViewQuinielaController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de quiniela inválido"})
		return
	}

	quiniela, err := vq.useCase.ExecuteByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiniela no encontrada"})
		return
	}

	c.JSON(http.StatusOK, quiniela)
}

func (vq *ViewQuinielaController) GetByEdicionID(c *gin.Context) {
	idStr := c.Param("idEdicion")
	idEdicion, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edición inválido"})
		return
	}

	quinielas, err := vq.useCase.ExecuteByEdicionID(int32(idEdicion))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quinielas)
}