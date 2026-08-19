# Post Article — Backend (Go + Fiber)

REST API for blog article management. Built with Go 1.26, Fiber, GORM, and MySQL.

---

## Prerequisites

| Tool | Version | Install |
|---|---|---|
| Go | 1.21+ | https://go.dev/dl/ |
| MySQL | 8.x | https://dev.mysql.com/downloads/ |
| golang-migrate CLI | latest | see below |

### Install the migrate CLI

```bash
# macOS / Linux (Homebrew)
brew install golang-migrate

# Or download a binary from https://github.com/golang-migrate/migrate/releases
# and place it on your PATH.

# for wsl
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
```

---

## Getting started

### 1. Clone and install dependencies

```bash
cd backend
go mod download
```

### 2. Configure environment

```bash
cp .env.example .env
# Edit .env with your MySQL credentials
```

`.env` values:

| Key | Default | Description |
|---|---|---|
| `DB_HOST` | `127.0.0.1` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `<db_user>` | MySQL user |
| `DB_PASSWORD` | `<db_password>` | MySQL password |
| `DB_NAME` | `sv_article` | Database name (created by migration) |
| `SERVER_PORT` | `8080` | Port the API listens on |

### 3. Run migrations

The migration creates the `sv_article` database and the `posts` table.

Replace `<db_user>` and `<db_password>` with the values from your `.env` (`DB_USER`, `DB_PASSWORD`).

```bash
# Apply (up)
migrate -path ./migrations \
        -database "mysql://<db_user>:<db_password>@tcp(127.0.0.1:3306)/sv_article" \
        up

# Rollback (down)
migrate -path ./migrations \
        -database "mysql://<db_user>:<db_password>@tcp(127.0.0.1:3306)/sv_article" \
        down
```

> If your MySQL password is empty, omit the colon: `mysql://<db_user>@tcp(...)/sv_article`

### 4. Start the server

```bash
go run cmd/api/main.go
```

The API is now available at `http://localhost:8080`.

---

## API reference

### POST `/article/` — Create an article

```bash
curl -X POST http://localhost:8080/article/ \
  -H "Content-Type: application/json" \
  -d '{
    "title":    "A title that is long enough to pass validation",
    "content":  "This content must be at least two hundred characters long. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam.",
    "category": "Technology",
    "status":   "publish"
  }'
```

Response: `{}` (201 Created)

---

### GET `/article/:limit/:offset` — List articles (paginated)

```bash
# First 10 articles
curl http://localhost:8080/article/10/0

# Only published (for Dashboard / Preview)
curl "http://localhost:8080/article/10/0?status=publish"

# Only drafts
curl "http://localhost:8080/article/10/0?status=draft"
```

The optional `?status` query param filters by `publish`, `draft`, or `trash`. When omitted, all statuses are returned.

---

### GET `/article/:id` — Get one article

```bash
curl http://localhost:8080/article/1
```

Response: full article JSON including `id`, `created_date`, `updated_date`.

---

### PUT `/article/:id` — Update an article

```bash
curl -X PUT http://localhost:8080/article/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title":    "Updated title that still meets the minimum length",
    "content":  "Updated content. This must still be at least two hundred characters long. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore. Ut enim ad minim veniam, quis nostrud.",
    "category": "Technology",
    "status":   "draft"
  }'
```

Response: `{}`

---

### DELETE `/article/:id` — Delete an article

```bash
curl -X DELETE http://localhost:8080/article/1
```

Response: `{}`

> **Note:** The Dashboard "Trash" action does **not** call DELETE. It calls PUT with `"status": "trash"`. DELETE is a hard delete.

---

## Validation rules

| Field | Rules |
|---|---|
| `title` | required, minimum 20 characters |
| `content` | required, minimum 200 characters |
| `category` | required, minimum 3 characters |
| `status` | required, must be `publish`, `draft`, or `trash` |

Validation failures return HTTP 400 with a per-field error object:

```json
{
  "errors": {
    "title": "Title must be at least 20 characters",
    "status": "Status must be one of: publish, draft, trash"
  }
}
```

---

## Import the Postman collection

Open Postman → Import → select `postman_collection.json` from the repo root. All five endpoints are included with sample bodies that pass validation.
