package service_test

import (
	"fmt"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/stretchr/testify/assert"
)

func TestCalcularFactorAgrupamiento(t *testing.T) {
	tablaTest := []service.EntradaTablaFactorAgrupamiento{
		{CantidadMin: 1, CantidadMax: 1, Factor: 1.00},
		{CantidadMin: 2, CantidadMax: 2, Factor: 0.80},
		{CantidadMin: 3, CantidadMax: 3, Factor: 0.70},
		{CantidadMin: 4, CantidadMax: 4, Factor: 0.65},
		{CantidadMin: 5, CantidadMax: 6, Factor: 0.60},
		{CantidadMin: 7, CantidadMax: 9, Factor: 0.50},
		{CantidadMin: 10, CantidadMax: 20, Factor: 0.45},
		{CantidadMin: 21, CantidadMax: 30, Factor: 0.40},
		{CantidadMin: 31, CantidadMax: 40, Factor: 0.35},
		{CantidadMin: 41, CantidadMax: -1, Factor: 0.30},
	}

	tests := []struct {
		cantidad       int
		expectedFactor float64
		wantErr        bool
	}{
		{1, 1.00, false},
		{2, 0.80, false},
		{3, 0.70, false},
		{4, 0.65, false},
		{5, 0.60, false},
		{6, 0.60, false},
		{7, 0.50, false},
		{10, 0.45, false},
		{21, 0.40, false},
		{31, 0.35, false},
		{41, 0.30, false},
		{50, 0.30, false},
		{0, 0, true},
		{-1, 0, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d_conductores", tt.cantidad), func(t *testing.T) {
			factor, err := service.CalcularFactorAgrupamiento(tt.cantidad, tablaTest)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.InDelta(t, tt.expectedFactor, factor, 0.001)
		})
	}
}
