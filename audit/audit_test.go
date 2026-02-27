package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogWritesAuditFile(t *testing.T) {
	outDir := t.TempDir()
	Log(outDir, "LOCK data/input.txt")

	path := filepath.Join(outDir, "audit.log")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read audit.log: %v", err)
	}
	if !strings.Contains(string(data), "LOCK data/input.txt") {
		t.Fatalf("audit content missing action: %q", string(data))
	}
}
