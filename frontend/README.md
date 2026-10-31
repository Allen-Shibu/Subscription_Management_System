# Frontend

React and Vite frontend for the Go subscription management API. Clerk handles sign-in and the dashboard user menu.

## Run locally

1. Run `npm install` in this directory.
2. Copy `.env.example` to `.env` and set `VITE_CLERK_PUBLISHABLE_KEY` to your Clerk publishable key.
3. Run `npm run dev`.

Run `npm run build` and `npm run lint` to check the frontend.

The dashboard and admin catalog currently use sample data. Catalog and staff edits live only in the current browser tab and disappear on refresh. The admin navigation and routes require an active Clerk organization with the `org:admin` role; no catalog or invoice API is connected yet.
