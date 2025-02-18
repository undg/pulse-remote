package utils

import (
	"net"
	"regexp"
	"testing"

	"github.com/undg/go-prapi/api/logger"
)

func TestIsLocalIP(t *testing.T) {
	testCases := []struct {
		name     string
		ip       string
		expected bool
		wantErr  bool
	}{
		// Happy path cases
		{"Loopback IPv4", "127.0.0.1", true, false},
		{"Loopback IPv6", "::1", true, false},
		{"Private IPv4 (10.x.x.x)", "10.0.0.1", true, false},
		{"Private IPv4 (172.16.x.x)", "172.16.0.1", true, false},
		{"Private IPv4 (192.168.x.x)", "192.168.0.1", true, false},
		{"Public IPv4", "8.8.8.8", false, false},
		{"Public IPv6", "2001:4860:4860::8888", false, false},

		// Edge cases
		{"Border of private range (172.15.255.255)", "172.15.255.255", false, false},
		{"Border of private range (172.32.0.0)", "172.32.0.0", false, false},
		{"Broadcast IPv4", "255.255.255.255", false, false},
		{"Unspecified IPv4", "0.0.0.0", false, false},
		{"Unspecified IPv6", "::", false, false},

		// Error cases
		{"Invalid IP", "invalid-ip", false, true},
		{"Empty string", "", false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ip := net.ParseIP(tc.ip)
			if tc.wantErr {
				if ip != nil {
					t.Errorf("Expected error for IP: %s, but got valid IP", tc.ip)
				}
			} else {
				if ip == nil {
					t.Fatalf("Failed to parse IP: %s", tc.ip)
				}
				result := IsLocalIP(ip)
				if result != tc.expected {
					t.Errorf("IsLocalIP(%s) = %v, want %v", tc.ip, result, tc.expected)
				}
			}
		})
	}
}

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()

	if err != nil {
		logger.Error().Err(err).Msg("problem with GetLocalIP()")
	}

	ipPattern := `^(\d{1,3}\.){3}\d{1,3}$`
	match, _ := regexp.MatchString(ipPattern, ip)
	if !match {
		t.Errorf("Invalid IP format: %s", ip)
	}
}
