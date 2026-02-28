// internal/calculos/application/usecase/calcular_charola_espaciado.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// CalcularCharolaEspaciadoUseCase calcula el dimensionamiento de charola con espaciado.
type CalcularCharolaEspaciadoUseCase struct {
	tablaRepo port.TablaNOMRepository
}

// NewCalcularCharolaEspaciadoUseCase creates a new instance.
func NewCalcularCharolaEspaciadoUseCase(tablaRepo port.TablaNOMRepository) *CalcularCharolaEspaciadoUseCase {
	return &CalcularCharolaEspaciadoUseCase{
		tablaRepo: tablaRepo,
	}
}

// Execute calcula el ancho requerido de charola para cables espaciados.
func (uc *CalcularCharolaEspaciadoUseCase) Execute(
	ctx context.Context,
	input dto.CharolaEspaciadoInput,
) (dto.CharolaEspaciadoOutput, error) {
	// 1. Validar input
	if err := input.Validate(); err != nil {
		return dto.CharolaEspaciadoOutput{}, fmt.Errorf("validación de entrada: %w", err)
	}

	// 2. Convertir primitivos a value objects
	conductorFase, err := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{
		DiametroMM: input.DiametroFaseMM,
	})
	if err != nil {
		return dto.CharolaEspaciadoOutput{}, fmt.Errorf("crear conductor fase: %w", err)
	}

	conductorTierra, err := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{
		DiametroMM: input.DiametroTierraMM,
	})
	if err != nil {
		return dto.CharolaEspaciadoOutput{}, fmt.Errorf("crear conductor tierra: %w", err)
	}

	// Parse sistema eléctrico
	sistema, err := entity.ParseSistemaElectrico(input.SistemaElectrico)
	if err != nil {
		return dto.CharolaEspaciadoOutput{}, fmt.Errorf("parsear sistema eléctrico: %w", err)
	}

	// Crear cable de control si se proporcionó
	var cablesControl []valueobject.CableControl
	if input.DiametroControlMM != nil && *input.DiametroControlMM > 0 {
		cableControl, err := valueobject.NewCableControl(valueobject.CableControlParams{
			Cantidad:   1,
			DiametroMM: *input.DiametroControlMM,
		})
		if err != nil {
			return dto.CharolaEspaciadoOutput{}, fmt.Errorf("crear cable control: %w", err)
		}
		cablesControl = append(cablesControl, cableControl)
	}

	// 3. Obtener tabla de charolas del repo
	tablaCharola, err := uc.tablaRepo.ObtenerTablaCharola(ctx, entity.TipoCanalizacionCharolaCableEspaciado)
	if err != nil {
		return dto.CharolaEspaciadoOutput{}, fmt.Errorf("obtener tabla charola: %w", err)
	}

	// 4. Llamar al servicio de dominio
	resultado, err := service.CalcularCharolaEspaciado(
		input.HilosPorFase,
		sistema,
		conductorFase,
		conductorTierra,
		tablaCharola,
		cablesControl,
	)
	if err != nil {
		return dto.CharolaEspaciadoOutput{}, fmt.Errorf("calcular charola espaciado: %w", err)
	}

	// 5. Convertir resultado domain a DTO output
	// El domain calcula internamente: numFases, tieneNeutro, hilosFaseTotal, totalHilos
	// y retorna el resultado con AnchoRequerido (que es el ancho total = EF + AF + EC + AC + tierra)
	// Para la memoria de cálculo necesitamos los valores intermedios que el domain no retorna,
	// así que los obtenemos directamente del input.
	// Nota: espacioControl = 1.0 * diametro coincide con la lógica del domain
	var espacioControl, anchoControl float64
	if input.DiametroControlMM != nil && *input.DiametroControlMM > 0 {
		espacioControl = *input.DiametroControlMM // 1.0 * diametro (coincide con domain)
		anchoControl = *input.DiametroControlMM   // diametro del cable (coincide con domain)
	}

	// Calcular hilos totales basado en sistema eléctrico - el domain hace lo mismo internamente
	numFases := 3 // default para Estrella
	tieneNeutro := true
	switch input.SistemaElectrico {
	case "MONOFASICO":
		numFases = 1
		tieneNeutro = true
	case "BIFASICO":
		numFases = 2
		tieneNeutro = true
	case "DELTA":
		numFases = 3
		tieneNeutro = false
	// ESTRELLA y otros: defaults arriba
	}

	hilosFaseTotal := numFases * input.HilosPorFase
	if tieneNeutro {
		hilosFaseTotal += input.HilosPorFase
	}

	// Valores de fuerza obtenidos directamente del input
	espacioFuerza := float64(hilosFaseTotal) * input.DiametroFaseMM
	anchoFuerza := espacioFuerza

	out := dto.CharolaEspaciadoOutput{
		Tipo:             string(resultado.Tipo),
		Tamano:           resultado.Tamano,
		TamanoPulgadas:   resultado.Tamano + "\"",
		AnchoRequerido:   resultado.AnchoRequerido,
		AnchoComercialMM: resultado.AnchoComercialMM,
		DiametroFaseMM:   input.DiametroFaseMM,
		DiametroTierraMM: input.DiametroTierraMM,
		NumHilosTotal:    hilosFaseTotal,
		EspacioFuerzaMM:  espacioFuerza,
		AnchoFuerzaMM:    anchoFuerza,
		EspacioControlMM: espacioControl,
		AnchoControlMM:   anchoControl,
		AnchoTierraMM:    input.DiametroTierraMM,
	}
	if input.DiametroControlMM != nil && *input.DiametroControlMM > 0 {
		out.DiametroControlMM = input.DiametroControlMM
	}
	return out, nil
}
