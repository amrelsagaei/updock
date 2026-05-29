# Contributing to upDock

Thanks for your interest in contributing! Here's how to get started.

## Development Setup

```bash
# Clone the repo
git clone https://github.com/amrelsagaei/updock.git
cd updock

# Install Go 1.23+
# https://go.dev/dl/

# Build
go build ./cmd/updock/

# Run tests
go test -race ./...

# Lint (install golangci-lint: https://golangci-lint.run/welcome/install/)
golangci-lint run ./...

# Format
golangci-lint fmt ./...
```

## Running a Single Test

```bash
go test -race -run TestRankOfficialFirst ./internal/hub/
```

## Integration Tests

Some tests drive the real `docker compose` CLI and need a running Docker daemon.
They are tagged `integration` and skipped by default. Run them with:

```bash
go test -race -tags integration ./internal/docker/
```

Each test creates a throwaway project, asserts behavior, and tears it down. If
Docker is not running, they self-skip.

## Adding a Recipe

Recipes live in `internal/recipe/recipes/`. To add one:

1. Create `internal/recipe/recipes/yourapp.yaml`.
2. Follow the schema - see `wordpress.yaml` as an example.
3. Every `${VAR}` in services must have a matching prompt entry.
4. Run `go test ./internal/recipe/` - it validates all embedded recipes.
5. Open a PR.

## Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - new feature
- `fix:` - bug fix
- `docs:` - documentation only
- `refactor:` - no behavior change
- `test:` - adding or updating tests
- `chore:` - maintenance, deps, CI
- `ci:` - CI/CD changes

Example: `feat: add redis recipe`

## Pull Request Process

1. Fork the repo and create a branch from `main`.
2. Make your changes.
3. Ensure tests pass: `go test -race ./...`
4. Ensure lint passes: `golangci-lint run ./...`
5. Open a PR with a clear description.
6. All PRs require one maintainer review before merge.

## Code Style

- All `.go` files must have the copyright header.
- No `//nolint` without a justification comment.
- Import ordering: stdlib, external, internal (enforced by linter).
- Table-driven tests for all validation and parsing logic.

## Security

If you find a security vulnerability, do NOT open a public issue. See [SECURITY.md](SECURITY.md) for how to report it privately.
