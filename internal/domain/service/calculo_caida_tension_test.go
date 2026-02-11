package service_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcularCaidaTension(t *testing.T) {
	tests := []struct {
		name           string
		corrienteA     float64
		distanciaM     float64
		calibre        string
		material       string
		seccionMM2     float64
		voltaje        int
		limitePorc     float64
		expectedPorc   float64
		expectedCumple bool
	}{
		{
			name:       "Cu 2AWG 30m at 120A 480V within 3%",
			corrienteA: 120,
			distanciaM: 30,
			calibre:    "2 AWG",
			material:   "Cu",
			seccionMM2: 33.62,
			voltaje:    480,
			limitePorc: 3.0,
			// VD = (√3 × 0.01724 × 30 × 120) / 33.62 = (107.476) / 33.62 = 3.197 V
			// VD% = (3.197 / 480) × 100 = 0.666%
			expectedPorc:   0.666,
			expectedCumple: true,
		},
		{
			name:       "Cu 12AWG 100m at 25A 220V exceeds 3%",
			corrienteA: 25,
			distanciaM: 100,
			calibre:    "12 AWG",
			material:   "Cu",
			seccionMM2: 3.31,
			voltaje:    220,
			limitePorc: 3.0,
			// VD = (√3 × 0.01724 × 100 × 25) / 3.31 = (74.63) / 3.31 = 22.55 V
			// VD% = (22.55 / 220) × 100 = 10.25%
			expectedPorc:   10.25,
			expectedCumple: false,
		},
		{
			name:       "Al conductor higher resistivity",
			corrienteA: 100,
			distanciaM: 20,
			calibre:    "4/0 AWG",
			material:   "Al",
			seccionMM2: 107.2,
			voltaje:    480,
			limitePorc: 3.0,
			// VD = (√3 × 0.02826 × 20 × 100) / 107.2 = (97.87) / 107.2 = 0.913 V
			// VD% = (0.913 / 480) × 100 = 0.190%
			expectedPorc:   0.190,
			expectedCumple: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corriente, err := valueobject.NewCorriente(tt.corrienteA)
			require.NoError(t, err)

			conductor, err := valueobject.NewConductor(valueobject.ConductorParams{
				Calibre:    tt.calibre,
				Material:   tt.material,
				SeccionMM2: tt.seccionMM2,
			})
			require.NoError(t, err)

			tension, err := valueobject.NewTension(tt.voltaje)
			require.NoError(t, err)

			porcentaje, cumple, err := service.CalcularCaidaTension(
				conductor, corriente, tt.distanciaM, tension, tt.limitePorc,
			)
			require.NoError(t, err)
			assert.InDelta(t, tt.expectedPorc, porcentaje, 0.01)
			assert.Equal(t, tt.expectedCumple, cumple)
		})
	}
}

func TestCalcularCaidaTension_DistanciaCero(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(100)
	conductor, _ := valueobject.NewConductor(valueobject.ConductorParams{
		Calibre: "2 AWG", Material: "Cu", SeccionMM2: 33.62,
	})
	tension, _ := valueobject.NewTension(480)

	_, _, err := service.CalcularCaidaTension(conductor, corriente, 0, tension, 3.0)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrDistanciaInvalida))
}

func TestCalcularCaidaTension_DistanciaNegativa(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(100)
	conductor, _ := valueobject.NewConductor(valueobject.ConductorParams{
		Calibre: "2 AWG", Material: "Cu", SeccionMM2: 33.62,
	})
	tension, _ := valueobject.NewTension(480)

	_, _, err := service.CalcularCaidaTension(conductor, corriente, -10, tension, 3.0)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrDistanciaInvalida))
}

func TestCalcularCaidaTension_ValidInputNoError(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(100)
	conductor, _ := valueobject.NewConductor(valueobject.ConductorParams{
		Calibre: "2 AWG", Material: "Cu", SeccionMM2: 33.62,
	})
	tension, _ := valueobject.NewTension(480)

	_, _, err := service.CalcularCaidaTension(conductor, corriente, 10, tension, 3.0)
	assert.NoError(t, err)
}
