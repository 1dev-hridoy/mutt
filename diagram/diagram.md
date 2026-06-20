# Mutt Architecture

> Open-source error tracking system — monitor crashes, handle recovery, send real-time alerts.

---

## System Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Mutt System                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────┐       ┌──────────────────┐       ┌────────────────────┐  │
│  │              │       │                  │       │                    │  │
│  │   Go SDK     │──────▶│   Backend API    │──────▶│    PostgreSQL      │  │
│  │  (mutt-go)   │       │   (Fiber)        │       │    (Neon)          │  │
│  │              │       │                  │       │                    │  │
│  └──────────────┘       └────────┬─────────┘       └────────────────────┘  │
│                                  │                                          │
│                                  │                                          │
│                          ┌───────▼────────┐                                │
│                          │                │                                │
│                          │     Redis      │                                │
│                          │   (Upstash)    │                                │
│                          │                │                                │
│                          └────────────────┘                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Data Flow: Error Ingestion

```
┌─────────┐                                    ┌──────────────┐
│         │   1. POST /api/v1/ingest           │              │
│         │   Header: X-Mutt-Key: <key>        │              │
│         │   Body: { title, log, stackTrace } │              │
│  SDK    │───────────────────────────────────▶│   Backend    │
│ (Go)    │                                    │              │
│         │   2. 202 Accepted                  │              │
│         │◀───────────────────────────────────│              │
│         │                                    │              │
└─────────┘                                    └──────┬───────┘
                                                      │
                           ┌──────────────────────────┼──────────────────┐
                           │                          │                  │
                           ▼                          ▼                  ▼
                    ┌─────────────┐           ┌─────────────┐    ┌─────────────┐
                    │  Hash API   │           │  Compute    │    │   Store     │
                    │  Key        │           │  Fingerprint│    │   Error     │
                    └──────┬──────┘           └──────┬──────┘    └──────┬──────┘
                           │                         │                  │
                           ▼                         ▼                  ▼
                    ┌─────────────┐           ┌─────────────┐    ┌─────────────┐
                    │  Lookup     │           │ Find/Create │    │  Increment  │
                    │  Project    │           │ ErrorGroup  │    │  Count      │
                    └─────────────┘           └─────────────┘    └─────────────┘
```

---

## Database Schema

```
┌─────────────────────┐       ┌─────────────────────────┐
│       users         │       │        projects          │
├─────────────────────┤       ├─────────────────────────┤
│ id            (PK)  │◀──────│ user_id           (FK)  │
│ username            │       │ id                (PK)  │
│ email         (UQ)  │       │ name                    │
│ password            │       │ api_key          (UQ)   │
│ phone         (UQ)  │       │ plan                    │
│ plan                │       │ notify                  │
│ created_at          │       │ created_at              │
│ updated_at          │       │ updated_at              │
│ deleted_at          │       │ deleted_at              │
└─────────────────────┘       └────────────┬────────────┘
                                           │
                                           │ 1:N
                                           ▼
                               ┌─────────────────────────┐
                               │     error_groups        │
                               ├─────────────────────────┤
                               │ id                (PK)  │
                               │ project_id        (FK)  │
                               │ fingerprint       (UQ)  │
                               │ title                   │
                               │ status                  │
                               │ count                   │
                               │ last_seen_at            │
                               │ created_at              │
                               │ updated_at              │
                               └────────────┬────────────┘
                                            │
                                            │ 1:N
                                            ▼
                               ┌─────────────────────────┐
                               │        errors           │
                               ├─────────────────────────┤
                               │ id                (PK)  │
                               │ error_group_id    (FK)  │
                               │ project_id        (FK)  │
                               │ log                     │
                               │ stack_trace             │
                               │ severity                │
                               │ notified                │
                               │ occurred_at             │
                               │ created_at              │
                               │ updated_at              │
                               └─────────────────────────┘
```

---

## Authentication Flows

### Dashboard Auth (JWT)

```
┌──────────┐                                    ┌──────────┐
│          │   POST /api/v1/auth/login           │          │
│          │   Body: { email, password }         │          │
│          │────────────────────────────────────▶│          │
│  Browser │                                    │  Backend │
│          │   200 OK                            │          │
│          │   { access_token, refresh_token }   │          │
│          │◀────────────────────────────────────│          │
│          │                                    │          │
│          │   GET /api/v1/projects              │          │
│          │   Header: Authorization: Bearer <token>       │
│          │────────────────────────────────────▶│          │
│          │                                    │          │
│          │   200 OK + data                    │          │
│          │◀────────────────────────────────────│          │
└──────────┘                                    └──────────┘
```

### SDK Ingestion (API Key)

```
┌──────────┐                                    ┌──────────┐
│          │   POST /api/v1/ingest               │          │
│          │   Header: X-Mutt-Key: <api_key>     │          │
│  App     │────────────────────────────────────▶│  Backend │
│  (Go)    │                                    │          │
│          │   202 Accepted                      │          │
│          │◀────────────────────────────────────│          │
└──────────┘                                    └──────────┘
```

---

## API Routes

```
/api/v1/
├── ping                          GET    (health check)
│
├── auth/
│   ├── signup                    POST   (no auth)
│   ├── login                     POST   (no auth)
│   ├── refresh                   POST   (no auth)
│   ├── logout                    POST   (JWT)
│   └── me                        GET    (JWT)
│
├── projects/                     (JWT + Rate Limit)
│   ├── POST                      Create project
│   ├── GET                       List user's projects
│   ├── /:id
│   │   ├── GET                   Get project details
│   │   ├── PATCH                 Update project
│   │   ├── DELETE                Delete project
│   │   ├── rotate-key POST       Rotate API key
│   │   └── errors/
│   │       ├── GET               List error groups
│   │       └── /:errorGroupId
│   │           ├── GET           Get error group + occurrences
│   │           ├── PATCH         Update status
│   │           └── DELETE        Delete error group
│
└── ingest                        POST   (API Key auth)
```

---

## Error Status Lifecycle

```
                    ┌───────────┐
                    │           │
          ┌────────│  critical │◀──────────┐
          │        │           │           │
          │        └─────┬─────┘           │
          │              │                 │
          │              │ new error       │ new error
          │              │ in same group   │ in same group
          │              ▼                 │
          │        ┌───────────┐           │
          │        │           │           │
          │        │ recovered │───────────┘
          │        │           │
          │        └─────┬─────┘
          │              │
          │              │ user confirms fix
          │              ▼
          │        ┌───────────┐
          │        │           │
          └───────▶│ resolved  │
                   │           │
                   └───────────┘

critical  = Error is active, occurring
recovered = Error stopped occurring (auto-detected)
resolved  = User manually marked as fixed
```

---

## Project Structure

```
mutt/
├── cmd/
│   └── main.go                    # Entry point, init, CORS
├── consts/
│   └── server.go                  # Constants (limits, sizes)
├── diagram/
│   └── diagram.md                 # This file
├── internal/
│   ├── config/
│   │   ├── connectToDB.go         # PostgreSQL connection
│   │   ├── loadEnv.go             # .env loading
│   │   ├── redis.go               # Redis connection
│   │   └── syncDatabase.go        # Auto-migration
│   ├── middleware/
│   │   ├── auth.go                # JWT Bearer auth
│   │   ├── apiKeyAuth.go          # API key auth (SDK)
│   │   └── rateLimit.go           # Redis rate limiter
│   └── service/
│       ├── apiKey.go              # API key generate/hash
│       ├── errorGroup.go          # Fingerprinting, grouping
│       ├── hash.go                # bcrypt password hashing
│       ├── notification.go        # Notification flag logic
│       ├── redis.go               # Token store, blacklist
│       └── token.go               # JWT generate/validate
├── models/
│   ├── user.go                    # User, auth DTOs
│   ├── project.go                 # Project, ErrorGroup, Error, DTOs
│   └── error.go                   # (reserved for future)
├── server/
│   ├── handler/
│   │   ├── AuthHandler.go         # Signup, Login, Logout, Refresh, Me
│   │   ├── ErrorHandler.go        # Ingest, error management
│   │   ├── ProjectHandler.go      # Project CRUD
│   │   └── ping.go                # Health check
│   └── routes/
│       └── main.go                # Route definitions
└── sdk/                           # (separate repo: github.com/dishan1223/mutt-go)
```

---

## Key Design Decisions

| Decision | Rationale |
|---|---|
| API keys stored as SHA-256 | Fast lookup (unlike bcrypt), keys are random not user-chosen |
| ErrorGroup fingerprinting | Clusters similar errors like Sentry issues |
| Redis for rate limiting | Distributed, fast, aligns with existing token store |
| Ownership checks on all queries | `WHERE user_id = ? AND id = ?` prevents IDOR |
| API key shown once | Security best practice — like Stripe secret keys |
| Soft delete on projects | Data recovery, audit trail |
