package generator

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestGenerate(t *testing.T) {
	var buf bytes.Buffer
	count := 5

	err := Generate(&buf, count)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) != count {
		t.Errorf("Expected %d lines, got %d", count, len(lines))
	}

	for i, line := range lines {
		var record Record
		err := json.Unmarshal(line, &record)
		if err != nil {
			t.Errorf("Line %d: invalid JSON: %v", i, err)
			continue
		}

		// Validate fields
		if record.Version != "1.0.0" {
			t.Errorf("Line %d: expected Version=1.0.0, got %s", i, record.Version)
		}
		if record.Name == "" {
			t.Errorf("Line %d: Name is empty", i)
		}
		if record.UUID == "" || len(record.UUID) < 10 {
			t.Errorf("Line %d: UUID is invalid: %s", i, record.UUID)
		}
		if record.Timestamp.IsZero() {
			t.Errorf("Line %d: Timestamp is zero", i)
		}
	}
}
