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

	"github.com/garfex/calculadora-filtros/internal/application/port"
	"github.com/garfex/calculadora-filtros/internal/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/infrastructure/repository"
	"github.com/garfex/calculadora-filtros/internal/presentation"
)

func main() {
	// Crear repositorios
	tablaRepo, err := repository.NewCSVTablaNOMRepository("data/tablas_nom")
	if err != nil {
		log.Fatalf("Error cargando tablas NOM: %v", err)
	}

	// TODO: Implementar PostgreSQL repository
	var equipoRepo port.EquipoRepository

	// Crear use case
	calcularMemoriaUC := usecase.NewCalcularMemoriaUseCase(tablaRepo, equipoRepo)

	// Crear router
	router := presentation.NewRouter(calcularMemoriaUC)

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
