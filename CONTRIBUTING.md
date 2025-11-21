# Contributing to xk6-parquet

Thank you for your interest in contributing to xk6-parquet! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/0/code_of_conduct/). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected vs actual behavior**
- **Environment details** (OS, Go version, k6 version)
- **Sample code** or Parquet files (if applicable)
- **Error messages** and stack traces

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description** of the proposed functionality
- **Explain why this enhancement would be useful**
- **List any alternative solutions** you've considered

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following the coding standards
3. **Add tests** for any new functionality
4. **Update documentation** as needed
5. **Ensure tests pass** (`go test ./...`)
6. **Submit a pull request**

## Development Setup

### Prerequisites

- Go 1.24 or later
- xk6 (`go install go.k6.io/xk6/cmd/xk6@latest`)
- Git

### Setting Up Your Development Environment

```bash
# Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/xk6-parquet.git
cd xk6-parquet

# Add upstream remote
git remote add upstream https://github.com/mmga-lab/xk6-parquet.git

# Create a feature branch
git checkout -b feature/your-feature-name

# Install dependencies
go mod download

# Build with xk6
xk6 build --with github.com/mmga-lab/xk6-parquet=.
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.txt ./...

# Run tests for a specific package
go test -v ./pkg/parquet

# Run a specific test
go test -v -run TestRead ./pkg/parquet
```

### Building

```bash
# Build with xk6
xk6 build --with github.com/mmga-lab/xk6-parquet=.

# Build for specific platform
GOOS=linux GOARCH=amd64 xk6 build --with github.com/mmga-lab/xk6-parquet=.
```

## Coding Standards

### Go Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Use `go vet` to check for common mistakes
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors appropriately

### Code Organization

```
pkg/parquet/
├── module.go      # k6 module registration
├── reader.go      # File reading logic
├── converter.go   # Type conversion
├── cache.go       # Caching layer
└── schema.go      # Schema & metadata
```

### Testing

- Write tests for new functionality
- Maintain or improve code coverage
- Use table-driven tests where appropriate
- Test error cases
- Use meaningful test names

Example test structure:

```go
func TestRead(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    int
        wantErr bool
    }{
        {
            name:    "valid file",
            input:   "testdata/sample.parquet",
            want:    100,
            wantErr: false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Documentation

- Add GoDoc comments for exported functions and types
- Update README.md for user-facing changes
- Add examples for new features
- Update CHANGELOG.md

## Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:

```
feat(reader): add support for nested data types

Add support for reading nested Parquet structures including
arrays, maps, and nested records.

Closes #123

fix(cache): prevent memory leak in cache cleanup

The cache was not properly releasing references to cached data.
This commit ensures proper cleanup when cache entries expire.

docs(api): update readChunked examples

Add more comprehensive examples showing different use cases
for chunked reading with error handling.
```

## Pull Request Process

1. **Update documentation** for any changes to the public API
2. **Add tests** for new functionality
3. **Update CHANGELOG.md** with your changes
4. **Ensure CI passes** (all tests and lints)
5. **Request review** from maintainers
6. **Address review feedback** promptly
7. **Squash commits** if requested

### PR Title Format

Use the same format as commit messages:

```
feat(reader): add column filtering support
fix(cache): resolve race condition in concurrent access
docs(readme): improve installation instructions
```

### PR Description Template

```markdown
## Description
Brief description of the changes

## Motivation and Context
Why is this change required? What problem does it solve?

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to change)
- [ ] Documentation update

## Testing
Describe the tests you ran and how to reproduce them

## Checklist
- [ ] My code follows the code style of this project
- [ ] I have added tests to cover my changes
- [ ] All new and existing tests pass
- [ ] I have updated the documentation accordingly
- [ ] I have updated CHANGELOG.md
```

## Release Process

Releases are managed by maintainers:

1. Update CHANGELOG.md
2. Create and push a version tag (e.g., `v0.2.0`)
3. GitHub Actions automatically builds and creates the release
4. Binaries are attached to the GitHub release

## Questions?

- Open an issue for questions
- Join discussions on GitHub Discussions
- Check existing documentation and examples

## License

By contributing to xk6-parquet, you agree that your contributions will be licensed under the Apache License 2.0.

---

Thank you for contributing to xk6-parquet! 🎉
