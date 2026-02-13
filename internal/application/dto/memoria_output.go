// internal/application/dto/memoria_output.go
package dto

import (
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
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
	TipoEquipo     entity.TipoEquipo
	Clave          string
	Tension        valueobject.Tension
	FactorPotencia float64

	// Paso 1: Corriente Nominal
	CorrienteNominal valueobject.Corriente

	// Paso 2: Ajuste de Corriente
	FactorAgrupamiento float64
	FactorTemperatura  float64
	FactorTotalAjuste  float64
	CorrienteAjustada  valueobject.Corriente
	HilosPorFase       int
	CorrientePorHilo   float64

	// Paso 3: Tipo de Canalización
	TipoCanalizacion entity.TipoCanalizacion

	// Paso 4: Conductor de Alimentación
	TemperaturaUsada      valueobject.Temperatura
	ConductorAlimentacion ResultadoConductor
	TablaAmpacidadUsada   string

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
