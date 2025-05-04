package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Route struct {
	PathPrefix string   `yaml:"path_prefix"`
	Backends   []string `yaml:"backends"`
	Methods    []string `yaml:"methods"`
}
type Config struct {
	Routes   []Route `yaml:"routes"`
	Port     string  `yaml:"port"`
	Env      string  `yaml:"env"`
	LogLevel string  `yaml:"log_level"`
	LogFile  string  `yaml:"log_file"`
}

func LoadConfig() ([]byte, error) {
	return os.ReadFile("configs/routes.yaml")
}

func ParseConfig(yamlFile []byte) (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
