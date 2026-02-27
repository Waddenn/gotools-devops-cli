package procops

import "testing"

func TestParseWindowsLine(t *testing.T) {
	line := `"Code.exe","1234","Console","1","12,000 K"`
	p := parseWindowsLine(line)
	if p == nil {
		t.Fatal("expected parsed process")
	}
	if p.PID != 1234 || p.Name != "Code.exe" {
		t.Fatalf("unexpected parse result: %+v", p)
	}
}

func TestParseUnixLine(t *testing.T) {
	line := "4321 /usr/bin/bash"
	p := parseUnixLine(line)
	if p == nil {
		t.Fatal("expected parsed process")
	}
	if p.PID != 4321 || p.Name != "/usr/bin/bash" {
		t.Fatalf("unexpected parse result: %+v", p)
	}
}

func TestParseProcessesSkipsHeaders(t *testing.T) {
	out := "  PID COMMAND\n123 sshd\n"
	procs := parseProcesses(out, "darwin")
	if len(procs) != 1 {
		t.Fatalf("expected 1 process, got %d", len(procs))
	}
	if procs[0].PID != 123 || procs[0].Name != "sshd" {
		t.Fatalf("unexpected process: %+v", procs[0])
	}
}

func TestIsConfirmed(t *testing.T) {
	cases := map[string]bool{
		"yes": true,
		"Y":   true,
		"oui": true,
		"no":  false,
	}
	for in, want := range cases {
		if got := isConfirmed(in); got != want {
			t.Fatalf("isConfirmed(%q) = %v, want %v", in, got, want)
		}
	}
}
