# Recur Backend

Backend API service for Brightloop Software Subscription Billing Platform.

## Requirements
- Go 1.22+ (tested on Go 1.27)
- Clerk account for authentication

## Configuration
Create a `.env` file in `backend/` (or copy from `.env.example`):

```bash
PORT=3000
GIN_MODE=debug
CLERK_SECRET_KEY=sk_test_...
CLERK_PUBLISHABLE_KEY=pk_test_...
```

## Running the Server
```bash
go run ./cmd/api
```

The server starts by default on `http://localhost:3000`.

## Running Tests
```bash
go test -v ./...
```

## Endpoints

| Method | Endpoint | Auth Required | Description |
|---|---|---|---|
| `GET` | `/health` | No | Service health check |
| `GET` | `/api/health` | No | API health check |
| `GET` | `/api/me` | Yes (Bearer Token) | Returns authenticated user claims & role |
| `GET` | `/api/admin/ping` | Yes (Admin Role) | Admin-only test endpoint |
