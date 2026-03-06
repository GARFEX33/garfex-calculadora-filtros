// internal/pdf/infrastructure/adapter/driven/gotenberg/http_client.go
package gotenberg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// HTTPClient define la interfaz para realizar solicitudes HTTP.
// Permite inyectar un cliente mock en tests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GotenbergClient es la interfaz que debe implementar cualquier cliente
// que interactúe con el servicio Gotenberg.
type GotenbergClient interface {
	HTTPClient
	PostMultipart(ctx context.Context, url string, contentType string, content []byte, maxRetries int) (*Response, error)
}

// defaultHTTPClient es la implementación por defecto usando net/http.
type defaultHTTPClient struct {
	client *http.Client
}

// newDefaultHTTPClient crea un cliente HTTP con timeouts configurados.
func newDefaultHTTPClient(timeout time.Duration) *defaultHTTPClient {
	return &defaultHTTPClient{
		client: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// No seguir redirects automáticamente
				return http.ErrUseLastResponse
			},
		},
	}
}

// Do implementa la interfaz HTTPClient de net/http.
// Es un wrapper que permite usar defaultHTTPClient como http.RoundTripper.
func (c *defaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// Response representa la respuesta de una llamada HTTP.
type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// PostMultipart realiza una solicitud POST multipart/form-data con reintentos.
// ctx se usa para el timeout y cancelación.
// content es el body multipart/form-data.
func (c *defaultHTTPClient) PostMultipart(
	ctx context.Context,
	url string,
	contentType string,
	content []byte,
	maxRetries int,
) (*Response, error) {

	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Crear nuevo reader para cada intento (el body se consume)
		bodyReader := bytes.NewReader(content)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("creando request: %w", err)
		}

		req.Header.Set("Content-Type", contentType)

		resp, err := c.client.Do(req)
		if err != nil {
			// Error de red - reintentar si quedan intentos
			log.Printf("[DEBUG] gotenberg/http_client: attempt %d/%d - connection error: %v", attempt+1, maxRetries+1, err)
			if attempt < maxRetries && isRetryableError(err) {
				lastErr = err
				time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
				continue
			}
			return nil, fmt.Errorf("ejecutando request: %w", err)
		}
		defer resp.Body.Close()

		// Leer el body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("leyendo response body: %w", err)
		}

		// Verificar códigos de error HTTP
		if !isSuccess(resp.StatusCode) {
			// Error 5xx del servidor - reintentar
			if attempt < maxRetries && resp.StatusCode >= 500 {
				lastErr = fmt.Errorf("servidor Gotenberg respondió con %d: %s", resp.StatusCode, string(body))
				time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
				continue
			}

			// Error 4xx - no reintentar (error del cliente)
			return nil, fmt.Errorf("Gotenberg respondió con código %d: %s", resp.StatusCode, string(body))
		}

		return &Response{
			StatusCode: resp.StatusCode,
			Body:       body,
			Header:     resp.Header,
		}, nil
	}

	return nil, fmt.Errorf("reintentos agotados: %w", lastErr)
}

// isSuccess retorna true si el código de estado es 2xx.
func isSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// isRetryableError retorna true si el error es transitorio y merece reintento.
func isRetryableError(err error) bool {
	if err == context.DeadlineExceeded || err == context.Canceled {
		return false
	}
	// Errores de red como connection refused, timeout, etc.
	return true
}
