package main

import (
	"strings"
	"testing"
)

func TestColorizeNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "xterm-256color")

	got := colorize("hello", clrBlue)
	if got != "hello" {
		t.Fatalf("colorize() = %q, want plain text", got)
	}
}

func TestPromptAndStatusMessages(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	if got := prompt("Choix"); got != "[Choix]" {
		t.Fatalf("prompt() = %q", got)
	}

	ok := success("done")
	if !strings.Contains(ok, "done") {
		t.Fatalf("success() missing payload: %q", ok)
	}

	err := failure("oops")
	if !strings.Contains(err, "oops") {
		t.Fatalf("failure() missing payload: %q", err)
	}
}
