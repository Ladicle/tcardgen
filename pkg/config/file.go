package config

import (
	"os"

	"github.com/ghodss/yaml"
)

func LoadConfig(filename string) (*DrawingConfig, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := &DrawingConfig{}
	if err := yaml.Unmarshal(f, c); err != nil {
		return nil, err
	}
	return c, nil
}
