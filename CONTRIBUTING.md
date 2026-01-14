# Contributing to KIE Go SDK

Thank you for your interest in contributing!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR-USERNAME/kie-go.git`
3. Create a branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Run tests: `go test -race ./...`
6. Submit a pull request

## Development Setup

```bash
# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run with race detection
go test -race ./...

# Check formatting
gofmt -d .

# Run linter
go vet ./...
```

## Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Run `gofmt` before committing
- Add godoc comments for exported types
- Keep line length under 100 characters

## Pull Request Guidelines

1. Write clear commit messages
2. Add tests for new features
3. Update documentation
4. Ensure all tests pass
5. Keep PRs focused and small

## Reporting Issues

- Use GitHub Issues
- Include Go version and OS
- Provide minimal reproduction code
- Include error messages

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
