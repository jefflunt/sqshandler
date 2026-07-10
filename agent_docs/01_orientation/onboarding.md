# Onboarding Guide

## Prerequisites
* Go 1.26 or newer.
* git.

## Environment Setup
1. Create the `~/.sqshandler/` configuration directory:
   ```bash
   mkdir -p ~/.sqshandler
   ```
2. Create and write `config.yml` inside that directory, configuring your SQS URL and AWS static credentials:
   ```yaml
   sqs:
     region: "us-east-1"
     queue_url: "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue"
     aws_access_key_id: "AKIA..."
     aws_secret_access_key: "wJalr..."
   cmd:
     DRAFT:
       path: "/bin/bash"
       args: ["-c", "echo 'Handling ticket {{value}}'"]
   ```

## Development Commands

### Compilation
To compile the local binary:
```bash
./script/build
```
This builds and places the executable binary inside `bin/sqshandler` with the correct commit-based version injected.

### Running Unit Tests
Execute the Go test suite:
```bash
go test -v ./...
```
