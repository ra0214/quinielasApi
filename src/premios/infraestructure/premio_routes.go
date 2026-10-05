package infraestructure

import (
	"quinielas/src/premios/application"
	"quinielas/src/premios/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouterPremios(r *gin.Engine, repo domain.IPremio, aporteRepo application.AportesDeQuiniela, quinielaRepo application.MetaDeQuiniela) {
	createUseCase := application.NewCreatePremio(repo)
	deleteUseCase := application.NewDeletePremio(repo)
	viewUseCase := application.NewViewPremio(repo)
	calcularRepartoUseCase := application.NewCalcularRepartoPremio(repo, aporteRepo, quinielaRepo)
	viewRepartoUseCase := application.NewViewRepartoPremio(repo)

	createController := NewCreatePremioController(createUseCase)
	deleteController := NewDeletePremioController(deleteUseCase)
	viewController := NewViewPremioController(viewUseCase)
	repartoController := NewRepartoPremioController(calcularRepartoUseCase, viewRepartoUseCase)

	api := r.Group("/premios")
	{
		api.POST("", createController.Execute)
		api.GET("", viewController.GetAll)
		api.GET("/quiniela/:idQuiniela", viewController.GetByQuinielaID)
		api.GET("/quiniela/:idQuiniela/reparto", repartoController.Get)
		api.POST("/quiniela/:idQuiniela/reparto", repartoController.Generar)
		api.DELETE("/:id", deleteController.Execute)
	}
}
