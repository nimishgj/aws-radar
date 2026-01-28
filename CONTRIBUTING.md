# Contributing to AWS Radar

Thank you for your interest in contributing to AWS Radar! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How to Contribute

### Reporting Issues

- Check if the issue already exists in the [issue tracker](https://github.com/nimishgj/aws-radar/issues)
- If not, create a new issue with a clear title and description
- Include steps to reproduce, expected behavior, and actual behavior
- Add relevant logs, screenshots, or error messages

### Suggesting Features

- Open an issue with the `enhancement` label
- Describe the feature and its use case
- Explain why it would be beneficial to the project

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Make your changes
4. Run tests and linting:
   ```bash
   make ci
   ```
5. Commit your changes with a clear message:
   ```bash
   git commit -m "Add feature: description of changes"
   ```
6. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
7. Open a Pull Request against `main`

## Development Setup

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Make

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/nimishgj/aws-radar.git
   cd aws-radar
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

3. Run the application:
   ```bash
   make run
   ```

4. Run tests:
   ```bash
   make test
   ```

### Running with Docker

```bash
make docker-up
```

## Coding Standards

### Go Code

- Follow standard Go conventions and idioms
- Use `gofmt` for formatting (run `make fmt`)
- Pass `go vet` checks (run `make lint`)
- Write clear, descriptive variable and function names
- Add comments for exported functions and types
- Handle errors appropriately

### Commit Messages

- Use clear, descriptive commit messages
- Start with a verb in present tense (Add, Fix, Update, Remove)
- Keep the first line under 72 characters
- Reference issues when applicable: `Fix #123: description`

### Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Update documentation if needed
- Add tests for new functionality
- Ensure all CI checks pass
- Request review from maintainers

## Adding New AWS Collectors

To add a new AWS service collector:

1. Create a new file in `internal/collector/`:
   ```go
   // internal/collector/newservice.go
   package collector

   type NewServiceCollector struct{}

   func NewNewServiceCollector() *NewServiceCollector {
       return &NewServiceCollector{}
   }

   func (c *NewServiceCollector) Name() string {
       return "newservice"
   }

   func (c *NewServiceCollector) Collect(ctx context.Context, cfg aws.Config, region string) error {
       // Implementation
   }
   ```

2. Register the collector in `cmd/aws-radar/main.go`

3. Add the metric definition in `internal/metrics/metrics.go`

4. Update documentation in README.md

## Testing

- Write unit tests for new functionality
- Use table-driven tests where appropriate
- Mock AWS API calls in tests
- Run the full test suite before submitting PRs:
  ```bash
  make test
  ```

## Documentation

- Update README.md for user-facing changes
- Add inline comments for complex logic
- Update config examples if configuration changes

## Getting Help

- Open an issue for questions
- Check existing issues and documentation first

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
