package infraestructure

import (
	"quinielas/src/quinielas/application"
	"quinielas/src/quinielas/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouterQuinielas(r *gin.Engine, repo domain.IQuiniela) {
	createUseCase := application.NewCreateQuiniela(repo)
	deleteUseCase := application.NewDeleteQuiniela(repo)
	updateUseCase := application.NewUpdateQuiniela(repo)
	viewUseCase := application.NewViewQuiniela(repo)

	createController := NewCreateQuinielaController(createUseCase)
	deleteController := NewDeleteQuinielaController(deleteUseCase)
	editController := NewEditQuinielaController(updateUseCase)
	viewController := NewViewQuinielaController(viewUseCase)

	api := r.Group("/quinielas")
	{
		api.POST("", createController.Execute)
		api.GET("", viewController.GetAll)
		api.GET("/edicion/:idEdicion", viewController.GetByEdicionID)
		api.GET("/:id", viewController.GetByID)
		api.PUT("/:id", editController.Execute)
		api.DELETE("/:id", deleteController.Execute)
	}
}