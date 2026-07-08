# Developer Notes

This document is for anyone writing code in this repository—not for end users
deploying the app. If you're looking for deployment instructions, see
[DEPLOY.md](./DEPLOY.md). If you're looking for the high-level architecture
pitch, see [README.md](./README.md).

> **Note on Terminology:** While the public-facing application is named **Obligo** and uses user-friendly terms like "giving targets" and "obligations," the underlying database schema and variable names retain their legacy nomenclature (e.g., `charity_owed`, `donation_percentage`).

The goal of this file is to write down the things that aren't obvious from
reading the code in isolation—invariants, gotchas, and "don't do this"
warnings that a future contributor (including future you) won't necessarily
rediscover just by reading a single function.

## The charity math is a snapshot, not a live calculation

This is the single most important thing to understand before touching the
`Ledgers` table, the ledger handlers, or anything in `queries.go` that sums
`charity_owed`.

### How it actually works

When a row is written to the `Ledgers` table, the `charity_owed` field is
**not** computed on the fly when you read the row, and it is **not** kept
in sync with the user's current settings dynamically. It is calculated exactly
once at the moment of insertion (`setEntry` in `ledger_post.go`), using
whatever `donation_percentage` the user happens to have set *at that instant*.

Concretely:

- When a **paycheck** entry is inserted, `charity_owed` is calculated as
  `amount * (current_donation_percentage / 100)`, and written into that row.

If the user changes their donation percentage later (via
`PATCH /users/settings`), nothing about previously-inserted rows changes.
Only new rows inserted after the change will use the new percentage. This is
intentional and is covered by our end-to-end tests.

### Why this becomes a problem: Editing historical entries

We have an update endpoint (`PATCH /entries/{id}`) that allows users to fix
typos or incorrect amounts on existing entries. While this preserves the original
`transaction_date`, editing the `amount` of a paycheck has a side effect that is
easy to miss:

**The charity math is recalculated using today's settings, not the
settings that were in effect when the original entry was created.**

If the user's `donation_percentage` has changed at any point between the
original insert and the edit, modifying the amount will cause `ledger_patch.go`
to recalculate the row's `charity_owed` against the *new* percentage—with no
warning, no audit trail, and no indication anywhere that the historical math shifted.

There's a second-order effect worth knowing about too: because `getAmountOwed`
and `getAmountFulfilled` in `queries.go` sum across *all* of a user's rows,
editing one old paycheck can shift the aggregate "total owed" figure the dashboard
shows—which in turn changes the *meaning* of the fulfilled percentages calculated
against that total.

### Why this matters for you, specifically

None of this shows up as a bug in code review or in casual testing, because
every individual operation does exactly what it's supposed to do. The problem
is emergent: it only appears when you combine "snapshot math" and the `PATCH`
recalculation across a timeline of changing user settings.

If you are adding a feature that touches ledger entries in bulk—like an admin
script to backfill data or modifying the CSV Import feature—you need to explicitly
preserve the original `charity_owed` values rather than letting a fresh insert
or patch recalculate them, or you risk silently rewriting a user's financial history.

## Other things worth knowing

### Two separate Postgres instances, on purpose

`docker-compose.yml` defines `db-dev` and `db-test` as separate Postgres
containers on separate ports (`DEV_PORT` / `TEST_PORT` in `.env`). This is
deliberate: the `db-test` instance exists **exclusively** for the Go backend
test suite (`go test`).

Tests call `clearDatabase()` (`TRUNCATE ... RESTART IDENTITY CASCADE`)
between runs, and you absolutely do not want that running against your dev
data. Always check which `DB_*_URL` a given Make target or test file is
pointed at before assuming it's safe to wipe.

### The schema lives in two places, and that's intentional

`backend/schema/000_official_look.sql` is the canonical schema file.
`backend/internal/DB/schema` is a **symlink** to that same directory, which
exists purely so Go's `//go:embed` directive (in `migration.go`) can pull the
SQL into the compiled binary without duplicating the file or restructuring
`backend/schema` to live inside the Go module tree. If you ever need to add a
new migration file, add it to `backend/schema/`, not anywhere under
`internal/DB/`. The symlink will pick it up automatically.

### CORS only fully trusts one origin

`CorsMiddleware` (in `corsMiddleware.go`) checks for an `Origin` header
matching exactly `http://localhost:5173` (the default Vite dev server port)
and only sends credentialed CORS headers for that exact match. Any other
origin gets a wildcard `*` header instead, which will not work for
authenticated, cookie-based requests. If you change the frontend dev port,
or deploy somewhere with a different origin, you'll need to update this
check—there isn't an env var for it currently.

### The `10` fallback in `getDonationPercent` is intentional

If you read `getDonationPercent` in `queries.go`, you'll see it defaults to
`10.0` if the DB lookup fails for any reason (no rows, query error). This
matches the schema default (`donation_percentage INT DEFAULT 10`) and is a
deliberate fail-safe, not a placeholder that was meant to be replaced. Don't
"fix" it by removing it.

## Roadmap / known gaps

These are known, deliberately deferred—not oversights:

- **No frontend automated tests.** The backend has a real test suite running
  in CI (`backend-tests.yml`); the frontend currently does not, and that's a
  deliberate scope decision for now rather than a gap waiting to be filled.
