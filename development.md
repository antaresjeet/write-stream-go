# Development Guide

## Setting up for Development

1. **Install Go**  
   Make sure Go is installed on your system. You can verify the installation with:

   ```bash
   go version
   ```

2. **Install Dependencies**  
   All required dependencies are listed in `go.mod`. You can download them using:

   ```bash
   go mod download
   ```

3. **Environment Variables**  
   Use the `.env.example` file as a reference to create your own `.env`:

   ```bash
   cp .env.example .env
   ```

   Ensure to fill in the correct values for your PostgreSQL connection and Google OAuth credentials.

4. **Database Setup**  
   Migrations are stored in the `migrations` folder. Use [golang-migrate](https://github.com/golang-migrate/migrate) to run them:

   ```bash
   migrate -database "postgres://<username>:<password>@<host>:<port>/<dbname>?sslmode=require" -path ./migrations up
   ```

5. **GraphQL Code Generation**  
   Whenever you make changes to the GraphQL schema, you should regenerate the GraphQL code using gqlgen:

   ```bash
   go run -mod=mpd github.com/99designs/gqlgen generate
   ```

## Common Commands

- Run the app:

  ```bash
  go run ./main.go
  ```

- Run the migrations:

  ```bash
  migrate -database "postgres://<username>:<password>@<host>:<port>/<dbname>?sslmode=require" -path ./migrations up
  ```

- Generate GraphQL files:

  ```bash
  go run github.com/99designs/gqlgen generate
  ```

## Debugging

- **Switch to Release Mode**  
  If you're in production, you should switch to GIN's release mode:

  ```bash
  export GIN_MODE=release
  ```

- **Handle Proxy Warning**  
  By default, GIN trusts all proxies. In a production environment, it is safer to set the trusted proxy manually.
