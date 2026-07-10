package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// SQSConfig stores configuration details for Amazon SQS polling.
type SQSConfig struct {
	Region              string `yaml:"region"`
	QueueURL            string `yaml:"queue_url"`
	MaxNumberOfMessages int32  `yaml:"max_number_of_messages"`
	WaitTimeSeconds     int32  `yaml:"wait_time_seconds"`
	AWSAccessKeyID      string `yaml:"aws_access_key_id"`
	AWSSecretAccessKey  string `yaml:"aws_secret_access_key"`
}

// CommandConfig defines the command to execute for a matching command key.
type CommandConfig struct {
	Path string   `yaml:"path"`
	Args []string `yaml:"args"`
}

// Config is the top-level structure of the ~/.sqshandler/config.yml file.
type Config struct {
	SQS     SQSConfig                `yaml:"sqs"`
	Cmd     map[string]CommandConfig `yaml:"cmd"`
	Extract []string                 `yaml:"extract"`
}

// LoadConfig loads the configuration from ~/.sqshandler/config.yml.
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(home, ".sqshandler", "config.yml")
	return LoadConfigFromFile(configPath)
}

// LoadConfigFromFile reads and parses a YAML configuration from a specified file path.
func LoadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	// Default configuration values
	if cfg.SQS.MaxNumberOfMessages == 0 {
		cfg.SQS.MaxNumberOfMessages = 10
	}
	if cfg.SQS.WaitTimeSeconds == 0 {
		cfg.SQS.WaitTimeSeconds = 20
	}

	// Validate required AWS credentials
	if cfg.SQS.AWSAccessKeyID == "" {
		return nil, fmt.Errorf("aws_access_key_id is required in SQS configuration")
	}
	if cfg.SQS.AWSSecretAccessKey == "" {
		return nil, fmt.Errorf("aws_secret_access_key is required in SQS configuration")
	}

	// Validate required extract configuration
	if len(cfg.Extract) == 0 {
		return nil, fmt.Errorf("extract is required in configuration")
	}
	for _, key := range cfg.Extract {
		if key == "" {
			return nil, fmt.Errorf("extract cannot contain empty keys")
		}
	}

	return &cfg, nil
}
