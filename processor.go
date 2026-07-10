package main

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// MessagePayload represents the expected JSON structure of SQS messages.
type MessagePayload struct {
	Cmd   string `json:"cmd"`
	Value string `json:"value"`
}

// SQSAPI defines the interface for SQS client operations to allow unit testing with mocks.
type SQSAPI interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

// Processor manages polling and dispatching messages.
type Processor struct {
	cfg       *Config
	sqsClient SQSAPI
	wg        sync.WaitGroup
}

// NewProcessor creates a new message processor.
func NewProcessor(cfg *Config, sqsClient SQSAPI) *Processor {
	return &Processor{
		cfg:       cfg,
		sqsClient: sqsClient,
	}
}

// Start begins polling the SQS queue until the context is canceled.
func (p *Processor) Start(ctx context.Context) {
	LogUTC("Starting SQS listener for queue: %s", p.cfg.SQS.QueueURL)

	for {
		select {
		case <-ctx.Done():
			LogUTC("Polling loop stopped. Waiting for active workers to complete...")
			p.wg.Wait()
			return
		default:
			input := &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(p.cfg.SQS.QueueURL),
				MaxNumberOfMessages: p.cfg.SQS.MaxNumberOfMessages,
				WaitTimeSeconds:     p.cfg.SQS.WaitTimeSeconds,
			}

			output, err := p.sqsClient.ReceiveMessage(ctx, input)
			if err != nil {
				// Don't log context cancellation as an error
				if errors.Is(err, context.Canceled) {
					continue
				}
				LogUTC("Error receiving SQS messages: %v", err)
				continue
			}

			for _, msg := range output.Messages {
				p.wg.Add(1)
				go func(m types.Message) {
					defer p.wg.Done()
					p.processMessage(ctx, m)
				}(msg)
			}
		}
	}
}

// Wait blocks until all active processing routines complete.
func (p *Processor) Wait() {
	p.wg.Wait()
}

func (p *Processor) processMessage(ctx context.Context, msg types.Message) {
	msgID := aws.ToString(msg.MessageId)
	body := aws.ToString(msg.Body)

	// Step 1: Validate payload using native JSON unmarshalling
	var payload MessagePayload
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		LogUTC("Message [%s] failed native JSON unmarshalling: %v. Body: %s", msgID, err, body)
		p.deleteMessage(ctx, msgID, msg.ReceiptHandle)
		return
	}

	// Validate required fields
	if payload.Cmd == "" || payload.Value == "" {
		LogUTC("Message [%s] validation failed: 'cmd' and 'value' must be non-empty. Body: %s", msgID, body)
		p.deleteMessage(ctx, msgID, msg.ReceiptHandle)
		return
	}

	// Step 2: Look up command mapping
	cmdConfig, exists := p.cfg.Cmd[payload.Cmd]
	if !exists {
		LogUTC("Message [%s] cmd '%s' has no configured mapping. Message discarded.", msgID, payload.Cmd)
		p.deleteMessage(ctx, msgID, msg.ReceiptHandle)
		return
	}

	// Step 3: Log command invocation and execute
	LogUTC("Command [%s] invoked: path=%s args=%v", payload.Cmd, cmdConfig.Path, cmdConfig.Args)
	
	exitStatus := p.runCommand(cmdConfig.Path, cmdConfig.Args)
	
	LogUTC("Command [%s] closed: exit status=%d", payload.Cmd, exitStatus)

	// Step 4: Delete message from SQS
	p.deleteMessage(ctx, msgID, msg.ReceiptHandle)
}

func (p *Processor) runCommand(path string, args []string) int {
	cmd := exec.Command(path, args...)
	// DO NOT log or redirect stdout/stderr to standard log output
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		// If command couldn't be started or failed for other reasons (e.g. file not found)
		LogUTC("Failed to run command: %v", err)
		return -1
	}

	return 0
}

func (p *Processor) deleteMessage(ctx context.Context, msgID string, receiptHandle *string) {
	_, err := p.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(p.cfg.SQS.QueueURL),
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		LogUTC("Failed to delete message [%s] from SQS: %v", msgID, err)
		return
	}
	LogUTC("Successfully deleted message [%s] from SQS.", msgID)
}
