// internal/pdf/infrastructure/adapter/driven/gotenberg/form_builder.go
package gotenberg

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

// FormBuilder construye el cuerpo multipart/form-data para Gotenberg.
type FormBuilder struct {
	buffer *bytes.Buffer
	writer *multipart.Writer
}

// NewFormBuilder crea un nuevo constructor de formularios multipart.
func NewFormBuilder() *FormBuilder {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	return &FormBuilder{
		buffer: buffer,
		writer: writer,
	}
}

// AddHTML añade el contenido HTML principal como "index.html".
func (fb *FormBuilder) AddHTML(htmlContent string) error {
	return fb.addFile("index.html", "text/html", strings.NewReader(htmlContent))
}

// AddHeader añade el contenido del header HTML.
// Si headerContent está vacío, no añade nada (header opcional).
func (fb *FormBuilder) AddHeader(headerContent string) error {
	if headerContent == "" {
		return nil
	}
	return fb.addFile("header.html", "text/html", strings.NewReader(headerContent))
}

// AddFooter añade el contenido del footer HTML.
// Si footerContent está vacío, no añade nada (footer opcional).
func (fb *FormBuilder) AddFooter(footerContent string) error {
	if footerContent == "" {
		return nil
	}
	return fb.addFile("footer.html", "text/html", strings.NewReader(footerContent))
}

// AddCSS añade CSS embebido que se aplicará a todas las páginas.
// Gotenberg soporta esto a través del campo "styles" en el formulario.
func (fb *FormBuilder) AddCSS(cssContent string) error {
	if cssContent == "" {
		return nil
	}

	return fb.writer.WriteField("styles", cssContent)
}

// AddOption añade una opción de conversión de Gotenberg.
// Referencia: https://gotenberg.dev/docs/modules/chromium#pdf-formats
func (fb *FormBuilder) AddOption(key, value string) error {
	return fb.writer.WriteField(key, value)
}

// addFile añade un archivo al formulario con el nombre de campo especificado.
func (fb *FormBuilder) addFile(fieldName, contentType string, reader io.Reader) error {
	part, err := fb.writer.CreateFormFile(fieldName, fieldName)
	if err != nil {
		return fmt.Errorf("creando form file %q: %w", fieldName, err)
	}

	if _, err := io.Copy(part, reader); err != nil {
		return fmt.Errorf("copiando contenido a %q: %w", fieldName, err)
	}

	return nil
}

// ContentType retorna el Content-Type del formulario (incluye boundary).
func (fb *FormBuilder) ContentType() string {
	return fb.writer.FormDataContentType()
}

// Build retorna el cuerpo completo del formulario como slice de bytes.
func (fb *FormBuilder) Build() ([]byte, error) {
	if err := fb.writer.Close(); err != nil {
		return nil, fmt.Errorf("cerrando multipart writer: %w", err)
	}

	return fb.buffer.Bytes(), nil
}
