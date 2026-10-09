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
	scheme := strings.ToLower(parsed.Scheme)
	if parsed.Hostname() == "" {
		parsed, err = neturl.Parse("//" + url)
		if err != nil {
			return true
		}
		scheme = ""
	}
	if scheme != "" && scheme != "http" && scheme != "https" {
		return true
	}
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	domainURL, err := neturl.Parse("//" + strings.TrimSpace(os.Getenv("DOMAIN")))
	if err != nil {
		return true
	}
	domain := strings.TrimPrefix(strings.ToLower(domainURL.Hostname()), "www.")
	return host == "" || domain == "" || host != domain
}
