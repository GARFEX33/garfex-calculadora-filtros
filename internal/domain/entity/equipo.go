// internal/domain/entity/equipo.go
package entity

import "github.com/garfex/calculadora-filtros/internal/domain/valueobject"

// CalculadorCorriente defines the contract for equipment that can calculate
// its nominal current. Each equipment type implements this differently.
type CalculadorCorriente interface {
	CalcularCorrienteNominal() (valueobject.Corriente, error)
}

// CalculadorPotencia defines the contract for equipment that can report its
// electrical power in the three standard components.
//
// FiltroActivo (PF=1): KVA = I×V×√3/1000, KW = KVA, KVAR = 0
// FiltroRechazo (purely reactive, PF=0): KVAR = given, KVA = KVAR, KW = 0
type CalculadorPotencia interface {
	PotenciaKVA() float64
	PotenciaKW() float64
	PotenciaKVAR() float64
}

// Equipo is the base struct embedded by all equipment types.
type Equipo struct {
	Clave   string
	Tipo    TipoFiltro
	Voltaje int
	ITM     int
	Bornes  int
}
