// internal/calculos/application/dto/equipo_input_test.go
package dto_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/stretchr/testify/assert"
)

func tension220(t *testing.T) float64 {
	t.Helper()
	return 220.0
}

func inputBaseManualAmperaje(t *testing.T) dto.EquipoInput {
	t.Helper()
	return dto.EquipoInput{
		Modo:             dto.ModoManualAmperaje,
		TipoEquipo:       "FILTRO_ACTIVO",
		AmperajeNominal:  100,
		Tension:          tension220(t),
		SistemaElectrico: dto.SistemaElectricoDelta,
		Estado:           "Nuevo Leon",
	}
}

func inputBaseListado(t *testing.T) dto.EquipoInput {
	t.Helper()
	return dto.EquipoInput{
		Modo: dto.ModoListado,
		Equipo: dto.DatosEquipo{
			Clave:    "TEST-001",
			Tipo:     dto.TipoFiltroA,
			Voltaje:  480,
			Amperaje: 100,
			ITM:      125,
		},
		Tension:          tension220(t),
		SistemaElectrico: dto.SistemaElectricoDelta,
		Estado:           "Nuevo Leon",
	}
}

// --- ModoListado ---

func TestEquipoInput_Validate_ModoListado_Valido(t *testing.T) {
	input := inputBaseListado(t)
	assert.NoError(t, input.Validate())
}

func TestEquipoInput_Validate_ModoListado_TipoVacio(t *testing.T) {
	input := inputBaseListado(t)
	input.Equipo.Tipo = ""
	assert.Error(t, input.Validate())
}

func TestEquipoInput_Validate_ModoListado_VoltajeCero(t *testing.T) {
	input := inputBaseListado(t)
	input.Equipo.Voltaje = 0
	assert.Error(t, input.Validate())
}

func TestEquipoInput_Validate_ModoListado_AmperajeCero(t *testing.T) {
	input := inputBaseListado(t)
	input.Equipo.Amperaje = 0
	assert.Error(t, input.Validate())
}

func TestEquipoInput_Validate_ModoListado_ITMCero(t *testing.T) {
	input := inputBaseListado(t)
	input.Equipo.ITM = 0
	assert.Error(t, input.Validate())
}

// --- ModoManualAmperaje ---

func TestEquipoInput_Validate_ModoManualAmperaje_Valido(t *testing.T) {
	assert.NoError(t, inputBaseManualAmperaje(t).Validate())
}

func TestEquipoInput_Validate_ModoManualAmperaje_AmperajeCero(t *testing.T) {
	input := inputBaseManualAmperaje(t)
	input.AmperajeNominal = 0
	assert.Error(t, input.Validate())
}

func TestEquipoInput_Validate_ModoManualAmperaje_SinTipoEquipo(t *testing.T) {
	input := inputBaseManualAmperaje(t)
	input.TipoEquipo = ""
	assert.Error(t, input.Validate())
}

// --- ModoManualPotencia ---

func TestEquipoInput_Validate_ModoManualPotencia_Valido(t *testing.T) {
	input := inputBaseManualAmperaje(t)
	input.Modo = dto.ModoManualPotencia
	input.PotenciaNominal = 75
	assert.NoError(t, input.Validate())
}

func TestEquipoInput_Validate_ModoManualPotencia_PotenciaCero(t *testing.T) {
	input := inputBaseManualAmperaje(t)
	input.Modo = dto.ModoManualPotencia
	input.PotenciaNominal = 0
	assert.Error(t, input.Validate())
}

// --- Modo inválido ---

func TestEquipoInput_Validate_ModoInvalido(t *testing.T) {
	input := inputBaseManualAmperaje(t)
	input.Modo = dto.ModoCalculo("INEXISTENTE")
	assert.ErrorIs(t, input.Validate(), dto.ErrModoInvalido)
}

// --- Tensión ---

func TestEquipoInput_Validate_TensionCero(t *testing.T) {
	input := inputBaseManualAmperaje(t)
	input.Tension = 0.0
	assert.Error(t, input.Validate())
}

// --- Estado ---

func TestEquipoInput_Validate_EstadoVacio(t *testing.T) {
	input := inputBaseManualAmperaje(t)
	input.Estado = ""
	assert.Error(t, input.Validate())
}

// --- SistemaElectrico ---

func TestEquipoInput_Validate_SistemaElectricoInvalido(t *testing.T) {
	input := inputBaseManualAmperaje(t)
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
			input := inputBaseManualAmperaje(t)
			input.SistemaElectrico = s
			assert.NoError(t, input.Validate())
		})
	}
}

// --- TipoFiltro mapeo ---

func TestTipoFiltro_ToTipoEquipo(t *testing.T) {
	tests := []struct {
		tipo       dto.TipoFiltro
		esperado   entity.TipoEquipo
		debeFallar bool
	}{
		{dto.TipoFiltroA, entity.TipoEquipoFiltroActivo, false},
		{dto.TipoFiltroKVA, entity.TipoEquipoTransformador, false},
		{dto.TipoFiltroKVAR, entity.TipoEquipoFiltroRechazo, false},
		{dto.TipoFiltro("INVALIDO"), "", true},
	}

	for _, tt := range tests {
		t.Run(string(tt.tipo), func(t *testing.T) {
			resultado, err := tt.tipo.ToTipoEquipo()
			if tt.debeFallar {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.esperado, resultado)
			}
		})
	}
}

// --- GetTipoEquipo ---

func TestEquipoInput_GetTipoEquipo_ModoListado(t *testing.T) {
	input := inputBaseListado(t)
	tipo, err := input.GetTipoEquipo()
	assert.NoError(t, err)
	assert.Equal(t, entity.TipoEquipoFiltroActivo, tipo)
}

func TestEquipoInput_GetTipoEquipo_ModoManual(t *testing.T) {
	input := inputBaseManualAmperaje(t)
	tipo, err := input.GetTipoEquipo()
	assert.NoError(t, err)
	assert.Equal(t, entity.TipoEquipoFiltroActivo, tipo)
}
