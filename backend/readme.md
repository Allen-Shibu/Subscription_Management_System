# Recur Backend

Backend API service for Brightloop Software Subscription Billing Platform.

## Requirements
- Go 1.27.3 or newer
- Clerk account for authentication
- A Neon PostgreSQL database and `psql` for migrations

## Configuration
Create a `.env` file in `backend/` (or copy from `.env.example`):

```bash
PORT=3000
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
GIN_MODE=debug
DATABASE_URL=postgresql://USER:PASSWORD@YOUR_DIRECT_ENDPOINT.neon.tech/DBNAME?sslmode=require
CLERK_SECRET_KEY=sk_test_...
CLERK_PUBLISHABLE_KEY=pk_test_...
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=...
SMTP_PASSWORD=...
SMTP_FROM="Brightloop Billing <billing@example.com>"
```

Set `CORS_ALLOWED_ORIGINS` to the exact frontend origins (comma-separated, without paths). For a deployed frontend, replace the local values with its HTTPS origin. Browser requests from other origins receive no CORS permission. The frontend sends a Clerk bearer token in the `Authorization` header; cookie credentials are not enabled.

Copy the **direct (unpooled)** connection string from the Neon console into `backend/.env` as `DATABASE_URL`; keep the TLS settings supplied by Neon. The same URL is used by the Go server and the migration script. Keep this credential on the backend only, never in a `VITE_` variable or a committed file. Neon provides the database, so Docker is not needed.

## Running the Server

```bash
make migrate-up   # apply migration/*_up.sql
make run          # or: go run ./cmd/api
```

Migrations are tracked in `schema_migrations`, so `make migrate-up` is safe to re-run.
`make migrate-down` reverts the latest applied migrations. Run migrations against a development Neon branch before using a production branch.

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
| `subscriptions` | Recurring subscription for a customer on a product tier. Price, currency and interval are snapshotted from the variant at creation. Supports `cancel_at_period_end`, `canceled_at`, `ended_at` and status `active` / `trialing` / `past_due` / `canceled` / `expired`. Admins can attach a `tax_rule_id`. |
| `invoices` | One row per billing period (`period_start` is unique per subscription). Auto number `INV-000001`, status `draft` / `open` / `paid` / `void` / `uncollectible`, plus `email_sent_at` for delivery tracking. Tax rate and amount are snapshotted. `total_cents = subtotal_cents - discount_cents + tax_cents`. |
| `subscription_discounts` | One-time discounts. `applied_invoice_id` stays `NULL` until an invoice consumes it. |
| `refunds` | Money returned against a paid invoice. |
| `tax_rules` | Admin-defined `name` + `rate_bps` (basis points, 1800 = 18%). |

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
| `GET` | `/api/staff/invoices/:invoiceID/pdf` | Download the invoice as a PDF. |
| `POST` | `/api/staff/invoices/:invoiceID/send` | Email the invoice PDF to the customer and record `emailSentAt`. |
| `POST` | `/api/staff/invoices/:invoiceID/pay` | Mark an open invoice as paid. |
| `POST` | `/api/staff/invoices/:invoiceID/refunds` | Refund a paid invoice (`amountCents`, `reason`), capped at the unrefunded total. |
| `GET` | `/api/staff/invoices/:invoiceID/refunds` | List refunds. |
| `PUT` | `/api/admin/subscriptions/:id/tax-rule` | Admin only. `{"taxRuleId": 1}` attaches a tax rule; `null` clears it. |

<<<<<<< HEAD
Invoices are numbered automatically (`INV-000001`). Issuing one applies the subscription's tax rule: tax is charged on `subtotal - discount`, and the rate is snapshotted on the invoice. Emailing invoices is not implemented yet; `emailSentAt` stays `null`.

## Tax rules

Admins create and change rules; any signed-in user can read them.

| Method | Endpoint | Auth Required |
|---|---|---|
| `POST` | `/api/admin/tax-rules` | Admin |
| `PUT` | `/api/admin/tax-rules/:id` | Admin |
| `DELETE` | `/api/admin/tax-rules/:id` | Admin |
| `GET` | `/api/tax-rules` | Any signed-in user |
| `GET` | `/api/tax-rules/:id` | Any signed-in user |

```json
{ "name": "GST 18%", "rateBps": 1800, "active": true }
```

## Staff accounts

`POST /api/admin/staff` with `{"email":"staff@example.com"}` sends a Clerk invitation carrying `public_metadata.role = "staff"`. The invitee sets their own password, so no credentials pass through this API. Admin role required.
=======
Invoices are numbered automatically (`INV-000001`). `PDF` is rendered with gofpdf.
`send` emails the PDF through SMTP using the customer's `users.email`; it returns `409`
if the invoice was already sent, `422` if the customer has no email, and `502` if
delivery fails. Sending requires `SMTP_HOST`, `SMTP_FROM`, and (usually) credentials.
Port `587` uses STARTTLS, port `465` uses implicit TLS.
>>>>>>> 9022f34 (feat(backend):add invoice service)

## Example

```bash
curl -X POST localhost:3000/api/admin/products \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"Starter","description":"Basic plan"}'

curl -X POST localhost:3000/api/admin/products/1/variants \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"Pro","priceCents":1999,"billingInterval":"month"}'
```
