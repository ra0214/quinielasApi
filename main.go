package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"quinielas/src/config"
	"quinielas/src/config/middleware"

	aportesInfra "quinielas/src/aportes/infraestructure"
	clientesInfra "quinielas/src/clientes/infraestructure"
	edicionesInfra "quinielas/src/ediciones/infraestructure"
	mantenimientoApp "quinielas/src/mantenimiento/application"
	mantenimientoInfra "quinielas/src/mantenimiento/infraestructure"
	movimientosApp "quinielas/src/movimientos/application"
	movimientosInfra "quinielas/src/movimientos/infraestructure"
	premiosInfra "quinielas/src/premios/infraestructure"
	quinielasApp "quinielas/src/quinielas/application"
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
	if err := dbConn.AutoAlterEdiciones(); err != nil {
		log.Printf("Advertencia al ejecutar auto-alter de ediciones: %v", err)
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
	aportesAdapter := quinielasInfra.NewAportesAdapter(aporteRepo)
	clientesInfra.SetupRouterClientes(r, clienteRepo, saldoRepo)
	edicionesInfra.SetupRouterEdiciones(r, edicionRepo)
	quinielasInfra.SetupRouterQuinielas(r, quinielaRepo, cobradorQuiniela, aportesAdapter, premioRepo)
	saldosInfra.SetupRouterSaldos(r, saldoRepo)
	aportesInfra.SetupRouterAportes(r, aporteRepo, quinielaRepo, cobradorAporte)
	movimientosInfra.SetupRouterMovimientos(r, movimientoRepo, saldoRepo)

	// Abonador de premios: registra cada parte del reparto como movimiento
	// PREMIO_ABONO (ajusta el saldo) y puede deshacerse si el reparto falla.
	abonadorPremio := movimientosApp.NewAbonadorPremio(createMovimiento, deleteMovimiento)
	premiosInfra.SetupRouterPremios(r, premioRepo, aporteRepo, quinielaRepo, abonadorPremio, clienteRepo)

	// Job en segundo plano: cierra las quinielas DESARROLLO vencidas aplicando
	// la consecuencia de su estado_fallo.
	vencerQuinielas := quinielasApp.NewVencerQuinielas(quinielaRepo, cobradorQuiniela, aportesAdapter)
	go iniciarVencimientos(vencerQuinielas)

	// Job en segundo plano: elimina los concursos cerrados hace mas de N meses
	// (retencion configurable con RETENCION_MESES; por defecto 6). No toca los saldos.
	purgarConcursos := mantenimientoApp.NewPurgarConcursos(mantenimientoInfra.NewMySQL(), mesesRetencion())
	go iniciarPurga(purgarConcursos)

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

// iniciarVencimientos revisa periódicamente las quinielas DESARROLLO vencidas.
// El intervalo se configura con VENCIMIENTOS_MIN (minutos); por defecto 1.
func iniciarVencimientos(v *quinielasApp.VencerQuinielas) {
	intervalo := time.Minute
	if valor := os.Getenv("VENCIMIENTOS_MIN"); valor != "" {
		if n, err := strconv.Atoi(valor); err == nil && n > 0 {
			intervalo = time.Duration(n) * time.Minute
		}
	}
	log.Printf("⏰ Revisión de quinielas vencidas cada %s", intervalo)

	procesar := func() {
		if err := v.Procesar(time.Now()); err != nil {
			log.Printf("Error al procesar quinielas vencidas: %v", err)
		}
	}

	// Primera pasada al arrancar, sin esperar al primer tick.
	procesar()

	ticker := time.NewTicker(intervalo)
	defer ticker.Stop()
	for range ticker.C {
		procesar()
	}
}

// iniciarPurga elimina periódicamente los concursos cerrados hace más de N
// meses. El intervalo se configura con PURGA_MIN (minutos); por defecto 1440
// (una vez al día).
func iniciarPurga(p *mantenimientoApp.PurgarConcursos) {
	intervalo := 24 * time.Hour
	if valor := os.Getenv("PURGA_MIN"); valor != "" {
		if n, err := strconv.Atoi(valor); err == nil && n > 0 {
			intervalo = time.Duration(n) * time.Minute
		}
	}
	log.Printf("🧹 Purga de concursos antiguos cada %s", intervalo)

	procesar := func() {
		if err := p.Procesar(time.Now()); err != nil {
			log.Printf("Error al purgar concursos antiguos: %v", err)
		}
	}

	// Primera pasada al arrancar, sin esperar al primer tick.
	procesar()

	ticker := time.NewTicker(intervalo)
	defer ticker.Stop()
	for range ticker.C {
		procesar()
	}
}

// mesesRetencion indica cuántos meses se conservan los concursos cerrados.
// Se configura con RETENCION_MESES; por defecto 6.
func mesesRetencion() int {
	if valor := os.Getenv("RETENCION_MESES"); valor != "" {
		if n, err := strconv.Atoi(valor); err == nil && n > 0 {
			return n
		}
	}
	return 6
}
