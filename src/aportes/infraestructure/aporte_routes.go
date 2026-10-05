package infraestructure

import (
	"quinielas/src/aportes/application"
	"quinielas/src/aportes/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouterAportes(r *gin.Engine, repo domain.IAporte) {
	createUseCase := application.NewCreateOrUpdateAporte(repo)
	deleteUseCase := application.NewDeleteAporte(repo)
	viewUseCase := application.NewViewAporte(repo)

	createController := NewCreateAporteController(createUseCase)
	deleteController := NewDeleteAporteController(deleteUseCase)
	viewController := NewViewAporteController(viewUseCase)

	api := r.Group("/aportes")
	{
		api.POST("", createController.Execute)
		api.GET("/quiniela/:idQuiniela", viewController.GetByQuiniela)
		api.GET("/cliente/:idCliente", viewController.GetByCliente)
		api.DELETE("/:id", deleteController.Execute)
	}
}