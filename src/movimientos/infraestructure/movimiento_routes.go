package infraestructure

import (
	"quinielas/src/movimientos/application"
	"quinielas/src/movimientos/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouterMovimientos(r *gin.Engine, repo domain.IMovimiento) {
	createUseCase := application.NewCreateMovimiento(repo)
	deleteUseCase := application.NewDeleteMovimiento(repo)
	viewUseCase := application.NewViewMovimiento(repo)

	createController := NewCreateMovimientoController(createUseCase)
	deleteController := NewDeleteMovimientoController(deleteUseCase)
	viewController := NewViewMovimientoController(viewUseCase)

	api := r.Group("/movimientos")
	{
		api.POST("", createController.Execute)
		api.GET("", viewController.GetAll)
		api.GET("/cliente/:idCliente", viewController.GetByClienteID)
		api.DELETE("/:id", deleteController.Execute)
	}
}