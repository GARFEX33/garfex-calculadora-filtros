// internal/pdf/application/port/html_renderer.go
package port

import (
	"github.com/garfex/calculadora-filtros/internal/pdf/application/dto"
)

// HtmlRenderer es el port driven para renderizar el template HTML de la memoria de cálculo.
// La implementación concreta vive en infrastructure/adapter/driven/template/.
type HtmlRenderer interface {
	// Render aplica los datos al template identificado por templateName y retorna el HTML renderizado.
	// templateName es el nombre del template a usar (ej: "memoria_calculo.html").
	// data contiene todos los datos necesarios para el renderizado.
	// Retorna el HTML como string o un error envuelto con ErrRenderizadoHtml.
	Render(templateName string, data dto.TemplateData) (string, error)
}
