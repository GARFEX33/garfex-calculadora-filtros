// internal/pdf/infrastructure/adapter/driven/gotenberg/pdf_generator.go
package gotenberg

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/garfex/calculadora-filtros/internal/pdf/application/port"
)

const (
	// ErrGeneracionPdf es el código de error para errores de generación de PDF.
	// Se usa para identificar errores específicos del generador.
	ErrGeneracionPdf = "ERR_GENERACION_PDF"

	// rutaCSS es la ruta al archivo CSS dentro del FS embebido.
	rutaCSS = "templates/styles/pdf.css"

	// rutaHeader es la ruta al template de header dentro del FS embebido.
	rutaHeader = "templates/partials/header.html"

	// rutaFooter es la ruta al template de footer dentro del FS embebido.
	rutaFooter = "templates/partials/footer.html"
)

// PdfGeneratorAdapter implementa port.PdfGenerator usando Gotenberg vía HTTP.
type PdfGeneratorAdapter struct {
	config      *Config
	httpClient  GotenbergClient
	templatesFS fs.FS
}

// NewPdfGenerator crea un nuevo PdfGeneratorAdapter.
// templatesFS es el fs.FS embebido que contiene los templates (header, footer, CSS).
func NewPdfGenerator(templatesFS fs.FS) (*PdfGeneratorAdapter, error) {
	config := NewConfig()

	// Crear cliente HTTP con el timeout de la configuración
	httpClient := newDefaultHTTPClient(config.Timeout)

	return &PdfGeneratorAdapter{
		config:      config,
		httpClient:  httpClient,
		templatesFS: templatesFS,
	}, nil
}

// NewPdfGeneratorWithConfig crea un PdfGeneratorAdapter con configuración custom.
// Útil para tests o configuración específica.
func NewPdfGeneratorWithConfig(templatesFS fs.FS, cfg *Config) *PdfGeneratorAdapter {
	return &PdfGeneratorAdapter{
		config:      cfg,
		httpClient:  newDefaultHTTPClient(cfg.Timeout),
		templatesFS: templatesFS,
	}
}

// Generate convierte el HTML a PDF usando el servicio Gotenberg.
// Implementa port.PdfGenerator.
func (g *PdfGeneratorAdapter) Generate(ctx context.Context, htmlContent string) ([]byte, error) {
	// 1. Leer el CSS una sola vez
	cssData, err := fs.ReadFile(g.templatesFS, rutaCSS)
	if err != nil {
		return nil, wrapError(err, "leyendo CSS")
	}
	cssContent := string(cssData)

	// 2. Extraer header del FS embebido (footer se maneja vía CSS @page margin-box)
	headerContent, err := g.extractTemplate(rutaHeader, cssContent)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrGeneracionPdf, err)
	}

	// NOTA: Ya no usamos footer.html personalizado.
	// La paginación se maneja vía CSS @page { @bottom-center } en pdf.css
	// Esto es más compatible con Chromium/Gotenberg que los footers HTML personalizados.

	// 3. Inyectar CSS en el HTML principal
	htmlWithCSS := injectCSSIntoHTML(htmlContent, cssContent)

	// 4. Construir el formulario multipart
	form := NewFormBuilder()
	if err := form.AddHTML(htmlWithCSS); err != nil {
		return nil, wrapError(err, "añadiendo HTML")
	}

	if err := form.AddHeader(headerContent); err != nil {
		return nil, wrapError(err, "añadiendo header")
	}

	// Footer eliminado: paginación vía CSS @page { @bottom-center }
	// Ver pdf.css para la configuración de paginación

	// Opciones de conversión para Gotenberg (formato Letter como wkhtmltopdf)
	form.AddOption("pdfFormat", "Letter")
	form.AddOption("marginTop", "10mm")
	form.AddOption("marginBottom", "25mm") // Mayor margen para acomodar footer CSS
	form.AddOption("marginLeft", "15mm")
	form.AddOption("marginRight", "10mm")

	// 5. Construir el body del request
	body, err := form.Build()
	if err != nil {
		return nil, wrapError(err, "construyendo formulario")
	}

	// 6. Realizar la petición HTTP con reintentos
	resp, err := g.httpClient.PostMultipart(
		ctx,
		g.config.URL,
		form.ContentType(),
		body,
		g.config.MaxRetries,
	)
	if err != nil {
		return nil, wrapError(err, "llamando a Gotenberg")
	}

	// 7. Validar que tenemos contenido PDF
	if len(resp.Body) == 0 {
		return nil, errors.New("Gotenberg generó un PDF vacío")
	}

	// Verificar que parece ser un PDF (magic numbers)
	if !isPDF(resp.Body) {
		return nil, fmt.Errorf("respuesta de Gotenberg no es un PDF válido: %s", string(resp.Body[:min(200, len(resp.Body))]))
	}

	return resp.Body, nil
}

// extractTemplate lee un template del FS y le inyecta el CSS.
// Retorna string vacío si el template no existe (es opcional).
func (g *PdfGeneratorAdapter) extractTemplate(templatePath, cssContent string) (string, error) {
	data, err := fs.ReadFile(g.templatesFS, templatePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// Template no existe - retornar string vacío (es opcional)
			return "", nil
		}
		return "", fmt.Errorf("leyendo template %q: %w", templatePath, err)
	}

	html := string(data)
	return injectCSSIntoHTML(html, cssContent), nil
}

// injectCSSIntoHTML inyecta CSS en la etiqueta <head> del HTML.
// Busca la etiqueta <link> de estilos y la reemplaza por <style> con el CSS embebido.
func injectCSSIntoHTML(html, cssContent string) string {
	styleTag := "<style>\n" + cssContent + "\n</style>"

	// Reemplazar el link a stylesheet por el estilo embebido
	html = strings.Replace(
		html,
		`<link rel="stylesheet" href="/style.css">`,
		styleTag,
		1,
	)

	// Si no encontró el link, intentar añadirlo en el head
	if !strings.Contains(html, "<style>") {
		// Buscar </head> y añadir el style antes
		html = strings.Replace(
			html,
			"</head>",
			styleTag+"\n</head>",
			1,
		)
	}

	return html
}

// isPDF verifica si los bytes parecen ser un PDF válido.
// Un PDF válido empieza con "%PDF-" (magic number).
func isPDF(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	return string(data[:4]) == "%PDF"
}

// wrapError envuelve un error con contexto y el código de error específico.
func wrapError(err error, context string) error {
	return fmt.Errorf("%s: %s: %w", ErrGeneracionPdf, context, err)
}

// Config retorna la configuración actual del adapter (útil para testing).
func (g *PdfGeneratorAdapter) Config() *Config {
	return g.config
}

// Ensure we implement the port interface
var _ port.PdfGenerator = (*PdfGeneratorAdapter)(nil)
