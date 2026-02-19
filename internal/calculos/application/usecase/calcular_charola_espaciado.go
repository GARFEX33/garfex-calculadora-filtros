// internal/calculos/application/usecase/calcular_charola_espaciado.go
package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	return dto.CharolaEspaciadoOutput{
		Tipo:           string(resultado.Tipo),
		Tamano:         resultado.Tamano,
		TamanoPulgadas: convertirTamanoAPulgadas(resultado.Tamano),
		AnchoRequerido: resultado.AnchoRequerido,
	}, nil
}

// convertirTamanoAPulgadas convierte el tamaño de mm a pulgadas.
// Ejemplo: "300mm" -> "12\""
func convertirTamanoAPulgadas(tamano string) string {
	// Remover "mm" del tamaño
	tamano = strings.TrimSuffix(tamano, "mm")
	mm, err := strconv.ParseFloat(tamano, 64)
	if err != nil {
		return tamano // Si no se puede convertir, devolver el original
	}
	pulgadas := mm / 25.4
	return fmt.Sprintf("%.2f\"", pulgadas)
}
