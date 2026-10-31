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
| `GET` | `/api/products` | Yes | List products |
| `GET` | `/api/products/:id` | Yes | Get a product |
| `GET` | `/api/products/:id/variants` | Yes | List a product's tiers |
| `GET` | `/api/variants/:variantID` | Yes | Get a tier |
| `POST` | `/api/admin/products` | Yes (Admin Role) | Create a product |
| `PUT` | `/api/admin/products/:id` | Yes (Admin Role) | Replace a product |
| `DELETE` | `/api/admin/products/:id` | Yes (Admin Role) | Delete a product and its tiers |
| `POST` | `/api/admin/products/:id/variants` | Yes (Admin Role) | Create a tier |
| `PUT` | `/api/admin/variants/:variantID` | Yes (Admin Role) | Replace a tier |
| `DELETE` | `/api/admin/variants/:variantID` | Yes (Admin Role) | Delete a tier |

Prices are stored in integer cents (`priceCents`) to avoid floating point money bugs.
`billingInterval` is `month` or `year`; `currency` is a 3-letter code (defaults to `USD`).

## Example

```bash
curl -X POST localhost:3000/api/admin/products \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"Starter","description":"Basic plan"}'

curl -X POST localhost:3000/api/admin/products/1/variants \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"Pro","priceCents":1999,"billingInterval":"month"}'
```
