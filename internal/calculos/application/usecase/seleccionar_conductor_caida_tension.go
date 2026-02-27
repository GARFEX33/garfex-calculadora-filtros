// internal/calculos/application/usecase/seleccionar_conductor_caida_tension.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// SeleccionarConductorPorCaidaTensionUseCase busca el calibre de conductor mínimo
// que cumple con la caída de tensión permitida según NOM-001-SEDE.
// Se utiliza como fallback cuando el conductor seleccionado por ampacidad no cumple.
type SeleccionarConductorPorCaidaTensionUseCase struct {
	calcularCaidaTensionUC *CalcularCaidaTensionUseCase
	tablaRepo              port.TablaNOMRepository
}

// NewSeleccionarConductorPorCaidaTensionUseCase crea una nueva instancia.
func NewSeleccionarConductorPorCaidaTensionUseCase(
	calcularCaidaTensionUC *CalcularCaidaTensionUseCase,
	tablaRepo port.TablaNOMRepository,
) *SeleccionarConductorPorCaidaTensionUseCase {
	return &SeleccionarConductorPorCaidaTensionUseCase{
		calcularCaidaTensionUC: calcularCaidaTensionUC,
		tablaRepo:              tablaRepo,
	}
}

// Execute busca el calibre mínimo que cumple con la caída de tensión permitida.
// Comienza desde el calibre seleccionado por ampacidad y prueba calibres superiores
// hasta encontrar uno que cumpla con el límite de caída de tensión.
func (uc *SeleccionarConductorPorCaidaTensionUseCase) Execute(
	ctx              context.Context,
	calibreAmpacidad string,              // calibre seleccionado por ampacidad (punto de partida)
	material         valueobject.MaterialConductor,
	corrienteNominal valueobject.Corriente, // NOM usa corriente nominal para caída de tensión
	longitud         float64,              // metros
	tension          valueobject.Tension,
	limiteCaida      float64,              // porcentaje máximo permitido
	tipoCanalizacion entity.TipoCanalizacion,
	sistemaElectrico entity.SistemaElectrico,
	tipoVoltaje      entity.TipoVoltaje,
	hilosPorFase     int,
	factorPotencia   float64,
	temperatura      valueobject.Temperatura, // para ObtenerCapacidadConductor
) (dto.ResultadoConductorCaidaTension, error) {
	// La tabla NOM tiene 19 calibres (14 AWG → 1000 MCM).
	// En el peor caso se necesitan hasta 17 saltos desde el calibre más pequeño.
	// Usar 18 garantiza recorrer toda la tabla sin importar el punto de partida.
	const maxIntentos = 18

	calibreActual := calibreAmpacidad
	var ultimoResultado dto.ResultadoCaidaTension
	var ultimoCalibre string

	for intento := 0; intento < maxIntentos; intento++ {
		// Obtener siguiente calibre superior
		calibreSiguiente, err := service.ObtenerCalibreSuperior(calibreActual)
		if err != nil {
			// Llegamos al máximo de la tabla NOM, no hay más calibres
			break
		}

		// Calcular caída de tensión con el nuevo calibre
		resultadoCaida, err := uc.calcularCaidaTensionUC.Execute(
			ctx,
			calibreSiguiente,
			material,
			corrienteNominal,
			longitud,
			tension,
			limiteCaida,
			tipoCanalizacion,
			sistemaElectrico,
			tipoVoltaje,
			hilosPorFase,
			factorPotencia,
		)
		if err != nil {
			return dto.ResultadoConductorCaidaTension{}, fmt.Errorf("calcular caída de tensión para calibre %s: %w", calibreSiguiente, err)
		}

		ultimoResultado = resultadoCaida
		ultimoCalibre = calibreSiguiente

		if resultadoCaida.Cumple {
			// Obtener datos físicos del nuevo calibre para el output
			seccion, err := uc.tablaRepo.ObtenerSeccionConductor(ctx, calibreSiguiente)
			if err != nil {
				return dto.ResultadoConductorCaidaTension{}, fmt.Errorf("obtener sección para calibre %s: %w", calibreSiguiente, err)
			}
			capacidad, err := uc.tablaRepo.ObtenerCapacidadConductor(ctx, tipoCanalizacion, material, temperatura, calibreSiguiente)
			if err != nil {
				return dto.ResultadoConductorCaidaTension{}, fmt.Errorf("obtener capacidad para calibre %s: %w", calibreSiguiente, err)
			}

			nota := fmt.Sprintf(
				"Calibre aumentado de %s a %s por verificación de caída de tensión (NOM-001-SEDE)",
				calibreAmpacidad,
				calibreSiguiente,
			)

			return dto.ResultadoConductorCaidaTension{
				CalibreOriginal:     calibreAmpacidad,
				CalibreSeleccionado: calibreSiguiente,
				SeccionMM2:          seccion,
				TipoAislamiento:     "THW", // TODO: obtener de tabla cuando esté disponible
				Capacidad:           capacidad,
				CaidaTension:        resultadoCaida,
				Nota:                nota,
				Cumple:              true,
				IntentosRealizados:  intento + 1,
			}, nil
		}

		calibreActual = calibreSiguiente
	}

	// Agotamos los intentos sin encontrar calibre válido
	// Retornamos el último resultado con Cumple=false (no es error fatal)
	intentosRealizados := maxIntentos
	if ultimoCalibre == "" {
		// No se pudo ni un intento (calibreAmpacidad ya era el máximo)
		ultimoCalibre = calibreAmpacidad
		intentosRealizados = 0
	}

	return dto.ResultadoConductorCaidaTension{
		CalibreOriginal:     calibreAmpacidad,
		CalibreSeleccionado: ultimoCalibre,
		CaidaTension:        ultimoResultado,
		Nota:                fmt.Sprintf("No se encontró calibre que cumpla la caída de tensión tras %d intentos", intentosRealizados),
		Cumple:              false,
		IntentosRealizados:  intentosRealizados,
	}, nil
}
