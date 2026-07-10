# sqshandler

`sqshandler` is a lightweight daemon written in Go that monitors an Amazon SQS queue, parses incoming messages natively, executes shell commands based on the message command (`cmd`), and deletes processed messages from the queue.

It is designed for security and speed, running message processing asynchronously in parallel goroutines.

---

## Configuration

The configuration file must be located at `~/.sqshandler/config.yml`. Below is a configuration template:

```yaml
sqs:
  region: "us-east-1"
  queue_url: "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue"
  max_number_of_messages: 10   # optional (default: 10)
  wait_time_seconds: 20        # optional (default: 20)
  
  # Optional: Embed static AWS credentials directly
  # If omitted, sqshandler falls back to standard AWS configuration chains (~/.aws/credentials, environment variables, etc.)
  aws_access_key_id: "AKIA..."
  aws_secret_access_key: "wJalr..."

# Commands mapping matching the 'cmd' key in JSON payloads
cmd:
  DRAFT:
    path: "/bin/bash"
    args: ["-c", "echo 'Drafting command executed'"]
  RELEASE:
    path: "/usr/bin/make"
    args: ["release"]
```

---

## Message Format

Incoming SQS messages must have the following JSON structure:

```json
{
  "cmd": "DRAFT",
  "value": "SQSH-2"
}
```

* Messages that fail JSON unmarshalling or have empty required fields (`cmd`, `value`) are immediately logged as errors and deleted from the SQS queue.
* Messages containing a `cmd` that has no configured mapping in `config.yml` are also logged and deleted.

---

## Build and Run

### Prerequisites
* Go 1.20 or newer.

### Compilation
Build the application binary:
```bash
go build -o sqshandler
```

### Execution
Run the daemon:
```bash
./sqshandler
```

### Graceful Shutdown
To stop the daemon cleanly, send a termination signal (`SIGINT` or `SIGTERM`). The program will stop polling and wait for all active command executions to complete before exiting.

---

## Logging Format

All log events are written to standard output. Every log line is prepended with a high-precision UTC ISO8601 timestamp (`YYYY-MM-DDTHH:MM:SS.ffffffZ`) and formatted as:

`<timestamp> <event> <message>`

### Event Types
* `INIT` — Daemon/Polling start events.
* `STOP` — Graceful shutdown and termination events.
* `SUCC` — Command execution starting.
* `CLSD` — Command execution finishing (includes exit status).
* `DELM` — Successful SQS message deletions.
* `DELF` — Failed SQS message deletions.
* `JSON` — Payload parsing and validation errors.
* `NMAP` — Unconfigured command mapping warnings.
* `SQSE` — AWS SDK/SQS errors.
* `ERRO` — Configuration/General errors.
