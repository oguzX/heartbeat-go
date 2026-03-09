# ZPulse

A lightweight heartbeat monitoring API built with Go. Services register themselves and periodically send heartbeat pings; ZPulse tracks their health status and automatically opens incidents when a service goes silent.

## Features

- Register and manage monitored services
- Ingest heartbeat pings with optional metadata
- Track service health status (`unknown`, `healthy`, `down`)
- Automatic incident detection via a background evaluator (runs every 15 seconds)
- Auto-resolve incidents when a heartbeat is received
- Client IP detection via `X-Forwarded-For`, `X-Real-IP`, or `RemoteAddr`
- Structured JSON logging via `slog`
- Graceful shutdown with configurable timeouts
- PostgreSQL backend via `pgx`

## Tech Stack

- **Language:** Go
- **Router:** [chi](https://github.com/go-chi/chi)
- **Database:** PostgreSQL ([pgx](https://github.com/jackc/pgx))
- **Config:** `.env` via [godotenv](https://github.com/joho/godotenv)

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL

### Configuration

Copy `.env.example` to `.env` and set the values:

```env
APP_ENV=development
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=devpulse
DB_USER=devpulse
DB_PASSWORD=devpulse
DB_SSLMODE=disable
DB_MAX_CONNS=10
```

### Run

```bash
go run ./cmd/api
```

## API

### Health

| Method | Path      | Description              |
|--------|-----------|--------------------------|
| GET    | `/`       | Liveness check           |
| GET    | `/health` | Server health            |
| GET    | `/ready`  | Readiness (DB ping)      |

### Services

| Method | Path               | Description        |
|--------|--------------------|--------------------|
| POST   | `/api/v1/services` | Register a service |
| GET    | `/api/v1/services` | List all services  |

#### Create Service — Request Body

```json
{
  "name": "my-worker",
  "expected_interval_seconds": 60,
  "grace_seconds": 30
}
```

#### Service — Response

```json
{
  "id": 1,
  "name": "my-worker",
  "slug": "my-worker",
  "api_key": "...",
  "expected_interval_seconds": 60,
  "grace_seconds": 30,
  "status": "unknown",
  "created_at": "...",
  "updated_at": "..."
}
```

### Heartbeats

| Method | Path                 | Description           |
|--------|----------------------|-----------------------|
| POST   | `/api/v1/heartbeats` | Ingest a heartbeat ping |

#### Request Body

```json
{
  "service_key": "<api_key>",
  "meta": { "version": "1.0.0" }
}
```

#### Response

```json
{
  "service": { ... },
  "heartbeat": {
    "id": 42,
    "service_id": 1,
    "received_at": "...",
    "source_ip": "1.2.3.4",
    "meta_json": { "version": "1.0.0" }
  }
}
```

## Project Structure

```
.
├── cmd/api/              # Entry point
├── internal/
│   ├── config/           # Env-based configuration
│   ├── db/               # PostgreSQL connection pool
│   ├── domain/           # Core types (Service, Heartbeat)
│   ├── repository/       # Database access layer
│   ├── service/          # Business logic
│   └── http/
│       ├── handlers/     # HTTP handlers
│       └── routes/       # Router setup
```
