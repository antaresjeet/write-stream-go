# Write Stream Go Project

This is a Go-based backend for a blog app featuring authentication with Google, GraphQL API, and user management. The project uses the following stack:

- **Gin**: A web framework to handle HTTP requests
- **Gqlgen**: GraphQL implementation in Go
- **GORM**: ORM library for Go
- **Go Migrate**: Database migrations for PostgreSQL

## Features

- Google OAuth2 authentication
- GraphQL API (with playground)
- User creation and management
- PostgreSQL database integration

## Prerequisites

Before setting up the project, ensure you have the following installed:

- Go 1.19+
- PostgreSQL
- Gqlgen
- Go Migrate

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/antares-jeet/write-stream-go.git
   cd write-stream-go
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Create a `.env` file based on `.env.example`:

   ```bash
   cp .env.example .env
   ```

4. Update the `.env` file with your environment variables, such as database credentials, OAuth client IDs, and secrets.

5. Apply database migrations:

   ```bash
   migrate -database "postgres://<DB_USER>:<DB_PASSWORD>@<DB_HOST>:<DB_PORT>/<DB_NAME>?sslmode=require" -path ./migrations up
   ```

6. Generate the GraphQL schema using Gqlgen:

   ```bash
   go run -mod=mod github.com/99designs/gqlgen generate
   ```

7. Run the application:
   ```bash
   go run ./main.go
   ```

The server will be up at `http://localhost:8080`.

## Running in Production

For production, switch to release mode by setting the following environment variable:

```bash
export GIN_MODE=release
```

You can also check [development.md](development.md) for more detailed information on setting up the project and its development environment.

## License

This project is licensed under the MIT License.
