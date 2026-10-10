package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"quinielas/src/config/middleware"

	aportesInfra "quinielas/src/aportes/infraestructure"
	clientesInfra "quinielas/src/clientes/infraestructure"
	edicionesInfra "quinielas/src/ediciones/infraestructure"
	movimientosApp "quinielas/src/movimientos/application"
	movimientosInfra "quinielas/src/movimientos/infraestructure"
	premiosInfra "quinielas/src/premios/infraestructure"
	quinielasInfra "quinielas/src/quinielas/infraestructure"
	saldosApp "quinielas/src/saldos/application"
	saldosInfra "quinielas/src/saldos/infraestructure"
)

func main() {
	r := gin.Default()

	r.Use(middleware.NewCorsMiddleware())

	// Inicializar Repositorios
	clienteRepo := clientesInfra.NewMySQL()
	edicionRepo := edicionesInfra.NewMySQL()
	quinielaRepo := quinielasInfra.NewMySQL()
	saldoRepo := saldosInfra.NewMySQL()
	aporteRepo := aportesInfra.NewMySQL()
	movimientoRepo := movimientosInfra.NewMySQL()
	premioRepo := premiosInfra.NewMySQL()

	// Registrar Rutas directamente en la misma instancia 'r'
	clientesInfra.SetupRouterClientes(r, clienteRepo, saldoRepo)
	edicionesInfra.SetupRouterEdiciones(r, edicionRepo)
	quinielasInfra.SetupRouterQuinielas(r, quinielaRepo)
	saldosInfra.SetupRouterSaldos(r, saldoRepo)
	aportesInfra.SetupRouterAportes(r, aporteRepo, quinielaRepo)
	movimientosInfra.SetupRouterMovimientos(r, movimientoRepo, saldoRepo)

	// Abonador de premios: registra cada parte del reparto como movimiento
	// PREMIO_ABONO (ajusta el saldo) y puede deshacerse si el reparto falla.
	ajustadorSaldo := saldosApp.NewAjustarSaldo(saldoRepo)
	abonadorPremio := movimientosApp.NewAbonadorPremio(
		movimientosApp.NewCreateMovimiento(movimientoRepo, ajustadorSaldo),
		movimientosApp.NewDeleteMovimiento(movimientoRepo, ajustadorSaldo),
	)
	premiosInfra.SetupRouterPremios(r, premioRepo, aporteRepo, quinielaRepo, abonadorPremio)

	// Obtener el puerto dinámico de Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor escuchando en el puerto :%s", port)

	// Escuchar en 0.0.0.0
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
