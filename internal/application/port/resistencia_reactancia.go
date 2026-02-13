// internal/application/port/resistencia_reactancia.go
package port

// ResistenciaReactancia holds the impedance values for voltage drop calculation.
type ResistenciaReactancia struct {
	R float64 // Ohms per km
	X float64 // Ohms per km
}
