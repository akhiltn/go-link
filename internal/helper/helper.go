package helper

import (
	neturl "net/url"
	"strings"
)

func EnforceHTTP(url string) string {
	lowercaseURL := strings.ToLower(url)
	if !strings.HasPrefix(lowercaseURL, "http://") && !strings.HasPrefix(lowercaseURL, "https://") {
		return "http://" + url
	}
	return url
}

func RemoveDomainError(rawURL, domain string) bool {
	parsed, err := neturl.Parse(rawURL)
	if err != nil {
		return true
	}
	scheme := strings.ToLower(parsed.Scheme)
	if parsed.Hostname() == "" {
		parsed, err = neturl.Parse("//" + rawURL)
		if err != nil {
			return true
		}
		scheme = ""
	}
	if scheme != "" && scheme != "http" && scheme != "https" {
		return true
	}
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	domainURL, err := neturl.Parse("//" + strings.TrimSpace(domain))
	if err != nil {
		return true
	}
	normalizedDomain := strings.TrimPrefix(strings.ToLower(domainURL.Hostname()), "www.")
	return host == "" || normalizedDomain == "" || host != normalizedDomain
}
