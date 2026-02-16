// internal/calculos/application/usecase/orquestador_memoria.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// OrquestadorMemoriaCalculo orquesta el cálculo completo de memoria.
// Coordina los micro-use cases y monta el output final.
type OrquestadorMemoriaCalculo struct {
	calcularCorriente       *CalcularCorrienteUseCase
	ajustarCorriente        *AjustarCorrienteUseCase
	seleccionarConductor    *SeleccionarConductorUseCase
	dimensionarCanalizacion *DimensionarCanalizacionUseCase
	calcularCaidaTension    *CalcularCaidaTensionUseCase
	tablaRepo               port.TablaNOMRepository
}

// NewOrquestadorMemoriaCalculo crea una nueva instancia del orquestador.
func NewOrquestadorMemoriaCalculo(
	calcularCorriente *CalcularCorrienteUseCase,
	ajustarCorriente *AjustarCorrienteUseCase,
	seleccionarConductor *SeleccionarConductorUseCase,
	dimensionarCanalizacion *DimensionarCanalizacionUseCase,
	calcularCaidaTension *CalcularCaidaTensionUseCase,
	tablaRepo port.TablaNOMRepository,
) *OrquestadorMemoriaCalculo {
	return &OrquestadorMemoriaCalculo{
		calcularCorriente:       calcularCorriente,
		ajustarCorriente:        ajustarCorriente,
		seleccionarConductor:    seleccionarConductor,
		dimensionarCanalizacion: dimensionarCanalizacion,
		calcularCaidaTension:    calcularCaidaTension,
		tablaRepo:               tablaRepo,
	}
}

// Execute orquesta el cálculo completo de memoria de cálculo.
func (o *OrquestadorMemoriaCalculo) Execute(ctx context.Context, input dto.EquipoInput) (dto.MemoriaOutput, error) {
	// Validar input
	if err := input.Validate(); err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("validación de input: %w", err)
	}

	// Aplicar defaults
	hilosPorFase := input.HilosPorFase
	if hilosPorFase <= 0 {
		hilosPorFase = 1
	}

	numTuberias := input.NumTuberias
	if numTuberias <= 0 {
		numTuberias = 1
	}

	limiteCaida := input.PorcentajeCaidaMaximo
	if limiteCaida <= 0 {
		limiteCaida = 3.0
	}

	// Inicializar output
	output := dto.MemoriaOutput{
		TipoEquipo:       input.TipoEquipo,
		Clave:            input.Clave,
		Tension:          input.Tension.Valor(),
		FactorPotencia:   input.FactorPotencia,
		TipoCanalizacion: input.TipoCanalizacion,
		ITM:              input.ITM,
		LongitudCircuito: input.LongitudCircuito,
		FillFactor:       0.40,
		Estado:           input.Estado,
		SistemaElectrico: input.SistemaElectrico,
		HilosPorFase:     hilosPorFase,
	}

	// Obtener temperatura ambiente
	tempAmbiente, err := o.tablaRepo.ObtenerTemperaturaPorEstado(ctx, input.Estado)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("obtener temperatura para estado %s: %w", input.Estado, err)
	}
	output.TemperaturaAmbiente = tempAmbiente

	// Obtener cantidad de conductores
	cantidadConductores := input.SistemaElectrico.CantidadConductores()
	output.CantidadConductores = cantidadConductores

	// Paso 1: Corriente Nominal
	resultadoCorriente, err := o.calcularCorriente.Execute(ctx, input)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 1 - corriente nominal: %w", err)
	}
	output.CorrienteNominal = resultadoCorriente.CorrienteNominal

	// Convert to domain object for next step
	corrienteNominalVO, err := valueobject.NewCorriente(resultadoCorriente.CorrienteNominal)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("corriente nominal inválida: %w", err)
	}

	// Paso 2: Ajuste de Corriente
	resultadoAjuste, err := o.ajustarCorriente.Execute(
		ctx,
		corrienteNominalVO,
		input.Estado,
		input.ToEntityTipoCanalizacion(),
		input.SistemaElectrico.ToEntity(),
		input.ToEntityTipoEquipo(),
		hilosPorFase,
		numTuberias,
	)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 2 - ajustar corriente: %w", err)
	}

	// Convert DTO to domain objects for subsequent use cases
	corrienteAjustadaVO, err := valueobject.NewCorriente(resultadoAjuste.CorrienteAjustada)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("corriente ajustada inválida: %w", err)
	}
	temperaturaVO := valueobject.Temperatura(resultadoAjuste.Temperatura)

	output.CorrienteAjustada = resultadoAjuste.CorrienteAjustada
	output.CorrientePorHilo = resultadoAjuste.CorrienteAjustada / float64(hilosPorFase)
	output.FactorTemperaturaCalculado = resultadoAjuste.FactorTemperatura
	output.FactorAgrupamientoCalculado = resultadoAjuste.FactorAgrupamiento
	output.FactorTemperatura = resultadoAjuste.FactorTemperatura
	output.FactorAgrupamiento = resultadoAjuste.FactorAgrupamiento
	output.FactorTotalAjuste = resultadoAjuste.FactorTotal
	output.TemperaturaUsada = resultadoAjuste.Temperatura

	// Material
	material := input.Material
	output.Material = material.String()

	// Pasos 4-5: Seleccionar Conductores
	resultadoConductores, err := o.seleccionarConductor.Execute(
		ctx,
		corrienteAjustadaVO,
		hilosPorFase,
		input.ITM,
		material,
		temperaturaVO,
		input.ToEntityTipoCanalizacion(),
	)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("pasos 4-5 - seleccionar conductores: %w", err)
	}
	output.ConductorAlimentacion = dto.ResultadoConductor{
		Calibre:         resultadoConductores.Alimentacion.Calibre(),
		Material:        resultadoConductores.Alimentacion.Material().String(),
		SeccionMM2:      resultadoConductores.Alimentacion.SeccionMM2(),
		TipoAislamiento: resultadoConductores.Alimentacion.TipoAislamiento(),
		Capacidad:       resultadoConductores.Capacidad,
	}
	output.TablaAmpacidadUsada = resultadoConductores.TablaUsada

	output.ConductorTierra = dto.ResultadoConductor{
		Calibre:    resultadoConductores.Tierra.Calibre(),
		Material:   resultadoConductores.Tierra.Material().String(),
		SeccionMM2: resultadoConductores.Tierra.SeccionMM2(),
	}

	// Paso 6: Dimensionar Canalización
	resultadoCanalizacion, err := o.dimensionarCanalizacion.Execute(
		ctx,
		resultadoConductores.Alimentacion.SeccionMM2(),
		resultadoConductores.Tierra.SeccionMM2(),
		hilosPorFase,
		input.ToEntityTipoCanalizacion(),
	)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 6 - dimensionar canalización: %w", err)
	}
	output.Canalizacion = dto.ResultadoCanalizacion{
		Tamano:           resultadoCanalizacion.Tamano,
		AreaTotalMM2:     resultadoCanalizacion.AreaTotalMM2,
		AreaRequeridaMM2: resultadoCanalizacion.AreaTotalMM2 / 0.40,
		NumeroDeTubos:    resultadoCanalizacion.NumeroDeTubos,
	}

	// Paso 7: Caída de Tensión
	resultadoCaida, err := o.calcularCaidaTension.Execute(
		ctx,
		resultadoConductores.Alimentacion,
		corrienteAjustadaVO,
		input.LongitudCircuito,
		input.Tension,
		limiteCaida,
		input.ToEntityTipoCanalizacion(),
		input.FactorPotencia,
		hilosPorFase,
	)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 7 - caída de tensión: %w", err)
	}
	output.CaidaTension = dto.ResultadoCaidaTension{
		Porcentaje:          resultadoCaida.Porcentaje,
		CaidaVolts:          resultadoCaida.CaidaVolts,
		Cumple:              resultadoCaida.Cumple,
		LimitePorcentaje:    limiteCaida,
		ResistenciaEfectiva: resultadoCaida.ResistenciaEfectiva,
	}

	// Generar observaciones - básico por ahora
	output.CumpleNormativa = resultadoCaida.Cumple

	return output, nil
}
