// internal/calculos/domain/service/calibre_superior.go
package service

import (
	"errors"
	"fmt"
	"strings"
)

// ErrCalibreNoReconocido se retorna cuando el calibre no existe en la lista NOM.
var ErrCalibreNoReconocido = errors.New("calibre no reconocido en la lista de calibres NOM")

// ErrNoExisteCalibreSuperior se retorna cuando se pide el calibre superior al máximo (1000 MCM).
var ErrNoExisteCalibreSuperior = errors.New("no existe calibre superior a 1000 MCM en la tabla NOM")

// calibresNOM es la secuencia estándar de calibres de acuerdo a NOM-001-SEDE.
// Orden ascendente: desde el calibre más pequeño (14) hasta el más grande (1000 MCM).
var calibresNOM = []string{
	"14",
	"12",
	"10",
	"8",
	"6",
	"4",
	"2",
	"1/0",
	"2/0",
	"3/0",
	"4/0",
	"250",
	"300",
	"350",
	"400",
	"500",
	"600",
	"750",
	"1000",
}

// normalizarCalibre elimina el sufijo " AWG" que agregan los CSVs de tablas NOM.
// Ejemplos: "2 AWG" → "2", "1/0 AWG" → "1/0", "250" → "250" (sin cambio)
func normalizarCalibre(calibre string) string {
	return strings.TrimSuffix(strings.TrimSpace(calibre), " AWG")
}

// ObtenerCalibreSuperior devuelve el siguiente calibre superior en la secuencia NOM.
// Acepta calibres con o sin sufijo " AWG" (ej: "2 AWG" o "2").
// Retorna el calibre en formato normalizado (sin sufijo " AWG").
// Retorna error si el calibre no existe en la lista o si es el calibre máximo (1000).
func ObtenerCalibreSuperior(calibreActual string) (string, error) {
	// Normalizar: "2 AWG" → "2"
	calibreNormalizado := normalizarCalibre(calibreActual)

	// Buscar el calibre actual en la lista
	indiceActual := -1
	for i, calibre := range calibresNOM {
		if calibre == calibreNormalizado {
			indiceActual = i
			break
		}
	}

	// Error: calibre no encontrado
	if indiceActual == -1 {
		return "", fmt.Errorf("%w: calibre '%s'", ErrCalibreNoReconocido, calibreActual)
	}

	// Error: es el último calibre
	if indiceActual == len(calibresNOM)-1 {
		return "", fmt.Errorf("%w", ErrNoExisteCalibreSuperior)
	}

	// Retornar el siguiente calibre
	return calibresNOM[indiceActual+1], nil
}
