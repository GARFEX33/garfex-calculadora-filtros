// internal/calculos/application/dto/tuberia_output.go
package dto

import (
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
)

// TuberiaOutput contiene el resultado del cálculo de tamaño de tubería.
// Es el DTO de salida para el use case CalcularTamanioTuberia.
type TuberiaOutput struct {
	AreaPorTuboMM2     float64 `json:"area_por_tubo_mm2"`
	TuberiaRecomendada string  `json:"tuberia_recomendada"`
	DesignacionMetrica string  `json:"designacion_metrica"`
	TipoCanalizacion   string  `json:"tipo_canalizacion"`
	NumTuberias        int     `json:"num_tuberias"`
}

// TuberiaOutputFromDomain convierte entity.ResultadoTamanioTuberia a TuberiaOutput.
func TuberiaOutputFromDomain(r entity.ResultadoTamanioTuberia) TuberiaOutput {
	return TuberiaOutput{
		AreaPorTuboMM2:     r.AreaPorTuboMM2(),
		TuberiaRecomendada: r.TuberiaRecomendada(),
		DesignacionMetrica: r.DesignacionMetrica(),
		TipoCanalizacion:   string(r.TipoCanalizacion()),
		NumTuberias:        r.NumTuberias(),
	}
}
