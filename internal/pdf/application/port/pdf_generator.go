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

	// GenerateWithFooter convierte el HTML a PDF con un footer ya renderizado.
	// footerHTML es el footer con las variables de empresa ya resueltas (ej: {{.Empresa.NombreCompleto}} → "Garfex S.A.").
	// Esto permite que el footer se renderice con datos dinámicos antes de pasarlo al generador.
	// Si footerHTML está vacío, se usa el comportamiento por defecto (footer estático o vacío).
	GenerateWithFooter(ctx context.Context, html, footerHTML string) ([]byte, error)
}
