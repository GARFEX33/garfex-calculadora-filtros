// internal/calculos/domain/entity/tamanio_tuberia.go
package entity

import "fmt"

// ResultadoTamanioTuberia holds the result of conduit sizing calculation.
// It represents the recommended conduit size based on conductor area occupation.
type ResultadoTamanioTuberia struct {
	areaPorTuboMM2     float64          // Total conductor area per tube in mm²
	tuberiaRecomendada string           // Trade size (e.g., "1", "1 1/2", "2")
	designacionMetrica string           // Metric designation (e.g., "27mm", "35mm", "41mm")
	tipoCanalizacion   TipoCanalizacion // Type of canalization (PVC, ACERO_PG, ACERO_PD)
	numTuberias        int              // Number of parallel conduits used
}

// NewResultadoTamanioTuberia constructs a validated ResultadoTamanioTuberia.
func NewResultadoTamanioTuberia(
	areaPorTuboMM2 float64,
	tuberiaRecomendada string,
	designacionMetrica string,
	tipoCanalizacion TipoCanalizacion,
	numTuberias int,
) (ResultadoTamanioTuberia, error) {
	if areaPorTuboMM2 <= 0 {
		return ResultadoTamanioTuberia{}, fmt.Errorf("areaPorTuboMM2 debe ser mayor que cero")
	}
	if tuberiaRecomendada == "" {
		return ResultadoTamanioTuberia{}, fmt.Errorf("tuberiaRecomendada no puede estar vacía")
	}
	if err := ValidarTipoCanalizacion(tipoCanalizacion); err != nil {
		return ResultadoTamanioTuberia{}, fmt.Errorf("NewResultadoTamanioTuberia: %w", err)
	}
	if numTuberias < 1 {
		return ResultadoTamanioTuberia{}, fmt.Errorf("numTuberias debe ser mayor o igual a 1: %d", numTuberias)
	}
	return ResultadoTamanioTuberia{
		areaPorTuboMM2:     areaPorTuboMM2,
		tuberiaRecomendada: tuberiaRecomendada,
		designacionMetrica: designacionMetrica,
		tipoCanalizacion:   tipoCanalizacion,
		numTuberias:        numTuberias,
	}, nil
}

// AreaPorTuboMM2 returns the total conductor area per tube.
func (r ResultadoTamanioTuberia) AreaPorTuboMM2() float64 {
	return r.areaPorTuboMM2
}

// TuberiaRecomendada returns the trade size designation.
func (r ResultadoTamanioTuberia) TuberiaRecomendada() string {
	return r.tuberiaRecomendada
}

// DesignacionMetrica returns the metric designation.
func (r ResultadoTamanioTuberia) DesignacionMetrica() string {
	return r.designacionMetrica
}

// TipoCanalizacion returns the canalization type used.
func (r ResultadoTamanioTuberia) TipoCanalizacion() TipoCanalizacion {
	return r.tipoCanalizacion
}

// NumTuberias returns the number of parallel conduits.
func (r ResultadoTamanioTuberia) NumTuberias() int {
	return r.numTuberias
}
