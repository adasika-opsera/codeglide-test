package config_test

import (
	"strings"
	"testing"

	"github.com/input-api/mcp-server/config"
)

func TestLoadAPIConfig_ValidSTDIO(t *testing.T) {
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("BEARER_TOKEN", "token-123")
	t.Setenv("TRANSPORT", "stdio")
	t.Setenv("API_KEY", "key-456")
	t.Setenv("BASIC_AUTH", "user:pass")
	t.Setenv("PORT", "3000")

	cfg, err := config.LoadAPIConfig()
	if err != nil {
		t.Fatalf("LoadAPIConfig() unexpected error: %v", err)
	}
	if cfg.BaseURL != "https://api.example.com" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://api.example.com")
	}
	if cfg.BearerToken != "token-123" {
		t.Errorf("BearerToken = %q, want %q", cfg.BearerToken, "token-123")
	}
	if cfg.APIKey != "key-456" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "key-456")
	}
	if cfg.BasicAuth != "user:pass" {
		t.Errorf("BasicAuth = %q, want %q", cfg.BasicAuth, "user:pass")
	}
	if cfg.Port != "3000" {
		t.Errorf("Port = %q, want %q", cfg.Port, "3000")
	}
}

func TestLoadAPIConfig_MissingBaseURL_STDIO(t *testing.T) {
	t.Setenv("TRANSPORT", "stdio")
	t.Setenv("API_BASE_URL", "")

	_, err := config.LoadAPIConfig()
	if err == nil {
		t.Fatal("LoadAPIConfig() error = nil, want error containing API_BASE_URL")
	}
	if !strings.Contains(err.Error(), "API_BASE_URL") {
		t.Errorf("LoadAPIConfig() error = %q, want substring %q", err.Error(), "API_BASE_URL")
	}
}

func TestLoadAPIConfig_HTTPMode_NoBaseURL(t *testing.T) {
	t.Setenv("TRANSPORT", "http")
	t.Setenv("PORT", "8080")
	t.Setenv("API_BASE_URL", "")

	cfg, err := config.LoadAPIConfig()
	if err != nil {
		t.Fatalf("LoadAPIConfig() unexpected error: %v", err)
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.BaseURL != "" {
		t.Errorf("BaseURL = %q, want empty string", cfg.BaseURL)
	}
}

func TestLoadAPIConfig_PortResolution(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("TRANSPORT", "stdio")

	cfg, err := config.LoadAPIConfig()
	if err != nil {
		t.Fatalf("LoadAPIConfig() unexpected error: %v", err)
	}
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
}

func TestLoadAPIConfig_LowercaseTransport(t *testing.T) {
	t.Setenv("transport", "http")
	t.Setenv("TRANSPORT", "")
	t.Setenv("API_BASE_URL", "")
	t.Setenv("PORT", "8080")

	cfg, err := config.LoadAPIConfig()
	if err != nil {
		t.Fatalf("LoadAPIConfig() with lowercase transport unexpected error: %v", err)
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
}

func TestLoadAPIConfig_EmptyTransport_RequiresBaseURL(t *testing.T) {
	t.Setenv("TRANSPORT", "")
	t.Setenv("transport", "")
	t.Setenv("API_BASE_URL", "")

	_, err := config.LoadAPIConfig()
	if err == nil {
		t.Fatal("LoadAPIConfig() with empty TRANSPORT error = nil, want STDIO-mode error")
	}
	if !strings.Contains(err.Error(), "API_BASE_URL") {
		t.Errorf("LoadAPIConfig() error = %q, want substring %q", err.Error(), "API_BASE_URL")
	}
}

func TestLoadAPIConfig_LowercasePort(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("port", "7070")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("TRANSPORT", "stdio")

	cfg, err := config.LoadAPIConfig()
	if err != nil {
		t.Fatalf("LoadAPIConfig() unexpected error: %v", err)
	}
	if cfg.Port != "7070" {
		t.Errorf("Port = %q, want %q (lowercase port fallback)", cfg.Port, "7070")
	}
}
