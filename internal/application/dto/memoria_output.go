// internal/application/dto/memoria_output.go
package dto

import (
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
)

// PasoMemoria representa un paso individual del cálculo.
type PasoMemoria struct {
	Numero      int
	Nombre      string
	Descripcion string
	Resultado   interface{}
}

// ResultadoConductor contiene la información del conductor seleccionado.
type ResultadoConductor struct {
	Calibre         string
	Material        string
	SeccionMM2      float64
	TipoAislamiento string
	Capacidad       float64 // Ampacidad de la tabla
}

// ResultadoCanalizacion contiene la información de la canalización.
type ResultadoCanalizacion struct {
	Tipo             entity.TipoCanalizacion
	Tamano           string
	AreaTotalMM2     float64
	AreaRequeridaMM2 float64
	NumeroDeTubos    int
}

// ResultadoCaidaTension contiene el resultado del cálculo de caída.
type ResultadoCaidaTension struct {
	Porcentaje          float64
	CaidaVolts          float64
	Cumple              bool
	LimitePorcentaje    float64
	ResistenciaEfectiva float64 // R·cosθ + X·senθ
}

// MemoriaOutput contiene el resultado completo de todos los pasos.
// Es el DTO de salida para el use case CalcularMemoria.
type MemoriaOutput struct {
	// Información del equipo
	TipoEquipo     entity.TipoEquipo `json:"tipo_equipo"`
	Clave          string            `json:"clave"`
	Tension        int               `json:"tension"`
	FactorPotencia float64           `json:"factor_potencia"`

	// NUEVOS: Información de cálculo
	Estado              string                  `json:"estado"`
	TemperaturaAmbiente int                     `json:"temperatura_ambiente"`
	SistemaElectrico    entity.SistemaElectrico `json:"sistema_electrico"`
	CantidadConductores int                     `json:"cantidad_conductores"`

	// Factores calculados
	FactorTemperaturaCalculado  float64 `json:"factor_temperatura_calculado"`
	FactorAgrupamientoCalculado float64 `json:"factor_agrupamiento_calculado"`

	// Paso 1: Corriente Nominal
	CorrienteNominal float64 `json:"corriente_nominal"`

	// Paso 2: Ajuste de Corriente
	FactorAgrupamiento float64 `json:"factor_agrupamiento"`
	FactorTemperatura  float64 `json:"factor_temperatura"`
	FactorTotalAjuste  float64 `json:"factor_total_ajuste"`
	CorrienteAjustada  float64 `json:"corriente_ajustada"`
	HilosPorFase       int     `json:"hilos_por_fase"`
	CorrientePorHilo   float64 `json:"corriente_por_hilo"`

	// Paso 3: Tipo de Canalización
	TipoCanalizacion entity.TipoCanalizacion `json:"tipo_canalizacion"`

	// Material del conductor
	Material string `json:"material"` // "Cu" o "Al"

	// Paso 4: Conductor de Alimentación
	TemperaturaUsada      int                `json:"temperatura_usada"`
	ConductorAlimentacion ResultadoConductor `json:"conductor_alimentacion"`
	TablaAmpacidadUsada   string             `json:"tabla_ampacidad_usada"`

	// Paso 5: Conductor de Tierra
	ConductorTierra ResultadoConductor
	ITM             int

	// Paso 6: Canalización
	Canalizacion ResultadoCanalizacion
	FillFactor   float64 // 40% para tubería

	// Paso 7: Caída de Tensión
	LongitudCircuito float64
	CaidaTension     ResultadoCaidaTension

	// Resumen de cumplimiento
	CumpleNormativa bool
	Observaciones   []string

	// Todos los pasos para el reporte
	Pasos []PasoMemoria
}
