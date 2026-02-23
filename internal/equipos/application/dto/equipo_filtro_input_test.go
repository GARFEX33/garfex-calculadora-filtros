// internal/equipos/application/dto/equipo_filtro_input_test.go
package dto_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/equipos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/equipos/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEquipoInput_Validate(t *testing.T) {
	clave := "FA-480-100"
	bornes := 6

	validInput := dto.CreateEquipoInput{
		Clave:    &clave,
		Tipo:     "A",
		Voltaje:  480,
		Amperaje: 100,
		ITM:      125,
		Bornes:   &bornes,
	}

	t.Run("input válido pasa validación", func(t *testing.T) {
		err := validInput.Validate()
		assert.NoError(t, err)
	})

	t.Run("tipo vacío falla", func(t *testing.T) {
		input := validInput
		input.Tipo = ""
		err := input.Validate()
		assert.ErrorIs(t, err, dto.ErrInputInvalido)
	})

	t.Run("tipo inválido falla", func(t *testing.T) {
		input := validInput
		input.Tipo = "FILTRO_ACTIVO"
		err := input.Validate()
		assert.ErrorIs(t, err, dto.ErrInputInvalido)
	})

	t.Run("voltaje cero falla", func(t *testing.T) {
		input := validInput
		input.Voltaje = 0
		err := input.Validate()
		assert.ErrorIs(t, err, dto.ErrInputInvalido)
	})

	t.Run("amperaje cero falla", func(t *testing.T) {
		input := validInput
		input.Amperaje = 0
		err := input.Validate()
		assert.ErrorIs(t, err, dto.ErrInputInvalido)
	})

	t.Run("ITM cero falla", func(t *testing.T) {
		input := validInput
		input.ITM = 0
		err := input.Validate()
		assert.ErrorIs(t, err, dto.ErrInputInvalido)
	})
}

func TestCreateEquipoInput_ToDomain(t *testing.T) {
	clave := "FR-220-50"

	t.Run("convierte correctamente a domain entity", func(t *testing.T) {
		input := dto.CreateEquipoInput{
			Clave:    &clave,
			Tipo:     "KVAR",
			Voltaje:  220,
			Amperaje: 50,
			ITM:      60,
			Bornes:   nil,
		}

		eq, err := input.ToDomain()
		require.NoError(t, err)
		assert.Equal(t, &clave, eq.Clave)
		assert.Equal(t, entity.TipoFiltroKVAR, eq.Tipo)
		assert.Equal(t, 220, eq.Voltaje)
		assert.Equal(t, 50, eq.Amperaje)
		assert.Equal(t, 60, eq.ITM)
		assert.Nil(t, eq.Bornes)
	})
}
