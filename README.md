# Chirpy

Chirpy is a lightweight RESTful API server built in Go for a micro-blogging platform. It handles user authentication, post creation, profanity filtering, and webhook-driven membership upgrades.

## Features

- **User Management**: Secure password hashing with `bcrypt` and JWT-based authentication.
- **Refresh Tokens**: Revocable refresh token flow for persistent user sessions.
- **Chirp Management**: Create, view, and delete short text posts (chirps) with profanity masking.
- **Database Storage**: PostgreSQL integration using `sqlc` for compile-time safe database queries.
- **Webhooks**: Integration endpoint to upgrade users to "Chirpy Red" status.
- **Metrics**: Administrative endpoints to track hits and server health..

## Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Query Builder**: sqlc
- **Migrations**: Goose

## Getting Started

### Prerequisites

- Go (1.22 or later)
- PostgreSQL
- Goose (for database migrations)

### Installation & Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/<your-username>/chirpy.git
   cd chirpy
   ```

2. **Configure environment variables:**
   Create a `.env` file in the root directory:

   ```env
   PORT=8080
   DB_URL=postgres://<username>:<password>@localhost:5432/chirpy?sslmode=disable
   PLATFORM=dev
   JWT_SECRET=your_jwt_secret_key
   POLKA_KEY=your_polka_api_key
   ```

3. **Run database migrations:**

   ```bash
   cd sql/schema
   goose postgres "postgres://<username>:<password>@localhost:5432/chirpy?sslmode=disable" up
   cd ../..
   ```

4. **Build and run the server:**
   ```bash
   go build -o out && ./out
   ```

The server should now be running on `http://localhost:8080`.

## API Endpoints Overview

- `GET /api/healthz` - Health check
- `POST /api/users` - Create a user
- `POST /api/login` - Authenticate user & receive tokens
- `POST /api/refresh` - Issue new access token
- `POST /api/revoke` - Revoke refresh token
- `GET /api/chirps` - Fetch chirps (supports author and sorting queries)
- `POST /api/chirps` - Create a new chirp
- `DELETE /api/chirps/{chirpID}` - Delete a chirp (authenticated)
- `POST /api/polka/webhooks` - Upgrade user to Chirpy Red
