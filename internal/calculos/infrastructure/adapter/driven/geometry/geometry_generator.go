// internal/calculos/infrastructure/adapter/driven/geometry/geometry_generator.go
package geometry

import (
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/pdf/geometry"
)

// GeometryGeneratorAdapter implementa el puerto GeometryGeneratorPort
// usando el paquete interno pdf/geometry.
type GeometryGeneratorAdapter struct{}

// NewGeometryGeneratorAdapter crea una nueva instancia del adapter.
func NewGeometryGeneratorAdapter() *GeometryGeneratorAdapter {
	return &GeometryGeneratorAdapter{}
}

// GenerarDiagramaCharola implementa GeometryGeneratorPort.
func (a *GeometryGeneratorAdapter) GenerarDiagramaCharola(
	diametroFaseMM float64,
	diametroTierraMM float64,
	diametroControlMM *float64,
	numHilosControl int,
	sistemaElectrico string,
	hilosPorFase int,
	anchoComercialMM float64,
	areaRequeridaMM2 float64,
	tipoCanalizacion string,
) (*port.GeometryDiagramaCharola, error) {

	// Parse sistema eléctrico
	sisElec, err := geometry.ParseSistemaElectrico(sistemaElectrico)
	if err != nil {
		return nil, err
	}

	// Factor triangular estándar NOM
	factorTriangular := 2.15

	// Seleccionar función según tipo de canalización
	var posiciones []geometry.ConductorPosicion
	var tipo string

	switch tipoCanalizacion {
	case "CHAROLA_CABLE_TRIANGULAR":
		// Usar distribución triangular
		params := geometry.ParametrosCharolaTriangular{
			DiametroFaseMM:    diametroFaseMM,
			DiametroTierraMM:  diametroTierraMM,
			DiametroControlMM: diametroControlMM,
			HilosPorFase:      hilosPorFase,
			FactorTriangular:  factorTriangular,
			AnchoComercialMM:  anchoComercialMM,
			SistemaElectrico:  sisElec,
		}
		posiciones = geometry.CalcularPosicionesCharolaTriangular(params)
		tipo = "triangular"

	case "CHAROLA_CABLE_ESPACIADO":
		fallthrough
	default:
		// Usar distribución espaciada (default)
		params := geometry.ParametrosCharolaBase{
			DiametroFaseMM:    diametroFaseMM,
			DiametroTierraMM:  diametroTierraMM,
			DiametroControlMM: diametroControlMM,
			NumHilosControl:   numHilosControl,
			SistemaElectrico:  sisElec,
			HilosPorFase:      hilosPorFase,
			AnchoComercialMM:  anchoComercialMM,
		}
		posiciones = geometry.CalcularPosicionesCharolaEspaciada(params)
		tipo = "espaciada"
	}

	// Calcular ancho ocupado
	anchoOcupado := geometry.CalcularAnchoOcupadoCharola(posiciones)

	// Generar SVG completo
	svg := geometry.GenerarSVGCompletoCharola(geometry.ParametrosSVGCharola{
		Posiciones:       posiciones,
		AnchoComercialMM: anchoComercialMM,
		PeralteMM:        geometry.PeralteCharolaMM,
		TipoDistribucion: tipo,
	})

	// Calcular viewBox y cotas
	viewBox := geometry.CalcularViewBox(anchoComercialMM, geometry.PeralteCharolaMM, 20)
	cotas := geometry.CalcularCotasCharola(anchoComercialMM, areaRequeridaMM2, geometry.PeralteCharolaMM)

	// Convertir tipos internos a tipos del puerto
	result := &port.GeometryDiagramaCharola{
		Posiciones:   make([]port.GeometryConductorPosicion, len(posiciones)),
		AnchoOcupado: anchoOcupado,
		ViewBox:      viewBox.ViewBox,
		Cotas:        make([]port.GeometryLineaCota, len(cotas)),
		SVG:          svg,
	}

	for i, pos := range posiciones {
		result.Posiciones[i] = port.GeometryConductorPosicion{
			CX:       pos.CX,
			CY:       pos.CY,
			Radio:    pos.Radio,
			Color:    pos.Color,
			Etiqueta: pos.Etiqueta,
			Tipo:     string(pos.Tipo),
		}
	}

	for i, cota := range cotas {
		result.Cotas[i] = port.GeometryLineaCota{
			X1:            cota.X1,
			Y1:            cota.Y1,
			X2:            cota.X2,
			Y2:            cota.Y2,
			Valor:         cota.Valor,
			Texto:         cota.Texto,
			PosicionTexto: cota.PosicionTexto,
		}
	}

	return result, nil
}

// GenerarDiagramaTuberia implementa GeometryGeneratorPort.
func (a *GeometryGeneratorAdapter) GenerarDiagramaTuberia(
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
) (*port.GeometryDiagramaTuberia, error) {

	// Parse sistema eléctrico
	sisElec, err := geometry.ParseSistemaElectrico(sistemaElectrico)
	if err != nil {
		return nil, err
	}

	// Construir parámetros para cálculo de posiciones
	params := geometry.ParametrosTuberia{
		DiametroInteriorMM: diametroInteriorMM,
		DiametroExteriorMM: diametroExteriorMM,
		AreaFaseMM2:        areaFaseMM2,
		AreaNeutroMM2:      areaNeutroMM2,
		AreaTierraMM2:      areaTierraMM2,
		NumFasesPorTubo:    numFasesPorTubo,
		NumNeutrosPorTubo:  numNeutrosPorTubo,
		NumTierras:         numTierras,
		SistemaElectrico:   sisElec,
	}

	// Calcular posiciones de conductores
	posiciones := geometry.CalcularPosicionesTuberia(params)

	// Generar SVG completo con el número de tubos
	svg := geometry.GenerarSVGCompletoTuberia(posiciones, diametroInteriorMM, diametroExteriorMM, numTubos, 30)

	// Calcular viewBox considerando el número de tubos
	viewBox := geometry.CalcularViewBox(diametroExteriorMM, diametroExteriorMM, 30, numTubos)

	// Convertir tipos internos a tipos del puerto
	result := &port.GeometryDiagramaTuberia{
		Posiciones:       make([]port.GeometryConductorPosicion, len(posiciones)),
		DiametroInterior: diametroInteriorMM,
		DiametroExterior: diametroExteriorMM,
		ViewBox:          viewBox.ViewBox,
		SVG:              svg,
	}

	for i, pos := range posiciones {
		result.Posiciones[i] = port.GeometryConductorPosicion{
			CX:       pos.CX,
			CY:       pos.CY,
			Radio:    pos.Radio,
			Color:    pos.Color,
			Etiqueta: pos.Etiqueta,
			Tipo:     string(pos.Tipo),
		}
	}

	return result, nil
}

// Verify interface implementation at compile time
var _ port.GeometryGeneratorPort = (*GeometryGeneratorAdapter)(nil)
