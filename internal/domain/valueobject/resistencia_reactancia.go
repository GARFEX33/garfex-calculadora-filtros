// internal/domain/valueobject/resistencia_reactancia.go
package valueobject

// ResistenciaReactancia holds the impedance values for voltage drop calculation.
type ResistenciaReactancia struct {
	R float64 // Ohms per km
	X float64 // Ohms per km
}
