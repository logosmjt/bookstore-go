## Golang 快速学习与实践

### 第一天

#### 熟悉Go编程的特性
- 静态类型和简洁性：Go 强调简洁性，同时通过静态类型保证代码的健壮性。
- 并发支持：Go 内置对并发的支持，通过 goroutines 和 channels 实现。
- 编译型语言：编译为机器码，执行速度快。
- 标准库：Go 提供的强大且全面的标准库。

#### 环境搭建

- 安装 Go：从 [golang.org](https://go.dev/) 下载并安装和配置。
- 安装 IDE：选择 Visual Studio Code，并安装格式 Go 插件以支持语法高亮、代码检查和调试等等。
- 编写 Hello World 程序：创建一个简单的 "Hello, World!" 程序以验证环境配置。
- optional：以TDD的方式编写德州扑克

### 第二天

#### 数据库模型的创建，连接与迁移。

1. 使用 [dbdiagram.io](https://dbdiagram.io/) 进行数据库设计，导出 PostgreSQL 的 SQL 代码。
2. 数据库迁移，学习并使用[golang-migrate/migrate](https://github.com/golang-migrate/migrate)实现数据库的迁移.
常用命令示例。
```
migrate create -ext sql -dir db/migration -seq $(name)
migrate -path "$(MIGRATION_PATH)" -database "$(DB_URL)" -verbose up
migrate -path "$(MIGRATION_PATH)" -database "$(DB_URL)" -verbose down
```

文件示例
```
000001_init_schema.up.sql
000001_init_schema.down.sql
000002_add_session.up.sql
000002_add_session.down.sql
```
3. 学习并使用PostgreSQL Go的工具[pgx](https://github.com/jackc/pgx)，完成与数据库的连接。
4. 编写CRUD操作的SQL查询，在sample中我将其添加在了db/query目录下。
```
-- name: CreateBook :one
INSERT INTO books (
  title,
  author,
  price,
  description,
  cover_image_url,
  published_date,
  user_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetBook :one
SELECT * FROM books
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  name,
  hashed_password,
  email,
  role
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: ListBooks :many
SELECT * FROM books
WHERE user_id = $1
ORDER BY published_date
LIMIT $2
OFFSET $3;
```
5. 学习并使用[sqlc](https://github.com/sqlc-dev/sqlc)来生成带有这些查询的类型安全接口的代码。
sqlc.yaml示例:
```
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
运行`sqlc generate`生成类型安全的查询函数。
示例:
```
const createBook = `-- name: CreateBook :one
INSERT INTO books (
  title,
  author,
  price,
  description,
  cover_image_url,
  published_date,
  user_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING id, title, author, price, description, cover_image_url, published_date, user_id, updated_at, created_at
`

type CreateBookParams struct {
	Title         string      `json:"title"`
	Author        string      `json:"author"`
	Price         int64       `json:"price"`
	Description   string      `json:"description"`
	CoverImageUrl string      `json:"cover_image_url"`
	PublishedDate time.Time   `json:"published_date"`
	UserID        pgtype.Int8 `json:"user_id"`
}

func (q *Queries) CreateBook(ctx context.Context, arg CreateBookParams) (Book, error) {
	row := q.db.QueryRow(ctx, createBook,
		arg.Title,
		arg.Author,
		arg.Price,
		arg.Description,
		arg.CoverImageUrl,
		arg.PublishedDate,
		arg.UserID,
	)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Author,
		&i.Price,
		&i.Description,
		&i.CoverImageUrl,
		&i.PublishedDate,
		&i.UserID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
```

### 第三天
#### 数据层编写测试
1. 学习并使用[testify](https://github.com/stretchr/testify)完成数据层的测试代码，验证数据层的代码满足期待。
```
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Name:           util.RandomUserName(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
		Role:           "seller",
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Role, user.Role)
	require.NotZero(t, user.CreatedAt)
```
2. 学习并使用[viper](https://github.com/spf13/viper),完成应用程序配置。
```
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
```

### 第四天
#### 路由与服务
1. 学习并使用[gin](https://github.com/gin-gonic/gin)创建HTTP路由。
```
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
  router.Run(address)
```
2. 学习[golang mock](https://github.com/golang/mock), 使用`mockgen`生成存储接口的mock。编写api的测试代码与生产代码。
3. 学习并使用[paseto](https://github.com/o1egl/paseto)，实现Security Token的创建与校验。
- Create token using symmetric key
```
symmetricKey := []byte("YELLOW SUBMARINE, BLACK WIZARDRY") // Must be 32 bytes
now := time.Now()
exp := now.Add(24 * time.Hour)
nbt := now

jsonToken := paseto.JSONToken{
        Audience:   "test",
        Issuer:     "test_service",
        Jti:        "123",
        Subject:    "test_subject",
        IssuedAt:   now,
        Expiration: exp,
        NotBefore:  nbt,
        }
// Add custom claim    to the token    
jsonToken.Set("data", "this is a signed message")
footer := "some footer"

// Encrypt data
token, err := paseto.Encrypt(symmetricKey, jsonToken, footer)
// token = "v2.local.E42A2iMY9SaZVzt-WkCi45_aebky4vbSUJsfG45OcanamwXwieieMjSjUkgsyZzlbYt82miN1xD-X0zEIhLK_RhWUPLZc9nC0shmkkkHS5Exj2zTpdNWhrC5KJRyUrI0cupc5qrctuREFLAvdCgwZBjh1QSgBX74V631fzl1IErGBgnt2LV1aij5W3hw9cXv4gtm_jSwsfee9HZcCE0sgUgAvklJCDO__8v_fTY7i_Regp5ZPa7h0X0m3yf0n4OXY9PRplunUpD9uEsXJ_MTF5gSFR3qE29eCHbJtRt0FFl81x-GCsQ9H9701TzEjGehCC6Bhw.c29tZSBmb290ZXI"

// Decrypt data
var newJsonToken paseto.JSONToken
var newFooter string
err := paseto.Decrypt(token, symmetricKey, &newJsonToken, &newFooter)
```
- Create token using asymetric key
```
b, _ := hex.DecodeString("b4cbfb43df4ce210727d953e4a713307fa19bb7d9f85041438d9e11b942a37741eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
privateKey := ed25519.PrivateKey(b)

b, _ = hex.DecodeString("1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
publicKey := ed25519.PublicKey(b)

// or create a new keypair 
// publicKey, privateKey, err := ed25519.GenerateKey(nil)

jsonToken := paseto.JSONToken{
        Expiration: time.Now().Add(24 * time.Hour),
        }
        
// Add custom claim    to the token    
jsonToken.Set("data", "this is a signed message")
footer := "some footer"

// Sign data
token, err := paseto.Sign(privateKey, jsonToken, footer)
// token = "v2.public.eyJkYXRhIjoidGhpcyBpcyBhIHNpZ25lZCBtZXNzYWdlIiwiZXhwIjoiMjAxOC0wMy0xMlQxOTowODo1NCswMTowMCJ9Ojv0uXlUNXSFhR88KXb568LheLRdeGy2oILR3uyOM_-b7r7i_fX8aljFYUiF-MRr5IRHMBcWPtM0fmn9SOd6Aw.c29tZSBmb290ZXI"

// Verify data
var newJsonToken paseto.JSONToken
var newFooter string
err := paseto.Verify(token, publicKey, &newJsonToken, &newFooter)
```
4. 使用postman，insomnia或者curl命令与gin service进行交互

### TBC grpc
