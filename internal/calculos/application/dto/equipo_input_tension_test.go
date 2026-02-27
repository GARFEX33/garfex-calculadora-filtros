package dto_test

import (
	"encoding/json"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEquipoInput_JSONTensionUnidad(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		wantTension float64
		wantUnidad  string
	}{
		{
			name:        "480V",
			jsonInput:   `{"tension":480,"tension_unidad":"V","modo":"MANUAL_POTENCIA","potencia_nominal":15,"potencia_unidad":"KW","factor_potencia":0.85,"itm":30,"sistema_electrico":"ESTRELLA","estado":"Ciudad de Mexico","tipo_canalizacion":"TUBERIA_PVC","longitud_circuito":50,"tipo_voltaje":"FASE_FASE"}`,
			wantTension: 480,
			wantUnidad:  "V",
		},
		{
			name:        "0.48kV",
			jsonInput:   `{"tension":0.48,"tension_unidad":"kV","modo":"MANUAL_POTENCIA","potencia_nominal":15,"potencia_unidad":"KW","factor_potencia":0.85,"itm":30,"sistema_electrico":"ESTRELLA","estado":"Ciudad de Mexico","tipo_canalizacion":"TUBERIA_PVC","longitud_circuito":50,"tipo_voltaje":"FASE_FASE"}`,
			wantTension: 0.48,
			wantUnidad:  "kV",
		},
		{
			name:        "no unidad - default V",
			jsonInput:   `{"tension":480,"modo":"MANUAL_POTENCIA","potencia_nominal":15,"potencia_unidad":"KW","factor_potencia":0.85,"itm":30,"sistema_electrico":"ESTRELLA","estado":"Ciudad de Mexico","tipo_canalizacion":"TUBERIA_PVC","longitud_circuito":50,"tipo_voltaje":"FASE_FASE"}`,
			wantTension: 480,
			wantUnidad:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input dto.EquipoInput
			err := json.Unmarshal([]byte(tt.jsonInput), &input)
			require.NoError(t, err)

			assert.Equal(t, tt.wantTension, input.Tension, "Tension mismatch")
			assert.Equal(t, tt.wantUnidad, input.TensionUnidad, "TensionUnidad mismatch")
		})
	}
}
