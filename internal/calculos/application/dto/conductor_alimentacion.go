// internal/calculos/application/dto/conductor_alimentacion.go
package dto

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// ConductorAlimentacionInput contiene los datos de entrada para seleccionar
// el conductor de alimentacion segun tablas NOM 310-15.
type ConductorAlimentacionInput struct {
	// CorrienteAjustada es el amperaje ya ajustado por factores NOM.
	CorrienteAjustada float64 `json:"corriente_ajustada"`
	// TipoCanalizacion es el tipo de canalizacion (TUBERIA_PVC, CHAROLA_ESPACIADO, etc).
	TipoCanalizacion string `json:"tipo_canalizacion"`
	// Material es el material del conductor ("Cu" o "Al").
	// Si esta vacio, se usa "Cu" por defecto.
	Material string `json:"material"`
	// Temperatura es la temperatura de operacion (60, 75, 90).
	// Si es nil, se aplica la regla NOM automatica.
	Temperatura *int `json:"temperatura"`
	// HilosPorFase es el numero de conductores por fase.
	// Si es 0 o menor, se usa 1 por defecto.
	HilosPorFase int `json:"hilos_por_fase"`
}

// Validate verifica que los campos requeridos sean validos.
func (i ConductorAlimentacionInput) Validate() error {
	if i.CorrienteAjustada <= 0 {
		return fmt.Errorf("%w: corriente_ajustada debe ser mayor que cero", ErrEquipoInputInvalido)
	}
	if i.TipoCanalizacion == "" {
		return fmt.Errorf("%w: tipo_canalizacion es requerido", ErrEquipoInputInvalido)
	}
	return nil
}

// ToDomainMaterial convierte el material del DTO a value object.
// Si esta vacio o es invalido, retorna cobre por defecto.
func (i ConductorAlimentacionInput) ToDomainMaterial() valueobject.MaterialConductor {
	if i.Material == "Al" {
		return valueobject.MaterialAluminio
	}
	return valueobject.MaterialCobre
}

// ConductorAlimentacionOutput contiene el resultado de seleccionar
// el conductor de alimentacion.
type ConductorAlimentacionOutput struct {
	// Calibre es el calibre del conductor seleccionado (ej: "4 AWG", "250 MCM").
	Calibre string `json:"calibre"`
	// Material es el material del conductor ("Cu" o "Al").
	Material string `json:"material"`
	// SeccionMM2 es la seccion transversal en mm2.
	SeccionMM2 float64 `json:"seccion_mm2"`
	// TipoAislamiento es el tipo de aislamiento (ej: "THHN/THHW").
	TipoAislamiento string `json:"tipo_aislamiento"`
	// CapacidadNominal es la ampacidad del conductor segun la tabla.
	CapacidadNominal float64 `json:"capacidad_nominal"`
	// TablaUsada es el nombre descriptivo de la tabla NOM utilizada.
	TablaUsada string `json:"tabla_usada"`
}
