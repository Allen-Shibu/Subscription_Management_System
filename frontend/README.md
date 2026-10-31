# Frontend

React and Vite frontend for the Go subscription management API. Clerk handles sign-in and the dashboard user menu.

## Run locally

1. Run `npm install` in this directory.
2. Copy `.env.example` to `.env` and set `VITE_CLERK_PUBLISHABLE_KEY` to your Clerk publishable key.
3. Run `npm run dev`.

Run `npm run build` and `npm run lint` to check the frontend.

By default (`VITE_USE_MOCK=true`), products and plans use a browser mock catalog whose edits persist in localStorage. The dashboard, subscriptions, taxes, discounts, quotation templates, and staff screens show sample data. Admin navigation and routes accept a Clerk `org:admin` role or a custom session `admin` role.

To connect to the existing Go API, set `VITE_USE_MOCK=false` and `VITE_API_URL` to its origin, then restart Vite. The frontend sends a Clerk bearer token. In this mode, customer product browsing uses the authenticated product GET routes. Admin product editing uses `POST/PUT/DELETE /api/admin/products`; billing tier editing uses the separate variant routes. Product search and pagination run in the browser because the Go list endpoint returns an array. Subscriptions use `GET /api/staff/subscriptions`, taxes use `/api/tax-rules` and `/api/admin/tax-rules`, and staff invitations use `POST /api/admin/staff`.

Plans remain a labeled browser mock in either mode because the Go API has no plan endpoints. Dashboard totals, discounts, and quotation templates remain sample data; billing remains a placeholder. The live product editor uses the server's `name`, `description`, `active`, and billing tier fields. The mock editor retains the separate sales price, cost price, and variant adjustment design, which the server cannot store. The Go server requires Clerk and PostgreSQL to be running.
