// internal/shared/kernel/valueobject/temperatura_test.go
package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemperatura_Valor(t *testing.T) {
	tests := []struct {
		name string
		temp Temperatura
		want int
	}{
		{"60C", Temp60, 60},
		{"75C", Temp75, 75},
		{"90C", Temp90, 90},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.temp.Valor())
		})
	}
}
