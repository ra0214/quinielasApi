package infraestructure

import (
	"quinielas/src/aportes/application"
	"quinielas/src/aportes/domain"
	quinielasDomain "quinielas/src/quinielas/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouterAportes(r *gin.Engine, repo domain.IAporte, quinielaRepo quinielasDomain.IQuiniela) {
	createUseCase := application.NewCreateOrUpdateAporte(repo, quinielaRepo)
	updateUseCase := application.NewUpdateAporte(repo, quinielaRepo)
	deleteUseCase := application.NewDeleteAporte(repo)
	viewUseCase := application.NewViewAporte(repo)

	createController := NewCreateAporteController(createUseCase)
	updateController := NewUpdateAporteController(updateUseCase)
	deleteController := NewDeleteAporteController(deleteUseCase)
	viewController := NewViewAporteController(viewUseCase)

	api := r.Group("/aportes")
	{
		api.POST("", createController.Execute)
		api.PUT("/:id", updateController.Execute)
		api.GET("/quiniela/:idQuiniela", viewController.GetByQuiniela)
		api.GET("/cliente/:idCliente", viewController.GetByCliente)
		api.DELETE("/:id", deleteController.Execute)
	}
}