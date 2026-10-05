package infraestructure

import (
	"errors"
	"net/http"
	"strconv"

	"quinielas/src/premios/application"
	"quinielas/src/premios/domain"

	"github.com/gin-gonic/gin"
)

type RepartoPremioController struct {
	calcularUseCase *application.CalcularRepartoPremio
	viewUseCase     *application.ViewRepartoPremio
}

func NewRepartoPremioController(calcularUseCase *application.CalcularRepartoPremio, viewUseCase *application.ViewRepartoPremio) *RepartoPremioController {
	return &RepartoPremioController{calcularUseCase: calcularUseCase, viewUseCase: viewUseCase}
}

// statusReparto traduce un error del dominio al codigo HTTP que corresponde.
// Un meta no alcanzada o un reparto ya generado no son fallos del servidor:
// son reglas de negocio que el cliente necesita conocer.
func statusReparto(err error) int {
	switch {
	case errors.Is(err, domain.ErrRepartoYaGenerado),
		errors.Is(err, domain.ErrMetaNoAlcanzada),
		errors.Is(err, domain.ErrAportesSinMonto):
		return http.StatusConflict
	case errors.Is(err, domain.ErrPremioNoRegistrado),
		errors.Is(err, application.ErrRepartoNoGenerado):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrSinAportes),
		errors.Is(err, domain.ErrMetaInvalida):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Generar calcula el reparto del premio y lo congela. A partir de esta llamada
// el monto a pagar queda fijo aunque se registren más aportes.
func (rp *RepartoPremioController) Generar(c *gin.Context) {
	idQuiniela, err := strconv.Atoi(c.Param("idQuiniela"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de quiniela inválido"})
		return
	}

	reparto, err := rp.calcularUseCase.Execute(int32(idQuiniela))
	if err != nil {
		c.JSON(statusReparto(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reparto del premio generado y congelado correctamente",
		"reparto": reparto,
	})
}

// Get devuelve el reparto ya congelado. No recalcula nada.
func (rp *RepartoPremioController) Get(c *gin.Context) {
	idQuiniela, err := strconv.Atoi(c.Param("idQuiniela"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de quiniela inválido"})
		return
	}

	reparto, err := rp.viewUseCase.ExecuteByQuiniela(int32(idQuiniela))
	if err != nil {
		c.JSON(statusReparto(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reparto)
}
