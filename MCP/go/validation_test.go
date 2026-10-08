package main

import (
	"net"
	"strings"
	"testing"
)

func stubPublicLookup(t *testing.T) {
	t.Helper()
	orig := lookupHost
	lookupHost = func(host string) ([]string, error) {
		return []string{"93.184.216.34"}, nil // public documentation IP
	}
	t.Cleanup(func() { lookupHost = orig })
}

func TestValidateBaseURL_ValidURL(t *testing.T) {
	stubPublicLookup(t)
	err := validateBaseURL("https://api.example.com/v1", []string{"api.example.com"})
	if err != nil {
		t.Fatalf("validateBaseURL() unexpected error: %v", err)
	}
}

func TestValidateBaseURL_InvalidScheme(t *testing.T) {
	cases := []string{
		"file:///etc/passwd",
		"gopher://internal",
		"ftp://files.com",
	}
	for _, raw := range cases {
		t.Run(raw, func(t *testing.T) {
			err := validateBaseURL(raw, []string{"files.com", "internal", "etc"})
			if err == nil {
				t.Fatal("validateBaseURL() error = nil, want scheme error")
			}
			if !strings.Contains(strings.ToLower(err.Error()), "scheme") {
				t.Errorf("validateBaseURL() error = %q, want substring %q", err.Error(), "scheme")
			}
		})
	}
}

func TestValidateBaseURL_NonAllowlistedHost(t *testing.T) {
	err := validateBaseURL("https://evil.com", []string{"api.example.com"})
	if err == nil {
		t.Fatal("validateBaseURL() error = nil, want allowlist error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "allowlist") {
		t.Errorf("validateBaseURL() error = %q, want substring %q", err.Error(), "allowlist")
	}
}

func TestValidateBaseURL_MissingHost(t *testing.T) {
	err := validateBaseURL("https:///no-host", []string{"api.example.com"})
	if err == nil {
		t.Fatal("validateBaseURL() error = nil, want missing host error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "host") {
		t.Errorf("validateBaseURL() error = %q, want substring %q", err.Error(), "host")
	}
}

func TestValidateBaseURL_EmptyAllowlist(t *testing.T) {
	err := validateBaseURL("https://api.example.com/v1", nil)
	if err == nil {
		t.Fatal("validateBaseURL() error = nil, want allowlist rejection for empty allowlist")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "allowlist") {
		t.Errorf("validateBaseURL() error = %q, want substring %q", err.Error(), "allowlist")
	}
}

func TestValidateBaseURL_PrivateIPs(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		host string
	}{
		{name: "loopback", raw: "http://127.0.0.1/", host: "127.0.0.1"},
		{name: "rfc1918_10", raw: "http://10.0.0.1/", host: "10.0.0.1"},
		{name: "rfc1918_172", raw: "http://172.16.0.1/", host: "172.16.0.1"},
		{name: "rfc1918_192", raw: "http://192.168.1.1/", host: "192.168.1.1"},
		{name: "link_local_metadata", raw: "http://169.254.169.254/", host: "169.254.169.254"},
		{name: "ipv6_loopback", raw: "http://[::1]/", host: "::1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateBaseURL(tc.raw, []string{tc.host})
			if err == nil {
				t.Fatal("validateBaseURL() error = nil, want private/blocked error")
			}
			msg := strings.ToLower(err.Error())
			if !strings.Contains(msg, "private") && !strings.Contains(msg, "blocked") {
				t.Errorf("validateBaseURL() error = %q, want substring %q or %q", err.Error(), "private", "blocked")
			}
		})
	}
}

func TestValidateBaseURL_WithPort(t *testing.T) {
	stubPublicLookup(t)
	err := validateBaseURL("https://api.example.com:8443/v1", []string{"api.example.com"})
	if err != nil {
		t.Fatalf("validateBaseURL() with port unexpected error: %v", err)
	}
}

func TestValidateBaseURL_UserinfoRejected(t *testing.T) {
	err := validateBaseURL("https://user:pass@api.example.com/v1", []string{"api.example.com"})
	if err == nil {
		t.Fatal("validateBaseURL() error = nil, want userinfo rejection")
	}
}

func TestValidateBaseURL_AllowlistCaseInsensitive(t *testing.T) {
	stubPublicLookup(t)
	err := validateBaseURL("https://API.Example.COM/v1", []string{"api.example.com"})
	if err != nil {
		t.Fatalf("validateBaseURL() case-insensitive allowlist unexpected error: %v", err)
	}
}

func TestValidateBaseURL_MixedPublicPrivateIPs(t *testing.T) {
	orig := lookupHost
	lookupHost = func(host string) ([]string, error) {
		return []string{"93.184.216.34", "10.0.0.1"}, nil
	}
	t.Cleanup(func() { lookupHost = orig })

	err := validateBaseURL("https://api.example.com/v1", []string{"api.example.com"})
	if err == nil {
		t.Fatal("validateBaseURL() error = nil, want rejection when any resolved IP is private")
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "private") && !strings.Contains(msg, "blocked") {
		t.Errorf("validateBaseURL() error = %q, want substring %q or %q", err.Error(), "private", "blocked")
	}
}

func TestValidateBaseURL_DNSFailure(t *testing.T) {
	orig := lookupHost
	lookupHost = func(host string) ([]string, error) {
		return nil, &net.DNSError{Err: "no such host", Name: host, IsNotFound: true}
	}
	t.Cleanup(func() { lookupHost = orig })

	err := validateBaseURL("https://missing.example.com/", []string{"missing.example.com"})
	if err == nil {
		t.Fatal("validateBaseURL() error = nil, want DNS failure")
	}
}

func TestValidateToken_Empty(t *testing.T) {
	err := validateToken("")
	if err == nil {
		t.Fatal("validateToken() error = nil, want empty-token error")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	err := validateToken("valid-token-123")
	if err != nil {
		t.Fatalf("validateToken() unexpected error: %v", err)
	}
}

func TestValidateToken_Oversized(t *testing.T) {
	token := strings.Repeat("a", maxTokenLength+1)
	err := validateToken(token)
	if err == nil {
		t.Fatal("validateToken() error = nil, want oversized-token error")
	}
}

func TestValidateHeaders_OversizedURL(t *testing.T) {
	err := validateHeaders(strings.Repeat("u", maxURLLength+1), "tok", "", "")
	if err == nil {
		t.Fatal("validateHeaders() error = nil, want oversized URL error")
	}
}

func TestValidateHeaders_NonPrintable(t *testing.T) {
	err := validateHeaders("https://api.example.com", "tok\x00en", "", "")
	if err == nil {
		t.Fatal("validateHeaders() error = nil, want non-printable error")
	}
}

func TestIsPrivateIP(t *testing.T) {
	cases := []struct {
		ip   string
		want bool
	}{
		{"127.0.0.1", true},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"192.168.1.1", true},
		{"169.254.169.254", true},
		{"::1", true},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
	}
	for _, tc := range cases {
		got := isPrivateIP(net.ParseIP(tc.ip))
		if got != tc.want {
			t.Errorf("isPrivateIP(%s) = %v, want %v", tc.ip, got, tc.want)
		}
	}
}
