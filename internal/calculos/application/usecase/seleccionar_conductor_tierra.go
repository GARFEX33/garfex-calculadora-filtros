// internal/calculos/application/usecase/seleccionar_conductor_tierra.go
package usecase

import (
	"context"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// SeleccionarConductorTierraUseCase executes the selection of ground conductor
// per NOM-250-122 based on ITM rating and conductor material.
type SeleccionarConductorTierraUseCase struct {
	tablaRepo port.TablaNOMRepository
}

// NewSeleccionarConductorTierraUseCase creates a new instance.
func NewSeleccionarConductorTierraUseCase(
	tablaRepo port.TablaNOMRepository,
) *SeleccionarConductorTierraUseCase {
	return &SeleccionarConductorTierraUseCase{
		tablaRepo: tablaRepo,
	}
}

// Execute selects the appropriate ground conductor from NOM-250-122 table.
func (uc *SeleccionarConductorTierraUseCase) Execute(
	ctx context.Context,
	itm int,
	material string,
) (dto.ConductorTierraOutput, error) {
	// Validate input
	input := dto.ConductorTierraInput{
		ITM:      itm,
		Material: material,
	}
	if err := input.Validate(); err != nil {
		return dto.ConductorTierraOutput{}, fmt.Errorf("validar input: %w", err)
	}

	// Get ground conductor table from repository
	tablaTierra, err := uc.tablaRepo.ObtenerTablaTierra(ctx)
	if err != nil {
		return dto.ConductorTierraOutput{}, fmt.Errorf("obtener tabla tierra: %w", err)
	}

	// Call domain service to select conductor
	domainMaterial := input.ToDomainMaterial()
	conductor, err := service.SeleccionarConductorTierra(itm, domainMaterial, tablaTierra)
	if err != nil {
		return dto.ConductorTierraOutput{}, fmt.Errorf("seleccionar conductor tierra: %w", err)
	}

	// Find ITMHasta from the matching table entry
	itmHasta := findITMHasta(conductor, domainMaterial, tablaTierra)

	// Map domain result to DTO
	return dto.ConductorTierraOutput{
		Calibre:    conductor.Calibre(),
		Material:   conductor.Material().String(),
		SeccionMM2: conductor.SeccionMM2(),
		ITMHasta:   itmHasta,
	}, nil
}

// findITMHasta finds the ITMHasta value from the matching table entry.
func findITMHasta(conductor valueobject.Conductor, material valueobject.MaterialConductor, tabla []valueobject.EntradaTablaTierra) int {
	for _, entrada := range tabla {
		if material == valueobject.MaterialAluminio && entrada.ConductorAl != nil {
			if conductor.Calibre() == entrada.ConductorAl.Calibre {
				return entrada.ITMHasta
			}
		}
		if conductor.Calibre() == entrada.ConductorCu.Calibre {
			return entrada.ITMHasta
		}
	}
	return 0
}
