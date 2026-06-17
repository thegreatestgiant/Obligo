# Charity-Tracker

Charity-Tracker is a self-hosted web app for tracking your income and your
charitable giving against each other. You log paychecks as they come in, log
donations as you make them, and the app tells you what percentage of your
earnings you've actually given — measured against a giving target you set
yourself.

It's built for anyone who gives a portion of their income to charity on a
recurring basis and wants a running, persistent answer to "am I actually
keeping up with my own goal?" without maintaining a spreadsheet.

## What it does

- **Log paychecks and donations.** Each entry is a simple amount, a type
  (paycheck or donation), an optional description, and a date.
- **Set a giving target.** By default, the app assumes you want to give 10%
  of your earnings to charity, but this is fully customizable per user from
  the Settings page.
- **Track progress automatically.** Every time you log a paycheck, the app
  calculates how much you now owe toward your giving target. Every time you
  log a donation, it calculates how much of that obligation you've just
  fulfilled. The dashboard shows your total earned, total donated, total
  owed, and what percentage of your obligation is fulfilled — updated live
  as you add entries.
- **See your history.** A paginated ledger lists every paycheck and donation
  you've logged, with the ability to delete entries you no longer want.

## Who it's for

This is built as a personal finance tool for a single household or
individual, not a multi-tenant SaaS product. Each user has their own account,
their own ledger, and their own giving target — there's no sharing or
admin/family-account model. It's meant to be self-hosted and run locally or
on infrastructure you control, not deployed as a public service.

## How to use it

1. **Register an account** and log in.
2. **Set your giving target** under Settings (or leave it at the 10%
   default).
3. **Log entries as they happen.** Whenever you get paid, log a paycheck
   entry for that amount. Whenever you donate, log a donation entry. The app
   does the percentage math for you.
4. **Check the dashboard** any time to see your running totals: how much
   you've earned, how much you owe toward your giving target, how much
   you've actually given, and what percentage of your obligation is
   fulfilled.

One thing worth knowing up front: the amount you "owe" on any given paycheck
is calculated using your giving target *at the moment you log that
paycheck*, not recalculated later if you change your target. If you change
your target percentage, it only affects new entries going forward — your
history doesn't retroactively change. See
[DEVELOPERS.md](./DEVELOPERS.md) if you want the full technical explanation
of why this matters, especially before deleting and re-adding an old entry.

## Running it

Charity-Tracker ships as a single Docker image with the database
provisioning handled for you via Docker Compose — see
[DEPLOY.md](./DEPLOY.md) for setup instructions.

If you want to work on the project itself rather than just run it, see
[DEVELOPERS.md](./DEVELOPERS.md) for architecture notes, known gaps, and
things to watch out for in the code.

## Architecture, in brief

- **Backend:** Go, serving both the API and the built frontend as a single
  binary.
- **Frontend:** React + TypeScript + Vite, built and embedded into the Go
  binary at Docker build time.
- **Database:** PostgreSQL, with schema migrations embedded in the Go binary
  and applied automatically on startup.
- **Deployment:** A single hardened, non-root Docker container, built via
  GitHub Actions and intended to run alongside a Postgres container via
  Docker Compose.

For the reasoning behind these choices, and details on things like how the
frontend gets its API URL at runtime, see [DEVELOPERS.md](./DEVELOPERS.md).

## Roadmap

- [x] Graceful shutdown and structured logging (`log/slog`)
- [x] Docker non-root hardening
- [x] Self-contained database migrations
- [ ] Implement CSV transaction Export/Import functionality
