// cmd/pdf_preview/main.go
// Servidor de preview para desarrollo de templates PDF.
// Renderiza HTML directamente en el navegador para feedback instantáneo.
// Uso: go run cmd/pdf_preview/main.go
//
// Endpoints:
//   - GET /                    : Renderiza la memoria completa
//   - GET /?empresa=garfex     : Cambiar empresa (garfex, summaa, siemens)
package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	calculosdto "github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	pdfdto "github.com/garfex/calculadora-filtros/internal/pdf/application/dto"
	pdfdomain "github.com/garfex/calculadora-filtros/internal/pdf/domain"
	pdftemplate "github.com/garfex/calculadora-filtros/internal/pdf/infrastructure/adapter/driven/template"
)

const (
	templatesBasePath = "internal/pdf"
	logosPath         = "internal/pdf/assets/logos"
	testDataPath      = "test_memoria.json"
	serverPort        = "3000"
)

var (
	empresaID = flag.String("empresa", "garfex", "Empresa por defecto (garfex, summaa, siemens)")
	noOpen    = flag.Bool("no-open", false, "No abrir navegador automáticamente")
	port      = flag.String("port", serverPort, "Puerto del servidor")
)

// TestData representa la estructura del JSON de prueba
type TestData struct {
	Memoria      calculosdto.MemoriaOutput `json:"memoria"`
	Presentacion *PresentacionInput        `json:"presentacion"`
}

// PresentacionInput contiene los datos de presentación (re-definido para JSON parsing)
type PresentacionInput struct {
	EmpresaID            string `json:"empresa_id"`
	NombreProyecto       string `json:"nombre_proyecto"`
	DireccionProyecto    string `json:"direccion_proyecto"`
	Responsable          string `json:"responsable"`
	NombreEquipoOverride string `json:"nombre_equipo_override,omitempty"`
}

func main() {
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("🌐 ")

	log.Printf("🎯 PDF Preview Server Started")
	log.Printf("   URL:      http://localhost:%s", *port)
	log.Printf("   Empresa:  %s", *empresaID)
	log.Printf("   Templates: %s", templatesBasePath)
	log.Printf("")
	log.Printf("   Presiona Ctrl+C para detener")
	log.Printf("")

	// Verificar empresa existe
	if _, ok := pdfdomain.BuscarEmpresaPorID(*empresaID); !ok {
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

	// Validar empresa usando el catálogo del dominio
	empresa, ok := pdfdomain.BuscarEmpresaPorID(empID)
	if !ok {
		http.Error(w, fmt.Sprintf("Empresa desconocida: %s. Opciones: garfex, summaa, siemens", empID), http.StatusBadRequest)
		return
	}

	// Renderizar HTML (re-crea renderer para hot-reload)
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

	err := filepath.Walk(templatesBasePath, func(path string, info os.FileInfo, err error) error {
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

func renderHTMLForEmpresa(empresa pdfdomain.EmpresaPresentacion) (string, error) {
	// Cargar datos de prueba
	testData, err := loadTestData(testDataPath)
	if err != nil {
		return "", fmt.Errorf("cargando datos de prueba: %w", err)
	}

	// Actualizar empresa_id en presentacion
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

	// Construir TemplateData usando el DTO compartido
	data := templateDataFromTestData(testData, empresa, logoBase64, logoLetraBase64, nombreEquipo)

	// Usar os.DirFS para hot-reload - crea nuevo renderer en cada request
	// Esto permite re-parsear templates desde disco en cada request
	diskFS := os.DirFS(templatesBasePath)
	renderer, err := pdftemplate.NewHtmlRenderer(diskFS)
	if err != nil {
		return "", fmt.Errorf("creando renderer HTML: %w", err)
	}

	// Render usando "memoria.html" que es el nombre del define en memoria.html
	html, err := renderer.Render("memoria.html", data)
	if err != nil {
		return "", fmt.Errorf("renderizando template: %w", err)
	}

	// Inject hot-reload script
	html = injectHotReloadScript(html)

	return html, nil
}

// templateDataFromTestData construye el DTO TemplateData desde los datos de prueba
func templateDataFromTestData(
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
		FechaGeneracion:   time.Now().Format("02/01/2006"),
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
