# Contributing to Better Auth SDK for Go

Thank you for your interest in contributing to Better Auth SDK for Go! We welcome contributions from the community.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow. Please be respectful and constructive in all interactions.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/better-auth-sdk-go.git
   cd better-auth-sdk-go
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/Zytera/better-auth-sdk-go.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Install Dependencies

```bash
make install
# or
go mod download
```

### Run Tests

```bash
make test
# or
go test -v -race ./...
```

### Run Linters

```bash
make lint
# or
golangci-lint run ./...
```

## How to Contribute

### Reporting Issues

Before creating an issue, please:

1. Check if the issue already exists
2. Provide a clear and descriptive title
3. Include steps to reproduce the issue
4. Provide your environment details (Go version, OS, etc.)
5. Include relevant code samples or error messages

### Suggesting Features

We love new ideas! When suggesting a feature:

1. Check if it has already been suggested
2. Provide a clear use case
3. Explain how it would benefit users
4. Consider backwards compatibility

### Contributing Code

1. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes** following our coding standards

3. **Write tests** for your changes

4. **Run tests** to ensure everything works:
   ```bash
   make test
   ```

5. **Format your code**:
   ```bash
   make fmt
   ```

6. **Run linters**:
   ```bash
   make vet
   ```

7. **Commit your changes** with a clear message:
   ```bash
   git commit -m "feat: add new feature"
   ```

8. **Push to your fork**:
   ```bash
   git push origin feature/my-new-feature
   ```

9. **Create a Pull Request** on GitHub

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Naming Conventions

- Use descriptive variable names
- Use camelCase for variables and functions
- Use PascalCase for exported functions and types
- Use ALL_CAPS for constants

### Code Organization

- Keep functions small and focused
- Group related functionality together
- Use meaningful package names
- Avoid circular dependencies

### Comments

- Write clear and concise comments
- Document all exported functions and types
- Use godoc format for documentation
- Explain "why" not "what" in comments

Example:
```go
// SignUp registers a new user with the provided credentials.
// It returns the created user and session, or an error if registration fails.
func (s *AuthService) SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
    // Implementation
}
```

### Error Handling

- Always check and handle errors
- Use custom error types when appropriate
- Provide context with error messages
- Don't panic in library code

Example:
```go
if err != nil {
    return nil, WrapError(ErrorTypeInternal, "failed to create user", err)
}
```

## Testing Guidelines

### Test Coverage

- Aim for at least 80% test coverage
- Write tests for all public APIs
- Include both positive and negative test cases
- Test error conditions

### Test Structure

Use table-driven tests when appropriate:

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name  string
        email string
        want  bool
    }{
        {"valid email", "test@example.com", true},
        {"invalid email", "invalid", false},
        {"empty email", "", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := ValidateEmail(tt.email)
            if got != tt.want {
                t.Errorf("ValidateEmail() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific tests
go test -v -run TestFunctionName ./...

# Run benchmarks
make bench
```

## Pull Request Process

1. **Update documentation** if you're changing functionality
2. **Add tests** for new features or bug fixes
3. **Ensure all tests pass** before submitting
4. **Keep PRs focused** - one feature/fix per PR
5. **Write a clear PR description** explaining:
   - What changes you made
   - Why you made them
   - How to test them
6. **Link related issues** using keywords like "Fixes #123"
7. **Be responsive** to review feedback

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(auth): add OAuth support for Twitter
fix(session): handle expired refresh tokens correctly
docs(readme): update installation instructions
```

## Reporting Bugs

When reporting bugs, please include:

1. **Clear title** describing the issue
2. **Steps to reproduce** the problem
3. **Expected behavior** vs actual behavior
4. **Code samples** demonstrating the issue
5. **Environment details**:
   - Go version
   - Operating system
   - SDK version
6. **Error messages** or stack traces

## Suggesting Features

When suggesting features, please include:

1. **Clear description** of the feature
2. **Use case** explaining why it's needed
3. **Proposed API** or usage example
4. **Alternatives considered**
5. **Additional context** or mockups

## Questions?

If you have questions about contributing, feel free to:

- Open an issue with the `question` label
- Reach out to the maintainers
- Check existing documentation and issues

## License

By contributing to Better Auth SDK for Go, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing! 🎉
