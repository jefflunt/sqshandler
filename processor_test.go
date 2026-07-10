package main

import (
	"context"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type mockSQSClient struct {
	deletedHandles []string
	mu             sync.Mutex
}

func (m *mockSQSClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	return nil, nil
}

func (m *mockSQSClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deletedHandles = append(m.deletedHandles, aws.ToString(params.ReceiptHandle))
	return &sqs.DeleteMessageOutput{}, nil
}

func TestProcessMessage_Valid(t *testing.T) {
	cfg := &Config{
		Cmd: map[string]CommandConfig{
			"TEST_CMD": {
				Path: "echo",
				Args: []string{"test-output"},
			},
		},
	}
	mockClient := &mockSQSClient{}
	p := NewProcessor(cfg, mockClient)

	msg := types.Message{
		MessageId:     aws.String("msg-1"),
		ReceiptHandle: aws.String("handle-1"),
		Body:          aws.String(`{"cmd":"TEST_CMD","value":"some-val"}`),
	}

	p.processMessage(context.Background(), msg)

	mockClient.mu.Lock()
	deletedCount := len(mockClient.deletedHandles)
	mockClient.mu.Unlock()

	if deletedCount != 1 {
		t.Errorf("expected 1 message to be deleted, got %d", deletedCount)
	}
}

func TestProcessMessage_InvalidJSON(t *testing.T) {
	cfg := &Config{}
	mockClient := &mockSQSClient{}
	p := NewProcessor(cfg, mockClient)

	msg := types.Message{
		MessageId:     aws.String("msg-2"),
		ReceiptHandle: aws.String("handle-2"),
		Body:          aws.String(`{invalid json}`),
	}

	p.processMessage(context.Background(), msg)

	mockClient.mu.Lock()
	deletedCount := len(mockClient.deletedHandles)
	mockClient.mu.Unlock()

	if deletedCount != 1 {
		t.Errorf("expected invalid message to be deleted, got %d", deletedCount)
	}
}

func TestProcessMessage_MissingFields(t *testing.T) {
	cfg := &Config{}
	mockClient := &mockSQSClient{}
	p := NewProcessor(cfg, mockClient)

	// Missing 'value' field
	msg := types.Message{
		MessageId:     aws.String("msg-3"),
		ReceiptHandle: aws.String("handle-3"),
		Body:          aws.String(`{"cmd":"TEST"}`),
	}

	p.processMessage(context.Background(), msg)

	mockClient.mu.Lock()
	deletedCount := len(mockClient.deletedHandles)
	mockClient.mu.Unlock()

	if deletedCount != 1 {
		t.Errorf("expected message with missing fields to be deleted, got %d", deletedCount)
	}
}

func TestProcessMessage_UnconfiguredCmd(t *testing.T) {
	cfg := &Config{
		Cmd: map[string]CommandConfig{},
	}
	mockClient := &mockSQSClient{}
	p := NewProcessor(cfg, mockClient)

	msg := types.Message{
		MessageId:     aws.String("msg-4"),
		ReceiptHandle: aws.String("handle-4"),
		Body:          aws.String(`{"cmd":"UNCONFIGURED","value":"val"}`),
	}

	p.processMessage(context.Background(), msg)

	mockClient.mu.Lock()
	deletedCount := len(mockClient.deletedHandles)
	mockClient.mu.Unlock()

	if deletedCount != 1 {
		t.Errorf("expected unconfigured cmd message to be deleted, got %d", deletedCount)
	}
}
