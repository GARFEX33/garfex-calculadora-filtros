// internal/domain/entity/canalizacion.go
package entity

// Canalizacion represents the conduit or cable tray selected for the installation.
type Canalizacion struct {
	Tipo      string  // "TUBERIA" | "CHAROLA"
	Tamano    string  // e.g., "1 1/2" (inches for tubería)
	AreaTotal float64 // total conductor area in mm²
}
