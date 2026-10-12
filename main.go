package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"quinielas/src/config"
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

	// Auto-alter para columnas nuevas de quinielas
	dbConn := config.GetDBPool()
	if dbConn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", dbConn.Err)
	}
	if err := dbConn.AutoAlterQuinielas(); err != nil {
		log.Printf("Advertencia al ejecutar auto-alter de quinielas: %v", err)
	}
	if err := dbConn.AutoAlterPremios(); err != nil {
		log.Printf("Advertencia al ejecutar auto-alter de premios: %v", err)
	}

	// Inicializar Repositorios
	clienteRepo := clientesInfra.NewMySQL()
	edicionRepo := edicionesInfra.NewMySQL()
	quinielaRepo := quinielasInfra.NewMySQL()
	saldoRepo := saldosInfra.NewMySQL()
	aporteRepo := aportesInfra.NewMySQL()
	movimientoRepo := movimientosInfra.NewMySQL()
	premioRepo := premiosInfra.NewMySQL()

	// Ajuste de saldo y cobro de quinielas: registrar movimientos por el pago
	// y el saldo fiado de cada quiniela para que el saldo del cliente cambie.
	ajustadorSaldo := saldosApp.NewAjustarSaldo(saldoRepo)
	createMovimiento := movimientosApp.NewCreateMovimiento(movimientoRepo, ajustadorSaldo)
	deleteMovimiento := movimientosApp.NewDeleteMovimiento(movimientoRepo, ajustadorSaldo)
	cobradorQuiniela := movimientosApp.NewCobradorQuiniela(movimientoRepo, createMovimiento, deleteMovimiento)
	cobradorAporte := movimientosApp.NewCobradorAporte(movimientoRepo, createMovimiento, deleteMovimiento)

	// Registrar Rutas directamente en la misma instancia 'r'
	clientesInfra.SetupRouterClientes(r, clienteRepo, saldoRepo)
	edicionesInfra.SetupRouterEdiciones(r, edicionRepo)
	quinielasInfra.SetupRouterQuinielas(r, quinielaRepo, cobradorQuiniela)
	saldosInfra.SetupRouterSaldos(r, saldoRepo)
	aportesInfra.SetupRouterAportes(r, aporteRepo, quinielaRepo, cobradorAporte)
	movimientosInfra.SetupRouterMovimientos(r, movimientoRepo, saldoRepo)

	// Abonador de premios: registra cada parte del reparto como movimiento
	// PREMIO_ABONO (ajusta el saldo) y puede deshacerse si el reparto falla.
	abonadorPremio := movimientosApp.NewAbonadorPremio(createMovimiento, deleteMovimiento)
	premiosInfra.SetupRouterPremios(r, premioRepo, aporteRepo, quinielaRepo, abonadorPremio, clienteRepo)

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
