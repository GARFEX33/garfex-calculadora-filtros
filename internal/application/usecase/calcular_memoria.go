// internal/application/usecase/calcular_memoria.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/application/dto"
	"github.com/garfex/calculadora-filtros/internal/application/port"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// CalcularMemoriaUseCase orquesta el cálculo completo de una memoria de cálculo.
// Sigue el flujo de 7 pasos definido en la documentación.
type CalcularMemoriaUseCase struct {
	tablaRepo  port.TablaNOMRepository
	equipoRepo port.EquipoRepository
}

// NewCalcularMemoriaUseCase crea una nueva instancia del use case.
func NewCalcularMemoriaUseCase(
	tablaRepo port.TablaNOMRepository,
	equipoRepo port.EquipoRepository,
) *CalcularMemoriaUseCase {
	return &CalcularMemoriaUseCase{
		tablaRepo:  tablaRepo,
		equipoRepo: equipoRepo,
	}
}

// Execute ejecuta el cálculo completo siguiendo los 7 pasos.
func (uc *CalcularMemoriaUseCase) Execute(ctx context.Context, input dto.EquipoInput) (dto.MemoriaOutput, error) {
	// Validar input
	if err := input.Validate(); err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("validación de input: %w", err)
	}

	output := dto.MemoriaOutput{
		TipoEquipo:       input.TipoEquipo,
		Clave:            input.Clave,
		Tension:          input.Tension,
		FactorPotencia:   input.FactorPotencia,
		TipoCanalizacion: input.TipoCanalizacion,
		ITM:              input.ITM,
		LongitudCircuito: input.LongitudCircuito,
		FillFactor:       0.40, // 40% para tubería según NOM
	}

	// Configurar defaults
	hilosPorFase := input.HilosPorFase
	if hilosPorFase <= 0 {
		hilosPorFase = 1
	}
	output.HilosPorFase = hilosPorFase

	factorAgrupamiento := input.FactorAgrupamiento
	if factorAgrupamiento <= 0 {
		factorAgrupamiento = 1.0
	}

	factorTemperatura := input.FactorTemperatura
	if factorTemperatura <= 0 {
		factorTemperatura = 1.0
	}

	output.FactorAgrupamiento = factorAgrupamiento
	output.FactorTemperatura = factorTemperatura
	output.FactorTotalAjuste = factorAgrupamiento * factorTemperatura

	limiteCaida := input.PorcentajeCaidaMaximo
	if limiteCaida <= 0 {
		limiteCaida = 3.0
	}

	// Preparar factores para AjustarCorriente
	factores := map[string]float64{
		"agrupamiento": factorAgrupamiento,
		"temperatura":  factorTemperatura,
	}

	// ============================================================================
	// PASO 1: Calcular Corriente Nominal
	// ============================================================================
	corrienteNominal, err := uc.calcularCorrienteNominal(ctx, input)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 1 - corriente nominal: %w", err)
	}
	output.CorrienteNominal = corrienteNominal

	// ============================================================================
	// PASO 2: Ajustar Corriente
	// ============================================================================
	corrienteAjustada, err := service.AjustarCorriente(corrienteNominal, factores)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 2 - ajustar corriente: %w", err)
	}
	output.CorrienteAjustada = corrienteAjustada
	output.CorrientePorHilo = corrienteAjustada.Valor() / float64(hilosPorFase)

	// ============================================================================
	// PASO 3: Seleccionar TipoCanalizacion (ya viene en input)
	// ============================================================================
	// La canalización ya está en el input, determina qué tabla NOM usar

	// ============================================================================
	// PASO 4: Seleccionar Conductor de Alimentación
	// ============================================================================
	// Determinar temperatura según reglas del AGENTS.md
	temperatura := uc.seleccionarTemperatura(corrienteAjustada, input)
	output.TemperaturaUsada = temperatura

	// Determinar material (usamos cobre por defecto, podría venir en input)
	material := valueobject.MaterialCobre

	// Obtener tabla de ampacidad
	tablaAmpacidad, err := uc.tablaRepo.ObtenerTablaAmpacidad(ctx, input.TipoCanalizacion, material, temperatura)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 4 - obtener tabla ampacidad: %w", err)
	}

	// Seleccionar conductor
	conductor, err := service.SeleccionarConductorAlimentacion(corrienteAjustada, hilosPorFase, tablaAmpacidad)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 4 - seleccionar conductor: %w", err)
	}

	output.ConductorAlimentacion = dto.ResultadoConductor{
		Calibre:         conductor.Calibre(),
		Material:        conductor.Material(),
		SeccionMM2:      conductor.SeccionMM2(),
		TipoAislamiento: conductor.TipoAislamiento(),
		Capacidad:       uc.buscarCapacidadEnTabla(tablaAmpacidad, conductor.Calibre()),
	}

	// ============================================================================
	// PASO 5: Seleccionar Conductor de Tierra
	// ============================================================================
	tablaTierra, err := uc.tablaRepo.ObtenerTablaTierra(ctx)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 5 - obtener tabla tierra: %w", err)
	}

	conductorTierra, err := service.SeleccionarConductorTierra(input.ITM, tablaTierra)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 5 - seleccionar conductor tierra: %w", err)
	}

	output.ConductorTierra = dto.ResultadoConductor{
		Calibre:    conductorTierra.Calibre(),
		Material:   conductorTierra.Material(),
		SeccionMM2: conductorTierra.SeccionMM2(),
	}

	// ============================================================================
	// PASO 6: Dimensionar Canalización
	// ============================================================================
	tablaCanalizacion, err := uc.tablaRepo.ObtenerTablaCanalizacion(ctx, input.TipoCanalizacion)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 6 - obtener tabla canalización: %w", err)
	}

	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: hilosPorFase * 3, SeccionMM2: conductor.SeccionMM2()}, // Fases
		{Cantidad: 1, SeccionMM2: conductorTierra.SeccionMM2()},          // Tierra
	}

	canalizacion, err := service.CalcularCanalizacion(conductores, string(input.TipoCanalizacion), tablaCanalizacion)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 6 - calcular canalización: %w", err)
	}

	output.Canalizacion = dto.ResultadoCanalizacion{
		Tipo:             input.TipoCanalizacion,
		Tamano:           canalizacion.Tamano,
		AreaTotalMM2:     canalizacion.AreaTotal,
		AreaRequeridaMM2: canalizacion.AreaTotal / 0.40,
	}

	// ============================================================================
	// PASO 7: Calcular Caída de Tensión
	// ============================================================================
	impedancia, err := uc.tablaRepo.ObtenerImpedancia(ctx, conductor.Calibre(), input.TipoCanalizacion, material)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 7 - obtener impedancia: %w", err)
	}

	entradaCaida := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: impedancia.R,
		ReactanciaOhmPorKm:  impedancia.X,
		TipoCanalizacion:    input.TipoCanalizacion,
		HilosPorFase:        hilosPorFase,
		FactorPotencia:      input.FactorPotencia,
	}

	resultadoCaida, err := service.CalcularCaidaTension(
		entradaCaida,
		corrienteAjustada,
		input.LongitudCircuito,
		input.Tension,
		limiteCaida,
	)
	if err != nil {
		return dto.MemoriaOutput{}, fmt.Errorf("paso 7 - calcular caída de tensión: %w", err)
	}

	output.CaidaTension = dto.ResultadoCaidaTension{
		Porcentaje:          resultadoCaida.Porcentaje,
		CaidaVolts:          resultadoCaida.CaidaVolts,
		Cumple:              resultadoCaida.Cumple,
		LimitePorcentaje:    limiteCaida,
		ResistenciaEfectiva: resultadoCaida.Impedancia,
	}

	// ============================================================================
	// Resumen de cumplimiento
	// ============================================================================
	output.CumpleNormativa = resultadoCaida.Cumple
	output.Observaciones = uc.generarObservaciones(output)

	return output, nil
}

// calcularCorrienteNominal calcula la corriente nominal según el modo.
func (uc *CalcularMemoriaUseCase) calcularCorrienteNominal(ctx context.Context, input dto.EquipoInput) (valueobject.Corriente, error) {
	switch input.Modo {
	case dto.ModoListado:
		// Buscar equipo en BD
		equipo, err := uc.equipoRepo.BuscarPorClave(ctx, input.Clave)
		if err != nil {
			return valueobject.Corriente{}, fmt.Errorf("buscar equipo: %w", err)
		}
		return service.CalcularCorrienteNominal(equipo)

	case dto.ModoManualAmperaje:
		// Usar el amperaje directo proporcionado
		return valueobject.NewCorriente(input.AmperajeNominal)

	case dto.ModoManualPotencia:
		// Para modo manual potencia, necesitamos crear un equipo apropiado
		// Esto requiere más información según el tipo
		return valueobject.Corriente{}, fmt.Errorf("modo MANUAL_POTENCIA requiere implementación adicional")

	default:
		return valueobject.Corriente{}, dto.ErrModoInvalido
	}
}

// seleccionarTemperatura determina la temperatura según las reglas del AGENTS.md.
func (uc *CalcularMemoriaUseCase) seleccionarTemperatura(corriente valueobject.Corriente, input dto.EquipoInput) valueobject.Temperatura {
	// Si hay override explícito, usarlo
	if input.TemperaturaOverride != nil {
		return *input.TemperaturaOverride
	}

	// Reglas según AGENTS.md de application
	if corriente.Valor() <= 100 {
		// <= 100A -> 60C (o 75C si charola triangular sin columna 60C)
		if input.TipoCanalizacion == entity.TipoCanalizacionCharolaCableTriangular {
			return valueobject.Temp75
		}
		return valueobject.Temp60
	}

	// > 100A -> 75C
	return valueobject.Temp75
}

// buscarCapacidadEnTabla busca la capacidad de un calibre en la tabla.
func (uc *CalcularMemoriaUseCase) buscarCapacidadEnTabla(tabla []valueobject.EntradaTablaConductor, calibre string) float64 {
	for _, entrada := range tabla {
		if entrada.Conductor.Calibre == calibre {
			return entrada.Capacidad
		}
	}
	return 0
}

// generarObservaciones genera observaciones sobre el cálculo.
func (uc *CalcularMemoriaUseCase) generarObservaciones(output dto.MemoriaOutput) []string {
	var obs []string

	if !output.CaidaTension.Cumple {
		obs = append(obs, fmt.Sprintf(
			"Caída de tensión %.2f%% excede el límite de %.2f%%",
			output.CaidaTension.Porcentaje,
			output.CaidaTension.LimitePorcentaje,
		))
	}

	if output.HilosPorFase > 1 {
		obs = append(obs, fmt.Sprintf(
			"Se usan %d hilos por fase en paralelo",
			output.HilosPorFase,
		))
	}

	return obs
}
