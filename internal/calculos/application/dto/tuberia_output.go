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
}
