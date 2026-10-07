package main

import (
	"testing"

	"github.com/input-api/mcp-server/config"
)

func TestGetAll_ReturnsEmptySlice(t *testing.T) {
	cfg := &config.APIConfig{
		BaseURL: "https://api.example.com",
	}

	result := GetAll(cfg)
	if result == nil {
		t.Fatal("GetAll() returned nil, want empty non-nil slice")
	}
	if len(result) != 0 {
		t.Errorf("GetAll() len = %d, want 0", len(result))
	}
}
