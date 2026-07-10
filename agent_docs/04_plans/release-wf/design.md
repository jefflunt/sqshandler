# Design: GitHub Actions Release Workflow (release-wf)

## User Story
As a maintainer of the `sqshandler` project, I want a GitHub Actions workflow that automatically builds the application binary for all major platforms (Linux, macOS, Windows) and architectures (amd64, arm64) on every branch push, and publishes a release asset when changes are merged to the `main` branch, so that distribution and testing are streamlined.

## Requirements
1. **Triggering:** Build on every push on any branch. Create/update a GitHub Release only on merge to `main`.
2. **Build Targets:** Linux, Darwin (macOS), Windows for both `amd64` and `arm64`.
3. **Build Target Directory:** Root path `.`.
4. **Binary Naming:** Prefix binaries with `sqshandler-`.

## Backlog
- [ ] Create GitHub Actions workflow file `.github/workflows/release.yml`
- [ ] Verify workflow YAML syntax
