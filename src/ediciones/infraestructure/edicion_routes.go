package infraestructure

import (
	"quinielas/src/ediciones/application"
	"quinielas/src/ediciones/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouterEdiciones(r *gin.Engine, repo domain.IEdicion) {
	createUseCase := application.NewCreateEdicion(repo)
	deleteUseCase := application.NewDeleteEdicion(repo)
	updateUseCase := application.NewUpdateEdicion(repo)
	viewUseCase := application.NewViewEdicion(repo)

	createController := NewCreateEdicionController(createUseCase)
	deleteController := NewDeleteEdicionController(deleteUseCase)
	editController := NewEditEdicionController(updateUseCase)
	viewController := NewViewEdicionController(viewUseCase)

	api := r.Group("/ediciones")
	{
		api.POST("", createController.Execute)
		api.GET("", viewController.GetAll)
		api.GET("/:id", viewController.GetByID)
		api.PUT("/:id", editController.Execute)
		api.DELETE("/:id", deleteController.Execute)
	}
}