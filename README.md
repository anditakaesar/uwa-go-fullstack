# Uwa Go Fullstack

## Architecture
Conceptually, I tried a Clean-architecture approach:
```
┌───────────────────────────────┐
│ cmd/                          │  → entrypoints (web, seeder)
└───────────────▲───────────────┘
                │
┌───────────────┴───────────────┐
│ internal/server               │  → wiring / composition root
└───────────────▲───────────────┘
                │
┌───────────────┴───────────────┐
│ internal/handler              │  → HTTP adapters, handlers, and renderer
│   ├── middleware              │
│   └── renderer                │
└───────────────▲───────────────┘
                │
┌───────────────┴───────────────┐
│ internal/service              │  → application logic (use cases)
└───────────────▲───────────────┘
                │
┌───────────────┴───────────────┐
│ internal/repo                 │  → infrastructure adapters
│ internal/infra                │  
└───────────────▲───────────────┘
                │
┌───────────────┴───────────────┐
│ internal/domain               │  → core model
└───────────────────────────────┘

```

## Setup
Use golang 1.25
Libraries:
- github.com/go-chi/chi/v5
- github.com/gorilla/csrf
- github.com/gorilla/sessions
- github.com/jackc/pgx/v5
- golang.org/x/crypto
- golang-migrate

Run `$ go mod vendor` to resolve references

### Database
Database can be deployed using podman/docker on `development` folder.
`~/development$ docker compose up -d`
Database that is used is `PostgreSQL`. PostgreSQL is used here with these reasons:
- I expect the schema of this project to evolve given time.
- I rely on complex queries such as to aggregate. This means that I pushed logic near the data related to it.
- This is the safer long-term database because its extensibility. Who knows I might need it to use extension like PostGIS (geospatial) or simply use its capability for full-text search and query-able JSONB data type.

### Migration
Check `Makefile`

- To create migration: `$ make create-migration name=new_migration`
- To migrate all pending migration: `$ make migrate`
- To migrate down one previous migration: `$ make migrate-down`
- To seed initial users: run `$ make seed-database` OR run manually
`$ DB_URL="postgres://here" go run ./cmd/seed/main.go`. Seed file can be done using separate `.csv` files:
    - `cmd/seed/users.csv` contains seed data for `users` table

### Env
Environment variables are located in file `.env-example`. This file eventualy will be `.env` on the deployemnt environment or loaded through different manners.
```
ENV="development"
PORT=":3000"
DB_URL="postgres://postgres:password@localhost:5433/backend_db?sslmode=disable"
COOKIE_SECRET="Rm57qySVRliOZg5WqJ5GyKHKY6f4sJ41"
CSRF_SECRET="4eQWYCt7WjxLwPmL06MhOW5FS96wxOk6"
JWT_SECRET="4eQWYCt7WjxLwPmL06MhOW5FS96wxOk6"
UPLOAD_DIR="../../uploads"
HOSTNAME="http://localhost:3000"
```

## Docs `/docs`
- TBD

## API
APIs are designed to utilize JSON:API specification [JSON:API Spec](https://jsonapi.org)
### GET /
Return the page to login via web.

### POST /login
Route to handle the authentication logic and return `cookies` and `JWT`. Primarily used by the form from the login page. The resulting token can be used to authenticate future endpoints.