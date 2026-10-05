package main

import (
	"log"
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
	r := gin.Default()

	// 1. Usar el Middleware de CORS configurado en src/config/middleware/middleware.go
	r.Use(middleware.NewCorsMiddleware())

	// 2. Inicializar Repositorios MySQL de cada módulo
	clienteRepo := clientesInfra.NewMySQL()
	edicionRepo := edicionesInfra.NewMySQL()
	quinielaRepo := quinielasInfra.NewMySQL()
	saldoRepo := saldosInfra.NewMySQL()
	aporteRepo := aportesInfra.NewMySQL()
	movimientoRepo := movimientosInfra.NewMySQL()
	premioRepo := premiosInfra.NewMySQL()

	// 3. Registrar Rutas de Clientes
	clienteRouter := clientesInfra.SetupRouter(clienteRepo)
	for _, route := range clienteRouter.Routes() {
		r.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	// 4. Registrar Rutas de Ediciones
	edicionRouter := gin.New()
	edicionesInfra.SetupRouterEdiciones(edicionRouter, edicionRepo)
	for _, route := range edicionRouter.Routes() {
		r.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	// 5. Registrar Rutas de Quinielas
	quinielaRouter := gin.New()
	quinielasInfra.SetupRouterQuinielas(quinielaRouter, quinielaRepo)
	for _, route := range quinielaRouter.Routes() {
		r.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	// 6. Registrar Rutas de Saldos
	saldoRouter := gin.New()
	saldosInfra.SetupRouterSaldos(saldoRouter, saldoRepo)
	for _, route := range saldoRouter.Routes() {
		r.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	// 7. Registrar Rutas de Aportes
	aporteRouter := gin.New()
	aportesInfra.SetupRouterAportes(aporteRouter, aporteRepo)
	for _, route := range aporteRouter.Routes() {
		r.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	// 8. Registrar Rutas de Movimientos
	movimientoRouter := gin.New()
	movimientosInfra.SetupRouterMovimientos(movimientoRouter, movimientoRepo)
	for _, route := range movimientoRouter.Routes() {
		r.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	// 9. Registrar Rutas de Premios
	premioRouter := gin.New()
	premiosInfra.SetupRouterPremios(premioRouter, premioRepo)
	for _, route := range premioRouter.Routes() {
		r.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	// 10. Configuración del servidor
	r.SetTrustedProxies([]string{"127.0.0.1"})

	log.Println("Servidor corriendo en http://localhost:8080")

	// Iniciar servidor
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
