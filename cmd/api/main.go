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

	calculosport "github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/calculos/infrastructure"
	"github.com/garfex/calculadora-filtros/internal/calculos/infrastructure/adapter/driven/csv"
)

func main() {
	// Crear repositorios
	tablaRepo, err := csv.NewCSVTablaNOMRepository("data/tablas_nom")
	if err != nil {
		log.Fatalf("Error cargando tablas NOM: %v", err)
	}

	// TODO: Implementar PostgreSQL repository
	var equipoRepo calculosport.EquipoRepository

	// Crear micro use cases
	calcularCorrienteUC := usecase.NewCalcularCorrienteUseCase(equipoRepo)
	ajustarCorrienteUC := usecase.NewAjustarCorrienteUseCase(tablaRepo)
	seleccionarConductorAlimentacionUC := usecase.NewSeleccionarConductorAlimentacionUseCase(tablaRepo)
	seleccionarConductorTierraUC := usecase.NewSeleccionarConductorTierraUseCase(tablaRepo)
	calcularTamanioTuberiaUC := usecase.NewCalcularTamanioTuberiaUseCase(tablaRepo)
	calcularCharolaEspaciadoUC := usecase.NewCalcularCharolaEspaciadoUseCase(tablaRepo)
	calcularCharolaTriangularUC := usecase.NewCalcularCharolaTriangularUseCase(tablaRepo)

	// Crear router
	router := infrastructure.NewRouter(
		calcularCorrienteUC,
		ajustarCorrienteUC,
		seleccionarConductorAlimentacionUC,
		seleccionarConductorTierraUC,
		calcularTamanioTuberiaUC,
		calcularCharolaEspaciadoUC,
		calcularCharolaTriangularUC,
	)

	// Configurar servidor HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Canal para señales de sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("Servidor iniciado en puerto %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando servidor: %v", err)
		}
	}()

	// Esperar señal de cierre
	<-quit
	log.Println("Cerrando servidor...")

	// Graceful shutdown con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error forzando cierre del servidor: %v", err)
	}

	log.Println("Servidor cerrado correctamente")
}
