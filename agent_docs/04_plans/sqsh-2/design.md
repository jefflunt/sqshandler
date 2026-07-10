# Design: Configurable Payload Data Extraction (SQSH-2)

## User Story
As an operator, I want to make the extraction of data from SQS message payloads configurable via the configuration file, so that I can extract arbitrary keys from incoming JSON messages and interpolate them into command paths and arguments.

## Requirements
1. Allow defining an `extract` section in the config file listing top-level JSON keys to extract.
2. If `extract` is not specified, default to extracting `cmd` and `value` to maintain backward compatibility.
3. Validate that the routing command key (`cmd`) is always provided in the payload and is non-empty.
4. Support dynamic interpolation of all extracted variables using `{{key}}` syntax in `cmd` configuration paths and arguments.
5. Log validation errors, delete the message, and abort execution if any configured key is missing or empty in the payload.

## Architecture

```
+------------------+     +-----------------------+     +-------------------+
|  Incoming SQS    |---->| Unmarshal to Map      +---->| Extract Configured|
|  JSON Payload    |     |                       |     | Keys & Validate   |
+------------------+     +-----------------------+     +---------+---------+
                                                                 |
                                                                 v
+------------------+     +-----------------------+     +---------+---------+
| Execute Command  |<----+ Interpolate variables |<----+ Match routing cmd |
| with args        |     | {{key}} -> val        |     | in Config.Cmd     |
+------------------+     +-----------------------+     +-------------------+
```

## Backlog
- [ ] Implement `extract` config parsing and defaults in `config.go` (`steps/01_implementation.md`)
- [ ] Implement dynamic JSON extraction and interpolation in `processor.go` (`steps/01_implementation.md`)
- [ ] Add unit tests verifying both new and legacy extraction behavior (`steps/01_implementation.md`)
