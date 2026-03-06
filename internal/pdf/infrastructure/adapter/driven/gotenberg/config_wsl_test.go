package gotenberg

import (
	"os"
	"testing"
)

func TestIsWSL(t *testing.T) {
	// Should detect WSL on this machine
	result := isWSL()
	t.Logf("isWSL() = %v", result)
}

func TestGetWindowsIP(t *testing.T) {
	ip := getWindowsIP()
	t.Logf("Windows IP = %v", ip)

	if ip == "" {
		t.Error("Expected to get a Windows IP in WSL environment")
	}

	// Should not be empty and should be a valid IP format
	if len(ip) < 7 || len(ip) > 15 {
		t.Errorf("Invalid IP format: %s", ip)
	}
}

func TestGetGotenbergURL(t *testing.T) {
	url := getGotenbergURL()
	t.Logf("Gotenberg URL = %s", url)

	// In WSL without env override, should contain the Windows IP
	if isWSL() && os.Getenv("GOTENBERG_URL") == "" {
		ip := getWindowsIP()
		expectedPrefix := "http://" + ip + ":3000"
		if url != expectedPrefix+"/forms/chromium/convert/html" {
			t.Errorf("Expected URL to start with %s, got %s", expectedPrefix, url)
		}
	}
}

func TestGetGotenbergURLWithEnvOverride(t *testing.T) {
	// This test would need to be run with GOTENBERG_URL set
	// Skipping actual env var test to not interfere with other tests
	t.Log("To test env override, run: GOTENBERG_URL=http://custom:4000 go test -v -run TestGetGotenbergURLWithEnvOverride")
}
