# Developer Notes

This document is for anyone writing code in this repository — not for end users
deploying the app. If you're looking for deployment instructions, see
[DEPLOY.md](./DEPLOY.md). If you're looking for the high-level architecture
pitch, see [README.md](./README.md).

The goal of this file is to write down the things that aren't obvious from
reading the code in isolation — invariants, gotchas, and "don't do this"
warnings that a future contributor (including future-you) won't necessarily
rediscover just by reading a single function.

## The charity math is a snapshot, not a live calculation

This is the single most important thing to understand before touching the
`Ledgers` table, the ledger handlers, or anything in `queries.go` that sums
`charity_owed` or `charity_fulfilled`.

### How it actually works

Every row in `Ledgers` has two fields: `charity_owed` and `charity_fulfilled`.
These are **not** computed on the fly when you read a row, and they are
**not** kept in sync with the user's current settings. They are calculated
exactly once, at the moment a row is inserted (`setEntry` in
`ledger_post.go`), using whatever `donation_percentage` the user happens to
have set *at that instant*. Once written, the value is frozen. It will never
change again unless that specific row is deleted and a new row is inserted
in its place.

Concretely:

- When a **paycheck** entry is inserted, `charity_owed` is calculated as
  `amount * (current_donation_percentage / 100)`, and written into that row.
- When a **donation** entry is inserted, `charity_fulfilled` is calculated as
  `(amount / total_owed_so_far) * 100`, where `total_owed_so_far` is the sum
  of `charity_owed` across all of the user's existing rows at that moment.

If the user changes their donation percentage later (via
`PATCH /users/settings`), nothing about previously-inserted rows changes.
Only new rows inserted after the change will use the new percentage. This is
intentional and is covered by a test (`TestUserJourneyEndToEnd` in
`flow_test.go` posts a paycheck at 10%, changes the setting to 20%, then
posts a second paycheck and asserts the second one owes double the first).

### Why this becomes a problem: delete-and-recreate

There is currently no "edit" or "update" endpoint for ledger entries. The
only way to "fix" an existing entry — a typo in the description, a wrong
amount, anything — is to delete it (`DELETE /entries/{id}`) and then create a
new one (`POST /entries`) to replace it.

Doing this has two side effects that are easy to miss:

1. **The transaction date changes.** `transaction_date` defaults to
   `CURRENT_TIMESTAMP` at insert time (see the schema in
   `000_official_look.sql`). The deleted row's original date is gone forever.
   The recreated row is dated "now," not whenever the original transaction
   actually happened. Since entries are listed `ORDER BY transaction_date
   DESC`, this can move the entry to a different position in the user's
   history than where it actually occurred.

2. **The charity math is recalculated using today's settings, not the
   settings that were in effect when the original entry was created.** If
   the user's `donation_percentage` has changed at any point between the
   original insert and the delete-and-recreate, the new row's `charity_owed`
   (or `charity_fulfilled`, if the entry was a donation) will silently differ
   from what the original row had — with no warning, no audit trail, and no
   indication anywhere that anything changed.

There's a second-order effect worth knowing about too: because
`getAmountOwed` and `getAmountFulfilled` in `queries.go` sum across *all* of
a user's rows, deleting and recreating one paycheck can shift the aggregate
"total owed" figure the dashboard shows — which in turn changes the
*meaning* of `charity_fulfilled` values on donation rows that were calculated
against the old total. Those donation rows aren't touched directly, but the
percentage they represent effectively moves out from under them.

### Why this matters for you, specifically

None of this shows up as a bug in code review or in casual testing, because
every individual operation (delete, insert) does exactly what it's supposed
to do and would pass its own isolated unit test. The problem is emergent: it
only appears when you combine "snapshot math," "no update endpoint," and
"transaction_date resets on insert" together, across multiple operations.

If you're adding a feature that touches ledger entries — a real PATCH/update
endpoint, a CSV import/export feature (this is on the README roadmap), a
bulk-edit tool, an admin script to backfill data, anything that deletes and
recreates rows under the hood — you need to either:

- Preserve the original `transaction_date` and the original
  `charity_owed`/`charity_fulfilled` values explicitly, rather than letting
  a fresh `INSERT` recalculate them, or
- Make it very clear to the end user (and to yourself, in the code) that
  what's happening is a recalculation against current settings, not a
  faithful edit.

If you build a real update endpoint, this whole problem mostly goes away —
an `UPDATE` that only touches `amount`/`description` and leaves
`charity_owed`, `charity_fulfilled`, and `transaction_date` untouched would
sidestep all of this. That's the cleanest long-term fix, but it's currently
out of scope (see Roadmap below).

In the meantime, the only place this is surfaced to an end user is a hover
tooltip on the delete button in `Ledger.tsx`. That's better than nothing, but
it's not something a developer working in `queries.go` or `ledger_post.go`
would ever see. This file is that missing piece.

## Other things worth knowing

### Two separate Postgres instances, on purpose

`docker-compose.yml` defines `db-dev` and `db-test` as separate Postgres
containers on separate ports (`DEV_PORT` / `TEST_PORT` in `.env`). This is
deliberate — tests call `clearDatabase()` (`TRUNCATE ... RESTART IDENTITY
CASCADE`) between runs, and you do not want that running against your dev
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
check — there isn't an env var for it currently.

### The `10` fallback in `getDonationPercent` is intentional

If you read `getDonationPercent` in `queries.go`, you'll see it defaults to
`10.0` if the DB lookup fails for any reason (no rows, query error). This
matches the schema default (`donation_percentage INT DEFAULT 10`) and is a
deliberate fail-safe, not a placeholder that was meant to be replaced. Don't
"fix" it by removing it.

## Roadmap / known gaps

These are known, deliberately deferred — not oversights:

- **No update/edit endpoint for ledger entries.** This is the root cause of
  the delete-and-recreate problem described above. Deferred, not yet
  scheduled.
- **No frontend automated tests.** The backend has a real test suite running
  in CI (`backend-tests.yml`); the frontend currently does not, and that's a
  deliberate scope decision for now rather than a gap waiting to be filled.
- **CSV export/import** is listed as a future item in the main README.
