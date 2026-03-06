// cmd/pdf_preview/main.go
// Servidor de preview para desarrollo de templates PDF.
// Renderiza HTML directamente en el navegador para feedback instantáneo.
// Uso: go run cmd/pdf_preview/main.go
//
// Endpoints:
//   - GET /                    : Renderiza la memoria completa
//   - GET /style.css           : Serve CSS stylesheet
//   - GET /?empresa=garfex     : Cambiar empresa (garfex, summaa, siemens)
package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	htmpl "html/template"
)

const (
	templatesPath = "internal/pdf/templates"
	logosPath     = "internal/pdf/assets/logos"
	testDataPath  = "test_memoria.json"
	templateName  = "memoria_calculo.html"
	fechaLayout   = "02/01/2006"
	serverPort    = "3000"
)

var (
	empresaID = flag.String("empresa", "garfex", "Empresa por defecto (garfex, summaa, siemens)")
	noOpen    = flag.Bool("no-open", false, "No abrir navegador automáticamente")
	port      = flag.String("port", serverPort, "Puerto del servidor")
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

// Cache para el HTML renderizado
var cachedHTML string
var lastModified time.Time

func main() {
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("🌐 ")

	log.Printf("🎯 PDF Preview Server Started")
	log.Printf("   URL:      http://localhost:%s", *port)
	log.Printf("   Empresa:  %s", *empresaID)
	log.Printf("   Templates: %s", templatesPath)
	log.Printf("")
	log.Printf("   Presiona Ctrl+C para detener")
	log.Printf("")

	// Verificar empresa existe
	if _, ok := empresas[*empresaID]; !ok {
		log.Fatalf("❌ Empresa desconocida: %s. Opciones: garfex, summaa, siemens", *empresaID)
	}

	// Abrir navegador automáticamente
	if !*noOpen {
		go func() {
			time.Sleep(500 * time.Millisecond)
			openBrowser("http://localhost:" + *port)
		}()
	}

	// Configurar rutas
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/health", handleHealth)
	// CSS embebido en memoria.html - no necesita ruta separada

	// Iniciar servidor
	addr := ":" + *port
	log.Printf("✅ Servidor escuchando en http://localhost%s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// Obtener empresa de query param o usar default
	empID := r.URL.Query().Get("empresa")
	if empID == "" {
		empID = *empresaID
	}

	// Validar empresa
	empresa, ok := empresas[empID]
	if !ok {
		http.Error(w, fmt.Sprintf("Empresa desconocida: %s. Opciones: garfex, summaa, siemens", empID), http.StatusBadRequest)
		return
	}

	// Renderizar HTML
	html, err := renderHTMLForEmpresa(empresa)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error renderizando template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	fmt.Fprint(w, html)
}

// handleHealth returns the last modified timestamp of template files for hot-reload
func handleHealth(w http.ResponseWriter, r *http.Request) {
	modTime, err := getTemplatesModTime()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obteniendo timestamp: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"last_modified":      modTime.Unix(),
		"last_modified_nano": modTime.UnixNano(),
		"formatted":          modTime.Format(time.RFC3339),
	})
}

// getTemplatesModTime returns the most recent modification time of all template files
func getTemplatesModTime() (time.Time, error) {
	var newest time.Time

	err := filepath.Walk(templatesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories and hidden files
		if info.IsDir() || strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}
		// Check .html, .css files in templates
		if !info.IsDir() && (strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".css")) {
			if info.ModTime().After(newest) {
				newest = info.ModTime()
			}
		}
		return nil
	})

	if err != nil {
		return time.Time{}, err
	}

	// If no files found, return current time to avoid issues
	if newest.IsZero() {
		return time.Now(), nil
	}

	return newest, nil
}

func renderHTMLForEmpresa(empresa EmpresaPresentacion) (string, error) {
	// Cargar datos de prueba
	testData, err := loadTestData(testDataPath)
	if err != nil {
		return "", fmt.Errorf("cargando datos de prueba: %w", err)
	}

	// Actualizar empresa_id
	if testData.Presentacion == nil {
		testData.Presentacion = &PresentacionInput{}
	}
	testData.Presentacion.EmpresaID = empresa.ID

	// Cargar logos
	logoBase64 := loadLogoBase64(filepath.Join(logosPath, empresa.LogoPath))
	var logoLetraBase64 string
	if empresa.ID == "garfex" {
		logoLetraBase64 = loadLogoBase64(filepath.Join(logosPath, "lg.png"))
	}

	// Determinar nombre del equipo
	nombreEquipo := testData.Presentacion.NombreEquipoOverride
	if nombreEquipo == "" {
		nombreEquipo = testData.Memoria.Equipo.Clave
	}

	// Construir TemplateData
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

	return renderHTML(data, templatesPath)
}

func loadTestData(path string) (*TestData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("leyendo archivo %s: %w", path, err)
	}

	var testData TestData
	if err := parseJSON(data, &testData); err != nil {
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

func renderHTML(data TemplateData, templatesDir string) (string, error) {
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

	// Inject hot-reload script
	html := buf.String()
	html = injectHotReloadScript(html)

	return html, nil
}

// injectHotReloadScript adds JavaScript for automatic page refresh when templates change
func injectHotReloadScript(html string) string {
	hotReloadScript := `
<!-- Hot Reload Script -->
<script>
(function() {
  'use strict';
  
  let lastModified = null;
  let isRefreshing = false;
  
  // Get current template modification time
  async function checkForChanges() {
    try {
      const response = await fetch('/health');
      if (!response.ok) throw new Error('Health check failed');
      
      const data = await response.json();
      const currentModTime = data.last_modified_nano;
      
      // First load - store initial timestamp
      if (lastModified === null) {
        lastModified = currentModTime;
        console.log('[Hot Reload] Watching for changes...');
        return;
      }
      
      // Template changed - refresh
      if (currentModTime !== lastModified) {
        if (!isRefreshing) {
          isRefreshing = true;
          console.log('[Hot Reload] Template changed! Refreshing...');
          
          // Add visual indicator
          const indicator = document.createElement('div');
          indicator.id = 'hot-reload-indicator';
          indicator.style.cssText = 'position:fixed;top:0;left:0;right:0;background:#10b981;color:white;padding:8px;text-align:center;z-index:99999;font-family:monospace;font-size:14px;';
          indicator.textContent = '🔄 Template changed - Refreshing...';
          document.body.appendChild(indicator);
          
          // Brief delay for visual feedback
          setTimeout(() => {
            window.location.reload();
          }, 300);
        }
      }
    } catch (err) {
      // Silently ignore errors during polling
    }
  }
  
  // Poll every 2 seconds
  setInterval(checkForChanges, 2000);
  
  console.log('[Hot Reload] Initialized');
})();
</script>
</body>`

	// Find </body> tag and insert script before it
	if strings.Contains(html, "</body>") {
		html = strings.Replace(html, "</body>", hotReloadScript, 1)
	} else {
		// If no </body> tag, append at end
		html += hotReloadScript
	}

	return html
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch strings.ToLower(os.Getenv("GOOS")) {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		viewers := []string{"xdg-open", "firefox", "chromium", "google-chrome"}
		for _, viewer := range viewers {
			if _, err := exec.LookPath(viewer); err == nil {
				cmd = exec.Command(viewer, url)
				break
			}
		}
		if cmd == nil {
			log.Printf("⚠️  No se encontró navegador. Abre manualmente: %s", url)
			return
		}
	}

	if err := cmd.Start(); err != nil {
		log.Printf("⚠️  Error abriendo navegador: %v", err)
		return
	}

	log.Printf("📂 Navegador abierto: %s", url)
}

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

func parseJSON(data []byte, v interface{}) error {
	// Wrapper simple - en Go 1.21+ usa json.Unmarshal directamente
	// Mantenemos la compatibilidad con el proyecto existente
	return json.Unmarshal(data, v)
}
