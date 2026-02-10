// internal/domain/entity/equipo.go
package entity

import "github.com/garfex/calculadora-filtros/internal/domain/valueobject"

// CalculadorCorriente defines the contract for equipment that can calculate
// its nominal current. Each equipment type implements this differently.
type CalculadorCorriente interface {
	CalcularCorrienteNominal() (valueobject.Corriente, error)
}

// Equipo is the base struct embedded by all equipment types.
type Equipo struct {
	Clave   string
	Tipo    TipoFiltro
	Voltaje int
	ITM     int
	Bornes  int
}
