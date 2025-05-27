# Contributing to MCP Server Library

Thank you for your interest in contributing to the MCP Server Library! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in the Issues section
2. Use the bug report template when creating a new issue
3. Include detailed steps to reproduce the bug
4. Include expected and actual behavior
5. Include relevant logs and error messages

### Suggesting Features

1. Check if the feature has already been suggested
2. Use the feature request template
3. Provide a clear description of the feature
4. Explain why this feature would be useful
5. Include any relevant examples

### Pull Requests

1. Fork the repository
2. Create a new branch for your changes
3. Make your changes
4. Add tests for your changes
5. Update documentation
6. Submit a pull request

## Development Setup

1. Install Go 1.21 or later
2. Fork and clone the repository
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run tests:
   ```bash
   go test ./...
   ```

## Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Run `golint` to check for style issues
- Write clear and concise commit messages

## Testing

- Write unit tests for new features
- Ensure all tests pass before submitting a PR
- Maintain or improve test coverage
- Include integration tests for complex features

## Documentation

- Update README.md if needed
- Add comments to exported functions and types
- Update examples if you change functionality
- Document any new configuration options

## Review Process

1. All PRs require at least one review
2. CI checks must pass
3. Code must be properly formatted
4. Tests must pass
5. Documentation must be updated

## Release Process

1. Version numbers follow [Semantic Versioning](https://semver.org/)
2. Major version changes require discussion
3. Release notes must be updated
4. Security updates are released as patch versions

## Security

- Report security vulnerabilities to security@serubanpetershan.com
- Do not include sensitive information in issues or PRs
- Follow security best practices in your code

## License

By contributing, you agree that your contributions will be licensed under the project's MIT License. 