// internal/calculos/application/usecase/calcular_num_hilos_tierra_test.go
package usecase

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestCalcularNumHilosTierra(t *testing.T) {
	tests := []struct {
		name             string
		tipoCanalizacion entity.TipoCanalizacion
		numTuberias      int
		expected         int
	}{
		// Casos de Charola - siempre 1 hilo
		{
			name:             "CharolaCableEspaciado_5tubos_retorna1",
			tipoCanalizacion: entity.TipoCanalizacionCharolaCableEspaciado,
			numTuberias:      5,
			expected:         1,
		},
		{
			name:             "CharolaCableTriangular_5tubos_retorna1",
			tipoCanalizacion: entity.TipoCanalizacionCharolaCableTriangular,
			numTuberias:      5,
			expected:         1,
		},
		{
			name:             "CharolaCableEspaciado_0tubos_retorna1",
			tipoCanalizacion: entity.TipoCanalizacionCharolaCableEspaciado,
			numTuberias:      0,
			expected:         1,
		},

		// Casos de Tubería PVC - 1 hilo de tierra por tubo (regla NOM)
		{
			name:             "TuberiaPVC_1tubo_retorna1",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaPVC,
			numTuberias:      1,
			expected:         1,
		},
		{
			name:             "TuberiaPVC_2tubos_retorna2",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaPVC,
			numTuberias:      2,
			expected:         2,
		},
		{
			name:             "TuberiaPVC_3tubos_retorna3",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaPVC,
			numTuberias:      3,
			expected:         3,
		},
		{
			name:             "TuberiaPVC_4tubos_retorna4",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaPVC,
			numTuberias:      4,
			expected:         4,
		},
		{
			name:             "TuberiaPVC_100tubos_retorna100",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaPVC,
			numTuberias:      100,
			expected:         100,
		},

		// Casos de Tubería Aluminio
		{
			name:             "TuberiaAluminio_1tubo_retorna1",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaAluminio,
			numTuberias:      1,
			expected:         1,
		},
		{
			name:             "TuberiaAluminio_3tubos_retorna3",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaAluminio,
			numTuberias:      3,
			expected:         3,
		},

		// Casos de Tubería Acero PG
		{
			name:             "TuberiaAceroPG_1tubo_retorna1",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaAceroPG,
			numTuberias:      1,
			expected:         1,
		},
		{
			name:             "TuberiaAceroPG_3tubos_retorna3",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaAceroPG,
			numTuberias:      3,
			expected:         3,
		},

		// Casos de Tubería Acero PD
		{
			name:             "TuberiaAceroPD_1tubo_retorna1",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaAceroPD,
			numTuberias:      1,
			expected:         1,
		},
		{
			name:             "TuberiaAceroPD_3tubos_retorna3",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaAceroPD,
			numTuberias:      3,
			expected:         3,
		},

		// Casos con valores inválidos (0 o negativos) - default a 1
		{
			name:             "TuberiaPVC_0tubos_retorna1_default",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaPVC,
			numTuberias:      0,
			expected:         1,
		},
		{
			name:             "TuberiaPVC_negativo_retorna1_default",
			tipoCanalizacion: entity.TipoCanalizacionTuberiaPVC,
			numTuberias:      -5,
			expected:         1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calcularNumHilosTierra(tt.tipoCanalizacion, tt.numTuberias)
			assert.Equal(t, tt.expected, result, "Para tipo %s con %d tubos, esperaba %d pero obtuvo %d",
				tt.tipoCanalizacion, tt.numTuberias, tt.expected, result)
		})
	}
}
