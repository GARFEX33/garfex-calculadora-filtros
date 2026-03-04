// internal/pdf/infrastructure/adapter/driven/template/html_renderer.go
// Package htmltemplate implementa port.HtmlRenderer usando html/template con embed.FS.
package htmltemplate

import (
	"bytes"
	"fmt"
	htmpl "html/template"
	"io/fs"
	"math"
	"strings"

	"github.com/garfex/calculadora-filtros/internal/pdf/application/dto"
)

// HtmlRendererAdapter implementa port.HtmlRenderer usando html/template con embed.FS.
// Los templates se parsean en la construcción (fail-fast) y se reutilizan en cada llamada.
type HtmlRendererAdapter struct {
	tmpl *htmpl.Template
}

// NewHtmlRenderer crea un HtmlRendererAdapter parseando todos los templates de templatesFS.
// Falla si algún template no puede ser parseado (fail-fast al inicio de la aplicación).
func NewHtmlRenderer(templatesFS fs.FS) (*HtmlRendererAdapter, error) {
	funcMap := htmpl.FuncMap{
		// upper convierte un string a mayúsculas
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		// capFirst retorna el primer carácter en mayúscula
		// uso: {{capFirst .Empresa.ID}} retorna "G" para "garfex"
		"capFirst": func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToUpper(s[:1])
		},
		// slice retorna un substring desde start hasta end (exclusivo)
		// uso: {{slice .Empresa.ID 0 1}} retorna el primer carácter
		"slice": func(s string, start, end int) string {
			if start < 0 || start > len(s) || end < start || end > len(s) {
				return ""
			}
			return s[start:end]
		},
		// formatFloat formatea un float64 con n decimales
		"formatFloat": func(f float64, decimals int) string {
			format := fmt.Sprintf("%%.%df", decimals)
			return fmt.Sprintf(format, f)
		},
		// formatFloat2 formatea con 2 decimales (shorthand común)
		"formatFloat2": func(f float64) string {
			return fmt.Sprintf("%.2f", f)
		},
		// formatFloat4 formatea con 4 decimales para impedancias
		"formatFloat4": func(f float64) string {
			return fmt.Sprintf("%.4f", f)
		},
		// formatInt convierte float64 a int para mostrar valores enteros
		"formatInt": func(f float64) string {
			return fmt.Sprintf("%.0f", f)
		},
		// mul multiplica dos float64 (útil para cálculos en template)
		"mul": func(a, b float64) float64 {
			return a * b
		},
		// div divide dos float64 con protección contra división por cero
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		// sub resta dos float64
		"sub": func(a, b float64) float64 {
			return a - b
		},
		// sqrt retorna la raíz cuadrada
		"sqrt": func(f float64) float64 {
			return math.Sqrt(f)
		},
		// safeHTML permite HTML sin escape (para contenido controlado)
		"safeHTML": func(s string) htmpl.HTML {
			return htmpl.HTML(s) //nolint:gosec // contenido interno controlado
		},
		// contains verifica si un string contiene un substring
		"contains": func(s, sub string) bool {
			return containsStr(s, sub)
		},
		// notNil verifica que un puntero no sea nil
		"notNil": func(v interface{}) bool {
			return v != nil
		},
		// esMaterial verifica el material del conductor para la etiqueta
		"esMaterial": func(material, expected string) bool {
			return material == expected
		},
		// itoa convierte int a string
		"itoa": func(i int) string {
			return fmt.Sprintf("%d", i)
		},
		// percent multiplica por 100 y formatea con 0 decimales
		"percent": func(f float64) string {
			return fmt.Sprintf("%.0f", f*100)
		},
		// not es la negación booleana (útil en condiciones de template)
		"not": func(b bool) bool {
			return !b
		},
		// mulIntFloat multiplica int × float64 → float64 (para HilosPorFase × SeccionMM2)
		"mulIntFloat": func(i int, f float64) float64 {
			return float64(i) * f
		},
		// intToFloat convierte int a float64
		"intToFloat": func(i int) float64 {
			return float64(i)
		},
		// toFloat64 convierte cualquier tipo numérico a float64 (int, uint, float32, float64)
		"toFloat64": func(v interface{}) float64 {
			return convToFloat64(v)
		},
		// formatNumeric acepta cualquier tipo numérico y devuelve un string formateado con 2 decimales
		"formatNumeric": func(v interface{}) string {
			return fmt.Sprintf("%.2f", convToFloat64(v))
		},
		// derefFloat desreferencia un puntero *float64 de forma segura (retorna 0 si nil)
		"derefFloat": func(f *float64) float64 {
			if f == nil {
				return 0
			}
			return *f
		},
	}

	// Parsear template principal + todos los partials
	tmpl, err := htmpl.New("").Funcs(funcMap).ParseFS(templatesFS,
		"templates/memoria.html",
		"templates/partials/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("parseando templates: %w", err)
	}

	return &HtmlRendererAdapter{tmpl: tmpl}, nil
}

// Render aplica los datos al template identificado por templateName y retorna el HTML completo.
// Implementa port.HtmlRenderer.
func (r *HtmlRendererAdapter) Render(templateName string, data dto.TemplateData) (string, error) {
	var buf bytes.Buffer
	if err := r.tmpl.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", fmt.Errorf("ejecutando template %q: %w", templateName, err)
	}
	return buf.String(), nil
}

// containsStr es una función auxiliar para verificar si s contiene sub.
func containsStr(s, sub string) bool {
	if len(sub) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// convToFloat64 es una función auxiliar que convierte cualquier tipo numérico a float64.
// Utilizada por las funciones de template toFloat64 y formatNumeric.
func convToFloat64(v interface{}) float64 {
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
