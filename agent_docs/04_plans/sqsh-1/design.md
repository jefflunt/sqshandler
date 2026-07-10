# Design: Monitor SQS Queue and Invoke Command (SQSH-1)

## User Story
As a developer/operator, I want a lightweight daemon written in Go that monitors a configured Amazon SQS queue, validates incoming messages by unmarshalling them natively, and executes specific shell commands based on the message command (`cmd`) value, so that tasks can be run asynchronously and securely.

## Requirements
1. **Language:** Go (1.20+).
2. **AWS SQS:** Long polling, retrieve messages, delete after execution/failure/discard.
3. **JSON Validation:** Native Go `encoding/json` package. If a message fails to unmarshal, it is deleted from the SQS queue and the event is logged.
4. **Command Execution:** Run configured command on match; no JSON interpolation. Log invocation start and exit code on completion.
5. **Config Location:** `~/.sqshandler/config.yml`.
6. **Logging:** Standardized logging where every line is prefixed with ISO8601 UTC microsecond timestamp.

## Architecture

```
                       +-------------------+
                       |    Config File    |
                       | (~/.sqshandler/   |
                       |    config.yml)    |
                       +---------+---------+
                                 | loads
                                 v
+------------------+     +-------+-------+     +-------------------+
|                  |     |               |     |                   |
|  Amazon SQS      +---->+  sqshandler   +---->+  Shell Command    |
|  (Message Queue) |     |  Daemon       |     |  Execution        |
|                  |     |               |     |                   |
+------------------+     +-------+-------+     +-------------------+
                                 | logs to stdout
                                 v
                         +-------+-------+
                         | UTC ISO8601   |
                         | Logger        |
                         +---------------+
```

## Backlog
- [ ] Initialize repository structure and go.mod (`steps/01_init.md`)
- [ ] Create logger package with ISO8601 UTC timestamp layout (`steps/02_logger.md`)
- [ ] Load and parse configuration file (`steps/03_config.md`)
- [ ] Implement validator and processor module (`steps/04_processor.md`)
- [ ] Write integration test & mock suite (`steps/05_testing.md`)
