// internal/pdf/infrastructure/adapter/driven/gotenberg/config.go
package gotenberg

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Config representa la configuración del cliente Gotenberg.
// Se carga desde variables de entorno.
type Config struct {
	// URL del servicio Gotenberg (incluye el endpoint de conversión).
	// Default: "http://localhost:3000/forms/chromium/convert/html"
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
//   - GOTENBERG_URL: URL del servicio Gotenberg (default: http://localhost:3000/forms/chromium/convert/html)
//     En WSL, se detecta automáticamente la IP de Windows si GOTENBERG_URL no está definido
//   - GOTENBERG_TIMEOUT: Timeout en formato parseable por time.ParseDuration (default: 60s)
//   - GOTENBERG_MAX_RETRIES: Número máximo de reintentos (default: 3)
func NewConfig() *Config {
	url := getGotenbergURL()
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

// isWSL checks if the current environment is running under WSL.
func isWSL() bool {
	// Check /proc/version for "Microsoft" or "WSL" marker
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft") ||
		strings.Contains(strings.ToLower(string(data)), "wsl")
}

// getWindowsIP extracts the Windows host IP from /etc/resolv.conf in WSL.
// WSL creates a symlink to /run/wsl/resolv.conf, and the nameserver entry
// points to the Windows host gateway.
func getWindowsIP() string {
	// Try /etc/resolv.conf first (WSL1 and WSL2)
	data, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		// Try alternative path for WSL2
		data, err = os.ReadFile("/run/wsl/resolv.conf")
		if err != nil {
			return ""
		}
	}

	// Parse nameserver line: nameserver 172.x.x.x
	re := regexp.MustCompile(`nameserver\s+(\d+\.\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) > 1 {
		ip := matches[1]
		// Verify it's not a local address (127.x.x.x or 0.x.x.x)
		if !strings.HasPrefix(ip, "127.") && !strings.HasPrefix(ip, "0.") {
			return ip
		}
	}
	return ""
}

// getGotenbergURL returns the Gotenberg URL, automatically detecting
// the Windows host IP when running under WSL.
func getGotenbergURL() string {
	// If user explicitly set GOTENBERG_URL, use it
	if url := os.Getenv("GOTENBERG_URL"); url != "" {
		return url
	}

	// Default URL (Gotenberg 8.x API)
	defaultURL := "http://localhost:3000/forms/chromium/convert/html"

	// If running in WSL, try to detect Windows host IP
	if isWSL() {
		windowsIP := getWindowsIP()
		if windowsIP != "" {
			return fmt.Sprintf("http://%s:3000/forms/chromium/convert/html", windowsIP)
		}
		// Log a warning if we couldn't detect the IP (optional)
	}

	return defaultURL
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
