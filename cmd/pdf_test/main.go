// cmd/pdf_test/main.go
// Herramienta standalone para generar PDFs de memorias de cálculo sin necesidad del servidor.
// Uso: go run cmd/pdf_test/main.go [-empresa=garfex|summaa|siemens]
//
// El tool lee templates directamente del disk (no embedded) para permitir
// desarrollo iterativo sin recompilar el servidor.
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	htmpl "html/template"
	"math"
	"os/exec"
)

const (
	// Rutas a templates en el filesystem (no embedded)
	templatesPath  = "internal/pdf/templates"
	stylesPath     = "internal/pdf/templates/styles"
	assetsPath     = "internal/pdf/assets"
	logosPath      = "internal/pdf/assets/logos"
	testDataPath   = "test_memoria.json"
	outputFileName = "test_output.pdf"
	templateName   = "memoria_calculo.html"
	fechaLayout    = "02/01/2006"
)

// Empresas disponibles
var empresas = map[string]EmpresaPresentacion{
	"garfex": {
		ID:              "garfex",
		NombreCompleto:  "Garfex",
		LogoPath:        "garfex.png",
		Direccion:       "Av. Insurgentes Sur 1234, Col. Del Valle, CDMX, C.P. 03100",
		Telefono:        "+52 55 1193-0515",
		Email:           "jcgarcia@garfex.mx",
		ColorPrimario:   "#7C0000",
		ColorSecundario: "#F4CF00",
	},
	"summaa": {
		ID:              "summaa",
		NombreCompleto:  "Summaa S.A. de C.V.",
		LogoPath:        "summaa.png",
		Direccion:       "Blvd. Manuel Ávila Camacho 800, Lomas de Chapultepec, CDMX, C.P. 11000",
		Telefono:        "+52 55 9876-5432",
		Email:           "ventas@summaa.com",
		ColorPrimario:   "#004A99",
		ColorSecundario: "#1B75BB",
	},
	"siemens": {
		ID:              "siemens",
		NombreCompleto:  "Siemens S.A. de C.V.",
		LogoPath:        "siemens.png",
		Direccion:       "Lago Alberto 319, Anáhuac I Secc, Miguel Hidalgo, CDMX, C.P. 11320",
		Telefono:        "+52 55 5229-3600",
		Email:           "contacto@siemens.com.mx",
		ColorPrimario:   "#009999",
		ColorSecundario: "#000000",
	},
}

// EmpresaPresentacion contiene los datos de la empresa
type EmpresaPresentacion struct {
	ID              string
	NombreCompleto  string
	LogoPath        string
	Direccion       string
	Telefono        string
	Email           string
	ColorPrimario   string
	ColorSecundario string
}

// TemplateData es el struct que alimenta el template HTML
type TemplateData struct {
	Empresa           EmpresaPresentacion
	LogoBase64        string
	LogoLetraBase64   string
	NombreProyecto    string
	DireccionProyecto string
	Responsable       string
	NombreEquipo      string
	Memoria           dto.MemoriaOutput
	FechaGeneracion   string
}

func main() {
	// Parsear flags
	empresaID := flag.String("empresa", "garfex", "ID de empresa (garfex, summaa, siemens)")
	flag.Parse()

	log.Printf("🎯 Generando PDF de prueba para empresa: %s", *empresaID)

	// 1. Verificar empresa existe
	empresa, ok := empresas[*empresaID]
	if !ok {
		log.Fatalf("❌ Empresa desconocida: %s. Opciones: garfex, summaa, siemens", *empresaID)
	}

	// 2. Cargar datos de prueba
	testData, err := loadTestData(testDataPath)
	if err != nil {
		log.Fatalf("❌ Error cargando datos de prueba: %v", err)
	}

	// Actualizar empresa_id en los datos si el usuario especificó una diferente
	if testData.Presentacion == nil {
		testData.Presentacion = &PresentacionInput{}
	}
	testData.Presentacion.EmpresaID = *empresaID

	// 3. Cargar logos - paths son relativos a internal/pdf/assets/logos/
	logoBase64 := loadLogoBase64(filepath.Join(logosPath, empresa.LogoPath))
	var logoLetraBase64 string
	if *empresaID == "garfex" {
		logoLetraBase64 = loadLogoBase64(filepath.Join(logosPath, "lg.png"))
	}

	// 4. Determinar nombre del equipo
	nombreEquipo := testData.Presentacion.NombreEquipoOverride
	if nombreEquipo == "" {
		nombreEquipo = testData.Memoria.Equipo.Clave
	}

	// 5. Construir TemplateData
	data := TemplateData{
		Empresa:           empresa,
		LogoBase64:        logoBase64,
		LogoLetraBase64:   logoLetraBase64,
		NombreProyecto:    testData.Presentacion.NombreProyecto,
		DireccionProyecto: testData.Presentacion.DireccionProyecto,
		Responsable:       testData.Presentacion.Responsable,
		NombreEquipo:      nombreEquipo,
		Memoria:           testData.Memoria,
		FechaGeneracion:   time.Now().Format(fechaLayout),
	}

	// 6. Renderizar HTML desde disk (no embedded)
	log.Printf("📄 Renderizando template desde disk: %s", templatesPath)
	html, err := renderHTML(data, templatesPath, stylesPath)
	if err != nil {
		log.Fatalf("❌ Error renderizando HTML: %v", err)
	}

	// 7. Generar PDF
	log.Printf("📄 Generando PDF con wkhtmltopdf...")
	pdfBytes, err := generatePDF(context.Background(), html, templatesPath)
	if err != nil {
		log.Fatalf("❌ Error generando PDF: %v", err)
	}

	// 8. Guardar PDF
	if err := os.WriteFile(outputFileName, pdfBytes, 0644); err != nil {
		log.Fatalf("❌ Error guardando PDF: %v", err)
	}

	log.Printf("✅ PDF generado exitosamente: %s", outputFileName)
	log.Printf("   Empresa: %s (%s)", empresa.NombreCompleto, *empresaID)
	log.Printf("   Proyecto: %s", data.NombreProyecto)
}

// PresentacionInput contiene los datos de presentación
type PresentacionInput struct {
	EmpresaID            string `json:"empresa_id"`
	NombreProyecto       string `json:"nombre_proyecto"`
	DireccionProyecto    string `json:"direccion_proyecto"`
	Responsable          string `json:"responsable"`
	NombreEquipoOverride string `json:"nombre_equipo_override,omitempty"`
}

// TestData representa la estructura del JSON de prueba
type TestData struct {
	Memoria      dto.MemoriaOutput  `json:"memoria"`
	Presentacion *PresentacionInput `json:"presentacion"`
}

// loadTestData carga los datos de prueba desde un archivo JSON
func loadTestData(path string) (*TestData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("leyendo archivo %s: %w", path, err)
	}

	var testData TestData
	if err := jsonUnmarshal(data, &testData); err != nil {
		return nil, fmt.Errorf("parseando JSON: %w", err)
	}

	return &testData, nil
}

// loadLogoBase64 carga un archivo de imagen y lo codifica en base64
func loadLogoBase64(path string) string {
	if path == "" {
		return ""
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("⚠️  Logo no encontrado: %s", path)
		return ""
	}

	return base64.StdEncoding.EncodeToString(data)
}

// loadCSS lee el archivo CSS y lo retorna como string
func loadCSS(stylesDir string) (string, error) {
	cssPath := filepath.Join(stylesDir, "pdf.css")
	data, err := os.ReadFile(cssPath)
	if err != nil {
		return "", fmt.Errorf("leyendo archivo CSS %s: %w", cssPath, err)
	}
	return string(data), nil
}

// renderHTML renderiza el template HTML usando fs.FS del filesystem
func renderHTML(data TemplateData, templatesDir string, stylesDir string) (string, error) {
	// Cargar CSS para embedding inline
	cssContent, err := loadCSS(stylesDir)
	if err != nil {
		return "", fmt.Errorf("cargando CSS: %w", err)
	}

	// Usar os.DirFS para leer desde disk (no embed)
	diskFS := os.DirFS(templatesDir)

	funcMap := htmpl.FuncMap{
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"capFirst": func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToUpper(s[:1])
		},
		"slice": func(s string, start, end int) string {
			if start < 0 || start > len(s) || end < start || end > len(s) {
				return ""
			}
			return s[start:end]
		},
		"formatFloat": func(f float64, decimals int) string {
			return fmt.Sprintf(fmt.Sprintf("%%.%df", decimals), f)
		},
		"formatFloat2": func(f float64) string {
			return fmt.Sprintf("%.2f", f)
		},
		"formatFloat4": func(f float64) string {
			return fmt.Sprintf("%.4f", f)
		},
		"formatInt": func(f float64) string {
			return fmt.Sprintf("%.0f", f)
		},
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"sub": func(a, b float64) float64 {
			return a - b
		},
		"sqrt": func(f float64) float64 {
			return math.Sqrt(f)
		},
		"safeHTML": func(s string) htmpl.HTML {
			return htmpl.HTML(s)
		},
		"contains": func(s, sub string) bool {
			return strings.Contains(s, sub)
		},
		"notNil": func(v interface{}) bool {
			return v != nil
		},
		"esMaterial": func(material, expected string) bool {
			return material == expected
		},
		"itoa": func(i int) string {
			return fmt.Sprintf("%d", i)
		},
		"percent": func(f float64) string {
			return fmt.Sprintf("%.0f", f*100)
		},
		"not": func(b bool) bool {
			return !b
		},
		"mulIntFloat": func(i int, f float64) float64 {
			return float64(i) * f
		},
		"intToFloat": func(i int) float64 {
			return float64(i)
		},
		"toFloat64": func(v interface{}) float64 {
			return toFloat64(v)
		},
		"formatNumeric": func(v interface{}) string {
			return fmt.Sprintf("%.2f", toFloat64(v))
		},
		"derefFloat": func(f *float64) float64 {
			if f == nil {
				return 0
			}
			return *f
		},
	}

	// Parsear templates desde disk FS
	tmpl, err := htmpl.New("").Funcs(funcMap).ParseFS(diskFS,
		"memoria.html",
		"partials/*.html",
	)
	if err != nil {
		return "", fmt.Errorf("parseando templates: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", fmt.Errorf("ejecutando template %q: %w", templateName, err)
	}

	// Embed CSS inline: reemplazar <link> con <style>
	htmlContent := buf.String()
	htmlContent = strings.Replace(
		htmlContent,
		`<link rel="stylesheet" href="/style.css">`,
		`<style>`+cssContent+`</style>`,
		1,
	)

	return htmlContent, nil
}

// generatePDF convierte el HTML a PDF usando wkhtmltopdf
func generatePDF(ctx context.Context, htmlContent string, templatesDir string) ([]byte, error) {
	// Crear directorio temporal
	tmpDir, err := os.MkdirTemp("", "garfex-pdf-test-*")
	if err != nil {
		return nil, fmt.Errorf("creando directorio temporal: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Escribir HTML principal
	htmlFile := filepath.Join(tmpDir, "memoria.html")
	if err := os.WriteFile(htmlFile, []byte(htmlContent), 0600); err != nil {
		return nil, fmt.Errorf("escribiendo HTML temporal: %w", err)
	}

	// Extraer header y footer
	headerFile := filepath.Join(tmpDir, "header.html")
	footerFile := filepath.Join(tmpDir, "footer.html")

	if err := copyFile(filepath.Join(templatesDir, "partials/header.html"), headerFile); err != nil {
		return nil, fmt.Errorf("copiando header: %w", err)
	}
	if err := copyFile(filepath.Join(templatesDir, "partials/footer.html"), footerFile); err != nil {
		return nil, fmt.Errorf("copiando footer: %w", err)
	}

	// PDF de salida
	pdfFile := filepath.Join(tmpDir, "memoria.pdf")

	// Construir comando wkhtmltopdf
	args := []string{
		"--page-size", "Letter",
		"--margin-top", "20mm",
		"--margin-bottom", "25mm",
		"--margin-left", "25mm",
		"--margin-right", "15mm",
		"--header-html", headerFile,
		"--footer-html", footerFile,
		"--enable-local-file-access",
		"--encoding", "utf-8",
		"--dpi", "300",
		"--disable-smart-shrinking",
		"--quiet",
		htmlFile,
		pdfFile,
	}

	// Verificar si wkhtmltopdf está disponible
	wkhtmltopdfPath := "wkhtmltopdf"
	if envPath := os.Getenv("WKHTMLTOPDF_PATH"); envPath != "" {
		wkhtmltopdfPath = envPath
	}

	cmd := exec.CommandContext(ctx, wkhtmltopdfPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("wkhtmltopdf falló: %w — output: %s", err, string(output))
	}

	// Leer PDF generado
	pdfBytes, err := os.ReadFile(pdfFile)
	if err != nil {
		return nil, fmt.Errorf("leyendo PDF generado: %w", err)
	}

	if len(pdfBytes) == 0 {
		return nil, fmt.Errorf("wkhtmltopdf generó un PDF vacío")
	}

	return pdfBytes, nil
}

// copyFile copia un archivo de origen a destino
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0600)
}

// toFloat64 convierte cualquier tipo numérico a float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	default:
		return 0
	}
}

// jsonUnmarshal es un wrapper para json.Unmarshal que usa el paquete estándar
func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
