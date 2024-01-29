# supabase-jwt-auth-demo

> [DepsHub.com](https://depshub.com) - The easiest way to stay up-to-date with your dependencies.

This demo shows how to use [Supabase](https://supabase.io) to authenticate users with JWTs using your own backend.

Read this tutorial to learn how to build this demo: [Using Supabase as an Auth Service](https://depshub.com/blog/using-supabase-auth-as-a-service-with-a-custom-backend/)

## Setup

You need 3 environment variables:

- `SUPABASE_JWT_SECRET` to run your backend. You can either set it in your `.env` file or in your terminal before running the backend.
- `VITE_SUPABASE_PROJECT_ID` and `VITE_SUPABASE_ANON_KEY` to build your frontend. To do that, simply create a `.env.local` file in the `frontend` directory and add your variables there. Learn more on the Vite [docs](https://vitejs.dev/guide/env-and-mode).

## Running

- Run `go run .` in the `backend` folder
- Run `yarn dev` in the `frontend` folder
