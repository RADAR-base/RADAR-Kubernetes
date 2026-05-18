package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return &cfg, nil
}

func WriteConfig(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func LoadSecrets(path string) (*Secrets, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading secrets %s: %w", path, err)
	}
	var sec Secrets
	if err := yaml.Unmarshal(data, &sec); err != nil {
		return nil, fmt.Errorf("parsing secrets %s: %w", path, err)
	}
	return &sec, nil
}

func WriteSecrets(sec *Secrets, path string) error {
	data, err := yaml.Marshal(sec)
	if err != nil {
		return fmt.Errorf("marshalling secrets: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}
