# Vinyhound

Small Go HTTP service that lets users sign up, log in, and view content that belongs only to them. Everything is kept in memory so you can explore the flow without a database. A lightweight Svelte SPA is included for interacting with the API.

## Prerequisites

- Go 1.21+
- Node 18+ (for the Svelte frontend)

## Run locally

`ash
go mod tidy
go run ./cmd/vinyhound
`

In another terminal for the web app:

`ash
cd web
npm install
npm run dev
`

The Svelte app will listen on http://localhost:5173 and proxies API calls to the Go server running at http://localhost:8080. Set PORT or run 
pm run dev -- --port <port> to customize.

## API

- POST /signup ? create an account. Body:
  `json
  { "username": "alice", "password": "secret", "content": ["First playlist", "Second playlist"] }
  `
- POST /login ? receive a bearer token for authenticated calls.
  `json
  { "username": "alice", "password": "secret" }
  `
  Response:
  `json
  { "token": "..." }
  `
- GET /me/content ? return the current user content. Requires Authorization: Bearer <token>.
- PUT /me/content ? replace the current content array. Same auth header.

## Frontend UX

The Svelte SPA provides:

- Account creation form with optional seed content (one entry per line).
- Login form (pre-filled with the "demo / demo123" credentials).
- Authenticated view showing your stored entries, a textarea editor, and update/refresh actions.
- Session persistence via localStorage so reloads keep you signed in.

## Notes

- Passwords are hashed with bcrypt before storage.
- Sessions are held in memory; restarting the server clears them.
- A demo account (demo / demo123) is created at startup with sample content.
