package util

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintHelp(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Reset stdout when test completes
	defer func() {
		os.Stdout = oldStdout
	}()

	// Call the function
	PrintHelp()

	// Restore stdout and get the output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify the output contains expected elements
	expectedStrings := []string{
		"Command Line Arguments:",
		"-e: [dev, prod]",
		"-p: (default: 3000)",
		"--cert-path:",
		"--key-path:",
		"--disable-public-fs:",
		"--db-adapter: [imdb, mongo]",
		"--db-host:",
		"--db-name:",
		"--db-user:",
		"--db-pass:",
		"--no-db:",
		"--help:",
		"Environment Variables:",
		"GWC_ENV:",
		"GWC_PORT:",
		"GWC_CERT_PATH:",
		"GWC_KEY_PATH:",
		"GWC_ENABLE_PUBLIC_FS:",
		"GWC_DB_ADAPTER:",
		"GWC_DB_HOSTNAME:",
		"GWC_DB_NAME:",
		"GWC_DB_USERNAME:",
		"GWC_DB_PASSWORD:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q, but it didn't", expected)
		}
	}

	// Check overall output length
	if len(output) < 500 {
		t.Errorf("Output seems too short: %d characters", len(output))
	}
}
