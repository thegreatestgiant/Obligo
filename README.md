# Charity-Tracker

A highly portable, self-contained financial ledger for tracking charitable donations. Built with a Go backend and a React frontend, compiled into a single hardened Docker container.

## Architecture & Technical Decisions

Charity-Tracker prioritizes ease of deployment and security over unnecessary complexity.

* **Single Binary Deployment:** The React frontend (`dist` folder) is built and injected directly into the Go server during the Docker build stage.
* **Dynamic Frontend Injection:** To ensure portability without rebuilding the React app, `window.API_URL` is dynamically injected into `index.html` at runtime by the Go server.
* **Embedded Migrations (Symlink Strategy):** Database initialization happens automatically on startup. The SQL schemas reside in `backend/schema/`, and a symlink at `backend/internal/DB/schema` allows Go's strict `//go:embed` to package the SQL seamlessly without cluttering the internal structure.
* **Graceful Shutdowns:** The HTTP server listens for Docker/OS signals (`SIGTERM`/`SIGINT`) and utilizes `context.WithTimeout` to allow active connections 5 seconds to finish before safely killing the process.
* **Production Hardened Docker:** The final Alpine Linux container runs strictly as an unprivileged, non-root user (`appuser`) for maximum security.

## Deployment

For end-user deployment instructions, please read [DEPLOY.md](./DEPLOY.md).

## Local Development Setup

**1. Run a local Postgres database:**

```bash
docker run --name charity-db -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=charity_db -p 5432:5432 -d postgres:15-alpine
```

**2. Start the Frontend:**

```bash
cd frontend
npm install
npm run dev
```

**3. Start the Backend:**
Ensure your `.env` is configured to point to `127.0.0.1:5432`.

```bash
cd backend
go run main.go
```

## Roadmap

* [x] Graceful shutdown and structured logging (`log/slog`)
* [x] Docker non-root hardening
* [x] Self-contained database migrations
* [ ] Implement CSV transaction Export/Import functionality
