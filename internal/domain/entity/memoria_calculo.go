// internal/domain/entity/memoria_calculo.go
package entity

import "github.com/garfex/calculadora-filtros/internal/domain/valueobject"

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
	CaidaTension          float64
	CumpleNormativa       bool
}
