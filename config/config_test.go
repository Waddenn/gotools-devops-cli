package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadJSONWithDefaults(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "config.json")
	if err := os.WriteFile(p, []byte(`{"base_dir":"mydata"}`), 0644); err != nil {
		t.Fatalf("write json: %v", err)
	}

	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("load json: %v", err)
	}
	if cfg.BaseDir != "mydata" {
		t.Fatalf("base_dir = %q, want mydata", cfg.BaseDir)
	}
	if cfg.DefaultFile == "" || cfg.OutDir == "" {
		t.Fatal("expected defaults for missing keys")
	}
}

func TestLoadTXT(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "config.txt")
	body := "# comment\ndefault_file=data/demo.txt\nprocess_top_n=25\n"
	if err := os.WriteFile(p, []byte(body), 0644); err != nil {
		t.Fatalf("write txt: %v", err)
	}

	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("load txt: %v", err)
	}
	if cfg.DefaultFile != "data/demo.txt" {
		t.Fatalf("default_file = %q", cfg.DefaultFile)
	}
	if cfg.ProcessTopN != 25 {
		t.Fatalf("process_top_n = %d", cfg.ProcessTopN)
	}
}
