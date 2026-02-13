package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

var ErrCharolaNoEncontrada = errors.New("no se encontró charola suficiente")

// CableControl representa un cable de control o comunicación que se transporta en charola.
type CableControl struct {
	Cantidad   int
	DiametroMM float64
}

type ConductorConDiametro struct {
	DiametroMM float64
}

func CalcularCharolaEspaciado(
	hilosPorFase int,
	sistema entity.SistemaElectrico,
	conductorFase ConductorConDiametro,
	conductorTierra ConductorConDiametro,
	tablaCharola []valueobject.EntradaTablaCanalizacion,
	cablesControl []CableControl,
) (entity.Canalizacion, error) {
	if hilosPorFase < 1 {
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaEspaciado: %w", errors.New("hilos por fase debe ser >= 1"))
	}

	// Determinar numero de fases segun el tipo de sistema
	var numFases int
	var tieneNeutro bool
	switch sistema {
	case entity.SistemaElectricoMonofasico:
		numFases = 1
		tieneNeutro = true
	case entity.SistemaElectricoBifasico:
		numFases = 2
		tieneNeutro = true
	case entity.SistemaElectricoDelta:
		numFases = 3
		tieneNeutro = false
	case entity.SistemaElectricoEstrella:
		numFases = 3
		tieneNeutro = true
	default:
		return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaEspaciado: %w", fmt.Errorf("sistema eléctrico no válido: %v", sistema))
	}

	// Calcular hilos de fase y neutro multiplicando por hilosPorFase (conductores en paralelo)
	hilosFaseTotal := numFases * hilosPorFase

	hilosNeutro := 0
	if tieneNeutro {
		hilosNeutro = hilosPorFase // El neutro se multiplica igual que las fases
	}

	totalHilos := hilosFaseTotal + hilosNeutro

	// Calcular ancho para cables de control
	var anchoControl float64
	for _, cable := range cablesControl {
		if cable.Cantidad > 0 && cable.DiametroMM > 0 {
			anchoControl += float64(cable.Cantidad) * cable.DiametroMM
		}
	}

	anchoRequerido := float64(totalHilos-1)*conductorFase.DiametroMM + conductorTierra.DiametroMM + anchoControl

	const alturaCharolaMM float64 = 50.0
	for _, entrada := range tablaCharola {
		anchoCharolaMM := entrada.AreaInteriorMM2 / alturaCharolaMM
		if anchoCharolaMM >= anchoRequerido {
			return entity.Canalizacion{
				Tipo:           string(entity.TipoCanalizacionCharolaCableEspaciado),
				Tamano:         entrada.Tamano,
				AnchoRequerido: anchoRequerido,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf("CalcularCharolaEspaciado: %w", ErrCharolaNoEncontrada)
}
