# Coding Conventions

This project follows Go idioms (as enforced by `gofmt` and `go vet`) with the following strict additions:

## 1. Logging Format and High-Precision Timestamps
* Every logged line must prefix the output with a UTC ISO8601 timestamp with microsecond precision:
  `YYYY-MM-DDTHH:MM:SS.ffffffZ` (e.g. `2026-07-10T02:00:00.123456Z`).
* Messages must be formatted as:
  `<timestamp> <event> <message>`
* Event prefix codes include:
  * `INIT`: Startup initialization.
  * `STOP`: Graceful shutdown exit steps.
  * `SUCC`: Successful invocation start of a local command.
  * `CLSD`: Completion of a command, outputting its exit status code.
  * `DELM`: Deletion of a message from the queue.
  * `DELF`: Failure to delete a message from the SQS queue.
  * `JSON`: Parser errors from message unmarshalling.
  * `NMAP`: Missing command configuration warning.
  * `SQSE`: AWS connection or polling errors.
  * `ERRO`: Internal, general runtime errors.

## 2. Dependency Management
* Clean imports and minimal dependencies.
* Go version is locked to `1.26.2`.
* Avoid external validation libraries in favour of Go's built-in `encoding/json` package.

## 3. Concurrency
* Spawning goroutines must be paired with safety mechanisms like `sync.WaitGroup` to allow graceful shutdown without leaking running threads.
