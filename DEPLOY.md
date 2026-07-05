# Deploying Obligo

Obligo is published as a single Docker image with the frontend
already built in. You don't need to clone the source code or build anything
yourself to run it — you only need the image and a small set of config files.

## Prerequisites

* Docker and Docker Compose installed on your host machine.

## What you need

You need exactly three things in a directory on your deploy machine:

1. A `docker-compose.yml` that defines the app container (pulling the
   published image) and a Postgres container for it to talk to.
2. A `.env` file with your own secrets and settings.
3. Nothing else — no source code, no `node_modules`, no Go toolchain.

> **Before you use this file:** the compose example below assumes
> `APP_PORT` is wired all the way through as a variable, on both the
> `environment` block *and* the left/right sides of `ports`. If your local
> `created-compose.yml` still has `APP_PORT` and the right-hand side of
> `ports` hardcoded to `"1234"`, update that file to match this pattern
> first, or the example below won't reflect what your compose file actually
> does.

### 1. Create your `docker-compose.yml`

```yaml
services:
  db:
    image: postgres:17-alpine
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASS}
      POSTGRES_DB: ${DEV_NAME}

  app:
    image: thegreatestgiant/obligo
    restart: always
    environment:
      # Use the service name 'db', not 'localhost' — Docker Compose
      # resolves service names to the right container on its internal network.
      DB_URL: "postgres://${DB_USER}:${DB_PASS}@db:5432/${DEV_NAME}?sslmode=disable"
      JWT_SECRET: ${JWT_SECRET}
      APP_PORT: ${APP_PORT}
      APP_URL: ${APP_URL}
    ports:
      - "8080:${APP_PORT}"
    depends_on:
      - db
```

A couple of things worth understanding about that `ports` line, since it's
easy to get backwards:

* The **left** side (`8080`) is the port *on your host machine* — this is
  what you put in your browser's address bar (via `APP_URL`).
* The **right** side (`${APP_PORT}`) is the port the app actually listens on
  *inside the container*. This has to match whatever value you set
  `APP_PORT` to in your `.env`, because the app reads that same variable at
  startup to decide what port to bind to. If these two ever get out of sync
  — say, `APP_PORT` gets changed in `.env` but the compose file still has a
  literal number hardcoded on the right side — the app will be listening on
  one port while Docker forwards traffic to a different one, and the
  deployment will silently fail to respond.
* The host-side `8080` doesn't have to match `APP_PORT` at all — that's just
  which port *you* pick to reach it from outside. You could just as easily
  make that side a variable too if you want full control over both ends, but
  it isn't required the way the right-hand side is.

This is intentionally a minimal, standalone compose file for deployment. It
is not the same `docker-compose.yml` used inside the project repo for local
development (that one spins up separate dev/test databases and an Adminer
instance for people actively writing code against the project — see
[DEVELOPERS.md](./DEVELOPERS.md) if that's what you're trying to do
instead).

### 2. Create your `.env`

In the same directory:

```env
DB_USER=charity_user
DB_PASS=super_secure_db_password
DEV_NAME=charity_db

APP_PORT=1234
APP_URL=http://localhost:8080

JWT_SECRET=your_super_long_random_string_here
```

Notes on each value:

* `DB_USER` / `DB_PASS` / `DEV_NAME` — credentials and database name for the
  Postgres container. These aren't pre-set to anything; pick your own.
* `APP_PORT` — the port the app listens on *inside* its container. Most
  deployments can leave this at `1234` and never think about it again; you'd
  only change it if something else inside the same container/network were
  already using port 1234.
* `APP_URL` — the externally-reachable URL where you'll access the app. The
  app injects this into the frontend at runtime so the browser knows where
  to send API requests, so it needs to match wherever you'll actually be
  browsing to (including the host-side port from your `docker-compose.yml`,
  if non-standard).
* `JWT_SECRET` — must be a long, random string used to sign login sessions.
  Generate one with:

  ```bash
  openssl rand -base64 64
  ```

## Launching it

From the directory containing your `docker-compose.yml` and `.env`:

```bash
docker compose up -d
```

Docker will pull the app image and the Postgres image, start both
containers, and the app will automatically run its database migrations on
first startup — there's no separate migration step to run yourself.

## Accessing the app

Navigate to whatever you set `APP_URL` to (default in the example above:
`http://localhost:8080`).

## Stopping it

```bash
docker compose down
```

Add `-v` if you also want to delete the Postgres data volume (this will
permanently delete all accounts and ledger data):

```bash
docker compose down -v
```
