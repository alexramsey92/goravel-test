# goravel-test

A SaaS scaffold built with [Goravel](https://goravel.dev) — a Go framework that mirrors Laravel's architecture. Built to evaluate whether Goravel can support the same patterns as a full Laravel SaaS (auth, CRM, CMS, REST API).

**Stack:** Go · Goravel · SQLite · Tailwind CSS (CDN) · JWT

**[Live demo page →](https://alexramsey92.github.io/goravel-test/)**

---

## What's included

| Area | Detail |
|---|---|
| **Auth** | Session-based login/register for the web UI; JWT Bearer for the API |
| **CRM** | Clients with status (lead / active / inactive), notes, company |
| **CMS** | Pages with title, slug, content, published toggle |
| **REST API** | `/api/auth/*` for token auth, `/api/v1/*` protected with Bearer JWT |
| **Roles** | `admin`, `user`, `client` — first registered user gets `admin` |

---

## Quick start

Requires **Go 1.22+**. No database server needed — SQLite file is created automatically.

```bash
# 1. Clone and install dependencies
git clone https://github.com/alexramsey92/goravel-test.git
cd goravel-test
go mod tidy

# 2. Set up environment
cp .env.example .env
go run . artisan key:generate
go run . artisan jwt:secret

# 3. Run migrations
go run . artisan migrate

# 4. Start the server
go run .
```

Open **http://localhost:3000** — register your first account (gets admin role).

### Hot reload with Air

```bash
go install github.com/air-verse/air@latest
air
```

---

## Routes

### Web UI

| Method | Path | Description |
|---|---|---|
| GET | `/login` | Login form |
| GET | `/register` | Register (first user = admin) |
| GET | `/dashboard` | Stats overview |
| GET | `/clients` | Client list |
| GET | `/clients/create` | Add client |
| GET | `/clients/:id` | View / edit / delete client |
| GET | `/pages` | CMS page list |
| GET | `/pages/create` | Create page |

### JSON API

All API responses are JSON. Protected routes require `Authorization: Bearer <token>`.

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/api/auth/register` | — | Register → returns JWT |
| POST | `/api/auth/login` | — | Login → returns JWT |
| GET | `/api/v1/me` | Bearer | Current user |
| GET | `/api/v1/clients` | Bearer | List clients |
| POST | `/api/v1/clients` | Bearer | Create client |
| GET | `/api/v1/clients/:id` | Bearer | Get client |
| DELETE | `/api/v1/clients/:id` | Bearer | Delete client |
| GET | `/api/v1/pages` | Bearer | List pages |
| POST | `/api/v1/pages` | Bearer | Create page |

#### Example

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -d "email=you@example.com&password=secret" | jq -r .token)

# Use it
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/me
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/clients
```

---

## Laravel → Goravel mapping

This project was built by someone with a Laravel background. Here's how the concepts translate:

| Laravel | Goravel |
|---|---|
| `Auth::user()` | `facades.Auth(ctx).User(&user)` |
| `Hash::make($pw)` | `facades.Hash().Make(pw)` |
| `Hash::check($pw, $hash)` | `facades.Hash().Check(pw, hash)` (returns `bool`) |
| `Model::query()->where(...)` | `facades.Orm().Query().Where(...)` |
| `->count()` | `->Count()` returns `(int64, error)` |
| `->create([...])` | `->Create(&model)` |
| `->save()` | `->Save(&model)` |
| `->delete()` | `->Delete(&model)` |
| `redirect('/dashboard')` | `ctx.Response().Redirect(302, "/dashboard")` |
| `view('clients.index', [...])` | `ctx.Response().View().Make("clients/index.tmpl", map[string]any{...})` |
| `Route::middleware('auth')->group(fn)` | `facades.Route().Middleware(mw.Handle).Group(func(r route.Router) {...})` |
| `Route::prefix('/api')->group(fn)` | `facades.Route().Prefix("/api").Group(func(r route.Router) {...})` |
| `request()->input('name')` | `ctx.Request().Input("name")` |
| `config/auth.php` | `config/auth.go` |
| `database/migrations/*.php` | `database/migrations/*.go` |
| `.env` | `.env` (identical format) |
| `php artisan migrate` | `go run . artisan migrate` |
| Laravel Mix / Vite | Tailwind CDN (no build step) |
| Blade `@if` / `@foreach` | Go `{{if}} {{range}}` |

### Gotchas

- **Templates** — Every `.tmpl` file must be wrapped in `{{define "path/name.tmpl"}}...{{end}}` so Go's `html/template` registers the full path as the template name (not just the base filename).
- **Session middleware** — You must apply `sessionmiddleware.StartSession()` to all web routes or `facades.Auth(ctx).Login()` will panic with "session driver is not set".
- **GroupFunc** — Route groups take `func(router route.Router)`, not `func()`.
- **Hash.Check** — Returns `bool`, not `(bool, error)` like you might expect.

---

## Switching databases

The default driver is SQLite (zero config). To switch:

### PostgreSQL

```bash
go get github.com/goravel/postgres
```

`bootstrap/providers.go` — swap `sqlite.ServiceProvider` for `postgres.ServiceProvider`

`config/database.go` — replace the sqlite connection block with:

```go
import postgresfacades "github.com/goravel/postgres/facades"

"default": "postgres",
"connections": map[string]any{
    "postgres": map[string]any{
        "host":     config.Env("DB_HOST", "127.0.0.1"),
        "port":     config.Env("DB_PORT", "5432"),
        "database": config.Env("DB_DATABASE", "goravel"),
        "username": config.Env("DB_USERNAME", "postgres"),
        "password": config.Env("DB_PASSWORD", ""),
        "sslmode":  "disable",
        "via": func() (driver.Driver, error) {
            return postgresfacades.Postgres("postgres")
        },
    },
},
```

### MySQL

Same pattern with `github.com/goravel/mysql` and `mysqlFacades.Mysql("mysql")`.

---

## Project structure

```
├── app/
│   ├── facades/          # Re-exports framework facades (Auth, Orm, Hash …)
│   ├── http/
│   │   ├── controllers/  # AuthController, DashboardController, ClientController, PageController
│   │   └── middleware/   # Auth (session), Admin (role check), ApiAuth (JWT Bearer)
│   └── models/           # User, Client, Page
├── bootstrap/            # App boot, service providers, migrations list
├── config/               # auth.go, database.go, session.go, jwt.go …
├── database/
│   └── migrations/       # Go migration files (auto-run on boot)
├── resources/
│   └── views/            # Go html/template files (.tmpl)
├── routes/
│   ├── web.go            # Session-based browser routes
│   └── api.go            # JWT-protected JSON API routes
└── main.go
```

---

## License

MIT
