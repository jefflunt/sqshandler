# Domain Reference

This document maps out the core data structures and interfaces defining the domain behavior of `sqshandler`.

## Core Configuration Structures

### 1. `SQSConfig`
Governs parameters for the SQS queue and connection credentials:
* `region` (string): AWS region (e.g. `us-east-1`).
* `queue_url` (string): Standard HTTP URL pointing to the SQS queue.
* `max_number_of_messages` (int32): Number of messages pulled per poll.
* `wait_time_seconds` (int32): Long polling duration (max 20 seconds).
* `aws_access_key_id` (string): Mandatory AWS credential.
* `aws_secret_access_key` (string): Mandatory AWS credential.

### 2. `CommandConfig`
Defines the executable script or command definition:
* `path` (string): Executable system file path (e.g. `/bin/bash`).
* `args` (slice of strings): Command arguments. Supports `{{value}}` placeholder interpolation.

---

## Payload Structure

### `MessagePayload`
Parsed from the raw body of an SQS message:
* `cmd` (string): Name of the command to execute (lookup key in config map).
* `value` (string): Context string to interpolate in arguments.

Both fields are mandatory.

---

## Client Mocking Interface

### `SQSAPI`
To run unit tests without live SQS connections, operations are decoupled into the following mockable interface:
```go
type SQSAPI interface {
    ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
    DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}
```
