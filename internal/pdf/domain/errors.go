// internal/pdf/domain/errors.go
package domain

import "errors"

var (
	// ErrEmpresaNoEncontrada se retorna cuando el ID de empresa no existe en el catálogo.
	ErrEmpresaNoEncontrada = errors.New("empresa no encontrada en el catálogo")

	// ErrGeneracionPdf se retorna cuando falla la conversión de HTML a PDF.
	ErrGeneracionPdf = errors.New("error al generar el PDF")

	// ErrRenderizadoHtml se retorna cuando falla el renderizado del template HTML.
	ErrRenderizadoHtml = errors.New("error al renderizar el HTML")
)
