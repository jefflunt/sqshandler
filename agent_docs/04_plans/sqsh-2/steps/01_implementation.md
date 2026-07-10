# Step 1: Configurable Payload Data Extraction

## Goals
1. Add `extract` to the config struct in `config.go` and implement defaults (i.e. default to `["cmd", "value"]` if empty). Validate that no empty strings are present in `extract`.
2. Replace static payload parsing in `processor.go` with dynamic JSON unmarshalling to `map[string]interface{}`.
3. Validate that the routing command key `cmd` is always present and non-empty.
4. Extract all configured keys under `extract`, convert them to string representations, and validate they are non-empty.
5. Support dynamic interpolation of all extracted keys using the `{{key}}` syntax.
6. Add unit tests for configuration parsing and processing logic to cover both new and legacy scenarios.

## Dependencies
None.

## Verification
- Run tests: `go test -v ./...`
- Verify that both configuration parsing tests and message processor tests pass cleanly.
