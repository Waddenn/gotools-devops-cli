package fileops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilterHeadTail(t *testing.T) {
	tmp := t.TempDir()
	in := filepath.Join(tmp, "input.txt")
	out := filepath.Join(tmp, "out")
	if err := os.MkdirAll(out, 0755); err != nil {
		t.Fatalf("mkdir out: %v", err)
	}

	content := "alpha one\nbeta two\nalpha three\ngamma\n"
	if err := os.WriteFile(in, []byte(content), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	if _, err := CountKeyword(in, "alpha"); err != nil {
		t.Fatalf("count keyword: %v", err)
	}
	if err := FilterKeyword(in, "alpha", out); err != nil {
		t.Fatalf("filter keyword: %v", err)
	}
	if err := Head(in, 2, out); err != nil {
		t.Fatalf("head: %v", err)
	}
	if err := Tail(in, 2, out); err != nil {
		t.Fatalf("tail: %v", err)
	}

	expected := []string{"filtered.txt", "filtered_not.txt", "head.txt", "tail.txt"}
	for _, name := range expected {
		if _, err := os.Stat(filepath.Join(out, name)); err != nil {
			t.Fatalf("missing output file %s: %v", name, err)
		}
	}
}
