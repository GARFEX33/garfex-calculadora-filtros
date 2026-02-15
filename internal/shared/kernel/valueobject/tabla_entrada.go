// internal/shared/kernel/valueobject/tabla_entrada.go
package valueobject

// EntradaTablaConductor represents one row from NOM table 310-15(b)(16).
// Must be sorted smallest-to-largest calibre (as in the NOM table).
// Conductor holds the full physical/electrical properties needed to construct
// a Conductor value object.
type EntradaTablaConductor struct {
	Capacidad float64 // ampacity in amperes
	Conductor ConductorParams
}

// EntradaTablaTierra represents one row from NOM table 250-122.
// Entries must be sorted by ITMHasta ascending.
type EntradaTablaTierra struct {
	ITMHasta    int
	ConductorCu ConductorParams
	ConductorAl *ConductorParams // nil = not permitted for this ITM range
}

// EntradaTablaCanalizacion represents one row from a conduit sizing table.
// Entries must be sorted by AreaInteriorMM2 ascending.
type EntradaTablaCanalizacion struct {
	Tamano          string
	AreaInteriorMM2 float64
}
