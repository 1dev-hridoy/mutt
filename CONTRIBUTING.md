# Contributing

Contributions are welcome from anyone. Feel free to open issues or submit pull requests.

## Getting Started

1. Fork the repo
2. Clone your fork
3. Create a branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Push and open a PR

## Development Setup

- Go 1.21+
- PostgreSQL (or use NeonDB)
- Redis (or use Upstash)

Copy `.env.example` to `.env` and fill in your credentials.

```bash
go run cmd/main.go
```

## Code Style

- Follow existing patterns in the codebase
- Keep functions short and focused
- No unnecessary abstractions

## Pull Requests

- Keep PRs small and focused on one change
- Describe what changed and why
- Reference related issues if any

## Architecture

To understand the system, data flows, and design decisions, see the [Architecture Diagram](diagram/diagram.md).

## License

By contributing, you agree that your contributions will be licensed under the [AGPL-3.0](LICENSE).
