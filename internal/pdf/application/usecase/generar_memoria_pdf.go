// internal/pdf/application/usecase/generar_memoria_pdf.go
package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/garfex/calculadora-filtros/internal/pdf"
	"github.com/garfex/calculadora-filtros/internal/pdf/application/dto"
	"github.com/garfex/calculadora-filtros/internal/pdf/application/port"
	"github.com/garfex/calculadora-filtros/internal/pdf/domain"
)

const (
	// templateName es el nombre del template HTML a usar para la memoria de cálculo.
	templateName = "memoria_calculo.html"

	// headerTemplateName es el nombre del template de header para Gotenberg.
	headerTemplateName = "gotenberg_header.html"

	// footerTemplateName es el nombre del template de footer.
	footerTemplateName = "gotenberg_footer.html"

	// fechaLayout es el formato de fecha para mostrar en el PDF.
	fechaLayout = "02/01/2006"
)

// GenerarMemoriaPdfUseCase orquesta la generación de la memoria de cálculo en PDF.
// Coordina el renderizado HTML y la conversión a PDF con control de concurrencia.
type GenerarMemoriaPdfUseCase struct {
	renderer  port.HtmlRenderer
	generator port.PdfGenerator
	semaforo  chan struct{}
}

// NewGenerarMemoriaPdf crea una nueva instancia del use case con control de concurrencia.
// maxConcurrent limita el número de generaciones de PDF simultáneas (recomendado: 3).
func NewGenerarMemoriaPdf(
	renderer port.HtmlRenderer,
	generator port.PdfGenerator,
	maxConcurrent int,
) *GenerarMemoriaPdfUseCase {
	if maxConcurrent <= 0 {
		maxConcurrent = 3
	}
	return &GenerarMemoriaPdfUseCase{
		renderer:  renderer,
		generator: generator,
		semaforo:  make(chan struct{}, maxConcurrent),
	}
}

// Execute genera la memoria de cálculo en PDF a partir del request.
// Flujo: resolver empresa → cargar logo → construir TemplateData → adquirir semáforo →
// renderizar HTML → generar PDF → liberar semáforo → retornar bytes.
func (uc *GenerarMemoriaPdfUseCase) Execute(
	ctx context.Context,
	req dto.PdfMemoriaRequest,
) ([]byte, error) {
	// 1. Resolver empresa del catálogo estático
	empresa, ok := domain.BuscarEmpresaPorID(req.Presentacion.EmpresaID)
	if !ok {
		return nil, fmt.Errorf("%w: id=%q", domain.ErrEmpresaNoEncontrada, req.Presentacion.EmpresaID)
	}

	// 2. Cargar logos desde filesystem y codificar en base64
	// Graceful degradation: si no se puede cargar, continuar sin logo
	// - LogoBase64: logo completo para el header principal (garfex.png, summa.png, siemens.png)
	// - LogoLetraBase64: solo para Garfex (lg.png), vacío para otras empresas
	logoBase64 := cargarLogoBase64(empresa.LogoPath)
	var logoLetraBase64 string
	if empresa.ID == "garfex" {
		logoLetraBase64 = cargarLogoBase64("assets/logos/lg.png")
	}

	// 3. Determinar nombre del equipo
	nombreEquipo := req.Presentacion.NombreEquipoOverride
	if nombreEquipo == "" {
		nombreEquipo = req.Memoria.Equipo.Clave
	}

	// 4. Construir TemplateData
	data := dto.TemplateData{
		Empresa:           empresa,
		LogoBase64:        logoBase64,
		LogoLetraBase64:   logoLetraBase64,
		NombreProyecto:    req.Presentacion.NombreProyecto,
		DireccionProyecto: req.Presentacion.DireccionProyecto,
		Responsable:       req.Presentacion.Responsable,
		NombreEquipo:      nombreEquipo,
		Memoria:           req.Memoria,
		FechaGeneracion:   time.Now().Format(fechaLayout),
	}

	// 5. Adquirir semáforo con timeout del contexto para limitar concurrencia
	select {
	case uc.semaforo <- struct{}{}:
		// semáforo adquirido
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout esperando turno de generación: %w", ctx.Err())
	}
	defer func() { <-uc.semaforo }()

	// 6. Renderizar template HTML (el CSS ya está embebido en el template con variables dinámicas)
	html, err := uc.renderer.Render(templateName, data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrRenderizadoHtml, err)
	}

	// 7. Renderizar footer con los datos de la empresa para que las variables
	// como {{.Empresa.NombreCompleto}} se resuelvan correctamente
	footerHTML, err := uc.renderer.Render(footerTemplateName, data)
	if err != nil {
		// Si falla el footer, continuar sin él (graceful degradation)
		footerHTML = ""
	}

	// 8. Renderizar header con los datos de la empresa para que las variables
	// como {{.Empresa.NombreCompleto}} se resuelvan correctamente
	headerHTML, err := uc.renderer.Render(headerTemplateName, data)
	if err != nil {
		// Si falla el header, continuar sin él (graceful degradation)
		headerHTML = ""
	}

	// 9. Generar PDF desde HTML con el header y footer renderizados
	pdfBytes, err := uc.generator.GenerateWithHeaderFooter(ctx, html, headerHTML, footerHTML)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrGeneracionPdf, err)
	}

	return pdfBytes, nil
}

// cargarLogoBase64 carga un archivo de imagen desde el filesystem embebido y lo codifica en base64.
// Retorna una cadena vacía si el archivo no existe o no se puede leer (graceful degradation).
func cargarLogoBase64(logoPath string) string {
	if logoPath == "" {
		return ""
	}

	data, err := pdf.AssetsFS.ReadFile(logoPath)
	if err != nil {
		// Graceful degradation: logo no encontrado → continuar sin logo
		return ""
	}

	return base64.StdEncoding.EncodeToString(data)
}
