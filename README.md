# Go JWT Session API

> Repository: `session-smuggler` · https://github.com/d28035203/session-smuggler

Go REST API with JWT authentication, GORM, and PostgreSQL. Includes Docker Compose and Kubernetes manifests.

## Features

- Register / login with bcrypt password hashing
- JWT sessions stored in PostgreSQL
- Health check endpoint for probes
- Multi-stage Docker build
- Postgres image with optional schema bootstrap
- Kubernetes Deployment + Service manifests

## Tech Stack

| Layer | Choice |
|-------|--------|
| Language | Go |
| HTTP | Fiber |
| ORM | GORM |
| Database | PostgreSQL |
| Auth | JWT (`golang-jwt`) + bcrypt |
| Ops | Docker Compose, Kubernetes |

## Project Structure

```
session-smuggler/
├── app/                 # Bootstrap, middleware, env validation
├── database/            # Connection helper, init SQL, DB Dockerfile
├── handlers/            # Auth + health handlers
├── kubernetes/          # App and Postgres manifests
├── models/              # User, sessions, response envelope
├── router/              # /api/v1 route registration
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── main.go
```

## API (`/api/v1`)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness probe |
| POST | `/register` | Create user `{ username, password }` |
| POST | `/login` | Returns JWT |
| POST | `/logout` | Invalidate session (Bearer token) |
| GET | `/authentication` | Validate current session |

## Quick Start

```bash
git clone https://github.com/d28035203/session-smuggler.git
cd session-smuggler
cp .env.example .env
# edit TOKEN_SECRET and Postgres settings

go mod tidy
go run .
```

### Docker Compose

```bash
docker compose up --build
```

API: `http://localhost:8085`

### Kubernetes

```bash
kubectl apply -f kubernetes/
```

## Environment

```env
APP_HOST=0.0.0.0
APP_PORT=8085
TOKEN_SECRET=change-me-to-a-long-random-string
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DBNAME=goproject
```

## License

MIT
