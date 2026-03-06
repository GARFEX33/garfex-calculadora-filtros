// cmd/pdf_test/main.go
// Herramienta standalone para generar PDFs de memorias de cálculo sin necesidad del servidor.
// Uso: go run cmd/pdf_test/main.go [-empresa=garfex|summaa|siemens] [-output=output.pdf]
//
// Requiere Gotenberg ejecutándose: docker-compose up gotenberg
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
	"time"

	calculosdto "github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	pdfdto "github.com/garfex/calculadora-filtros/internal/pdf/application/dto"
	pdfdomain "github.com/garfex/calculadora-filtros/internal/pdf/domain"
	pdfgotenberg "github.com/garfex/calculadora-filtros/internal/pdf/infrastructure/adapter/driven/gotenberg"
	pdftemplate "github.com/garfex/calculadora-filtros/internal/pdf/infrastructure/adapter/driven/template"
)

const (
	// Rutas a templates en el filesystem (no embedded)
	templatesBasePath = "internal/pdf"
	logosPath         = "internal/pdf/assets/logos"
	testDataPath      = "test_memoria.json"
	fechaLayout       = "02/01/2006"
)

var (
	outputFileName = flag.String("output", "test_output.pdf", "Archivo de salida del PDF")
	empresaID      = flag.String("empresa", "garfex", "ID de empresa (garfex, summaa, siemens)")
)

// PresentacionInput contiene los datos de presentación (re-definido para JSON parsing)
type PresentacionInput struct {
	EmpresaID            string `json:"empresa_id"`
	NombreProyecto       string `json:"nombre_proyecto"`
	DireccionProyecto    string `json:"direccion_proyecto"`
	Responsable          string `json:"responsable"`
	NombreEquipoOverride string `json:"nombre_equipo_override,omitempty"`
}

// TestData representa la estructura del JSON de prueba
type TestData struct {
	Memoria      calculosdto.MemoriaOutput `json:"memoria"`
	Presentacion *PresentacionInput        `json:"presentacion"`
}

func main() {
	flag.Parse()

	log.Printf("🎯 Generando PDF de prueba para empresa: %s", *empresaID)

	// 1. Validar empresa existe usando el catálogo del dominio
	empresa, ok := pdfdomain.BuscarEmpresaPorID(*empresaID)
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

	// 5. Construir TemplateData usando el DTO compartido
	data := buildTemplateData(testData, empresa, logoBase64, logoLetraBase64, nombreEquipo)

	// 6. Crear HtmlRendererAdapter con os.DirFS para leer templates desde disk
	// (igual que pdf_preview para permitir desarrollo iterativo sin recompilar)
	log.Printf("📄 Creando renderer HTML con templates desde disk: %s", templatesBasePath)
	diskFS := os.DirFS(templatesBasePath)
	renderer, err := pdftemplate.NewHtmlRenderer(diskFS)
	if err != nil {
		log.Fatalf("❌ Error inicializando renderer HTML: %v", err)
	}

	// 7. Crear PdfGeneratorAdapter (Gotenberg)
	// Pasar diskFS como fallback por si el renderer falla y necesita extraer templates
	log.Printf("📄 Creando generador PDF con Gotenberg...")
	generator, err := pdfgotenberg.NewPdfGenerator(diskFS)
	if err != nil {
		log.Fatalf("❌ Error inicializando generador PDF (Gotenberg): %v", err)
	}
	log.Printf("   Gotenberg URL: %s", generator.Config().URL)

	// 8. Renderizar body HTML
	log.Printf("📄 Renderizando template body: memoria_calculo.html")
	bodyHTML, err := renderer.Render("memoria_calculo.html", data)
	if err != nil {
		log.Fatalf("❌ Error renderizando body HTML: %v", err)
	}

	// 9. Renderizar header HTML (para Gotenberg native header)
	log.Printf("📄 Renderizando template header: gotenberg_header.html")
	headerHTML, err := renderer.Render("gotenberg_header.html", data)
	if err != nil {
		log.Printf("⚠️  Error renderizando header: %v, continuando sin header", err)
		headerHTML = ""
	}

	// 10. Renderizar footer HTML (para Gotenberg native footer con paginación)
	log.Printf("📄 Renderizando template footer: gotenberg_footer.html")
	footerHTML, err := renderer.Render("gotenberg_footer.html", data)
	if err != nil {
		log.Printf("⚠️  Error renderizando footer: %v, continuando sin footer", err)
		footerHTML = ""
	}

	// 11. Generar PDF con Gotenberg
	log.Printf("📄 Generando PDF con Gotenberg...")
	ctx := context.Background()
	pdfBytes, err := generator.GenerateWithHeaderFooter(ctx, bodyHTML, headerHTML, footerHTML)
	if err != nil {
		log.Fatalf("❌ Error generando PDF: %v", err)
	}

	// 12. Guardar PDF
	if err := os.WriteFile(*outputFileName, pdfBytes, 0644); err != nil {
		log.Fatalf("❌ Error guardando PDF: %v", err)
	}

	log.Printf("✅ PDF generado exitosamente: %s", *outputFileName)
	log.Printf("   Empresa: %s (%s)", empresa.NombreCompleto, *empresaID)
	log.Printf("   Proyecto: %s", data.NombreProyecto)
	log.Printf("   Gotenberg: %s", generator.Config().URL)
}

// buildTemplateData construye el DTO TemplateData desde los datos de prueba
func buildTemplateData(
	testData *TestData,
	empresa pdfdomain.EmpresaPresentacion,
	logoBase64 string,
	logoLetraBase64 string,
	nombreEquipo string,
) pdfdto.TemplateData {
	return pdfdto.TemplateData{
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
}

func loadTestData(path string) (*TestData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("leyendo archivo %s: %w", path, err)
	}

	var testData TestData
	if err := json.Unmarshal(data, &testData); err != nil {
		return nil, fmt.Errorf("parseando JSON: %w", err)
	}

	return &testData, nil
}

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
