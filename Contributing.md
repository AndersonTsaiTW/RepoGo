# 🛠️ Contributing to RepoGo

Welcome to RepoGo! We're thrilled you're interested in contributing. This guide outlines how to set up your environment, follow our coding standards, and submit issues or pull requests. Whether you're fixing bugs, adding features, or improving documentation—your contributions are valued.

---

## 🚀 Getting Started

### 1. Clone the Repository
```bash
git clone https://github.com/AndersonTsaiTW/RepoGo.git
cd RepoGo
```

### 2. Install Dependencies
Ensure you have Go installed (version ≥ 1.20 recommended).
```bash
go mod download
```

### 3. Build the Project
```bash
make build
```

### 4. Run Locally
```bash
go run main.go
```

### 5. Run Tests
```bash
go test ./...
```

---

## 🎨 Code Style & Conventions

To maintain consistency and readability:

- **Language**: Go (Golang)
- **Formatting**: Use `gofmt` or `goimports` before committing
- **Naming**: Use descriptive names for variables, functions, and packages
- **Comments**: Write clear, concise comments for complex logic
- **Error Handling**: Always check and handle errors gracefully
- **Modularity**: Keep functions small and focused
- **Tests**: Include unit tests for new features and bug fixes

---

## 🐛 Submitting Issues

If you encounter a bug or have a feature request:

1. **Search existing issues** to avoid duplicates.
2. **Create a new issue** with a clear title and description.
3. Include:
   - Steps to reproduce (if applicable)
   - Expected vs actual behavior
   - Screenshots or logs (if helpful)
   - Environment details (OS, Go version)

---

## 📦 Submitting Pull Requests

We welcome PRs for bug fixes, features, documentation, and refactoring.

### Steps:

1. **Fork the repo** and create your branch:
   ```bash
   git checkout -b feature/my-awesome-feature
   ```

2. **Make your changes** and ensure tests pass.

3. **Run linters and formatters**:
   ```bash
   gofmt -w .
   ```

4. **Commit with a clear message**:
   ```bash
   git commit -m "Add feature: smart token estimation"
   ```

5. **Push and open a PR**:
   ```bash
   git push origin feature/my-awesome-feature
   ```

6. In your PR description, include:
   - What the change does
   - Why it’s needed
   - Any related issues (e.g., `Fixes #42`)

---

## 🔍 Review Process

Once you submit a PR:

- A maintainer will review your code.
- You may be asked to make changes or clarify intent.
- Once approved, your PR will be merged.

We aim to review PRs within **3–5 business days**.

---

## 📜 Code of Conduct

All contributors are expected to follow our [Code of Conduct](CODE_OF_CONDUCT.md). Be respectful, inclusive, and constructive in all interactions.

---

## 🙌 Thank You

Your contributions make RepoGo better for everyone. Whether you're fixing typos or building new features, we appreciate your effort and creativity.

Let’s build something great together!
