<p align="center">
  <img src="assets/logo.svg" alt="Mutt" width="120" />
</p>

<h1 align="center">Mutt</h1>

<p align="center">
  Open-source error tracking with built-in OpenTelemetry observability.
</p>

---

## Overview

Mutt is a self-hosted error tracking system for monitoring crashes, handling recovery, and sending real-time alerts — with full observability baked in via OpenTelemetry.

## Features

- **Error Grouping** — Clusters similar errors using SHA-256 fingerprint hashing
- **Real-time Ingestion** — Capture errors from your apps via SDK
- **Status Tracking** — Mark errors as critical, recovered, or resolved
- **Per-Project Alerts** — Toggle notifications on/off per project
- **API Key Auth** — Secure SDK ingestion with hashed API keys
- **Backup & Restore** — Export and import project data (JSON / gzip)
- **Rate Limiting** — Redis-backed rate limits on all endpoints

## OpenTelemetry

Mutt ships with OTel tracing out of the box:

- **HTTP request tracing** via Fiber OTel middleware
- **Database query tracing** via `otelgorm` (auto-instruments all GORM calls)
- **OTLP export** — sends traces to any OTel-compatible collector (Jaeger, Grafana Tempo, etc.)

Configure with standard OTel env vars:

```
OTEL_SERVICE_NAME=mutt
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
```

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go |
| Framework | Fiber |
| Database | PostgreSQL (NeonDB) + GORM |
| Cache | Redis (Upstash) |
| Auth | JWT (access/refresh tokens) |
| Observability | OpenTelemetry (OTLP/HTTP) |

## Status

> **In Development** — Core backend built. SDKs for Go, JS, and other languages coming soon.

## License

[AGPL-3.0](LICENSE)

## Contributing

Contributions welcome — open issues or submit PRs.

- [Contributing Guide](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Issue Template](ISSUE_TEMPLATE.md)
- [PR Template](PR_TEMPLATE.md)
- [Architecture Diagram](diagram/diagram.md)
