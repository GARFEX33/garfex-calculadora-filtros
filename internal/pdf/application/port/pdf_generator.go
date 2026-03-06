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

	// GenerateWithHeaderFooter convierte el HTML a PDF con header y footer ya renderizados.
	// headerHTML es el header con las variables de empresa ya resueltas (ej: {{.Empresa.NombreCompleto}} → "Garfex S.A.").
	// footerHTML es el footer con las variables de empresa ya resueltas.
	// Esto permite que header/footer se rendericen con datos dinámicos antes de pasarlos al generador.
	// Si headerHTML está vacío, se usa el comportamiento por defecto (sin header o header estático).
	// Si footerHTML está vacío, se usa el comportamiento por defecto (footer estático o vacío).
	GenerateWithHeaderFooter(ctx context.Context, html, headerHTML, footerHTML string) ([]byte, error)
}
