# Testing Patterns

We implement a test-first approach to protect configuration loaders, parsers, SQS processing paths, and command execution logic from regressions.

## 1. Decoupled AWS SQS Testing
* To avoid calling real AWS SQS queues during tests, the processor requires an SQS operations interface (`SQSAPI`) instead of a concrete client pointer.
* Use a lightweight local test mock struct (`mockSQSClient`) satisfying this interface to simulate SQS Receive and Delete requests.
* Validate that SQS delete handles are correctly populated for all message processing outcomes.

## 2. Config File Loading and Validation
* Unit tests for configuration loading must use temporary files (generated via `t.TempDir()`).
* Write dummy YAML content with and without parameters (such as SQS credentials), verifying that validation error blocks trigger correctly.

## 3. Command Execution Assertion
* Use lightweight system shell calls (such as `echo`) to assert command execution paths without breaking across testing host architectures.
