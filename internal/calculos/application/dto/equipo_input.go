// internal/calculos/application/dto/equipo_input.go
package dto

import (
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// ModoCalculo define cómo se proporcionan los datos del equipo.
type ModoCalculo string

const (
	ModoListado        ModoCalculo = "LISTADO"
	ModoManualAmperaje ModoCalculo = "MANUAL_AMPERAJE"
	ModoManualPotencia ModoCalculo = "MANUAL_POTENCIA"
)

// SistemaElectrico representa el tipo de sistema eléctrico (DTO - sin métodos de negocio).
type SistemaElectrico string

const (
	SistemaElectricoDelta      SistemaElectrico = "DELTA"
	SistemaElectricoEstrella   SistemaElectrico = "ESTRELLA"
	SistemaElectricoBifasico   SistemaElectrico = "BIFASICO"
	SistemaElectricoMonofasico SistemaElectrico = "MONOFASICO"
)

// ToEntity convierte el DTO a la entidad del domain.
func (s SistemaElectrico) ToEntity() entity.SistemaElectrico {
	return entity.SistemaElectrico(s)
}

// CantidadConductores returns el número de conductores según el sistema.
func (s SistemaElectrico) CantidadConductores() int {
	switch s {
	case SistemaElectricoDelta, SistemaElectricoBifasico:
		return 3
	case SistemaElectricoEstrella:
		return 4
	case SistemaElectricoMonofasico:
		return 2
	default:
		return 3
	}
}

// EquipoInput contiene todos los datos necesarios para calcular una memoria.
// Es el DTO de entrada para el use case CalcularMemoria.
type EquipoInput struct {
	// Modo indica cómo se proporcionan los datos del equipo
	Modo ModoCalculo

	// Clave del equipo (requerido si Modo = LISTADO)
	Clave string

	// Datos del equipo (requeridos si Modo = MANUAL_*)
	TipoEquipo      string  // Solo para MANUAL_AMPERAJE / MANUAL_POTENCIA
	AmperajeNominal float64 // Solo para MANUAL_AMPERAJE
	PotenciaNominal float64 // Solo para MANUAL_POTENCIA (KVAR o KVA)
	Tension         valueobject.Tension
	FactorPotencia  float64 // Solo para CARGA en MANUAL_POTENCIA
	ITM             int

	// Parámetros de instalación
	TipoCanalizacion      string                        // "TUBERIA_PVC", "CHAROLA_CABLE_ESPACIADO", etc.
	TemperaturaOverride   *valueobject.Temperatura      // nil = usar lógica por defecto
	HilosPorFase          int                           // default: 1
	Material              valueobject.MaterialConductor `json:"material"` // "Cu" o "Al"; si vacío, default Cu
	LongitudCircuito      float64                       // metros, para caída de tensión
	PorcentajeCaidaMaximo float64                       // default: 3.0%

	// Sistema eléctrico (DTO con tipos primitivos, no entity)
	SistemaElectrico SistemaElectrico `json:"sistema_electrico" binding:"required"`
	Estado           string           `json:"estado" binding:"required"`
}

// Validate verifica que el input tenga los campos requeridos según el modo.
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

	// Validar sistema eléctrico
	entitySistema := e.SistemaElectrico.ToEntity()
	if err := entity.ValidarSistemaElectrico(entitySistema); err != nil {
		return err
	}

	return nil
}

// ToEntityTipoEquipo convierte el DTO string a entity.TipoEquipo.
func (e EquipoInput) ToEntityTipoEquipo() entity.TipoEquipo {
	return entity.TipoEquipo(e.TipoEquipo)
}

// ToEntityTipoCanalizacion convierte el DTO string a entity.TipoCanalizacion.
func (e EquipoInput) ToEntityTipoCanalizacion() entity.TipoCanalizacion {
	return entity.TipoCanalizacion(e.TipoCanalizacion)
}
