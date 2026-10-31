# Clerk React app

A React and Vite frontend with Clerk sign-in and a user menu after sign-in. The Go backend is separate from this folder.

## Run locally

1. Run `npm install`.
2. Copy `.env.example` to `.env` and set `VITE_CLERK_PUBLISHABLE_KEY` to the publishable key from your [Clerk dashboard](https://dashboard.clerk.com/).
3. Run `npm run dev` and open the URL Vite prints.

The local `.env` is ignored by Git. This workspace already has a test publishable key in `.env`; replace it if you want to use a different Clerk instance. Configure the corresponding instance's sign-in methods in the Clerk dashboard.

Run `npm run build` to verify the production build.
