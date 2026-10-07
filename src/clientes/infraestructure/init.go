// init.go
package infraestructure

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Init() {
	deps := InitDependencies()
	router := gin.Default()
	SetupRouterClientes(router, deps.ClienteRepo)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}