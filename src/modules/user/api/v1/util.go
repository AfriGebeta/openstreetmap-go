package v1

import (
	"net"
	"net/http"
	"regexp"
	"strings"
)

func SanitizeForRegex(displayName string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]")
	return reg.ReplaceAllString(displayName, "")
}

func getClientIP(r *http.Request) string {

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")

		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])

			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
