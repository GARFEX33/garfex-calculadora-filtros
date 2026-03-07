// internal/pdf/infrastructure/adapter/driven/gotenberg/pdf_generator.go
package gotenberg

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/garfex/calculadora-filtros/internal/pdf/application/port"
)

const (
	// ErrGeneracionPdf es el código de error para errores de generación de PDF.
	// Se usa para identificar errores específicos del generador.
	ErrGeneracionPdf = "ERR_GENERACION_PDF"

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
	return g.GenerateWithHeaderFooter(ctx, htmlContent, "", "")
}

// GenerateWithHeaderFooter convierte el HTML a PDF usando el servicio Gotenberg con header y footer renderizados.
// Si headerHTML está vacío, usa el header del FS embebido (comportamiento legacy).
// Si footerHTML está vacío, usa el footer del FS embebido (comportamiento legacy).
// Implementa port.PdfGenerator.
func (g *PdfGeneratorAdapter) GenerateWithHeaderFooter(ctx context.Context, htmlContent, headerHTML, footerHTML string) ([]byte, error) {
	// 1. El HTML ya contiene CSS embebido desde memoria.html con variables dinámicas

	// 2. Usar header renderizado si se provee, sino extraer del FS
	// El headerHTML ya viene renderizado con los datos del proyecto

	// 3. Usar footer renderizado si se provee, sino extraer del FS
	// El footerHTML ya viene renderizado con los datos del proyecto

	// 4. El HTML ya tiene CSS embebido - no necesita inyección adicional

	// 5. Construir el formulario multipart
	form := NewFormBuilder()
	if err := form.AddHTML(htmlContent); err != nil {
		return nil, wrapError(err, "añadiendo HTML")
	}

	// Header nativo de Gotenberg (se renderiza en cada página)
	if err := form.AddHeader(headerHTML); err != nil {
		return nil, wrapError(err, "añadiendo header")
	}

	// Footer con paginación nativa de Chromium (pageNumber/totalPages)
	if err := form.AddFooter(footerHTML); err != nil {
		return nil, wrapError(err, "añadiendo footer")
	}

	// Opciones de conversión para Gotenberg
	form.AddOption("marginTop", "15mm")
	form.AddOption("marginBottom", "15mm")
	form.AddOption("marginLeft", "12mm")
	form.AddOption("marginRight", "12mm")

	// Sin waitDelay - MathJax fue removido del template
	// La página carga inmediatamente sin scripts externos que renderizar

	// 6. Construir el body del request
	body, err := form.Build()
	if err != nil {
		return nil, wrapError(err, "construyendo formulario")
	}

	// 7. Realizar la petición HTTP con reintentos
	resp, err := g.httpClient.PostMultipart(
		ctx,
		g.config.URL,
		form.ContentType(),
		body,
		g.config.MaxRetries,
	)
	if err != nil {
		log.Printf("[ERROR] gotenberg: failed to call Gotenberg at %s: %v", g.config.URL, err)
		return nil, wrapError(err, "llamando a Gotenberg")
	}

	// 8. Validar que tenemos contenido PDF
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
