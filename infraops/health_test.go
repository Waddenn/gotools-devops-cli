package infraops

import (
	"math"
	"testing"
)

func nearlyEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.01
}

func TestParseUnixDFUsedPercent(t *testing.T) {
	in := "Filesystem 1K-blocks Used Available Use% Mounted on\n/dev/sda1 100000 45000 55000 45% /\n"
	used, err := parseUnixDFUsedPercent(in)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !nearlyEqual(used, 45) {
		t.Fatalf("used = %v, want 45", used)
	}
}

func TestParseWMICUsedPercent(t *testing.T) {
	in := "Node,FreeSpace,Size\nPC,500,1000\n"
	used, err := parseWMICUsedPercent(in)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !nearlyEqual(used, 50) {
		t.Fatalf("used = %v, want 50", used)
	}
}

func TestParsePowerShellUsedPercent(t *testing.T) {
	in := "FreeSpace      Size\n---------      ----\n500 1000\n"
	used, err := parsePowerShellUsedPercent(in)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !nearlyEqual(used, 50) {
		t.Fatalf("used = %v, want 50", used)
	}
}
