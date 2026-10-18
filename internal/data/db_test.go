package data

import (
	"context"
	"strings"
	"testing"
)

func TestNewRejectsUnsupportedDriver(t *testing.T) {
	_, err := New(context.Background(), "unknown", "")
	if err == nil || !strings.Contains(err.Error(), `unsupported database driver "unknown"`) {
		t.Fatalf("New() error = %v; want unsupported database driver error", err)
	}
}
