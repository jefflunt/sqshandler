# Step 1: Initialize Project, Setup Logging, Config Parsing, and Execution

## Goals
1. Setup Go module and install dependencies (`github.com/aws/aws-sdk-go-v2`, `github.com/xeipuuv/gojsonschema`, `gopkg.in/yaml.v3`).
2. Implement ISO8601 UTC microsecond precision custom logger.
3. Parse the YAML configuration file dynamically.
4. Implement validation against JSON schema, routing, process execution in parallel, and SQS deletion/logging.
5. Create tests and verify.

## Dependencies
None.

## Verification
- Run tests: `go test ./...`.
- Verify logging format matching exact ISO8601 UTC timestamp pattern.
