package helper

import (
	neturl "net/url"
	"os"
	"strings"
)

func EnforceHTTP(url string) string {
	if !strings.HasPrefix(url, "http") {
		return "http://" + url
	}
	return url
}

func RemoveDomainError(url string) bool {
	parsed, err := neturl.Parse(url)
	if err != nil {
		return true
	}
	if parsed.Hostname() == "" && parsed.Scheme == "" {
		parsed, err = neturl.Parse("//" + url)
		if err != nil {
			return true
		}
	}
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	domain := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(os.Getenv("DOMAIN"))), "www.")
	return host == "" || domain == "" || host != domain
}
