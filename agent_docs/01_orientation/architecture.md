# System Architecture

`sqshandler` is a daemon designed to read commands from an Amazon SQS queue, perform native validation checks on the message payload, execute the corresponding command locally, and clean up SQS messages in all scenarios.

## Components and Flows

```
+-------------------------------------------------------------+
|                     sqs.ReceiveMessage                      |
+------------------------------+------------------------------+
                               |
                               v
                       [Message Poller]
                               |
                   (Spawns parallel goroutine)
                               |
                               v
               [Native JSON Payload Parser]
                               |
                     (Unmarshal & Validate)
                     /                    \
            (Success)                      (Failure)
               /                                \
              v                                  v
      [Command Router]                   [Error Log (JSON)]
              |                                  |
     (Match configured cmd)                      |
     /                    \                      |
 (Match)              (No Match)                 |
   /                        \                    |
  v                          v                   |
[Invoke Command]        [Warn Log (NMAP)]        |
  | (Log SUCC/CLSD)          |                   |
  | (Capture exit status)    |                   |
  \                          /                   |
   \                        /                    |
    v                      v                     v
+-------------------------------------------------------------+
|                     sqs.DeleteMessage                       |
+------------------------------+------------------------------+
                               |
                               v
                [Log deletion result (DELM)]
```

## Concurrency Model
The daemon retrieves messages in batches (up to `max_number_of_messages`) via SQS long-polling. For each message, it launches a Go routine, ensuring high throughput. Graceful shutdown listens for `SIGINT`/`SIGTERM` to safely wait for all active tasks to complete.

## High-Precision Logger
All stdout logs are prepended with an ISO8601 UTC microsecond precision timestamp (`YYYY-MM-DDTHH:MM:SS.ffffffZ`).
Log lines match `<timestamp> <event> <message>`, where event type codes define execution stages (e.g. `SUCC`, `CLSD`, `DELM`, `JSON`, `NMAP`).
