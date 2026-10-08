package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"unicode"
)

const (
	maxURLLength   = 2048
	maxTokenLength = 8192
)

// lookupHost resolves a hostname to address strings. Overridable in tests.
var lookupHost = net.LookupHost

// validateBaseURL parses rawURL and rejects unsafe targets that could enable SSRF.
// Checks: scheme allowlist (http/https), host allowlist, and private/reserved IP ranges.
func validateBaseURL(rawURL string, allowedHosts []string) error {
	if strings.TrimSpace(rawURL) == "" {
		return fmt.Errorf("url is empty: missing host")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("unsupported scheme %q: only http and https are allowed", parsed.Scheme)
	}

	if parsed.User != nil {
		return fmt.Errorf("url must not contain userinfo")
	}

	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("url is missing host")
	}

	if !hostAllowed(host, allowedHosts) {
		return fmt.Errorf("host %q is not in the allowlist", host)
	}

	ips, err := resolveHostIPs(host)
	if err != nil {
		return fmt.Errorf("dns resolution failed for host %q: %w", host, err)
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return fmt.Errorf("host %q resolves to private or blocked address %s", host, ip.String())
		}
	}

	return nil
}

func hostAllowed(host string, allowedHosts []string) bool {
	if len(allowedHosts) == 0 {
		return false
	}
	for _, allowed := range allowedHosts {
		if strings.EqualFold(host, allowed) {
			return true
		}
	}
	return false
}

func resolveHostIPs(host string) ([]net.IP, error) {
	if ip := net.ParseIP(host); ip != nil {
		return []net.IP{ip}, nil
	}

	addrs, err := lookupHost(host)
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("no addresses found")
	}

	ips := make([]net.IP, 0, len(addrs))
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip == nil {
			continue
		}
		ips = append(ips, ip)
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no valid IP addresses found")
	}
	return ips, nil
}

// isPrivateIP reports whether ip is loopback, private, link-local, or otherwise blocked for outbound calls.
func isPrivateIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
		return true
	}

	// Explicit cloud-metadata / link-local range (also covered by IsLinkLocalUnicast for IPv4).
	_, linkLocal, err := net.ParseCIDR("169.254.0.0/16")
	if err == nil && linkLocal.Contains(ip) {
		return true
	}

	return false
}

// validateToken rejects empty tokens and tokens longer than 8192 characters.
func validateToken(token string) error {
	if len(token) == 0 {
		return fmt.Errorf("token is empty")
	}
	if len(token) > maxTokenLength {
		return fmt.Errorf("token exceeds maximum length of %d characters", maxTokenLength)
	}
	return nil
}

// validateHeaders enforces length and printable-character rules on auth-related header values.
func validateHeaders(baseURL string, bearerToken string, apiKey string, basicAuth string) error {
	if err := validateHeaderValue("baseURL", baseURL, maxURLLength); err != nil {
		return err
	}
	if err := validateHeaderValue("bearerToken", bearerToken, maxTokenLength); err != nil {
		return err
	}
	if err := validateHeaderValue("apiKey", apiKey, maxTokenLength); err != nil {
		return err
	}
	if err := validateHeaderValue("basicAuth", basicAuth, maxTokenLength); err != nil {
		return err
	}
	return nil
}

func validateHeaderValue(name, value string, maxLen int) error {
	if value == "" {
		return nil
	}
	if len(value) > maxLen {
		return fmt.Errorf("%s exceeds maximum length of %d characters", name, maxLen)
	}
	for _, r := range value {
		if !unicode.IsPrint(r) {
			return fmt.Errorf("%s contains non-printable characters", name)
		}
	}
	return nil
}
