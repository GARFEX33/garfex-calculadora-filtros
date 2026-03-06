// internal/pdf/infrastructure/adapter/driven/gotenberg/config.go
package gotenberg

import (
	"os"
	"time"
)

// Config representa la configuración del cliente Gotenberg.
// Se carga desde variables de entorno.
type Config struct {
	// URL del servicio Gotenberg (incluye el endpoint de conversión).
	// Default: "http://localhost:3000/convert/html"
	URL string

	// Timeout para la generación de un PDF individual.
	// Default: 60 segundos
	Timeout time.Duration

	// MaxRetries número máximo de reintentos en caso de error transitorio.
	// Default: 3
	MaxRetries int
}

// NewConfig crea una nueva configuración leyendo las variables de entorno.
// Variables de entorno:
//   - GOTENBERG_URL: URL del servicio Gotenberg (default: http://localhost:3000/convert/html)
//   - GOTENBERG_TIMEOUT: Timeout en formato parseable por time.ParseDuration (default: 60s)
//   - GOTENBERG_MAX_RETRIES: Número máximo de reintentos (default: 3)
func NewConfig() *Config {
	url := getEnv("GOTENBERG_URL", "http://localhost:3000/convert/html")
	timeoutStr := getEnv("GOTENBERG_TIMEOUT", "60s")
	maxRetries := 3

	if retriesStr := os.Getenv("GOTENBERG_MAX_RETRIES"); retriesStr != "" {
		if retries := parseIntOrDefault(retriesStr, 3); retries > 0 {
			maxRetries = retries
		}
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		timeout = 60 * time.Second
	}

	return &Config{
		URL:        url,
		Timeout:    timeout,
		MaxRetries: maxRetries,
	}
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseIntOrDefault parses a string to int and returns default on error.
func parseIntOrDefault(s string, defaultValue int) int {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return defaultValue
		}
		n = n*10 + int(c-'0')
	}
	return n
}
