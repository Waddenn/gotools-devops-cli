package fileops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateReportIndexMerge(t *testing.T) {
	tmp := t.TempDir()
	dataDir := filepath.Join(tmp, "data")
	outDir := filepath.Join(tmp, "out")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatalf("mkdir data: %v", err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("mkdir out: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dataDir, "a.txt"), []byte("hello world\n"), 0644); err != nil {
		t.Fatalf("write a.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "b.txt"), []byte("42 data\n"), 0644); err != nil {
		t.Fatalf("write b.txt: %v", err)
	}

	if err := GenerateReport(dataDir, outDir); err != nil {
		t.Fatalf("report: %v", err)
	}
	if err := GenerateIndex(dataDir, outDir); err != nil {
		t.Fatalf("index: %v", err)
	}
	if err := MergeFiles(dataDir, outDir); err != nil {
		t.Fatalf("merge: %v", err)
	}

	expected := []string{"report.txt", "index.txt", "merged.txt"}
	for _, name := range expected {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("missing output file %s: %v", name, err)
		}
	}
}
