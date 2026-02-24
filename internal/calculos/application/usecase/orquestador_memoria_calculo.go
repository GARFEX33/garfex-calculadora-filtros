// internal/calculos/application/usecase/orquestador_memoria_calculo.go
package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// OrquestadorMemoriaCalculoUseCase executes the complete memory calculation pipeline.
// It chains all steps sequentially where each step's output feeds the next.
type OrquestadorMemoriaCalculoUseCase struct {
	// Use cases for each step
	calcularCorrienteUC         *CalcularCorrienteUseCase
	ajustarCorrienteUC          *AjustarCorrienteUseCase
	seleccionarConductorUC      *SeleccionarConductorUseCase
	calcularTamanioTuberiaUC    *CalcularTamanioTuberiaUseCase
	calcularCharolaEspaciadoUC  *CalcularCharolaEspaciadoUseCase
	calcularCharolaTriangularUC *CalcularCharolaTriangularUseCase
	calcularCaidaTensionUC      *CalcularCaidaTensionUseCase

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
	tablaRepo port.TablaNOMRepository,
) *OrquestadorMemoriaCalculoUseCase {
	return &OrquestadorMemoriaCalculoUseCase{
		calcularCorrienteUC:         calcularCorrienteUC,
		ajustarCorrienteUC:          ajustarCorrienteUC,
		seleccionarConductorUC:      seleccionarConductorUC,
		calcularTamanioTuberiaUC:    calcularTamanioTuberiaUC,
		calcularCharolaEspaciadoUC:  calcularCharolaEspaciadoUC,
		calcularCharolaTriangularUC: calcularCharolaTriangularUC,
		calcularCaidaTensionUC:      calcularCaidaTensionUC,
		tablaRepo:                   tablaRepo,
	}
}

// calcularNumHilosTierra calcula el número de hilos de tierra según las reglas de la NOM.
//
// Reglas de negocio:
//   - Charola (cualquier tipo) → siempre 1 hilo de tierra
//   - Tubería con ≤2 tubos → 1 hilo de tierra
//   - Tubería con >2 tubos → 2 hilos de tierra
//
// Ejemplos de uso:
//
//	calcularNumHilosTierra(entity.TipoCanalizacionCharolaCableEspaciado, 5) // returns 1
//	calcularNumHilosTierra(entity.TipoCanalizacionTuboPVC, 2)              // returns 1
//	calcularNumHilosTierra(entity.TipoCanalizacionTuboPVC, 3)              // returns 2
func calcularNumHilosTierra(tipoCanalizacion entity.TipoCanalizacion, numTuberias int) int {
	// Charola siempre tiene 1 hilo de tierra
	if tipoCanalizacion.EsCharola() {
		return 1
	}

	// Default para valores inválidos
	if numTuberias <= 0 {
		return 1
	}

	// Tubería: ≤2 tubos = 1 hilo, >2 tubos = 2 hilos
	if numTuberias <= 2 {
		return 1
	}
	return 2
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
	if tipoCanalizacion.EsCharola() {
		// For CHAROLA: need to look up diameters first
		diametroFase, err := uc.tablaRepo.ObtenerDiametroConductor(
			ctx,
			output.ConductorAlimentacion.Calibre,
			material.String(),
			true, // with insulation
		)
		if err != nil {
			return dto.MemoriaOutput{}, fmt.Errorf("obtener diámetro fase: %w", err)
		}

		diametroTierra, err := uc.tablaRepo.ObtenerDiametroConductor(
			ctx,
			output.ConductorTierra.Calibre,
			material.String(),
			false, // ground is bare
		)
		if err != nil {
			return dto.MemoriaOutput{}, fmt.Errorf("obtener diámetro tierra: %w", err)
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
			// For triangular, use triangular use case (similar input structure)
			triangularInput := dto.CharolaTriangularInput{
				HilosPorFase:     input.HilosPorFase,
				DiametroFaseMM:   diametroFase,
				DiametroTierraMM: diametroTierra,
			}
			resultadoTriangular, err := uc.calcularCharolaTriangularUC.Execute(ctx, triangularInput)
			if err != nil {
				return dto.MemoriaOutput{}, fmt.Errorf("paso 4 (charola triangular): %w", err)
			}
			// Map to compatible output
			resultadoCanalizacion = dto.CharolaEspaciadoOutput{
				Tipo:           resultadoTriangular.Tipo,
				Tamano:         resultadoTriangular.Tamano,
				TamanoPulgadas: resultadoTriangular.TamanoPulgadas,
				AnchoRequerido: resultadoTriangular.AnchoRequerido,
			}

		case entity.TipoCanalizacionCharolaCableEspaciado:
			resultadoCharola, err := uc.calcularCharolaEspaciadoUC.Execute(ctx, charolaInput)
			if err != nil {
				return dto.MemoriaOutput{}, fmt.Errorf("paso 4 (charola espaciado): %w", err)
			}
			resultadoCanalizacion = resultadoCharola

		default:
			return dto.MemoriaOutput{}, fmt.Errorf("tipo de canalización no soportado para charola: %s", tipoCanalizacion)
		}

		// Map charola result to output
		output.Canalizacion = dto.ResultadoCanalizacion{
			Tamano:           resultadoCanalizacion.Tamano,
			AreaRequeridaMM2: resultadoCanalizacion.AnchoRequerido,
			NumeroDeTubos:    1,
		}

		output.Pasos = append(output.Pasos, dto.PasoMemoria{
			Numero:      4,
			Nombre:      "Dimensionamiento de Charola",
			Descripcion: "Cálculo de tamaño de charola según configuración",
			Resultado:   resultadoCanalizacion,
		})

	} else {
		// For TUBERIA types: use tuberia use case

		// Calcular número de hilos de tierra según normativa NOM
		numTierras := calcularNumHilosTierra(tipoCanalizacion, input.NumTuberias)

		tuberiaInput := dto.TuberiaInput{
			NumFases:         sistemaElectrico.CantidadFases(),
			CalibreFase:      output.ConductorAlimentacion.Calibre,
			NumNeutros:       numNeutros,
			CalibreNeutro:    output.ConductorAlimentacion.Calibre, // Same as fase for now
			CalibreTierra:    output.ConductorTierra.Calibre,
			TipoCanalizacion: input.TipoCanalizacion,
			NumTuberias:      input.NumTuberias,
			NumTierras:       numTierras,
		}

		resultadoTuberia, err := uc.calcularTamanioTuberiaUC.Execute(ctx, tuberiaInput)
		if err != nil {
			return dto.MemoriaOutput{}, fmt.Errorf("paso 4 (tubería): %w", err)
		}

		output.Canalizacion = dto.ResultadoCanalizacion{
			Tamano:           resultadoTuberia.TuberiaRecomendada,
			AreaTotalMM2:     resultadoTuberia.AreaPorTuboMM2,
			AreaRequeridaMM2: 0, // Not directly provided by this use case
			NumeroDeTubos:    resultadoTuberia.NumTuberias,
		}

		output.Pasos = append(output.Pasos, dto.PasoMemoria{
			Numero:      4,
			Nombre:      "Dimensionamiento de Tubería",
			Descripcion: "Cálculo de tamaño de tubería según área de conductores",
			Resultado:   resultadoTuberia,
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
	// Get voltage reference for calculation
	// According to NOM: MONOFASICO, BIFASICO, ESTRELLA use Vfn; DELTA uses Vff
	// If user provided Vff but system needs Vfn, convert
	tensionReferencia := float64(tension.Valor())

	// Adjust voltage reference based on tipo voltaje and system type
	if sistemaElectrico == entity.SistemaElectricoDelta && tipoVoltaje.EsFaseNeutro() {
		// User provided Vfn but DELTA needs Vff: Vff = Vfn * √3
		tensionReferencia = tensionReferencia * math.Sqrt(3)
	} else if (sistemaElectrico == entity.SistemaElectricoMonofasico ||
		sistemaElectrico == entity.SistemaElectricoBifasico ||
		sistemaElectrico == entity.SistemaElectricoEstrella) && tipoVoltaje.EsFaseFase() {
		// User provided Vff but system needs Vfn: Vfn = Vff / √3
		tensionReferencia = tensionReferencia / math.Sqrt(3)
	}

	tensionReferenciaVO, err := valueobject.NewTension(tensionReferencia, "V")
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("tensión de referencia inválida: %w", err)
	}

	resultadoCaidaTension, err := uc.calcularCaidaTensionUC.Execute(
		ctx,
		output.ConductorAlimentacion.Calibre,
		material,
		corrienteAjustada,
		input.LongitudCircuito, // already in meters from input
		tensionReferenciaVO,
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
	}

	output.Pasos = append(output.Pasos, dto.PasoMemoria{
		Numero:      5,
		Nombre:      "Caída de Tensión",
		Descripcion: "Cálculo de caída de tensión según NOM-001",
		Resultado:   resultadoCaidaTension,
	})

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
