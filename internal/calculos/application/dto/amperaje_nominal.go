// internal/calculos/application/dto/amperaje_nominal.go
package dto

import (
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// TipoCargaDTO representa el tipo de carga eléctrica (DTO).
type TipoCargaDTO string

const (
	TipoCargaDTOGenerica   TipoCargaDTO = "GENERICA"
	TipoCargaDTOMonofasica TipoCargaDTO = "MONOFASICA"
	TipoCargaDTOTrifasica  TipoCargaDTO = "TRIFASICA"
)

// ToEntity convierte el DTO al tipo del dominio.
func (t TipoCargaDTO) ToEntity() service.TipoCarga {
	switch t {
	case TipoCargaDTOMonofasica:
		return service.TipoCargaMonofasica
	case TipoCargaDTOTrifasica:
		return service.TipoCargaTrifasica
	default:
		// Por defecto, trifásico (caso más común en instalaciones industriales)
		return service.TipoCargaTrifasica
	}
}

// SistemaElectricoDTO representa el tipo de sistema eléctrico (DTO).
type SistemaElectricoDTO string

const (
	SistemaElectricoDTOEstrella SistemaElectricoDTO = "ESTRELLA"
	SistemaElectricoDTODelta    SistemaElectricoDTO = "DELTA"
)

// ToEntity convierte el DTO al tipo del dominio.
func (s SistemaElectricoDTO) ToEntity() service.SistemaElectrico {
	switch s {
	case SistemaElectricoDTOEstrella:
		return service.SistemaElectricoEstrella
	case SistemaElectricoDTODelta:
		return service.SistemaElectricoDelta
	default:
		// Por defecto, estrella
		return service.SistemaElectricoEstrella
	}
}

// AmperajeNominalInput es el DTO de entrada para calcular amperaje nominal desde potencia.
type AmperajeNominalInput struct {
	// PotenciaWatts es la potencia activa en Watts (requerido, > 0)
	PotenciaWatts float64 `json:"potencia_watts" binding:"required,gt=0"`

	// Tension es la tensión del circuito en volts (requerido)
	// Valores válidos: 127, 220, 240, 277, 440, 480, 600
	Tension int `json:"tension" binding:"required"`

	// TipoCarga indica el tipo de carga eléctrica (requerido)
	// Valores: "MONOFASICA" | "TRIFASICA"
	TipoCarga TipoCargaDTO `json:"tipo_carga" binding:"required"`

	// SistemaElectrico indica el tipo de sistema eléctrico (requerido)
	// Valores: "ESTRELLA" | "DELTA"
	SistemaElectrico SistemaElectricoDTO `json:"sistema_electrico" binding:"required"`

	// FactorPotencia es el factor de potencia (requerido, > 0 y <= 1)
	FactorPotencia float64 `json:"factor_potencia" binding:"required,gt=0,lte=1"`
}

// Validate verifica que el input tenga valores válidos.
func (i AmperajeNominalInput) Validate() error {
	// Validar tensión con value object
	_, err := valueobject.NewTension(i.Tension)
	if err != nil {
		return err
	}

	// Validar tipo de carga
	if i.TipoCarga == "" {
		return ErrEquipoInputInvalido
	}

	// Validar sistema eléctrico
	if i.SistemaElectrico == "" {
		return ErrEquipoInputInvalido
	}

	return nil
}

// AmperajeNominalOutput es el DTO de salida del cálculo de amperaje nominal.
type AmperajeNominalOutput struct {
	// Amperaje es el valor de corriente calculada en Amperes
	Amperaje float64 `json:"amperaje"`

	// Unidad es la unidad de medida ("A")
	Unidad string `json:"unidad"`
}
