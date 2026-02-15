// internal/calculos/domain/service/calcular_factor_temperatura_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestCalcularFactorTemperatura(t *testing.T) {
	// Tabla de test
	tablaTest := []service.EntradaTablaFactorTemperatura{
		{RangoTempC: "10-15", Factor60C: 1.20, Factor75C: 1.15, Factor90C: 1.12},
		{RangoTempC: "16-20", Factor60C: 1.15, Factor75C: 1.11, Factor90C: 1.09},
		{RangoTempC: "21-25", Factor60C: 1.10, Factor75C: 1.07, Factor90C: 1.05},
		{RangoTempC: "26-30", Factor60C: 1.05, Factor75C: 1.03, Factor90C: 1.02},
		{RangoTempC: "31-35", Factor60C: 1.00, Factor75C: 1.00, Factor90C: 1.00},
		{RangoTempC: "36-40", Factor60C: 0.94, Factor75C: 0.95, Factor90C: 0.96},
		{RangoTempC: "41-45", Factor60C: 0.88, Factor75C: 0.90, Factor90C: 0.91},
		{RangoTempC: "46-50", Factor60C: 0.82, Factor75C: 0.85, Factor90C: 0.87},
		{RangoTempC: "51-55", Factor60C: 0.75, Factor75C: 0.80, Factor90C: 0.82},
		{RangoTempC: "56-60", Factor60C: 0.67, Factor75C: 0.74, Factor90C: 0.77},
		{RangoTempC: "61-70", Factor60C: 0.58, Factor75C: 0.67, Factor90C: 0.71},
		{RangoTempC: "71-80", Factor60C: 0.47, Factor75C: 0.58, Factor90C: 0.63},
	}

	tests := []struct {
		name           string
		tempAmbiente   int
		tempConductor  valueobject.Temperatura
		expectedFactor float64
		wantErr        bool
	}{
		{"21°C + 75C conductor", 21, valueobject.Temp75, 1.07, false},
		{"26°C + 75C conductor", 26, valueobject.Temp75, 1.03, false},
		{"36°C + 75C conductor", 36, valueobject.Temp75, 0.95, false},
		{"21°C + 60C conductor", 21, valueobject.Temp60, 1.10, false},
		{"31°C exacto + 75C", 31, valueobject.Temp75, 1.00, false},
		{"Temperatura negativa", -5, valueobject.Temp75, 0, true},
		{"40°C + 90C conductor", 40, valueobject.Temp90, 0.96, false},
		{"45°C + 60C conductor", 45, valueobject.Temp60, 0.88, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factor, err := service.CalcularFactorTemperatura(tt.tempAmbiente, tt.tempConductor, tablaTest)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.InDelta(t, tt.expectedFactor, factor, 0.001)
		})
	}
}
