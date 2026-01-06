# Contributing to k8s-switch

First off, thank you for considering contributing to `k8s-switch`! It's people like you who make open-source tools
better for everyone.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Please read it to
understand the expectations for our community.

## How Can I Contribute?

### Reporting Bugs

* Check the [Issues](https://github.com/tuplle/k8s-switch/issues) to see if the bug has already been reported.
* If not, open a new issue. Include a clear title, a description of the problem, and steps to reproduce the issue.

### Suggesting Enhancements

* Open an issue with the "enhancement" tag.
* Describe the feature you'd like to see and why it would be useful.

### Pull Requests

1. Fork the repository.
2. Create a new branch for your feature or fix (`git checkout -b feature/amazing-feature`).
3. Make your changes.
4. Ensure your code follows Go best practices and is formatted (`go fmt ./...`).
5. Commit your changes (`git commit -m 'Add some amazing feature'`).
6. Push to the branch (`git push origin feature/amazing-feature`).
7. Open a Pull Request.

## Development Setup

### Prerequisites

- Go 1.25 or later.
- `make` installed on your system.

### Building Locally

Use the provided `Makefile` to build and test your changes:

```bash
# Build the binary into the bin/ directory
make build

# Install your modified version locally to test
make install
```

## Project Structure

- `main.go`: Entry point of the application.
- `cmd/`: Contains Cobra command definitions and primary logic.
- `internal/`: Shared utility functions and internal logic.

## Style Guide

- Follow standard Go conventions.
- Use meaningful variable and function names.
- Keep functions small and focused on a single task.

Thank you for your contributions!
