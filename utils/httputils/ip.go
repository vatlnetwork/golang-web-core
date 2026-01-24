package httputils

import (
	"net"
	"net/http"
	"strings"
)

// GetRealIP extracts the real IP address from an HTTP request.
// It checks various headers that reverse proxies (like nginx) use to forward the original client IP.
// The order of precedence is:
// 1. X-Real-IP (commonly used by nginx)
// 2. X-Forwarded-For (standard header, may contain multiple IPs)
// 3. CF-Connecting-IP (Cloudflare)
// 4. True-Client-IP (Akamai and Cloudflare)
// 5. RemoteAddr (fallback to direct connection IP)
func GetRealIP(r *http.Request) (string, error) {
	// Check X-Real-IP header (nginx default)
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		ip := strings.TrimSpace(realIP)
		if parsedIP := net.ParseIP(ip); parsedIP != nil {
			return ip, nil
		}
	}

	// Check X-Forwarded-For header (may contain multiple IPs)
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		// X-Forwarded-For can be a comma-separated list: "client, proxy1, proxy2"
		// The first IP is the original client
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if parsedIP := net.ParseIP(ip); parsedIP != nil {
				return ip, nil
			}
		}
	}

	// Check Cloudflare specific header
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		ip := strings.TrimSpace(cfIP)
		if parsedIP := net.ParseIP(ip); parsedIP != nil {
			return ip, nil
		}
	}

	// Check True-Client-IP (Akamai and Cloudflare)
	if trueClientIP := r.Header.Get("True-Client-IP"); trueClientIP != "" {
		ip := strings.TrimSpace(trueClientIP)
		if parsedIP := net.ParseIP(ip); parsedIP != nil {
			return ip, nil
		}
	}

	// Fallback to RemoteAddr
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	return remoteIP, nil
}

// GetUserAgent extracts the User-Agent header from an HTTP request.
func GetUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}
