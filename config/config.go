package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DefaultFile string `json:"default_file"`
	BaseDir     string `json:"base_dir"`
	OutDir      string `json:"out_dir"`
	DefaultExt  string `json:"default_ext"`
	WikiLang    string `json:"wiki_lang"`
	ProcessTopN int    `json:"process_top_n"`
}

func DefaultConfig() *Config {
	return &Config{
		DefaultFile: "data/input.txt",
		BaseDir:     "data",
		OutDir:      "out",
		DefaultExt:  ".txt",
		WikiLang:    "fr",
		ProcessTopN: 10,
	}
}

// Load detecte le format (json ou txt) et charge la config
func Load(path string) (*Config, error) {
	if strings.HasSuffix(path, ".json") {
		return loadJSON(path)
	}
	return loadTXT(path)
}

func loadJSON(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire %s: %w", path, err)
	}
	// on part des valeurs par defaut, comme ca si une clé manque dans le json c'est pas grave
	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("JSON invalide dans %s: %w", path, err)
	}
	return cfg, nil
}

func loadTXT(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire %s: %w", path, err)
	}
	defer f.Close()

	cfg := DefaultConfig()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "default_file":
			cfg.DefaultFile = val
		case "base_dir":
			cfg.BaseDir = val
		case "out_dir":
			cfg.OutDir = val
		case "default_ext":
			cfg.DefaultExt = val
		case "wiki_lang":
			cfg.WikiLang = val
		case "process_top_n":
			if n, err := strconv.Atoi(val); err == nil && n >= 0 {
				cfg.ProcessTopN = n
			}
		}
	}
	return cfg, scanner.Err()
}

func (c *Config) EnsureOutDir() error {
	return os.MkdirAll(c.OutDir, 0755)
}
