# Vinyhound

Small Go HTTP service that lets users sign up, log in, and view content that belongs only to them. Data now lives in Postgres instead of an in-memory map, and a lightweight Svelte SPA is included for interacting with the API.

## Prerequisites

- Go 1.21+
- Node 18+ (for the Svelte frontend)
- Postgres 14+ (or compatible hosted option)

## Database setup

1. Create a Postgres database that the app can reach.
2. Apply the schema located at `db/schema.sql`. Example:
   ```bash
   psql $DATABASE_URL -f db/schema.sql
   ```

`DATABASE_URL` should be a standard Postgres connection string such as `postgres://user:pass@localhost:5432/vinyhound?sslmode=disable`.

## Run locally

```bash
go mod tidy
go run ./cmd/vinyhound
```

The backend expects `DATABASE_URL` in the environment before it starts. In another terminal for the web app:

```bash
cd web
npm install
npm run dev
```

The Svelte app listens on http://localhost:5173 and proxies API calls to the Go server running at http://localhost:8080. Set `PORT` or run `npm run dev -- --port <port>` to customize.

## API

- POST /signup – create an account. Body:
  ```json
  { "username": "alice", "password": "secret", "content": ["First playlist", "Second playlist"] }
  ```
- POST /login – receive a bearer token for authenticated calls.
  ```json
  { "username": "alice", "password": "secret" }
  ```
  Response:
  ```json
  { "token": "..." }
  ```
- GET /me/content – return the current user content. Requires `Authorization: Bearer <token>`.
- PUT /me/content – replace the current content array. Same auth header.

## Frontend UX

The Svelte SPA provides:

- Account creation form with optional seed content (one entry per line).
- Login form (pre-filled with the "demo / demo123" credentials).
- Authenticated view showing your stored entries, a textarea editor, and update/refresh actions.
- Session persistence via localStorage so reloads keep you signed in.

## Notes

- Passwords are hashed with bcrypt before storage.
- Session tokens persist in the `sessions` table; delete rows or drop the table to invalidate them.
- A demo account (`demo` / `demo123`) is created at startup with sample content if it is not already present.
