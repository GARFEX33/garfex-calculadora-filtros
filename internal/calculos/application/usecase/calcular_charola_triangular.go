// internal/calculos/application/usecase/calcular_charola_triangular.go
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

// CalcularCharolaTriangularUseCase calcula el dimensionamiento de charola con configuración triangular.
type CalcularCharolaTriangularUseCase struct {
	tablaRepo port.TablaNOMRepository
}

// NewCalcularCharolaTriangularUseCase creates a new instance.
func NewCalcularCharolaTriangularUseCase(tablaRepo port.TablaNOMRepository) *CalcularCharolaTriangularUseCase {
	return &CalcularCharolaTriangularUseCase{
		tablaRepo: tablaRepo,
	}
}

// Execute calcula el ancho requerido de charola para cables en configuración triangular.
func (uc *CalcularCharolaTriangularUseCase) Execute(
	ctx context.Context,
	input dto.CharolaTriangularInput,
) (dto.CharolaTriangularOutput, error) {
	// 1. Validar input
	if err := input.Validate(); err != nil {
		return dto.CharolaTriangularOutput{}, fmt.Errorf("validación de entrada: %w", err)
	}

	// 2. Convertir primitivos a value objects
	conductorFase, err := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{
		DiametroMM: input.DiametroFaseMM,
	})
	if err != nil {
		return dto.CharolaTriangularOutput{}, fmt.Errorf("crear conductor fase: %w", err)
	}

	conductorTierra, err := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{
		DiametroMM: input.DiametroTierraMM,
	})
	if err != nil {
		return dto.CharolaTriangularOutput{}, fmt.Errorf("crear conductor tierra: %w", err)
	}

	// Crear cable de control si se proporcionó
	var cablesControl []valueobject.CableControl
	if input.DiametroControlMM != nil && *input.DiametroControlMM > 0 {
		cableControl, err := valueobject.NewCableControl(valueobject.CableControlParams{
			Cantidad:   1,
			DiametroMM: *input.DiametroControlMM,
		})
		if err != nil {
			return dto.CharolaTriangularOutput{}, fmt.Errorf("crear cable control: %w", err)
		}
		cablesControl = append(cablesControl, cableControl)
	}

	// 3. Obtener tabla de charolas del repo
	tablaCharola, err := uc.tablaRepo.ObtenerTablaCharola(ctx, entity.TipoCanalizacionCharolaCableTriangular)
	if err != nil {
		return dto.CharolaTriangularOutput{}, fmt.Errorf("obtener tabla charola: %w", err)
	}

	// 4. Llamar al servicio de dominio
	resultado, err := service.CalcularCharolaTriangular(
		input.HilosPorFase,
		conductorFase,
		conductorTierra,
		tablaCharola,
		cablesControl,
	)
	if err != nil {
		return dto.CharolaTriangularOutput{}, fmt.Errorf("calcular charola triangular: %w", err)
	}

	// 5. Convertir resultado domain a DTO output con valores intermedios para la memoria
	// El domain calcula: espacioControl = 1.0 * diametro, anchoControl = diametro
	// Usamos directamente los valores del input que coinciden con la lógica del domain
	var espacioControl, anchoControl float64
	if input.DiametroControlMM != nil && *input.DiametroControlMM > 0 {
		espacioControl = *input.DiametroControlMM // 1.0 * diametro (coincide con domain)
		anchoControl = *input.DiametroControlMM    // diametro del cable (coincide con domain)
	}

	// Valores de fuerza obtenidos directamente del input (el domain no los retorna)
	anchoPotencia := 2.0 * input.DiametroFaseMM * float64(input.HilosPorFase)
	espacioFuerza := float64(input.HilosPorFase-1) * 2.15 * input.DiametroFaseMM
	factorTriangular := 2.15

	out := dto.CharolaTriangularOutput{
		Tipo:             string(resultado.Tipo),
		Tamano:           resultado.Tamano,
		TamanoPulgadas:   resultado.Tamano + "\"",
		AnchoRequerido:   resultado.AnchoRequerido,
		AnchoComercialMM: resultado.AnchoComercialMM,
		DiametroFaseMM:   input.DiametroFaseMM,
		DiametroTierraMM: input.DiametroTierraMM,
		AnchoPotenciaMM:  anchoPotencia,
		EspacioFuerzaMM:  espacioFuerza,
		EspacioControlMM: espacioControl,
		AnchoControlMM:   anchoControl,
		AnchoTierraMM:    input.DiametroTierraMM,
		FactorTriangular: factorTriangular,
		FactorControl:    1.0,
	}
	if input.DiametroControlMM != nil && *input.DiametroControlMM > 0 {
		out.DiametroControlMM = input.DiametroControlMM
	}
	return out, nil
}
