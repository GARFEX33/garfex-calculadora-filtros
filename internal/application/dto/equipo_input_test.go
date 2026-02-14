// internal/application/dto/equipo_input_test.go
package dto_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/application/dto"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tension220(t *testing.T) valueobject.Tension {
	t.Helper()
	v, err := valueobject.NewTension(220)
	require.NoError(t, err)
	return v
}

func inputBase(t *testing.T) dto.EquipoInput {
	t.Helper()
	return dto.EquipoInput{
		Modo:             dto.ModoManualAmperaje,
		AmperajeNominal:  100,
		Tension:          tension220(t),
		ITM:              125,
		Estado:           "Nuevo Leon",
		SistemaElectrico: dto.SistemaElectricoDelta,
	}
}

// --- ModoListado ---

func TestEquipoInput_Validate_ModoListado_ClavePresente(t *testing.T) {
	input := inputBase(t)
	input.Modo = dto.ModoListado
	input.Clave = "FA-001"
	assert.NoError(t, input.Validate())
}

func TestEquipoInput_Validate_ModoListado_SinClave(t *testing.T) {
	input := inputBase(t)
	input.Modo = dto.ModoListado
	input.Clave = ""
	assert.ErrorIs(t, input.Validate(), dto.ErrEquipoInputInvalido)
}

// --- ModoManualAmperaje ---

func TestEquipoInput_Validate_ModoManualAmperaje_Valido(t *testing.T) {
	assert.NoError(t, inputBase(t).Validate())
}

func TestEquipoInput_Validate_ModoManualAmperaje_AmperajeCero(t *testing.T) {
	input := inputBase(t)
	input.AmperajeNominal = 0
	assert.ErrorIs(t, input.Validate(), dto.ErrEquipoInputInvalido)
}

func TestEquipoInput_Validate_ModoManualAmperaje_AmperajeNegativo(t *testing.T) {
	input := inputBase(t)
	input.AmperajeNominal = -10
	assert.ErrorIs(t, input.Validate(), dto.ErrEquipoInputInvalido)
}

// --- ModoManualPotencia ---

func TestEquipoInput_Validate_ModoManualPotencia_Valido(t *testing.T) {
	input := inputBase(t)
	input.Modo = dto.ModoManualPotencia
	input.PotenciaNominal = 75
	assert.NoError(t, input.Validate())
}

func TestEquipoInput_Validate_ModoManualPotencia_PotenciaCero(t *testing.T) {
	input := inputBase(t)
	input.Modo = dto.ModoManualPotencia
	input.PotenciaNominal = 0
	assert.ErrorIs(t, input.Validate(), dto.ErrEquipoInputInvalido)
}

// --- Modo inválido ---

func TestEquipoInput_Validate_ModoInvalido(t *testing.T) {
	input := inputBase(t)
	input.Modo = dto.ModoCalculo("INEXISTENTE")
	assert.ErrorIs(t, input.Validate(), dto.ErrModoInvalido)
}

// --- Tensión ---

func TestEquipoInput_Validate_TensionCero(t *testing.T) {
	input := inputBase(t)
	input.Tension = valueobject.Tension{} // zero value → Valor() == 0
	assert.ErrorIs(t, input.Validate(), dto.ErrEquipoInputInvalido)
}

// --- ITM ---

func TestEquipoInput_Validate_ITMCero(t *testing.T) {
	input := inputBase(t)
	input.ITM = 0
	assert.ErrorIs(t, input.Validate(), dto.ErrEquipoInputInvalido)
}

func TestEquipoInput_Validate_ITMNegativo(t *testing.T) {
	input := inputBase(t)
	input.ITM = -1
	assert.ErrorIs(t, input.Validate(), dto.ErrEquipoInputInvalido)
}

// --- Estado ---

func TestEquipoInput_Validate_EstadoVacio(t *testing.T) {
	input := inputBase(t)
	input.Estado = ""
	assert.ErrorIs(t, input.Validate(), dto.ErrEquipoInputInvalido)
}

// --- SistemaElectrico ---

func TestEquipoInput_Validate_SistemaElectricoInvalido(t *testing.T) {
	input := inputBase(t)
	input.SistemaElectrico = dto.SistemaElectrico("INVALIDO")
	assert.ErrorIs(t, input.Validate(), entity.ErrSistemaElectricoInvalido)
}

func TestEquipoInput_Validate_TodosSistemas(t *testing.T) {
	sistemas := []dto.SistemaElectrico{
		dto.SistemaElectricoDelta,
		dto.SistemaElectricoEstrella,
		dto.SistemaElectricoBifasico,
		dto.SistemaElectricoMonofasico,
	}
	for _, s := range sistemas {
		t.Run(string(s), func(t *testing.T) {
			input := inputBase(t)
			input.SistemaElectrico = s
			assert.NoError(t, input.Validate())
		})
	}
}
