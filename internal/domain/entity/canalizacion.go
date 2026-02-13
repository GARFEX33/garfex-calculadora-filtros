// internal/domain/entity/canalizacion.go
package entity

// Canalizacion represents the conduit or cable tray selected for the installation.
type Canalizacion struct {
	Tipo           string  // "TUBERIA" | "CHAROLA"
	Tamano         string  // e.g., "1 1/2" (inches for tubería)
	AnchoRequerido float64 // for charola: required width in mm; for tubería: total conductor area in mm²
	NumeroDeTubos  int     // number of parallel conduits; 1 = single conduit installation
}
