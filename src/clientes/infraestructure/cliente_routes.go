package infraestructure

import (
	"quinielas/src/clientes/application"
	"quinielas/src/clientes/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouter(repo domain.ICliente) *gin.Engine {
	r := gin.Default()

	// Inicializar Casos de Uso
	createUseCase := application.NewCreateCliente(repo)
	deleteUseCase := application.NewDeleteCliente(repo)
	updateUseCase := application.NewUpdateCliente(repo)
	viewUseCase := application.NewViewCliente(repo)
	searchUseCase := application.NewSearchCliente(repo)

	// Inicializar Controladores
	createController := NewCreateClienteController(createUseCase)
	deleteController := NewDeleteClienteController(deleteUseCase)
	editController := NewEditClienteController(updateUseCase)
	viewController := NewViewClienteController(viewUseCase)
	searchController := NewSearchClienteController(searchUseCase)

	// Rutas de la API de Clientes
	api := r.Group("/clientes")
	{
		api.POST("", createController.Execute)
		api.GET("", viewController.GetAll)
		api.GET("/search", searchController.Execute)
		api.GET("/:id", viewController.GetByID)
		api.PUT("/:id", editController.Execute)
		api.DELETE("/:id", deleteController.Execute)
	}

	return r
}