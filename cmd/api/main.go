// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	calculosport "github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/infrastructure"
	"github.com/garfex/calculadora-filtros/internal/calculos/infrastructure/adapter/driven/csv"
	calculospostgres "github.com/garfex/calculadora-filtros/internal/calculos/infrastructure/adapter/driven/postgres"
	calcmock "github.com/garfex/calculadora-filtros/internal/calculos/infrastructure/adapter/driven/mock"

	equiposusecase "github.com/garfex/calculadora-filtros/internal/equipos/application/usecase"
	equiposport "github.com/garfex/calculadora-filtros/internal/equipos/application/port"
	equiposinfra "github.com/garfex/calculadora-filtros/internal/equipos/infrastructure"
	equipospostgres "github.com/garfex/calculadora-filtros/internal/equipos/infrastructure/adapter/driven/postgres"
	equipohttp "github.com/garfex/calculadora-filtros/internal/equipos/infrastructure/adapter/driver/http"
	mockequipos "github.com/garfex/calculadora-filtros/internal/equipos/infrastructure/adapter/driven/mock"
	sharedpostgres "github.com/garfex/calculadora-filtros/internal/shared/infrastructure/postgres"
)

func main() {
	// Cargar variables de entorno desde .env (solo en desarrollo, ignora error si no existe)
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// ─── Configuración de modo mock ─────────────────────────────────────────────
	mockMode := os.Getenv("MOCK_MODE") == "true"
	environment := os.Getenv("ENVIRONMENT")

	// Validación de seguridad: MOCK_MODE no permitido en producción
	if mockMode && environment == "production" {
		log.Fatal("❌ MOCK_MODE no está permitido en entorno de producción")
	}

	// ─── Tablas NOM (CSV) ─────────────────────────────────────────────────────

	tablaRepo, err := csv.NewCSVTablaNOMRepository("data/tablas_nom")
	if err != nil {
		log.Fatalf("Error cargando tablas NOM: %v", err)
	}

	// ─── Repositorios: PostgreSQL o Mock según MOCK_MODE ─────────────────────
	// En MOCK_MODE=true no se intenta conectar a PostgreSQL. Ambos repositorios
	// (equipos y cálculos) usan implementaciones en memoria.

	var calcEquipoRepo calculosport.EquipoRepository
	var equipoFiltroRepo equiposport.EquipoFiltroRepository

	if mockMode {
		log.Println("⚠️  MOCK_MODE activo — usando repositorios en memoria (sin PostgreSQL)")
		calcEquipoRepo = calcmock.NewCalcEquipoMockRepository()
		equipoFiltroRepo = mockequipos.NewMockEquipoFiltroRepository()
	} else {
		dbCfg, err := sharedpostgres.LoadDBConfigFromEnv()
		if err != nil {
			log.Fatalf("Error cargando configuración de base de datos: %v", err)
		}

		pool, err := equipospostgres.NewPool(dbCfg)
		if err != nil {
			log.Fatalf("Error conectando a PostgreSQL: %v", err)
		}
		defer pool.Close()
		log.Printf("✅ Conectado a PostgreSQL en %s:%s", dbCfg.Host, dbCfg.Port)

		calcEquipoRepo = calculospostgres.NewCalcEquipoFiltroRepository(pool)
		equipoFiltroRepo = equipospostgres.NewPostgresEquipoFiltroRepository(pool)
	}

	// ─── Calculos: use cases ──────────────────────────────────────────────────

	calcularCorrienteUC := usecase.NewCalcularCorrienteUseCase(calcEquipoRepo)
	ajustarCorrienteUC := usecase.NewAjustarCorrienteUseCase(tablaRepo)
	seleccionarConductorUC := usecase.NewSeleccionarConductorUseCase(tablaRepo)
	seleccionarConductorAlimentacionUC := usecase.NewSeleccionarConductorAlimentacionUseCase(tablaRepo)
	seleccionarConductorTierraUC := usecase.NewSeleccionarConductorTierraUseCase(tablaRepo)
	calcularTamanioTuberiaUC := usecase.NewCalcularTamanioTuberiaUseCase(tablaRepo)
	calcularCharolaEspaciadoUC := usecase.NewCalcularCharolaEspaciadoUseCase(tablaRepo)
	calcularCharolaTriangularUC := usecase.NewCalcularCharolaTriangularUseCase(tablaRepo)
	calcularCaidaTensionUC := usecase.NewCalcularCaidaTensionUseCase(tablaRepo)
	seleccionarConductorCaidaTensionUC := usecase.NewSeleccionarConductorPorCaidaTensionUseCase(calcularCaidaTensionUC, tablaRepo)

	orquestadorMemoriaUC := usecase.NewOrquestadorMemoriaCalculoUseCase(
		calcularCorrienteUC,
		ajustarCorrienteUC,
		seleccionarConductorUC,
		calcularTamanioTuberiaUC,
		calcularCharolaEspaciadoUC,
		calcularCharolaTriangularUC,
		calcularCaidaTensionUC,
		seleccionarConductorCaidaTensionUC,
		tablaRepo,
	)

	// ─── Equipos: use cases ───────────────────────────────────────────────────

	crearEquipoUC := equiposusecase.NewCrearEquipoUseCase(equipoFiltroRepo)
	obtenerEquipoUC := equiposusecase.NewObtenerEquipoUseCase(equipoFiltroRepo)
	listarEquiposUC := equiposusecase.NewListarEquiposUseCase(equipoFiltroRepo)
	actualizarEquipoUC := equiposusecase.NewActualizarEquipoUseCase(equipoFiltroRepo)
	eliminarEquipoUC := equiposusecase.NewEliminarEquipoUseCase(equipoFiltroRepo)

	equipoHandler := equipohttp.NewEquipoHandler(
		crearEquipoUC,
		obtenerEquipoUC,
		listarEquiposUC,
		actualizarEquipoUC,
		eliminarEquipoUC,
	)

	// ─── Router principal ────────────────────────────────────────────────────

	router := infrastructure.NewRouter(
		calcularCorrienteUC,
		ajustarCorrienteUC,
		seleccionarConductorUC,
		seleccionarConductorAlimentacionUC,
		seleccionarConductorTierraUC,
		calcularTamanioTuberiaUC,
		calcularCharolaEspaciadoUC,
		calcularCharolaTriangularUC,
		calcularCaidaTensionUC,
		orquestadorMemoriaUC,
	)

	// Montar rutas de equipos bajo /api/v1
	v1 := router.Group("/api/v1")
	equiposinfra.RegisterEquiposRoutes(v1, equipoHandler)

	// ─── Servidor HTTP ───────────────────────────────────────────────────────

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Servidor iniciado en puerto %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando servidor: %v", err)
		}
	}()

	<-quit
	log.Println("Cerrando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error forzando cierre del servidor: %v", err)
	}

	log.Println("Servidor cerrado correctamente")
}
