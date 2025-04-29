package processor

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcess(t *testing.T) {
	// Skip if integration tests are not enabled
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set ENABLE_INTEGRATION_TESTS=true to run")
	}

	// Setup test output directory
	testOutputDir := "test_output_" + t.Name()
	os.Setenv("OUTPUT_DIR", testOutputDir)
	defer os.RemoveAll(testOutputDir)

	// Sample input
	input := `
{"version":"1.0.0","name":"Test Boat 1","uuid":"urn:mrn:signalk:uuid:12345678-1234-1234-1234-123456789012","latitude":37.7,"longitude":-122.3,"altitude":0.0,"course":123.4,"speed":5.0,"timestamp":"2025-04-28T12:34:56Z"}
{"version":"1.0.0","name":"Test Boat 2","uuid":"urn:mrn:signalk:uuid:abcdef12-abcd-abcd-abcd-abcdef123456","latitude":37.8,"longitude":-122.4,"altitude":0.0,"course":234.5,"speed":6.0,"timestamp":"2025-04-28T12:34:56Z"}
`
	reader := strings.NewReader(input)
	var output bytes.Buffer

	// Process the input
	err := Process(reader, &output)
	if err != nil {
		t.Fatalf("Process() failed: %v", err)
	}

	// Check if output files were created
	files, err := filepath.Glob(filepath.Join(testOutputDir, "*/*/*/*/*"))
	if err != nil {
		t.Fatalf("Failed to glob output files: %v", err)
	}

	if len(files) == 0 {
		t.Errorf("No output files were created")
	}

	// Output should indicate success
	outputStr := output.String()
	if !strings.Contains(outputStr, "Created partition file") {
		t.Errorf("Expected output to contain 'Created partition file', got: %s", outputStr)
	}
}