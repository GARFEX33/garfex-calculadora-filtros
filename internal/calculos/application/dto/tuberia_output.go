// internal/calculos/application/dto/tuberia_output.go
package dto

// TuberiaOutput contiene el resultado del cálculo de tamaño de tubería.
// Es el DTO de salida para el use case CalcularTamanioTuberia.
type TuberiaOutput struct {
	AreaPorTuboMM2     float64 `json:"area_por_tubo_mm2"`
	TuberiaRecomendada string  `json:"tuberia_recomendada"`
	DesignacionMetrica string  `json:"designacion_metrica"`
	TipoCanalizacion   string  `json:"tipo_canalizacion"`
	NumTuberias        int     `json:"num_tuberias"`

	// Nuevos campos para el desarrollo detallado
	AreaFaseMM2          float64  `json:"area_fase_mm2"`
	AreaNeutroMM2        *float64 `json:"area_neutro_mm2,omitempty"` // nil si no hay neutro (DELTA)
	AreaTierraMM2        float64  `json:"area_tierra_mm2"`
	NumFasesPorTubo      int      `json:"num_fases_por_tubo"`
	NumNeutrosPorTubo    int      `json:"num_neutros_por_tubo"`
	NumTierras           int      `json:"num_tierras"`
	AreaOcupacionTuboMM2 float64  `json:"area_ocupacion_tubo_mm2"`
	FillFactor           float64  `json:"fill_factor"`
}
