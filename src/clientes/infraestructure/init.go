// init.go
package infraestructure

import (
	"log"
)

func Init() {
	deps := InitDependencies()
	router := SetupRouter(deps.ClienteRepo)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}