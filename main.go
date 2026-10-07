package main

import (
	"log"
	"os"
	"quinielas/src/config/middleware"

	aportesInfra "quinielas/src/aportes/infraestructure"
	clientesInfra "quinielas/src/clientes/infraestructure"
	edicionesInfra "quinielas/src/ediciones/infraestructure"
	movimientosInfra "quinielas/src/movimientos/infraestructure"
	premiosInfra "quinielas/src/premios/infraestructure"
	quinielasInfra "quinielas/src/quinielas/infraestructure"
	saldosInfra "quinielas/src/saldos/infraestructure"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Crear una SOLA instancia de Gin
	r := gin.Default()

	// 2. Usar el Middleware de CORS configurado en src/config/middleware/middleware.go
	r.Use(middleware.NewCorsMiddleware())

	// 3. Inicializar Repositorios MySQL de cada módulo
	clienteRepo := clientesInfra.NewMySQL()
	edicionRepo := edicionesInfra.NewMySQL()
	quinielaRepo := quinielasInfra.NewMySQL()
	saldoRepo := saldosInfra.NewMySQL()
	aporteRepo := aportesInfra.NewMySQL()
	movimientoRepo := movimientosInfra.NewMySQL()
	premioRepo := premiosInfra.NewMySQL()

	// 4. Pasar la MISMA instancia `r` a todos los routers
	clientesInfra.SetupRouterClientes(r, clienteRepo)
	edicionesInfra.SetupRouterEdiciones(r, edicionRepo)
	quinielasInfra.SetupRouterQuinielas(r, quinielaRepo)
	saldosInfra.SetupRouterSaldos(r, saldoRepo)
	aportesInfra.SetupRouterAportes(r, aporteRepo)
	movimientosInfra.SetupRouterMovimientos(r, movimientoRepo)
	premiosInfra.SetupRouterPremios(r, premioRepo, aporteRepo, quinielaRepo)

	// 5. Configuración del servidor
	r.SetTrustedProxies([]string{"127.0.0.1"})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor iniciando en el puerto :%s", port)

	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
