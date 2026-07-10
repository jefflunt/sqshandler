package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

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
	LogUTC("INIT Starting SQS listener for queue: %s", p.cfg.SQS.QueueURL)

	for {
		select {
		case <-ctx.Done():
			LogUTC("STOP Polling loop stopped. Waiting for active workers to complete...")
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
				LogUTC("SQSE Error receiving SQS messages: %v", err)
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

	// Step 1: Validate payload using native JSON unmarshalling to map
	var rawPayload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &rawPayload); err != nil {
		LogUTC("JSON unmarshalling: %v", err)
		p.deleteMessage(ctx, msgID, msg.ReceiptHandle)
		return
	}

	// Extract cmd key for routing
	cmdStr := ""
	if cmdVal, exists := rawPayload["cmd"]; exists && cmdVal != nil {
		if s, ok := cmdVal.(string); ok {
			cmdStr = s
		} else {
			cmdStr = fmt.Sprintf("%v", cmdVal)
		}
	}

	// Extract each key configured under Config.Extract without failing on missing/empty values
	extracted := make(map[string]string)
	for _, key := range p.cfg.Extract {
		val, exists := rawPayload[key]
		if !exists || val == nil {
			extracted[key] = ""
			continue
		}
		switch v := val.(type) {
		case string:
			extracted[key] = v
		default:
			extracted[key] = fmt.Sprintf("%v", v)
		}
	}

	// Step 2: Look up command mapping
	cmdConfig, exists := p.cfg.Cmd[cmdStr]
	if !exists {
		LogUTC("NMAP [%s] has no configured mapping", cmdStr)
		p.deleteMessage(ctx, msgID, msg.ReceiptHandle)
		return
	}

	// Step 3: Log command invocation and execute
	LogUTC("SUCC [%s] invoked", cmdStr)

	// Dynamically interpolate extracted keys in path and args
	path := cmdConfig.Path
	for k, v := range extracted {
		path = strings.ReplaceAll(path, "{{"+k+"}}", v)
	}

	args := make([]string, len(cmdConfig.Args))
	for i, arg := range cmdConfig.Args {
		argVal := arg
		for k, v := range extracted {
			argVal = strings.ReplaceAll(argVal, "{{"+k+"}}", v)
		}
		args[i] = argVal
	}
	
	exitStatus := p.runCommand(path, args)
	
	LogUTC("CLSD [%s] closed: exit status=%d", cmdStr, exitStatus)

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
		LogUTC("ERRO Failed to run command: %v", err)
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
		LogUTC("DELF [%s] delete failed: %v", msgID, err)
		return
	}
	LogUTC("DELM [%s] deleted from SQS", msgID)
}
