// internal/application/dto/equipo_input.go
package dto

import (
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// ModoCalculo define cómo se proporcionan los datos del equipo.
type ModoCalculo string

const (
	ModoListado        ModoCalculo = "LISTADO"
	ModoManualAmperaje ModoCalculo = "MANUAL_AMPERAJE"
	ModoManualPotencia ModoCalculo = "MANUAL_POTENCIA"
)

// EquipoInput contiene todos los datos necesarios para calcular una memoria.
// Es el DTO de entrada para el use case CalcularMemoria.
type EquipoInput struct {
	// Modo indica cómo se proporcionan los datos del equipo
	Modo ModoCalculo

	// Clave del equipo (requerido si Modo = LISTADO)
	Clave string

	// Datos del equipo (requeridos si Modo = MANUAL_*)
	TipoEquipo      entity.TipoEquipo
	AmperajeNominal float64 // Solo para MANUAL_AMPERAJE
	PotenciaNominal float64 // Solo para MANUAL_POTENCIA (KVAR o KVA)
	Tension         valueobject.Tension
	FactorPotencia  float64 // Solo para CARGA en MANUAL_POTENCIA
	ITM             int

	// Parámetros de instalación
	TipoCanalizacion      entity.TipoCanalizacion
	TemperaturaOverride   *valueobject.Temperatura      // nil = usar lógica por defecto
	HilosPorFase          int                           // default: 1
	Material              valueobject.MaterialConductor // "Cu" o "Al"; si vacío, default Cu
	LongitudCircuito      float64                       // metros, para caída de tensión
	PorcentajeCaidaMaximo float64                       // default: 3.0%

	// NUEVO: Reemplaza factor_agrupamiento y factor_temperatura
	Estado           string                  `json:"estado" binding:"required"`
	SistemaElectrico entity.SistemaElectrico `json:"sistema_electrico" binding:"required"`
}

// Validate verifica que el input tenga los campos requeridos según el modo.
// Esta validación es básica; el use case hace validaciones más específicas.
func (e EquipoInput) Validate() error {
	switch e.Modo {
	case ModoListado:
		if e.Clave == "" {
			return ErrEquipoInputInvalido
		}
	case ModoManualAmperaje:
		if e.AmperajeNominal <= 0 {
			return ErrEquipoInputInvalido
		}
	case ModoManualPotencia:
		if e.PotenciaNominal <= 0 {
			return ErrEquipoInputInvalido
		}
	default:
		return ErrModoInvalido
	}

	if e.Tension.Valor() <= 0 {
		return ErrEquipoInputInvalido
	}

	if e.ITM <= 0 {
		return ErrEquipoInputInvalido
	}

	if e.Estado == "" {
		return ErrEquipoInputInvalido
	}

	if err := entity.ValidarSistemaElectrico(e.SistemaElectrico); err != nil {
		return err
	}

	return nil
}
