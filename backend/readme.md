# Recur Backend

Backend API service for Brightloop Software Subscription Billing Platform.

## Requirements
- Go 1.22+ (tested on Go 1.27)
- Clerk account for authentication
- Docker (for local Postgres) and `psql`

## Configuration
Create a `.env` file in `backend/` (or copy from `.env.example`):

```bash
PORT=3000
GIN_MODE=debug
DATABASE_URL=postgres://recur:recur@localhost:5432/recur?sslmode=disable
CLERK_SECRET_KEY=sk_test_...
CLERK_PUBLISHABLE_KEY=pk_test_...
```

## Running the Server

```bash
make db-up        # start local Postgres in Docker
make migrate-up   # apply migration/*_up.sql
make run          # or: go run ./cmd/api
```

Migrations are tracked in `schema_migrations`, so `make migrate-up` is safe to re-run.
`make migrate-down` reverts the latest applied migrations. No Docker? Point `DATABASE_URL` at any Postgres and run the same commands.

The server starts by default on `http://localhost:3000`.

## Roles

Roles live in Clerk; the local `users.role` column is only a queryable copy.
Use one of `admin`, `staff`, `customer`.

In the Clerk dashboard add a custom session claim (Sessions → Customize session token):

```json
{ "role": "{{user.public_metadata.role}}" }
```

Organization roles (`org:admin`, `org:staff`, `org:customer`) are accepted too.

## Running Tests
```bash
go test ./...
```

## Endpoints

| Method | Endpoint | Auth Required | Description |
|---|---|---|---|
| `GET` | `/health` | No | Service health check |
| `GET` | `/api/health` | No | API health check |
| `GET` | `/api/me` | Yes (Bearer Token) | Returns authenticated user claims & role |
| `GET` | `/api/admin/ping` | Yes (Admin Role) | Admin-only test endpoint |
