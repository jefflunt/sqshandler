package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary YAML config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yml")

	yamlContent := `
sqs:
  region: "us-west-2"
  queue_url: "https://sqs.us-west-2.amazonaws.com/12345/my-queue"
  max_number_of_messages: 5
  wait_time_seconds: 15
  aws_access_key_id: "test-key"
  aws_secret_access_key: "test-secret"

cmd:
  BUILD:
    path: "/usr/bin/make"
    args: ["build"]
  TEST:
    path: "/usr/bin/go"
    args: ["test", "./..."]
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	cfg, err := LoadConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadConfigFromFile failed: %v", err)
	}

	if cfg.SQS.Region != "us-west-2" {
		t.Errorf("expected SQS.Region to be 'us-west-2', got '%s'", cfg.SQS.Region)
	}
	if cfg.SQS.QueueURL != "https://sqs.us-west-2.amazonaws.com/12345/my-queue" {
		t.Errorf("unexpected QueueURL: %s", cfg.SQS.QueueURL)
	}
	if cfg.SQS.MaxNumberOfMessages != 5 {
		t.Errorf("expected MaxNumberOfMessages to be 5, got %d", cfg.SQS.MaxNumberOfMessages)
	}
	if cfg.SQS.WaitTimeSeconds != 15 {
		t.Errorf("expected WaitTimeSeconds to be 15, got %d", cfg.SQS.WaitTimeSeconds)
	}

	buildCmd, ok := cfg.Cmd["BUILD"]
	if !ok {
		t.Fatal("expected 'BUILD' command mapping to exist")
	}
	if buildCmd.Path != "/usr/bin/make" || len(buildCmd.Args) != 1 || buildCmd.Args[0] != "build" {
		t.Errorf("unexpected BUILD command configuration: %+v", buildCmd)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yml")

	yamlContent := `
sqs:
  region: "us-east-1"
  queue_url: "https://sqs.us-east-1.amazonaws.com/12345/my-queue"
  aws_access_key_id: "test-key"
  aws_secret_access_key: "test-secret"
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	cfg, err := LoadConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadConfigFromFile failed: %v", err)
	}

	// Verify defaults
	if cfg.SQS.MaxNumberOfMessages != 10 {
		t.Errorf("expected default MaxNumberOfMessages 10, got %d", cfg.SQS.MaxNumberOfMessages)
	}
	if cfg.SQS.WaitTimeSeconds != 20 {
		t.Errorf("expected default WaitTimeSeconds 20, got %d", cfg.SQS.WaitTimeSeconds)
	}
}

func TestLoadConfigMissingCredentials(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yml")

	yamlContent := `
sqs:
  region: "us-east-1"
  queue_url: "https://sqs.us-east-1.amazonaws.com/12345/my-queue"
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	_, err := LoadConfigFromFile(configPath)
	if err == nil {
		t.Error("expected LoadConfigFromFile to fail when credentials are missing")
	}
}
