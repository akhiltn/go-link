package helper

import (
	"os"
	"strings"
)

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func RemoveDomainError(url string) bool {
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(url, "https://", "", 1)
	newURL = strings.Replace(url, "www.", "", 1)
	newURL = strings.Replace(url, "www.", "", 1)
	newURL = strings.Split(url, "/")[0]
	if strings.ToUpper(newURL) == strings.ToUpper(os.Getenv("DOMAIN")) {
		return false
	}
	return true
}
