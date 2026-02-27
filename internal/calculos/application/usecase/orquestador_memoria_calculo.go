// internal/calculos/application/usecase/orquestador_memoria_calculo.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// OrquestadorMemoriaCalculoUseCase executes the complete memory calculation pipeline.
// It chains all steps sequentially where each step's output feeds the next.
type OrquestadorMemoriaCalculoUseCase struct {
	// Use cases for each step
	calcularCorrienteUC                  *CalcularCorrienteUseCase
	ajustarCorrienteUC                   *AjustarCorrienteUseCase
	seleccionarConductorUC               *SeleccionarConductorUseCase
	calcularTamanioTuberiaUC             *CalcularTamanioTuberiaUseCase
	calcularCharolaEspaciadoUC           *CalcularCharolaEspaciadoUseCase
	calcularCharolaTriangularUC          *CalcularCharolaTriangularUseCase
	calcularCaidaTensionUC               *CalcularCaidaTensionUseCase
	seleccionarConductorCaidaTensionUC   *SeleccionarConductorPorCaidaTensionUseCase

	// Repository for diameter lookups (needed for charola)
	tablaRepo port.TablaNOMRepository
}

// NewOrquestadorMemoriaCalculoUseCase creates a new orchestrator instance.
func NewOrquestadorMemoriaCalculoUseCase(
	calcularCorrienteUC *CalcularCorrienteUseCase,
	ajustarCorrienteUC *AjustarCorrienteUseCase,
	seleccionarConductorUC *SeleccionarConductorUseCase,
	calcularTamanioTuberiaUC *CalcularTamanioTuberiaUseCase,
	calcularCharolaEspaciadoUC *CalcularCharolaEspaciadoUseCase,
	calcularCharolaTriangularUC *CalcularCharolaTriangularUseCase,
	calcularCaidaTensionUC *CalcularCaidaTensionUseCase,
	seleccionarConductorCaidaTensionUC *SeleccionarConductorPorCaidaTensionUseCase,
	tablaRepo port.TablaNOMRepository,
) *OrquestadorMemoriaCalculoUseCase {
	return &OrquestadorMemoriaCalculoUseCase{
		calcularCorrienteUC:                calcularCorrienteUC,
		ajustarCorrienteUC:                ajustarCorrienteUC,
		seleccionarConductorUC:            seleccionarConductorUC,
		calcularTamanioTuberiaUC:          calcularTamanioTuberiaUC,
		calcularCharolaEspaciadoUC:        calcularCharolaEspaciadoUC,
		calcularCharolaTriangularUC:       calcularCharolaTriangularUC,
		calcularCaidaTensionUC:            calcularCaidaTensionUC,
		seleccionarConductorCaidaTensionUC: seleccionarConductorCaidaTensionUC,
		tablaRepo:                         tablaRepo,
	}
}

// calcularNumHilosTierra calcula el número de hilos de tierra según las reglas de la NOM.
//
// Reglas de negocio:
//   - Charola (cualquier tipo) → siempre 1 hilo de tierra
//   - Tubería → 1 hilo de tierra por tubo (= numTuberias)
//
// Ejemplos de uso:
//
//	calcularNumHilosTierra(entity.TipoCanalizacionCharolaCableEspaciado, 5) // returns 1
//	calcularNumHilosTierra(entity.TipoCanalizacionTuboPVC, 1)              // returns 1
//	calcularNumHilosTierra(entity.TipoCanalizacionTuboPVC, 2)              // returns 2
//	calcularNumHilosTierra(entity.TipoCanalizacionTuboPVC, 3)              // returns 3
func calcularNumHilosTierra(tipoCanalizacion entity.TipoCanalizacion, numTuberias int) int {
	// Charola siempre tiene 1 hilo de tierra
	if tipoCanalizacion.EsCharola() {
		return 1
	}

	// Default para valores inválidos
	if numTuberias <= 0 {
		return 1
	}

	// Tubería: 1 tierra por tubo
	return numTuberias
}

// Execute runs the complete memory calculation pipeline.
// It orchestrates all 6 steps sequentially.
func (uc *OrquestadorMemoriaCalculoUseCase) Execute(
	ctx context.Context,
	input dto.EquipoInput,
) (dto.MemoriaOutput, error) {
	// Step 0: Apply defaults and validate
	input.ApplyDefaults()
	if err := input.ValidateForMemoria(); err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("validación de entrada: %w", err)
	}

	// Convert DTOs to domain entities/types for downstream use
	tipoCanalizacion, err := entity.ParseTipoCanalizacion(input.TipoCanalizacion)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("tipo canalización inválido: %w", err)
	}

	sistemaElectrico := input.SistemaElectrico.ToEntity()

	tipoVoltaje, err := entity.ParseTipoVoltaje(input.TipoVoltaje)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("tipo voltaje inválido: %w", err)
	}

	material, err := input.ToDomainMaterial()
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("material inválido: %w", err)
	}

	tension, err := input.ToDomainTension()
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("tensión inválida: %w", err)
	}

	// Get TipoEquipo according to mode
	tipoEquipo, err := input.GetTipoEquipo()
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("tipo de equipo inválido: %w", err)
	}

	// Get ITM according to mode
	itm := input.Equipo.ITM

	// Prepare output structure
	output := dto.MemoriaOutput{
		Equipo:           input.Equipo,
		TipoEquipo:       string(tipoEquipo),
		Tension:          tension.Valor(),
		FactorPotencia:   input.FactorPotencia,
		Estado:           input.Estado,
		SistemaElectrico: input.SistemaElectrico,
		HilosPorFase:     input.HilosPorFase,
		TipoCanalizacion: input.TipoCanalizacion,
		Material:         input.Material,
		LongitudCircuito: input.LongitudCircuito,
		ITM:              itm,
	}

	// Determine neutral count from system type
	numNeutros := sistemaElectrico.CantidadNeutros()

	// ============================================================
	// STEP 1: Calculate Nominal Current
	// ============================================================
	resultadoCorriente, err := uc.calcularCorrienteUC.Execute(ctx, input)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 1 (corriente nominal): %w", err)
	}
	output.CorrienteNominal = resultadoCorriente.CorrienteNominal
	output.Pasos = append(output.Pasos, dto.PasoMemoria{
		Numero:      1,
		Nombre:      "Corriente Nominal",
		Descripcion: "Cálculo de corriente nominal desde potencia o amperaje",
		Resultado:   resultadoCorriente,
	})

	// ============================================================
	// STEP 2: Adjust Current (temperature, grouping, usage factors)
	// ============================================================
	corrienteNominalVO, err := valueobject.NewCorriente(output.CorrienteNominal)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("corriente nominal inválida: %w", err)
	}

	resultadoAjuste, err := uc.ajustarCorrienteUC.Execute(
		ctx,
		corrienteNominalVO,
		input.Estado,
		tipoCanalizacion,
		sistemaElectrico,
		tipoEquipo,
		input.HilosPorFase,
		input.NumTuberias,
	)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 2 (ajuste de corriente): %w", err)
	}
	output.CorrienteAjustada = resultadoAjuste.CorrienteAjustada
	output.CorrientePorHilo = resultadoAjuste.CorrienteAjustada / float64(input.HilosPorFase)
	output.FactorTemperatura = resultadoAjuste.FactorTemperatura
	output.FactorAgrupamiento = resultadoAjuste.FactorAgrupamiento
	output.FactorTotalAjuste = resultadoAjuste.FactorTotal
	output.TemperaturaAmbiente = resultadoAjuste.TemperaturaAmbiente
	output.CantidadConductores = resultadoAjuste.CantidadConductoresTotal
	output.ConductoresPorTubo = resultadoAjuste.ConductoresPorTubo
	output.Pasos = append(output.Pasos, dto.PasoMemoria{
		Numero:      2,
		Nombre:      "Ajuste de Corriente",
		Descripcion: "Aplicación de factores de temperatura, agrupamiento y uso",
		Resultado:   resultadoAjuste,
	})

	// ============================================================
	// STEP 3: Select Feed and Ground Conductors
	// ============================================================
	corrienteAjustada, err := valueobject.NewCorriente(resultadoAjuste.CorrienteAjustada)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("corriente ajustada inválida: %w", err)
	}

	// Get temperature used (from adjustment step or default)
	temperaturaUsada := valueobject.Temperatura(resultadoAjuste.Temperatura)

	resultadoConductores, err := uc.seleccionarConductorUC.Execute(
		ctx,
		corrienteAjustada,
		input.HilosPorFase,
		itm,
		material,
		temperaturaUsada,
		tipoCanalizacion,
	)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 3 (seleccionar conductor): %w", err)
	}
	output.ConductorAlimentacion = resultadoConductores.Alimentacion
	output.ConductorTierra = resultadoConductores.Tierra
	output.TablaAmpacidadUsada = resultadoConductores.TablaUsada
	output.TemperaturaUsada = temperaturaUsada.Valor()
	output.Pasos = append(output.Pasos, dto.PasoMemoria{
		Numero:      3,
		Nombre:      "Selección de Conductores",
		Descripcion: "Selección de conductor de alimentación y tierra",
		Resultado:   resultadoConductores,
	})

	// ============================================================
	// STEP 4: Size Conduit/Tray (branch by canalization type)
	// ============================================================
	canalizacion, detalleCharola, detalleTuberia, fillFactor, err := uc.calcularCanalizacion(
		ctx,
		output.ConductorAlimentacion.Calibre,
		output.ConductorTierra.Calibre,
		material,
		tipoCanalizacion,
		sistemaElectrico,
		input,
		numNeutros,
	)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 4 (canalización): %w", err)
	}
	output.Canalizacion = canalizacion
	output.DetalleCharola = detalleCharola
	output.DetalleTuberia = detalleTuberia
	output.FillFactor = fillFactor

	// Append the appropriate paso based on canalization type
	if tipoCanalizacion.EsCharola() {
		output.Pasos = append(output.Pasos, dto.PasoMemoria{
			Numero:      4,
			Nombre:      "Dimensionamiento de Charola",
			Descripcion: "Cálculo de tamaño de charola según configuración",
			Resultado:   detalleCharola,
		})
	} else {
		output.Pasos = append(output.Pasos, dto.PasoMemoria{
			Numero:      4,
			Nombre:      "Dimensionamiento de Tubería",
			Descripcion: "Cálculo de tamaño de tubería según área de conductores",
			Resultado:   detalleTuberia,
		})
	}

	// ============================================================
	// STEP 4b: Asignar número de hilos de tierra al conductor
	// ============================================================
	// El número de hilos de tierra se calcula según el tipo de canalización:
	// - Charola: siempre 1 hilo
	// - Tubería ≤2 tubos: 1 hilo
	// - Tubería >2 tubos: 2 hilos
	numHilosTierra := calcularNumHilosTierra(tipoCanalizacion, input.NumTuberias)
	output.ConductorTierra.NumHilos = numHilosTierra

	// ============================================================
	// STEP 5: Calculate Voltage Drop
	// ============================================================
	// Pasar la tensión tal como la ingresó el usuario (sin convertir).
	// La conversión Vfn↔Vff la realiza internamente el domain service
	// calcularVoltajeReferencia usando tipoVoltaje + sistemaElectrico.
	// Hacer la conversión aquí causaría doble conversión → %caída incorrecto.
	resultadoCaidaTension, err := uc.calcularCaidaTensionUC.Execute(
		ctx,
		output.ConductorAlimentacion.Calibre,
		material,
		corrienteNominalVO, // NOM-001-SEDE: caída de tensión usa corriente nominal, no ajustada
		input.LongitudCircuito, // already in meters from input
		tension,
		input.PorcentajeCaidaMaximo,
		tipoCanalizacion,
		sistemaElectrico,
		tipoVoltaje,
		input.HilosPorFase,
		input.FactorPotencia,
	)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 5 (caída de tensión): %w", err)
	}

	output.CaidaTension = dto.ResultadoCaidaTension{
		Porcentaje:       resultadoCaidaTension.Porcentaje,
		CaidaVolts:       resultadoCaidaTension.CaidaVolts,
		Cumple:           resultadoCaidaTension.Cumple,
		LimitePorcentaje: resultadoCaidaTension.LimitePorcentaje,
		Impedancia:       resultadoCaidaTension.Impedancia,
		Resistencia:      resultadoCaidaTension.Resistencia,
		Reactancia:       resultadoCaidaTension.Reactancia,
	}

	output.Pasos = append(output.Pasos, dto.PasoMemoria{
		Numero:      5,
		Nombre:      "Caída de Tensión",
		Descripcion: "Cálculo de caída de tensión según NOM-001",
		Resultado:   resultadoCaidaTension,
	})

	// ============================================================
	// STEP 5b: Recálculo por caída de tensión
	// Si el conductor de ampacidad no cumple la caída de tensión,
	// buscar el calibre superior mínimo que cumpla (NOM-001-SEDE).
	// ============================================================
	if !output.CaidaTension.Cumple {
		resultadoRecalc, err := uc.seleccionarConductorCaidaTensionUC.Execute(
			ctx,
			output.ConductorAlimentacion.Calibre,
			material,
			corrienteNominalVO,
			input.LongitudCircuito,
			tension,
			input.PorcentajeCaidaMaximo,
			tipoCanalizacion,
			sistemaElectrico,
			tipoVoltaje,
			input.HilosPorFase,
			input.FactorPotencia,
			temperaturaUsada,
		)
		if err != nil {
			output.Observaciones = append(output.Observaciones,
				fmt.Sprintf("No se pudo recalcular calibre por caída de tensión: %v", err))
		}
		if err == nil && resultadoRecalc.Cumple {
			// Override conductor de alimentación con el calibre superior
			calibreOriginal := output.ConductorAlimentacion.Calibre
			output.ConductorAlimentacion.Calibre = resultadoRecalc.CalibreSeleccionado
			output.ConductorAlimentacion.SeccionMM2 = resultadoRecalc.SeccionMM2
			output.ConductorAlimentacion.Capacidad = resultadoRecalc.Capacidad
			output.ConductorAlimentacion.TipoAislamiento = resultadoRecalc.TipoAislamiento
			output.ConductorAlimentacion.SeleccionPorCaidaTension = true
			output.ConductorAlimentacion.CalibreOriginalAmpacidad = calibreOriginal
			output.ConductorAlimentacion.NotaSeleccion = resultadoRecalc.Nota

			// Override caída de tensión con el resultado del nuevo calibre
			output.CaidaTension = dto.ResultadoCaidaTension{
				Porcentaje:       resultadoRecalc.CaidaTension.Porcentaje,
				CaidaVolts:       resultadoRecalc.CaidaTension.CaidaVolts,
				Cumple:           resultadoRecalc.CaidaTension.Cumple,
				LimitePorcentaje: resultadoRecalc.CaidaTension.LimitePorcentaje,
				Impedancia:       resultadoRecalc.CaidaTension.Impedancia,
				Resistencia:      resultadoRecalc.CaidaTension.Resistencia,
				Reactancia:       resultadoRecalc.CaidaTension.Reactancia,
			}

			// Re-ejecutar canalización con el nuevo calibre
			canalizacionRecalc, detalleCharolaRecalc, detalleTuberiaRecalc, fillFactorRecalc, errCanal := uc.calcularCanalizacion(
				ctx,
				resultadoRecalc.CalibreSeleccionado,
				output.ConductorTierra.Calibre,
				material,
				tipoCanalizacion,
				sistemaElectrico,
				input,
				numNeutros,
			)
			if errCanal != nil {
				// No es error fatal: mantener canalización original
				output.Observaciones = append(output.Observaciones,
					fmt.Sprintf("No se pudo recalcular canalización con calibre %s: %v",
						resultadoRecalc.CalibreSeleccionado, errCanal))
			} else {
				output.Canalizacion = canalizacionRecalc
				output.DetalleCharola = detalleCharolaRecalc
				output.DetalleTuberia = detalleTuberiaRecalc
				output.FillFactor = fillFactorRecalc
			}

			output.Pasos = append(output.Pasos, dto.PasoMemoria{
				Numero:      6,
				Nombre:      "Recálculo por Caída de Tensión",
				Descripcion: "Calibre aumentado al siguiente superior para cumplir caída de tensión NOM-001-SEDE",
				Resultado:   resultadoRecalc,
			})
		}
		// Si resultadoRecalc.Cumple == false: se agotaron calibres, mantener original con Cumple=false
	}

	// ============================================================
	// Final Assembly: CumpleNormativa and Observaciones
	// ============================================================
	output.CumpleNormativa = output.CaidaTension.Cumple

	// Generate observations
	output.Observaciones = uc.generarObservaciones(output)

	return output, nil
}

// generarObservaciones creates a list of observations based on the calculation results.
func (uc *OrquestadorMemoriaCalculoUseCase) generarObservaciones(memoria dto.MemoriaOutput) []string {
	var obs []string

	// Observation about voltage drop
	if memoria.CaidaTension.Cumple {
		obs = append(obs, fmt.Sprintf("La caída de tensión (%.2f%%) cumple con el límite de %.1f%%",
			memoria.CaidaTension.Porcentaje, memoria.CaidaTension.LimitePorcentaje))
	} else {
		obs = append(obs, fmt.Sprintf("ADVERTENCIA: La caída de tensión (%.2f%%) excede el límite de %.1f%%. Considere aumentar el calibre del conductor.",
			memoria.CaidaTension.Porcentaje, memoria.CaidaTension.LimitePorcentaje))
	}

	// Observation about conductor selection
	obs = append(obs, fmt.Sprintf("Conductor de alimentación: %s %s (%s, %.2f mm²)",
		memoria.ConductorAlimentacion.Material,
		memoria.ConductorAlimentacion.Calibre,
		memoria.ConductorAlimentacion.TipoAislamiento,
		memoria.ConductorAlimentacion.SeccionMM2))

	// Observation about ground conductor
	obs = append(obs, fmt.Sprintf("Conductor de tierra: %s %s (%.2f mm²)",
		memoria.ConductorTierra.Material,
		memoria.ConductorTierra.Calibre,
		memoria.ConductorTierra.SeccionMM2))

	// Observation about canalization
	obs = append(obs, fmt.Sprintf("Canalización recomendada: %s", memoria.Canalizacion.Tamano))

	// Observation about factors
	if memoria.FactorTotalAjuste < 1.0 {
		obs = append(obs, fmt.Sprintf("Factores aplicados: Temperatura=%.2f, Agrupamiento=%.2f, Uso=%.2f",
			memoria.FactorTemperatura, memoria.FactorAgrupamiento, memoria.FactorTotalAjuste))
	}

	return obs
}

// calcularCanalizacion ejecuta el paso 4 de dimensionamiento de canalización.
// Es llamado tanto en el flujo principal como en el recálculo por caída de tensión.
func (uc *OrquestadorMemoriaCalculoUseCase) calcularCanalizacion(
	ctx              context.Context,
	calibreFase      string,
	calibreTierra    string,
	material         valueobject.MaterialConductor,
	tipoCanalizacion entity.TipoCanalizacion,
	sistemaElectrico entity.SistemaElectrico,
	input            dto.EquipoInput,
	numNeutros       int,
) (canalizacion dto.ResultadoCanalizacion, detalleCharola *dto.DetalleCharola, detalleTuberia *dto.DetalleTuberia, fillFactor float64, err error) {

	if tipoCanalizacion.EsCharola() {
		// For CHAROLA: need to look up diameters first
		diametroFase, err := uc.tablaRepo.ObtenerDiametroConductor(
			ctx,
			calibreFase,
			material.String(),
			true, // with insulation
		)
		if err != nil {
			return dto.ResultadoCanalizacion{}, nil, nil, 0, fmt.Errorf("obtener diámetro fase: %w", err)
		}

		diametroTierra, err := uc.tablaRepo.ObtenerDiametroConductor(
			ctx,
			calibreTierra,
			material.String(),
			false, // ground is bare
		)
		if err != nil {
			return dto.ResultadoCanalizacion{}, nil, nil, 0, fmt.Errorf("obtener diámetro tierra: %w", err)
		}

		// Build charola input
		charolaInput := dto.CharolaEspaciadoInput{
			HilosPorFase:     input.HilosPorFase,
			SistemaElectrico: string(sistemaElectrico),
			DiametroFaseMM:   diametroFase,
			DiametroTierraMM: diametroTierra,
		}

		// Handle optional control diameter
		if input.DiametroControlMM != nil && *input.DiametroControlMM > 0 {
			charolaInput.DiametroControlMM = input.DiametroControlMM
		}

		var resultadoCanalizacion dto.CharolaEspaciadoOutput

		switch tipoCanalizacion {
		case entity.TipoCanalizacionCharolaCableTriangular:
			// For triangular, use triangular use case
			triangularInput := dto.CharolaTriangularInput{
				HilosPorFase:     input.HilosPorFase,
				DiametroFaseMM:   diametroFase,
				DiametroTierraMM: diametroTierra,
			}
			if input.DiametroControlMM != nil && *input.DiametroControlMM > 0 {
				triangularInput.DiametroControlMM = input.DiametroControlMM
			}
			resultadoTriangular, err := uc.calcularCharolaTriangularUC.Execute(ctx, triangularInput)
			if err != nil {
				return dto.ResultadoCanalizacion{}, nil, nil, 0, fmt.Errorf("charola triangular: %w", err)
			}
			// Map to compatible output
			resultadoCanalizacion = dto.CharolaEspaciadoOutput{
				Tipo:             resultadoTriangular.Tipo,
				Tamano:           resultadoTriangular.Tamano,
				TamanoPulgadas:   resultadoTriangular.TamanoPulgadas,
				AnchoRequerido:   resultadoTriangular.AnchoRequerido,
				AnchoComercialMM: resultadoTriangular.AnchoComercialMM,
			}
			// Poblar detalle con valores intermedios del triangular
			detalleCharola = &dto.DetalleCharola{
				DiametroFaseMM:    resultadoTriangular.DiametroFaseMM,
				DiametroTierraMM:  resultadoTriangular.DiametroTierraMM,
				DiametroControlMM: resultadoTriangular.DiametroControlMM,
				AnchoPotenciaMM:   resultadoTriangular.AnchoPotenciaMM,
				EspacioFuerzaMM:   resultadoTriangular.EspacioFuerzaMM,
				EspacioControlMM:  resultadoTriangular.EspacioControlMM,
				AnchoControlMM:    resultadoTriangular.AnchoControlMM,
				AnchoTierraMM:     resultadoTriangular.AnchoTierraMM,
				FactorTriangular:  resultadoTriangular.FactorTriangular,
			}

		case entity.TipoCanalizacionCharolaCableEspaciado:
			resultadoCharola, err := uc.calcularCharolaEspaciadoUC.Execute(ctx, charolaInput)
			if err != nil {
				return dto.ResultadoCanalizacion{}, nil, nil, 0, fmt.Errorf("charola espaciado: %w", err)
			}
			resultadoCanalizacion = resultadoCharola
			// Poblar detalle con valores intermedios del espaciado
			detalleCharola = &dto.DetalleCharola{
				DiametroFaseMM:    resultadoCharola.DiametroFaseMM,
				DiametroTierraMM:  resultadoCharola.DiametroTierraMM,
				DiametroControlMM: resultadoCharola.DiametroControlMM,
				NumHilosTotal:     resultadoCharola.NumHilosTotal,
				EspacioFuerzaMM:   resultadoCharola.EspacioFuerzaMM,
				AnchoFuerzaMM:     resultadoCharola.AnchoFuerzaMM,
				EspacioControlMM:  resultadoCharola.EspacioControlMM,
				AnchoControlMM:    resultadoCharola.AnchoControlMM,
				AnchoTierraMM:     resultadoCharola.AnchoTierraMM,
			}

		default:
			return dto.ResultadoCanalizacion{}, nil, nil, 0, fmt.Errorf("tipo de canalización no soportado para charola: %s", tipoCanalizacion)
		}

		// Map charola result to output
		canalizacion = dto.ResultadoCanalizacion{
			Tamano:           resultadoCanalizacion.Tamano,
			AnchoComercialMM: resultadoCanalizacion.AnchoComercialMM,
			AreaRequeridaMM2: resultadoCanalizacion.AnchoRequerido,
			NumeroDeTubos:    1,
		}

		return canalizacion, detalleCharola, nil, 0, nil

	} else {
		// For TUBERIA types: use tuberia use case

		// Calcular número de hilos de tierra según normativa NOM
		numTierras := calcularNumHilosTierra(tipoCanalizacion, input.NumTuberias)

		tuberiaInput := dto.TuberiaInput{
			NumFases:         sistemaElectrico.CantidadFases(),
			CalibreFase:      calibreFase,
			NumNeutros:       numNeutros,
			CalibreNeutro:    calibreFase, // Same as fase for now
			CalibreTierra:    calibreTierra,
			TipoCanalizacion: input.TipoCanalizacion,
			NumTuberias:      input.NumTuberias,
			NumTierras:       numTierras,
			HilosPorFase:     input.HilosPorFase, // Necesario para calcular conductores por tubo
		}

		resultadoTuberia, err := uc.calcularTamanioTuberiaUC.Execute(ctx, tuberiaInput)
		if err != nil {
			return dto.ResultadoCanalizacion{}, nil, nil, 0, fmt.Errorf("tubería: %w", err)
		}

		// Poblar DetalleTuberia con valores intermedios
		detalleTuberia = &dto.DetalleTuberia{
			AreaFaseMM2:          resultadoTuberia.AreaFaseMM2,
			AreaNeutroMM2:        resultadoTuberia.AreaNeutroMM2,
			AreaTierraMM2:        resultadoTuberia.AreaTierraMM2,
			NumFasesPorTubo:      resultadoTuberia.NumFasesPorTubo,
			NumNeutrosPorTubo:    resultadoTuberia.NumNeutrosPorTubo,
			NumTierras:           resultadoTuberia.NumTierras,
			AreaOcupacionTuboMM2: resultadoTuberia.AreaOcupacionTuboMM2,
			DesignacionMetrica:   resultadoTuberia.DesignacionMetrica,
			FillFactor:           resultadoTuberia.FillFactor,
		}

		canalizacion = dto.ResultadoCanalizacion{
			Tamano:           resultadoTuberia.TuberiaRecomendada,
			AreaTotalMM2:     resultadoTuberia.AreaPorTuboMM2,
			AreaRequeridaMM2: 0, // Not directly provided by this use case
			NumeroDeTubos:    resultadoTuberia.NumTuberias,
		}

		return canalizacion, nil, detalleTuberia, resultadoTuberia.FillFactor, nil
	}
}
