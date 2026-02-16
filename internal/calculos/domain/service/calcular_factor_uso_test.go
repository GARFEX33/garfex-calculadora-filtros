// internal/calculos/domain/service/calcular_factor_uso_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/stretchr/testify/assert"
)

func TestCalcularFactorUso(t *testing.T) {
	tests := []struct {
		name          string
		tipoEquipo    entity.TipoEquipo
		expectedValue float64
		wantErr       bool
	}{
		{
			name:          "FILTRO_ACTIVO retorna 1.35",
			tipoEquipo:    entity.TipoEquipoFiltroActivo,
			expectedValue: 1.35,
			wantErr:       false,
		},
		{
			name:          "FILTRO_RECHAZO retorna 1.35",
			tipoEquipo:    entity.TipoEquipoFiltroRechazo,
			expectedValue: 1.35,
			wantErr:       false,
		},
		{
			name:          "TRANSFORMADOR retorna 1.25",
			tipoEquipo:    entity.TipoEquipoTransformador,
			expectedValue: 1.25,
			wantErr:       false,
		},
		{
			name:          "CARGA retorna 1.25",
			tipoEquipo:    entity.TipoEquipoCarga,
			expectedValue: 1.25,
			wantErr:       false,
		},
		{
			name:          "Tipo de equipo inválido retorna error",
			tipoEquipo:    "TIPO_INVALIDO",
			expectedValue: 0,
			wantErr:       true,
		},
		{
			name:          "Tipo de equipo vacío retorna error",
			tipoEquipo:    "",
			expectedValue: 0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factor, err := service.CalcularFactorUso(tt.tipoEquipo)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "tipo de equipo no válido")
				return
			}

			assert.NoError(t, err)
			assert.InDelta(t, tt.expectedValue, factor, 0.001)
		})
	}
}
