# Chirpy

Chirpy is a lightweight RESTful API server built in Go for a micro-blogging platform. It handles user authentication, post creation, profanity filtering, and webhook-driven membership upgrades.

## Features

- **User Management**: Secure password hashing with `argon2id` and JWT-based authentication.
- **Refresh Tokens**: Revocable refresh token flow for persistent user sessions.
- **Chirp Management**: Create, view, and delete short text posts (chirps) with profanity masking.
- **Database Storage**: PostgreSQL integration using `sqlc` for compile-time safe database queries.
- **Webhooks**: Integration endpoint to upgrade users to "Chirpy Red" status.
- **Metrics**: Administrative endpoints to track hits and server health.

## Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Query Builder**: sqlc
- **Migrations**: Goose

## Getting Started

### Prerequisites

- Go (1.26.0 or later)
- PostgreSQL
- Goose (for database migrations)

### Installation & Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/sandro-gg/chirpy.git
   cd chirpy
   ```

2. **Configure environment variables:**
   Create your `.env` file in the root directory by copying the example template:

   ```bash
   cp .env.example .env
   ```

   ```env
   DB_URL=your_db_url_here
   PLATFORM=dev
   JWT_SECRET=your_secret_here
   POLKA_KEY=your_key_here
   PORT=your_port_here
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

The server should now be running on your configured port (default 8080).

## Example Usage

Here is how you can test the core workflows locally using `curl`:

1. **Register a new user:**

   ```bash
   curl -X POST http://localhost:8080/api/users \
     -H "Content-Type: application/json" \
     -d '{"email": "sandro@example.com", "password": "mypassword123"}'
   ```

2. **Log in to receive your access token:**

   ```bash
   curl -X POST http://localhost:8080/api/login \
      -H "Content-Type: application/json" \
      -d '{"email": "sandro@example.com", "password": "mypassword123"}'
   ```

3. **Create an authenticated chirp (replace <YOUR_ACCESS_TOKEN> with the token from login):**
   ```bash
   curl -X POST http://localhost:8080/api/chirps \
      -H "Authorization: Bearer <YOUR_ACCESS_TOKEN>" \
      -H "Content-Type: application/json" \
      -d '{"body": "This is my first portfolio chirp!"}'
   ```

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
