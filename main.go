package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	LogUTC("INIT Initializing sqshandler...")

	// Load configuration
	cfg, err := LoadConfig()
	if err != nil {
		LogUTC("ERRO Failed to load configuration: %v", err)
		os.Exit(1)
	}

	// Initialize SQS client
	var optFns []func(*awsConfig.LoadOptions) error
	optFns = append(optFns, awsConfig.WithRegion(cfg.SQS.Region))

	if cfg.SQS.AWSAccessKeyID != "" && cfg.SQS.AWSSecretAccessKey != "" {
		optFns = append(optFns, awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.SQS.AWSAccessKeyID, cfg.SQS.AWSSecretAccessKey, ""),
		))
	}

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(), optFns...)
	if err != nil {
		LogUTC("ERRO Failed to load AWS SDK config: %v", err)
		os.Exit(1)
	}

	sqsClient := sqs.NewFromConfig(awsCfg)

	// Create processor
	processor := NewProcessor(cfg, sqsClient)

	// Set up context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		LogUTC("STOP Received shutdown signal: %v. Initiating graceful shutdown...", sig)
		cancel()
	}()

	// Start processor
	processor.Start(ctx)
	LogUTC("STOP sqshandler daemon exited.")
}
