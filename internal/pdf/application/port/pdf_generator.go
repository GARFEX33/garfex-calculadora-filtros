// internal/pdf/application/port/pdf_generator.go
package port

import "context"

// PdfGenerator es el port driven para convertir HTML a bytes de PDF.
// La implementación concreta usa wkhtmltopdf vía infrastructure/adapter/driven/wkhtmltopdf/.
type PdfGenerator interface {
	// Generate convierte el HTML proporcionado a bytes de PDF.
	// ctx permite cancelar la operación (útil para timeouts).
	// html es el HTML completo a convertir.
	// Retorna los bytes del PDF generado o un error envuelto con ErrGeneracionPdf.
	Generate(ctx context.Context, html string) ([]byte, error)
}
