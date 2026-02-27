package webops

import "testing"

func TestSafeFilePart(t *testing.T) {
	cases := map[string]string{
		"Go_(langage)":   "Go__langage",
		"  docker stats": "docker_stats",
		"###":            "article",
	}
	for in, want := range cases {
		if got := safeFilePart(in); got != want {
			t.Fatalf("safeFilePart(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestExtractWordsFromText(t *testing.T) {
	words := extractWordsFromText("Go 123 Docker 42, Kubernetes")
	if len(words) != 3 {
		t.Fatalf("len(words) = %d, want 3", len(words))
	}
	if words[0] != "Go" || words[1] != "Docker" || words[2] != "Kubernetes" {
		t.Fatalf("unexpected words: %#v", words)
	}
}
