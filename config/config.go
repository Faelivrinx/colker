package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of the YAML file
type Config struct {
	Messages   Messages    `yaml:"messages"`
	Secret     Secret      `yaml:"secret"`
	Containers []Container `yaml:"containers"`
	Webhooks   []Hook      `yaml:"webhooks"`
}

type Container struct {
	Name      string `yaml:"name"`
	StatusURL string `yaml:"status_url"`
}

type Messages struct {
	StartMessage string `yaml:"start_message"`
	StopMessage  string `yaml:"stop_message"`
	FinalMessage string `yaml:"final_message"`
}

type Secret struct {
	SecretValue string `yaml:"secret_value"`
}

type Hook struct {
	Url     string `yaml:"url"`
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
