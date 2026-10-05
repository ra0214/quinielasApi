package infraestructure

import (
	"quinielas/src/saldos/application"
	"quinielas/src/saldos/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouterSaldos(r *gin.Engine, repo domain.ISaldo) {
	viewUseCase := application.NewViewSaldo(repo)
	updateUseCase := application.NewUpdateSaldo(repo)

	viewController := NewViewSaldoController(viewUseCase)
	editController := NewEditSaldoController(updateUseCase)

	api := r.Group("/saldos")
	{
		api.GET("", viewController.GetAll)
		api.GET("/cliente/:idCliente", viewController.GetByClienteID)
		api.PUT("/cliente/:idCliente", editController.Execute)
	}
}