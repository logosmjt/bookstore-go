# Bootcamp Golang

## Day 1

### Understanding Go Programming Features

- **Static Typing and Simplicity**: Go emphasizes simplicity while ensuring code robustness through static typing.
- **Concurrency Support**: Built-in concurrency support using goroutines and channels.
- **Compiled Language**: Compiled to machine code, offering fast execution speed.
- **Standard Library**: Comes with a powerful and comprehensive standard library.

### Environment Setup

- **Install Go**: Download and install from [golang.org](https://go.dev/), and set up the environment.
- **Install IDE**: Use Visual Studio Code with the Go extension for syntax highlighting, linting, debugging, etc.
- **Write Hello World**: Create a simple "Hello, World!" program to verify your environment.
- **Optional**: Write a simple poker game using TDD.

---

## Day 2

### Database Modeling, Connection, and Migration

1. Use [dbdiagram.io](https://dbdiagram.io/) for database design and export PostgreSQL SQL statements.
2. Learn and use [golang-migrate/migrate](https://github.com/golang-migrate/migrate) for database migrations.

Example commands:
```bash
migrate create -ext sql -dir db/migration -seq $(name)
migrate -path "$(MIGRATION_PATH)" -database "$(DB_URL)" -verbose up
migrate -path "$(MIGRATION_PATH)" -database "$(DB_URL)" -verbose down
```

Migration file examples:
```
000001_init_schema.up.sql
000001_init_schema.down.sql
000002_add_session.up.sql
000002_add_session.down.sql
```

3. Learn and use [pgx](https://github.com/jackc/pgx) to connect Go with PostgreSQL.
4. Write SQL queries for CRUD operations (e.g., under `db/query`):

```sql
-- name: CreateBook :one
INSERT INTO books (...) VALUES (...) RETURNING *;

-- name: GetBook :one
SELECT * FROM books WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (...) VALUES (...) RETURNING *;

-- name: ListBooks :many
SELECT * FROM books WHERE user_id = $1 ORDER BY published_date LIMIT $2 OFFSET $3;
```

5. Use [sqlc](https://github.com/sqlc-dev/sqlc) to generate type-safe Go code from SQL.

Example `sqlc.yaml`:
```yaml
version: "2"
sql:
- schema: "db/migration"
  queries: "db/query"
  engine: "postgresql"
  gen:
    go: 
      package: "db"
      out: "db/sqlc"
      sql_package: "pgx/v5"
      emit_json_tags: true
      emit_interface: true
      emit_empty_slices: true
      overrides:
        - db_type: "timestamptz"
          go_type: "time.Time"
        - db_type: "uuid"
          go_type: "github.com/google/uuid.UUID"
```

Run `sqlc generate` to generate code:

```go
func (q *Queries) CreateBook(ctx context.Context, arg CreateBookParams) (Book, error) {
    row := q.db.QueryRow(ctx, createBook, ...)
    var i Book
    err := row.Scan(...)
    return i, err
}
```

---

## Day 3

### Testing the Data Layer

1. Use [testify](https://github.com/stretchr/testify) to test your database logic:
```go
user, err := testStore.CreateUser(ctx, arg)
require.NoError(t, err)
require.Equal(t, arg.Email, user.Email)
```

2. Use [viper](https://github.com/spf13/viper) to load and manage app configurations:
```go
viper.AddConfigPath(path)
viper.SetConfigName("app")
viper.SetConfigType("env")
viper.AutomaticEnv()
viper.ReadInConfig()
viper.Unmarshal(&config)
```

---

## Day 4

### Routing and API Services

1. Use [gin](https://github.com/gin-gonic/gin) to build HTTP routes:
```go
router := gin.Default()
router.POST("/users", server.createUser)
router.POST("/users/login", server.loginUser)
router.Run(address)
```

2. Use [golang/mock](https://github.com/golang/mock) and `mockgen` to generate mocks for interfaces. Write tests for your APIs.
3. Use [paseto](https://github.com/o1egl/paseto) for secure token generation and validation.

**Symmetric encryption example**:
```go
token, err := paseto.Encrypt(symmetricKey, jsonToken, footer)
err := paseto.Decrypt(token, symmetricKey, &newJsonToken, &newFooter)
```

**Asymmetric example**:
```go
token, err := paseto.Sign(privateKey, jsonToken, footer)
err := paseto.Verify(token, publicKey, &newJsonToken, &newFooter)
```

4. Test your APIs using Postman, Insomnia, or `curl`.

---

## Day 5

### gRPC Simple Service

1. Follow [the gRPC Go quickstart guide](https://grpc.io/docs/languages/go/quickstart/):
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

2. Define your service under the `proto` directory.

3. Generate code:
```bash
protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
proto/*.proto
```

4. Implement and start your service:
- [`server.go`](https://github.com/logosmjt/bookstore-go/blob/main/gapi/server.go)
- [`createUser.go`](https://github.com/logosmjt/bookstore-go/blob/main/gapi/createUser.go)

5. Test with [evans](https://github.com/ktr0731/evans):
```
evans -r repl
> call CreateUser
```

---

### gRPC Gateway

[gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway) converts HTTP calls into gRPC requests via a reverse proxy.

1. Install gRPC-Gateway and its dependencies.
2. Clone [googleapis](https://github.com/googleapis/googleapis) and copy required `google/api` files into your `proto` folder.
3. Update `Makefile` with:
```bash
--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative
```

4. Annotate proto files with HTTP and OpenAPI options:
```proto
rpc CreateUser (CreateUserRequest) returns (CreateUserResponse) {
  option (google.api.http) = {
    post: "/v1/create_user"
    body: "*"
  };
  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
    description: "Use this API to create a new user";
    summary: "Create new user";
  };
}
```

5. Use [swagger-ui](https://github.com/swagger-api/swagger-ui) to serve and view the API docs.
6. Generate static files with [statik](https://github.com/rakyll/statik):
```bash
statik -src=./doc/swagger -dest=./doc
```

7. Visit `http://localhost:8080/swagger/` to view your interactive docs.

---

### Final Notes

gRPC APIs for login, token validation, and book-related features are not yet implemented. These are left as exercises for you to explore and practice further.