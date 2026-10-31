# Frontend

React and Vite frontend for the Go subscription management API. Clerk handles sign-in and the dashboard user menu.

## Run locally

1. Run `npm install` in this directory.
2. Copy `.env.example` to `.env` and set `VITE_CLERK_PUBLISHABLE_KEY` to your Clerk publishable key.
3. Run `npm run dev`.

Run `npm run build` and `npm run lint` to check the frontend.

The dashboard, taxes, discounts, quotation templates, and staff screens use sample data. Products and plans use a browser mock catalog by default (`VITE_USE_MOCK=true`); its edits persist in localStorage. The admin navigation and routes require an active Clerk organization with the `org:admin` role.

To use a server, set `VITE_USE_MOCK=false` and `VITE_API_URL` to its origin. The frontend sends a Clerk bearer token. Product reads use `GET /api/products` and `GET /api/products/:id`; product writes use `POST/PUT/DELETE /api/admin/products`. Plan reads and writes use `/api/plans` and `/api/admin/plans` respectively. List endpoints must accept `q`, `page`, and `pageSize` and return `{items,total,page,pageSize}`. Product payloads must include `name`, `type`, `salesPrice`, `costPrice`, `notes`, and `variants: [{attribute,value,extraPrice}]`; plan payloads include the fields shown in the form. Validation failures should return HTTP 400 with `{error,fields:{fieldName:message}}`. The current Go server exposes different product fields and no plan endpoints, so real mode reports that mismatch and prevents writes until its contract is extended.
