// internal/calculos/application/port/geometry_generator.go
package port

// GeometryGeneratorPort es el puerto para generar diagramas SVG de canalización.
// La implementación está en infrastructure y usa el paquete pdf/geometry.
type GeometryGeneratorPort interface {
	// GenerarDiagramaCharola genera el SVG para una charola con distribución de cables.
	// tipoCanalizacion debe ser "CHAROLA_CABLE_ESPACIADO" o "CHAROLA_CABLE_TRIANGULAR".
	GenerarDiagramaCharola(
		diametroFaseMM float64,
		diametroTierraMM float64,
		diametroControlMM *float64,
		numHilosControl int,
		sistemaElectrico string,
		hilosPorFase int,
		anchoComercialMM float64,
		areaRequeridaMM2 float64,
		tipoCanalizacion string,
	) (*GeometryDiagramaCharola, error)

	// GenerarDiagramaTuberia genera el SVG para una tubería con conductores.
	GenerarDiagramaTuberia(
		areaFaseMM2 float64,
		areaNeutroMM2 *float64,
		areaTierraMM2 float64,
		numFasesPorTubo int,
		numNeutrosPorTubo int,
		numTierras int,
		sistemaElectrico string,
		diametroInteriorMM float64,
		diametroExteriorMM float64,
		numTubos int,
	) (*GeometryDiagramaTuberia, error)
}

// GeometryDiagramaCharola contiene el resultado de generar un diagrama de charola.
type GeometryDiagramaCharola struct {
	Posiciones   []GeometryConductorPosicion
	AnchoOcupado float64
	ViewBox      string
	Cotas        []GeometryLineaCota
	SVG          string
}

// GeometryDiagramaTuberia contiene el resultado de generar un diagrama de tubería.
type GeometryDiagramaTuberia struct {
	Posiciones       []GeometryConductorPosicion
	DiametroInterior float64
	DiametroExterior float64
	ViewBox          string
	SVG              string
}

// GeometryConductorPosicion representa la posición de un conductor en el diagrama.
type GeometryConductorPosicion struct {
	CX       float64
	CY       float64
	Radio    float64
	Color    string
	Etiqueta string
	Tipo     string
}

// GeometryLineaCota representa una línea de dimensión.
type GeometryLineaCota struct {
	X1            float64
	Y1            float64
	X2            float64
	Y2            float64
	Valor         float64
	Texto         string
	PosicionTexto string
}
