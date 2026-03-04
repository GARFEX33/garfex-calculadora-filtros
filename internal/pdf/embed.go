// internal/pdf/embed.go
// Este archivo registra los embed.FS para los templates HTML y assets del módulo PDF.
// Los archivos se incrustan en el binario Go en tiempo de compilación.
package pdf

import "embed"

// TemplatesFS contiene todos los templates HTML de la memoria de cálculo.
// Incluye el template principal (memoria.html) y todos los partials (partials/*.html).
//
//go:embed templates
var TemplatesFS embed.FS

// AssetsFS contiene los assets estáticos del módulo PDF (logos, imágenes).
// Los logos se leen y codifican en base64 en el use case para incrustarlos en el HTML.
//
//go:embed assets
var AssetsFS embed.FS
