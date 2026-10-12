package infraestructure

import (
	"quinielas/src/aportes/application"
	"quinielas/src/aportes/domain"
	quinielasDomain "quinielas/src/quinielas/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouterAportes(r *gin.Engine, repo domain.IAporte, quinielaRepo quinielasDomain.IQuiniela, cobro domain.CobroAporte) {
	createUseCase := application.NewCreateOrUpdateAporte(repo, quinielaRepo, cobro)
	updateUseCase := application.NewUpdateAporte(repo, quinielaRepo, cobro)
	deleteUseCase := application.NewDeleteAporte(repo, cobro)
	viewUseCase := application.NewViewAporte(repo)

	createController := NewCreateAporteController(createUseCase)
	updateController := NewUpdateAporteController(updateUseCase)
	deleteController := NewDeleteAporteController(deleteUseCase)
	viewController := NewViewAporteController(viewUseCase)

	api := r.Group("/aportes")
	{
		api.POST("", createController.Execute)
		api.GET("", viewController.GetAll)
		api.PUT("/:id", updateController.Execute)
		api.GET("/quiniela/:idQuiniela", viewController.GetByQuiniela)
		api.GET("/cliente/:idCliente", viewController.GetByCliente)
		api.DELETE("/:id", deleteController.Execute)
	}
}