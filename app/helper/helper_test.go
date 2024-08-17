package helper

import (
	"os"
	"testing"
)

func TestEnforceHTTP(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"EnforceHTTP", args{"http://localhost:3000"}, "http://localhost:3000"},
		{"EnforceHTTP", args{"https://localhost:3000"}, "https://localhost:3000"},
		{"EnforceHTTP", args{"localhost:3000"}, "http://localhost:3000"},
		{"EnforceHTTP", args{"localhost"}, "http://localhost"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EnforceHTTP(tt.args.url); got != tt.want {
				t.Errorf("EnforceHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDomainError(t *testing.T) {
	type args struct {
		url    string
		domain string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"RemoveDomainError", args{"http://localhost:3000", "localhost"}, true},
		{"RemoveDomainError", args{"https://localhost:3000", "localhost"}, true},
		{"RemoveDomainError", args{"localhost:3000", ""}, true},
		{"RemoveDomainError", args{"localhost", "localhost"}, false},
		{"RemoveDomainError", args{"www.localhost", ""}, true},
		{"RemoveDomainError", args{"localhost.com", "localhost"}, true},
		{"RemoveDomainError", args{"localhost.com:3000", ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DOMAIN", tt.args.domain)
			if got := RemoveDomainError(tt.args.url); got != tt.want {
				t.Errorf("RemoveDomainError() = %v, want %v", got, tt.want)
				t.Errorf("Test Case: %v", tt)
			}
		})
	}
}
