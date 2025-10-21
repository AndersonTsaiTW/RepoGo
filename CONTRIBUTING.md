# Contributing to RepoGo

Thank you for your interest in contributing to RepoGo! 🎉

This document provides guidelines and instructions for contributing to this project. Whether you're fixing a bug, adding a feature, or improving documentation, your contributions are welcome!

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Code Style Guidelines](#code-style-guidelines)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Getting Help](#getting-help)

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. Please be respectful and considerate in your interactions with others.

**Key principles:**
- Be welcoming and inclusive
- Be respectful of differing viewpoints
- Accept constructive criticism gracefully
- Focus on what is best for the community

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.22 or higher**: [Download Go](https://golang.org/dl/)
- **Git**: [Download Git](https://git-scm.com/downloads)
- **Make** (optional but recommended): For using the Makefile
  - Windows: `choco install make` or `scoop install make`
  - macOS: Already included with Xcode Command Line Tools
  - Linux: `apt-get install build-essential` or equivalent

### Fork and Clone

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/RepoGo.git
   cd RepoGo
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/AndersonTsaiTW/RepoGo.git
   ```

## Development Setup

### Build the Project

```bash
# Using make
make build

# Or using go directly
go build -o bin/repogo ./cmd/repogo
```

### Run the Application

```bash
# Scan current directory
./bin/repogo

# Show help
./bin/repogo --help

# Scan with options
./bin/repogo . -o output.md --include "*.go,*.md"
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/scanner/
```

**Testing Policy:**
- Tests are **encouraged but not required** for contributions
- For **bug fixes**, please include a test that reproduces the issue
- We provide example tests in `internal/analyzer/analyzer_test.go` to help you get started
- Go's testing is simple - just create a `*_test.go` file next to your code!

## How to Contribute

### 1. Find or Create an Issue

- Check [existing issues](https://github.com/AndersonTsaiTW/RepoGo/issues) for something to work on
- Look for issues labeled `good first issue` or `help wanted`
- If you have a new idea, create an issue first to discuss it

### 2. Create a Branch

Create a feature branch from `main`:

```bash
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name
```

**Branch naming conventions:**
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Adding or updating tests

### 3. Make Your Changes

- Write clean, readable code
- Follow the code style guidelines (see below)
- Update documentation as needed

### 4. Commit Your Changes

```bash
git add .
git commit -m "type: brief description"
```

See [Commit Message Guidelines](#commit-message-guidelines) for details.

### 5. Push and Create a Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Code Style Guidelines

### Go Code Style

We follow standard Go conventions:

- **Use `gofmt`**: Format your code with `gofmt` before committing
  ```bash
  gofmt -w .
  ```

- **Run `go vet`**: Check for common mistakes
  ```bash
  go vet ./...
  ```

- **Follow Go best practices**:
  - Use meaningful variable and function names
  - Keep functions small and focused
  - Add comments for exported functions and types
  - Use package-level documentation

### Package Organization

Our project structure:

```
cmd/repogo/          # Application entry point
internal/            # Internal packages (not importable by other projects)
├── models/         # Data structures
├── config/         # Configuration
├── scanner/        # File scanning
├── git/            # Git integration
├── analyzer/       # File analysis
└── renderer/       # Output rendering
```

**Guidelines:**
- Each package should have a single, well-defined purpose
- Avoid circular dependencies
- Use `internal/` for code that shouldn't be imported externally
- Add package-level comments explaining the package purpose

### Documentation

- **Public functions and types** must have comments
  ```go
  // GetInfo retrieves Git repository information from the specified root directory.
  // Returns nil if the directory is not a Git repository.
  func GetInfo(root string) (*models.GitInfo, error) {
      // ...
  }
  ```

- **Package documentation**: Add a `doc.go` or package comment
  ```go
  // Package scanner provides file system scanning and filtering functionality.
  package scanner
  ```

## Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semicolons, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks (updating dependencies, etc.)

### Examples

```
feat(scanner): add support for .gitignore patterns

Implement gitignore pattern matching to automatically exclude
files that are in .gitignore.

Closes #42
```

```
fix(analyzer): correct token estimation for multi-byte characters

The previous implementation incorrectly counted bytes instead of
characters, leading to inaccurate token estimates for non-ASCII text.

Fixes #38
```

```
docs: update README with new CLI options
```

### Rules

- Use imperative mood: "add feature" not "added feature"
- Start lowercase, no period at the end
- Keep it short and clear
- Reference issues: `Fixes #123` or `Closes #456`

## Pull Request Process

### Before Submitting

- [ ] Code follows the style guidelines
- [ ] All tests pass (`go test ./...`)
- [ ] Tests added (encouraged, especially for bug fixes)
- [ ] Documentation is updated
- [ ] Commit messages follow the guidelines
- [ ] Branch is up to date with `main`

### PR Template

When you create a PR, please fill out the template completely:

- **Description**: What does this PR do?
- **Type of change**: Bug fix, new feature, etc.
- **Testing**: How was this tested?
- **Checklist**: Complete all items

### Review Process

1. **Automated checks**: CI will run tests and linting
2. **Code review**: Maintainers will review your code
3. **Feedback**: Address any requested changes
4. **Approval**: Once approved, your PR will be merged

### After Your PR is Merged

- Delete your feature branch
- Pull the latest `main` branch
- Celebrate! 🎉

## Issue Guidelines

### Creating Issues

When creating an issue, please:

- **Use templates**: Select the appropriate issue template
- **Search first**: Check if a similar issue already exists
- **Be specific**: Provide detailed information
- **Include examples**: Code snippets, error messages, screenshots

### Issue Types

- **Bug Report**: Something isn't working
- **Feature Request**: Suggest a new feature
- **Documentation**: Improvements or additions to documentation
- **Question**: Ask for help or clarification

### Good First Issues

New to the project? Look for issues labeled `good first issue`. These are:
- Well-defined and scoped
- Good for learning the codebase
- Have clear acceptance criteria
- Usually don't require deep knowledge of the project

## Getting Help

### Communication Channels

- **GitHub Issues**: For bugs, features, and questions
- **GitHub Discussions**: For general discussions (if enabled)
- **Pull Request Comments**: For specific code questions

### Resources

- [Project README](README.md)
- [Go Documentation](https://golang.org/doc/)
- [GitHub Flow](https://guides.github.com/introduction/flow/)

### Questions?

If you have questions:
1. Check existing issues and documentation
2. Create a new issue with the `question` label
3. Be patient and respectful

## Recognition

Contributors will be:
- Listed in the project's contributors page
- Acknowledged in release notes (for significant contributions)
- Appreciated and thanked! 🙏

## License

By contributing to RepoGo, you agree that your contributions will be licensed under the same license as the project (see [LICENSE](LICENSE)).

---

Thank you for contributing to RepoGo! Your efforts help make this project better for everyone. 🚀

**Happy Coding!** 💻
