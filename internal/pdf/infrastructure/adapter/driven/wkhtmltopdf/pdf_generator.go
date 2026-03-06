// internal/pdf/infrastructure/adapter/driven/wkhtmltopdf/pdf_generator.go
package wkhtmltopdf

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// envVarBinaryPath es la variable de entorno para especificar la ruta del binario wkhtmltopdf.
	envVarBinaryPath = "WKHTMLTOPDF_PATH"

	// defaultBinaryName es el nombre por defecto del binario wkhtmltopdf en el PATH del sistema.
	defaultBinaryName = "wkhtmltopdf"

	// tempFilePattern es el patrón para los archivos temporales de HTML.
	tempFilePattern = "garfex-memoria-*.html"

	// tempHeaderPattern es el patrón para el archivo temporal de header.
	tempHeaderPattern = "garfex-header-*.html"

	// tempFooterPattern es el patrón para el archivo temporal de footer.
	tempFooterPattern = "garfex-footer-*.html"

	// tempPdfPattern es el patrón para el archivo temporal de PDF.
	tempPdfPattern = "garfex-pdf-*.pdf"
)

// PdfGeneratorAdapter implementa port.PdfGenerator usando wkhtmltopdf vía exec.CommandContext.
// El binario se busca primero en WKHTMLTOPDF_PATH, luego en el PATH del sistema.
type PdfGeneratorAdapter struct {
	binaryPath  string
	templatesFS fs.FS
}

// NewPdfGenerator crea un PdfGeneratorAdapter.
// Verifica que el binario wkhtmltopdf esté disponible al momento de la construcción.
// templatesFS es el embed.FS que contiene los templates de header y footer.
func NewPdfGenerator(templatesFS fs.FS) (*PdfGeneratorAdapter, error) {
	binaryPath := resolveBinaryPath()

	// Verificar que el binario existe y es ejecutable
	if _, err := exec.LookPath(binaryPath); err != nil {
		// Si el PATH configurado no funciona, intentar el fallback
		if binaryPath != defaultBinaryName {
			if _, fallbackErr := exec.LookPath(defaultBinaryName); fallbackErr == nil {
				binaryPath = defaultBinaryName
			} else {
				return nil, fmt.Errorf("wkhtmltopdf no encontrado en %q ni en PATH del sistema: %w", binaryPath, err)
			}
		} else {
			return nil, fmt.Errorf("wkhtmltopdf no encontrado en PATH del sistema: %w", err)
		}
	}

	return &PdfGeneratorAdapter{
		binaryPath:  binaryPath,
		templatesFS: templatesFS,
	}, nil
}

// Generate convierte el HTML a PDF usando wkhtmltopdf.
// Flujo: escribir HTML temp → extraer header/footer → ejecutar wkhtmltopdf → leer PDF → limpiar temps.
// Implementa port.PdfGenerator.
func (g *PdfGeneratorAdapter) Generate(ctx context.Context, htmlContent string) ([]byte, error) {
	// Crear directorio temporal para todos los archivos de esta generación
	tmpDir, err := os.MkdirTemp("", "garfex-pdf-*")
	if err != nil {
		return nil, fmt.Errorf("creando directorio temporal: %w", err)
	}
	defer os.RemoveAll(tmpDir) // limpiar siempre al finalizar

	// El HTML ya contiene CSS embebido desde memoria.html con variables dinámicas
	htmlFile := filepath.Join(tmpDir, "memoria.html")
	if err := os.WriteFile(htmlFile, []byte(htmlContent), 0600); err != nil {
		return nil, fmt.Errorf("escribiendo HTML temporal: %w", err)
	}
	// El HTML ya contiene CSS embebido desde memoria.html - no necesita inyección

	// 2. Extraer header.html del embed.FS y escribirlo como archivo temporal
	headerFile := filepath.Join(tmpDir, "header.html")
	if err := extractTemplate(g.templatesFS, "templates/partials/header.html", headerFile); err != nil {
		return nil, fmt.Errorf("extrayendo header template: %w", err)
	}
	// El header usa las variables CSS embebidas en memoria.html

	// 3. Extraer footer.html del embed.FS y escribirlo como archivo temporal
	footerFile := filepath.Join(tmpDir, "footer.html")
	if err := extractTemplate(g.templatesFS, "templates/partials/footer.html", footerFile); err != nil {
		return nil, fmt.Errorf("extrayendo footer template: %w", err)
	}
	// El footer se maneja vía CSS @page en memoria.html

	// 4. Archivo de salida PDF
	pdfFile := filepath.Join(tmpDir, "memoria.pdf")

	// 5. Construir comando wkhtmltopdf con los flags del PRD
	args := buildWkhtmltopdfArgs(htmlFile, headerFile, footerFile, pdfFile)

	// 6. Ejecutar wkhtmltopdf con el contexto (respeta timeout/cancelación)
	cmd := exec.CommandContext(ctx, g.binaryPath, args...) //nolint:gosec // path validado en NewPdfGenerator
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Verificar si el error es por contexto cancelado
		if ctx.Err() != nil {
			return nil, fmt.Errorf("generación de PDF cancelada/timeout: %w", ctx.Err())
		}

		// Verificar si el binario no fue encontrado
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("wkhtmltopdf falló con código %d: %s", exitErr.ExitCode(), string(output))
		}

		// Error de tipo "ejecutable no encontrado"
		if errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("wkhtmltopdf no encontrado: %w", err)
		}

		return nil, fmt.Errorf("ejecutando wkhtmltopdf: %w — output: %s", err, string(output))
	}

	// 7. Leer el PDF generado
	pdfBytes, err := os.ReadFile(pdfFile)
	if err != nil {
		return nil, fmt.Errorf("leyendo PDF generado: %w", err)
	}

	if len(pdfBytes) == 0 {
		return nil, fmt.Errorf("wkhtmltopdf generó un PDF vacío")
	}

	return pdfBytes, nil
}

// BinaryPath retorna la ruta del binario wkhtmltopdf que se está usando.
// Útil para logging y diagnóstico.
func (g *PdfGeneratorAdapter) BinaryPath() string {
	return g.binaryPath
}

// resolveBinaryPath determina la ruta del binario wkhtmltopdf.
// Prioridad: variable de entorno WKHTMLTOPDF_PATH → "wkhtmltopdf" en PATH.
func resolveBinaryPath() string {
	if envPath := os.Getenv(envVarBinaryPath); envPath != "" {
		return envPath
	}
	return defaultBinaryName
}

// buildWkhtmltopdfArgs construye el slice de argumentos para wkhtmltopdf.
// Usa los flags definidos en el PRD sección 4.2 para tamaño Letter con márgenes NOM.
func buildWkhtmltopdfArgs(htmlFile, headerFile, footerFile, outputFile string) []string {
	return []string{
		// Configuración de página
		"--page-size", "Letter",
		"--margin-top", "10mm",
		"--margin-bottom", "15mm",
		"--margin-left", "15mm",
		"--margin-right", "10mm",

		// Header y footer como archivos HTML separados
		"--header-html", headerFile,
		"--footer-html", footerFile,

		// Permitir acceso a archivos locales (necesario para logos en base64 y CSS)
		"--enable-local-file-access",

		// Encoding y resolución para mejor calidad
		"--encoding", "utf-8",
		"--dpi", "300",

		// Desactivar inteligencia de páginas para evitar páginas en blanco extra
		"--disable-smart-shrinking",

		// Silenciar output de progreso (solo errores reales)
		"--quiet",

		// Archivos de entrada y salida
		htmlFile,
		outputFile,
	}
}

// extractTemplate lee un template del embed.FS y lo escribe en destPath.
func extractTemplate(templatesFS fs.FS, templatePath, destPath string) error {
	content, err := fs.ReadFile(templatesFS, templatePath)
	if err != nil {
		return fmt.Errorf("leyendo %q del FS embebido: %w", templatePath, err)
	}

	if err := os.WriteFile(destPath, content, 0600); err != nil {
		return fmt.Errorf("escribiendo %q: %w", destPath, err)
	}

	return nil
}

// injectCSS inyecta el contenido CSS dentro de una etiqueta <style> en el <head> del archivo HTML.
func injectCSS(filePath, cssContent string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("leyendo archivo %q: %w", filePath, err)
	}

	html := string(data)
	styleTag := "<style>" + cssContent + "</style>"
	html = strings.Replace(
		html,
		`<link rel="stylesheet" href="/style.css">`,
		styleTag,
		1,
	)

	if err := os.WriteFile(filePath, []byte(html), 0600); err != nil {
		return fmt.Errorf("escribiendo archivo %q con CSS: %w", filePath, err)
	}

	return nil
}
