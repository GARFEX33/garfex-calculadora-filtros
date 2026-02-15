// internal/calculos/application/dto/errors.go
package dto

import (
	"errors"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
)

// Errores de validación de input.
var (
	ErrEquipoInputInvalido = errors.New("datos de equipo inválidos")
	ErrModoInvalido        = errors.New("modo de cálculo inválido")
	ErrEquipoNoEncontrado  = errors.New("equipo no encontrado")
)

// Re-exportar errores de domain/service para que presentation no importe domain directamente.
// Esto mantiene la arquitectura hexagonal: presentation -> application -> domain.
var (
	ErrConductorNoEncontrado    = service.ErrConductorNoEncontrado
	ErrCanalizacionNoDisponible = service.ErrCanalizacionNoDisponible
	ErrDistanciaInvalida        = service.ErrDistanciaInvalida
	ErrHilosPorFaseInvalido     = service.ErrHilosPorFaseInvalido
	ErrFactorPotenciaInvalido   = service.ErrFactorPotenciaInvalido
)
