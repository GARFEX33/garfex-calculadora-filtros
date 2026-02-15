// internal/calculos/domain/entity/itm_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewITM_Valid(t *testing.T) {
	itm, err := entity.NewITM(125, 3, 3, 480)
	require.NoError(t, err)

	assert.Equal(t, 125, itm.Amperaje)
	assert.Equal(t, 3, itm.Polos)
	assert.Equal(t, 3, itm.Bornes)
	assert.Equal(t, 480, itm.Voltaje)
}

func TestNewITM_Invalido(t *testing.T) {
	tests := []struct {
		name                             string
		amperaje, polos, bornes, voltaje int
	}{
		{"amperaje cero", 0, 3, 3, 480},
		{"amperaje negativo", -1, 3, 3, 480},
		{"polos cero", 125, 0, 3, 480},
		{"bornes cero", 125, 3, 0, 480},
		{"voltaje cero", 125, 3, 3, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entity.NewITM(tt.amperaje, tt.polos, tt.bornes, tt.voltaje)
			assert.Error(t, err)
		})
	}
}
