// internal/domain/entity/memoria_calculo.go
package entity

import "github.com/garfex/calculadora-filtros/internal/domain/valueobject"

// ResultadoCaidaTension holds the complete voltage drop calculation results.
// R_ef, X_ef, and the IEEE-141 effective term are exposed for the memoria de cálculo report.
type ResultadoCaidaTension struct {
	Porcentaje  float64 // %VD = (VD / V) × 100
	CaidaVolts  float64 // VD in volts
	Cumple      bool    // %VD ≤ limiteNOM
	Impedancia  float64 // Término efectivo (Ω/km) = R·cosθ + X·senθ  (IEEE-141)
	Resistencia float64 // R_ef = ResistenciaOhmPorKm / HilosPorFase
	Reactancia  float64 // X_ef = ReactanciaOhmPorKm / HilosPorFase  (Tabla 9)
}

// MemoriaCalculo holds the complete results of all calculation steps
// for an electrical installation memory.
type MemoriaCalculo struct {
	Equipo                CalculadorCorriente
	CorrienteNominal      valueobject.Corriente
	CorrienteAjustada     valueobject.Corriente
	FactoresAjuste        map[string]float64
	PotenciaKVA           float64
	PotenciaKW            float64
	PotenciaKVAR          float64
	ConductorAlimentacion valueobject.Conductor
	HilosPorFase          int
	ConductorTierra       valueobject.Conductor
	Canalizacion          Canalizacion
	TemperaturaUsada      int                   // Actual NOM temperature column used (60, 75, or 90°C)
	CaidaTension          ResultadoCaidaTension // Complete voltage drop result including R_ef, X_ef, effective term
	CumpleNormativa       bool
}
