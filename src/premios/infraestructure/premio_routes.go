package infraestructure

import (
	"quinielas/src/premios/application"
	"quinielas/src/premios/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouterPremios(r *gin.Engine, repo domain.IPremio) {
	createUseCase := application.NewCreatePremio(repo)
	deleteUseCase := application.NewDeletePremio(repo)
	viewUseCase := application.NewViewPremio(repo)

	createController := NewCreatePremioController(createUseCase)
	deleteController := NewDeletePremioController(deleteUseCase)
	viewController := NewViewPremioController(viewUseCase)

	api := r.Group("/premios")
	{
		api.POST("", createController.Execute)
		api.GET("", viewController.GetAll)
		api.GET("/quiniela/:idQuiniela", viewController.GetByQuinielaID)
		api.DELETE("/:id", deleteController.Execute)
	}
}