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

## Billing schema

| Table | Purpose |
|---|---|
| `subscriptions` | Recurring subscription for a customer on a product tier. Price, currency and interval are snapshotted from the variant at creation. Supports `cancel_at_period_end`, `canceled_at`, `ended_at` and status `active` / `trialing` / `past_due` / `canceled` / `expired`. |
| `invoices` | One row per billing period (`period_start` is unique per subscription). Auto number `INV-000001`, status `draft` / `open` / `paid` / `void` / `uncollectible`, plus `email_sent_at` for delivery tracking. `total_cents = subtotal_cents - discount_cents`. |
| `subscription_discounts` | One-time discounts. `applied_invoice_id` stays `NULL` until an invoice consumes it. |
| `refunds` | Money returned against a paid invoice. |

Customers are rows in `users`, so a subscription references `users.id`.

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

## Subscription endpoints

All of these require the staff or admin role.

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/staff/subscriptions` | Create a subscription (`customerId`, `variantId`, optional `quantity`, `startAt`). Price and interval are snapshotted from the variant. |
| `GET` | `/api/staff/subscriptions` | List subscriptions, filter with `?customerId=` and `?status=`. |
| `GET` | `/api/staff/subscriptions/:id` | Get a subscription. |
| `PATCH` | `/api/staff/subscriptions/:id` | Change `quantity`, `currentPeriodStart`, `currentPeriodEnd`, `cancelAtPeriodEnd`. |
| `POST` | `/api/staff/subscriptions/:id/cancel` | `{"atPeriodEnd":true}` flags cancellation; `false` cancels immediately. |
| `POST` | `/api/staff/subscriptions/:id/discounts` | Add a one-time discount (`amountCents`, `reason`). |
| `GET` | `/api/staff/subscriptions/:id/discounts` | List discounts and whether each was applied. |
| `POST` | `/api/staff/subscriptions/:id/invoices` | Issue an invoice for the current period, consume pending discounts, and advance the subscription to the next period. |
| `GET` | `/api/staff/subscriptions/:id/invoices` | List invoices. |
| `GET` | `/api/staff/invoices/:invoiceID` | Get one invoice. |
| `POST` | `/api/staff/invoices/:invoiceID/pay` | Mark an open invoice as paid. |
| `POST` | `/api/staff/invoices/:invoiceID/refunds` | Refund a paid invoice (`amountCents`, `reason`), capped at the unrefunded total. |
| `GET` | `/api/staff/invoices/:invoiceID/refunds` | List refunds. |

Invoices are numbered automatically (`INV-000001`). Emailing them is not implemented yet; `emailSentAt` stays `null`.

## Example

```bash
curl -X POST localhost:3000/api/admin/products \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"Starter","description":"Basic plan"}'

curl -X POST localhost:3000/api/admin/products/1/variants \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"Pro","priceCents":1999,"billingInterval":"month"}'
```
