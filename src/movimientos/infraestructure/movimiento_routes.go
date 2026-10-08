package infraestructure

import (
	"quinielas/src/movimientos/application"
	"quinielas/src/movimientos/domain"
	saldosApp "quinielas/src/saldos/application"
	saldosDomain "quinielas/src/saldos/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouterMovimientos(r *gin.Engine, repo domain.IMovimiento, saldoRepo saldosDomain.ISaldo) {
	ajustador := saldosApp.NewAjustarSaldo(saldoRepo)

	createUseCase := application.NewCreateMovimiento(repo, ajustador)
	deleteUseCase := application.NewDeleteMovimiento(repo, ajustador)
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