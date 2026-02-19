// internal/calculos/domain/service/calcular_tamanio_tuberia.go
package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// ErrTuberiaNoEncontrada is returned when no conduit size fits the required area.
var ErrTuberiaNoEncontrada = errors.New("no se encontró tamaño de tubería con área suficiente")

// ErrAreaRequeridaInvalida is returned when areaRequerida is not positive.
var ErrAreaRequeridaInvalida = errors.New("el área requerida debe ser mayor que cero")

// ErrTablaOcupacionVacia is returned when the occupation table is empty.
var ErrTablaOcupacionVacia = errors.New("tabla de ocupación vacía")

// TablaOcupacionRepository defines the contract for reading conduit occupation tables.
// This is a domain-level interface that abstracts the data source.
type TablaOcupacionRepository interface {
	// ObtenerTablaOcupacion returns occupation table entries for the given canalization type.
	// The table should contain entries sorted by area_ocupacion ascending.
	ObtenerTablaOcupacion(canalizacion entity.TipoCanalizacion) ([]valueobject.EntradaTablaOcupacion, error)
}

// CalcularAreaPorTubo calculates the total conductor area per tube in mm².
// This method implements the distribution rules:
// - Phase conductors are distributed evenly across tubes
// - Neutral conductors are distributed evenly across tubes
// - Ground conductors are NOT distributed - each tube gets the full ground area
func CalcularAreaPorTubo(
	fases int,
	neutros int,
	tierras int,
	areaFase float64,
	areaNeutral float64,
	areaTierra float64,
	numTuberias int,
) float64 {
	// Phase area per tube: distribute evenly
	fasesPorTubo := float64(fases) / float64(numTuberias)
	areaFasesPorTubo := fasesPorTubo * areaFase

	// Neutral area per tube: distribute evenly
	neutrosPorTubo := float64(neutros) / float64(numTuberias)
	areaNeutrosPorTubo := neutrosPorTubo * areaNeutral

	// Ground area per tube: NOT distributed - full ground goes in each tube
	areaTierraPorTubo := float64(tierras) * areaTierra * float64(numTuberias)

	return areaFasesPorTubo + areaNeutrosPorTubo + areaTierraPorTubo
}

// CalcularAreaRequerida calculates the required conduit interior area
// based on conductor area and fill factor (40% for >2 conductors per NOM).
func CalcularAreaRequerida(areaConectores float64) float64 {
	const fillFactor = 0.40 // 40% fill for more than 2 conductors per NOM Chapter 9
	return areaConectores / fillFactor
}

// BuscarTamanioTuberia finds the smallest conduit whose occupation area
// accommodates the required area. The table should be already loaded.
//
// NOTE: The occupation tables already incorporate the 40% fill factor per NOM Chapter 9.
func BuscarTamanioTuberia(
	areaRequerida float64,
	tipoCanalizacion entity.TipoCanalizacion,
	tablaOcupacion []valueobject.EntradaTablaOcupacion,
) (entity.ResultadoTamanioTuberia, error) {
	if areaRequerida <= 0 {
		return entity.ResultadoTamanioTuberia{}, fmt.Errorf("BuscarTamanioTuberia: %w", ErrAreaRequeridaInvalida)
	}
	if err := entity.ValidarTipoCanalizacion(tipoCanalizacion); err != nil {
		return entity.ResultadoTamanioTuberia{}, fmt.Errorf("BuscarTamanioTuberia: %w", err)
	}
	if len(tablaOcupacion) == 0 {
		return entity.ResultadoTamanioTuberia{}, fmt.Errorf("BuscarTamanioTuberia: %w", ErrTablaOcupacionVacia)
	}

	// The occupation tables already have 40% fill factor applied
	// Search directly where area_ocupacion >= area_conductores
	for _, entrada := range tablaOcupacion {
		if entrada.AreaOcupacionMM2 >= areaRequerida {
			return entity.NewResultadoTamanioTuberia(
				areaRequerida,
				entrada.Tamano,
				entrada.DesignacionMetrica,
				tipoCanalizacion,
				1, // Default to single conduit; caller should adjust if needed
			)
		}
	}

	return entity.ResultadoTamanioTuberia{}, fmt.Errorf(
		"%w: área requerida %.2f mm² excede máxima disponible %.2f mm²",
		ErrTuberiaNoEncontrada, areaRequerida, tablaOcupacion[len(tablaOcupacion)-1].AreaOcupacionMM2,
	)
}

// CalcularTamanioTuberiaWithMultiplePipes calculates conduit size considering multiple parallel pipes.
// It distributes conductors evenly across pipes and finds the appropriate size.
//
// NOTE: The occupation tables already incorporate the 40% fill factor per NOM Chapter 9.
// Therefore, we search directly where area_ocupacion >= area_conductores (no division by 0.40 needed).
func CalcularTamanioTuberiaWithMultiplePipes(
	fases int,
	neutros int,
	tierras int,
	areaFase float64,
	areaNeutral float64,
	areaTierra float64,
	numTuberias int,
	tipoCanalizacion entity.TipoCanalizacion,
	tablaOcupacion []valueobject.EntradaTablaOcupacion,
) (entity.ResultadoTamanioTuberia, error) {
	if numTuberias < 1 {
		return entity.ResultadoTamanioTuberia{}, fmt.Errorf("CalcularTamanioTuberiaWithMultiplePipes: %w", ErrNumeroDeTubosInvalido)
	}

	// Calculate total conductor area per tube (already represents the area that needs to fit)
	areaConectores := CalcularAreaPorTubo(fases, neutros, tierras, areaFase, areaNeutral, areaTierra, numTuberias)

	// The occupation tables already have 40% fill factor applied (area_ocupacion = area_interior * 0.40)
	// So we search directly where area_ocupacion >= area_conductores
	// (no need to divide by 0.40 again)
	for _, entrada := range tablaOcupacion {
		if entrada.AreaOcupacionMM2 >= areaConectores {
			return entity.NewResultadoTamanioTuberia(
				areaConectores,
				entrada.Tamano,
				entrada.DesignacionMetrica,
				tipoCanalizacion,
				numTuberias,
			)
		}
	}

	return entity.ResultadoTamanioTuberia{}, fmt.Errorf(
		"%w: área requerida %.2f mm² excede máxima disponible %.2f mm² (numTuberias=%d)",
		ErrTuberiaNoEncontrada, areaConectores, tablaOcupacion[len(tablaOcupacion)-1].AreaOcupacionMM2, numTuberias,
	)
}

// DiseñoMetricoConduit returns the metric designation for a given trade size.
// This is a helper for converting between trade sizes and metric designations.
func DiseñoMetricoConduit(tamano string) (string, error) {
	conversions := map[string]string{
		"1/2":   "13mm",
		"3/4":   "19mm",
		"1":     "25mm",
		"1 1/4": "32mm",
		"1 1/2": "38mm",
		"2":     "51mm",
		"2 1/2": "64mm",
		"3":     "75mm",
		"3 1/2": "89mm",
		"4":     "100mm",
		"5":     "125mm",
		"6":     "150mm",
	}

	if val, ok := conversions[tamano]; ok {
		return val, nil
	}
	return "", fmt.Errorf("tamaño de tubería no reconocido: %s", tamano)
}
