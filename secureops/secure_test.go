package secureops

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLockUnlockLifecycle(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "sample.txt")
	if err := os.WriteFile(file, []byte("hello"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	r := bufio.NewReader(strings.NewReader("yes\nyes\n"))
	if err := LockFile(file, tmp, r); err != nil {
		t.Fatalf("lock: %v", err)
	}
	if !IsLocked(file, tmp) {
		t.Fatal("expected locked file")
	}

	if err := UnlockFile(file, tmp, r); err != nil {
		t.Fatalf("unlock: %v", err)
	}
	if IsLocked(file, tmp) {
		t.Fatal("expected unlocked file")
	}
}

func TestSetReadOnlyReadWrite(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "rw.txt")
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if err := SetReadOnly(file, tmp); err != nil {
		t.Fatalf("set readonly: %v", err)
	}
	if err := SetReadWrite(file, tmp); err != nil {
		t.Fatalf("set readwrite: %v", err)
	}
}
