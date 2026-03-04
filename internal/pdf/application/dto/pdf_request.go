// internal/pdf/application/dto/pdf_request.go
package dto

import (
	calculosdto "github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/pdf/domain"
)

// PresentacionInput contiene los datos de presentación ingresados por el usuario
// en el formulario de configuración de PDF.
type PresentacionInput struct {
	// EmpresaID es el identificador de la empresa del catálogo estático (ej: "garfex", "summa", "siemens").
	EmpresaID string `json:"empresa_id"`

	// NombreProyecto es el nombre del proyecto o instalación.
	NombreProyecto string `json:"nombre_proyecto"`

	// DireccionProyecto es la dirección donde se realiza la instalación.
	DireccionProyecto string `json:"direccion_proyecto"`

	// Responsable es el nombre del ingeniero o técnico responsable del cálculo.
	Responsable string `json:"responsable"`

	// NombreEquipoOverride permite sobreescribir el nombre del equipo en la memoria.
	// Si está vacío, se usará el nombre obtenido de la MemoriaOutput.
	NombreEquipoOverride string `json:"nombre_equipo_override,omitempty"`
}

// PdfMemoriaRequest combina el resultado de cálculo con los datos de presentación
// para generar la memoria de cálculo en PDF.
type PdfMemoriaRequest struct {
	// Memoria contiene el resultado completo del cálculo eléctrico.
	Memoria calculosdto.MemoriaOutput `json:"memoria"`

	// Presentacion contiene los datos de presentación para el PDF.
	Presentacion PresentacionInput `json:"presentacion"`
}

// TemplateData es el struct que alimenta el template HTML de la memoria de cálculo.
// Agrupa todos los datos necesarios para renderizar el PDF.
type TemplateData struct {
	// Empresa contiene los datos de la empresa presentadora.
	Empresa domain.EmpresaPresentacion

	// LogoBase64 es el logo de la empresa codificado en base64 para incrustar en el HTML.
	// Si el logo no se puede cargar, este campo queda vacío (graceful degradation).
	LogoBase64 string

	// NombreProyecto es el nombre del proyecto o instalación.
	NombreProyecto string

	// DireccionProyecto es la dirección de la instalación.
	DireccionProyecto string

	// Responsable es el nombre del responsable del cálculo.
	Responsable string

	// NombreEquipo es el nombre del equipo (puede ser override o derivado de la memoria).
	NombreEquipo string

	// Memoria contiene el resultado completo del cálculo eléctrico.
	Memoria calculosdto.MemoriaOutput

	// FechaGeneracion es la fecha de generación del documento (formato: "02/01/2006").
	FechaGeneracion string
}
