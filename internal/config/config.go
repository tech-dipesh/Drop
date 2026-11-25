package config

import (
	"cmp"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Dir  string `json:"dir"`
}

func Load() (Config, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "drop", "config.json"))
	var c Config
	if err == nil && json.Unmarshal(data, &c) == nil && c.ID != "" {
		return c, nil
	}
	id := make([]byte, 8)
	if _, err := rand.Read(id); err != nil {
		return Config{}, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	host, _ := os.Hostname()
	c = Config{ID: hex.EncodeToString(id), Name: cmp.Or(host, "drop device"), Dir: filepath.Join(home, "Downloads")}
	return c, Save(c)
}

func Save(c Config) error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	dir = filepath.Join(dir, "drop")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, _ := json.MarshalIndent(c, "", "  ")
	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0o644)
}
