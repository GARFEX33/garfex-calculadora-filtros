// internal/calculos/application/dto/conductor_tierra.go
package dto

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// ConductorTierraInput contiene los datos de entrada para seleccionar
// el conductor de tierra según NOM-250-122.
type ConductorTierraInput struct {
	// ITM es el Interruptor Termomagnético en amperes.
	ITM int `json:"itm" binding:"required"`
	// Material es el material del conductor ("Cu" o "Al").
	// Si está vacío, se usa "Cu" por defecto.
	Material string `json:"material"`
}

// Validate verifica que los campos requeridos sean válidos.
func (i ConductorTierraInput) Validate() error {
	if i.ITM <= 0 {
		return fmt.Errorf("%w: ITM debe ser mayor que cero", ErrEquipoInputInvalido)
	}
	return nil
}

// ToDomainMaterial convierte el material del DTO a value object.
// Si está vacío o es inválido, retorna cobre por defecto.
func (i ConductorTierraInput) ToDomainMaterial() valueobject.MaterialConductor {
	switch i.Material {
	case "Al":
		return valueobject.MaterialAluminio
	default:
		return valueobject.MaterialCobre
	}
}

// ConductorTierraOutput contiene el resultado de seleccionar
// el conductor de tierra según NOM-250-122.
type ConductorTierraOutput struct {
	// Calibre es el calibre del conductor seleccionado (ej: "12 AWG", "2 AWG").
	Calibre string `json:"calibre"`
	// Material es el material del conductor ("Cu" o "Al").
	Material string `json:"material"`
	// SeccionMM2 es la sección transversal en mm².
	SeccionMM2 float64 `json:"seccion_mm2"`
	// ITMHasta es el ITM máximo para el cual aplica este calibre.
	ITMHasta int `json:"itm_hasta"`
}
