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

### 第五天 

#### gRPC simple service
1. 按照[指南](https://grpc.io/docs/languages/go/quickstart/)开始Go中的gRPC旅程
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
2. 更新gRPC service，参考文档或[proto](https://github.com/logosmjt/bookstore-go/tree/main/proto)下的文件

3. 生成代码参考[文档](https://grpc.io/docs/languages/go/generated-code/)或Makefile中proto
```
rm -f pb/*.go
protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
proto/*.proto
```
4. 更新service，启动服务
[server.go](https://github.com/logosmjt/bookstore-go/blob/main/gapi/server.go)
[createUser.go](https://github.com/logosmjt/bookstore-go/blob/main/gapi/createUser.go)
[runGrpcServer](https://github.com/logosmjt/bookstore-go/blob/main/main.go#L85C6-L85C19)

5. 安装[evans](https://github.com/ktr0731/evans) 完成测试

```
localhost:9090> package pb
pb@localhost:9090> show service
pb@localhost:9090> service BookStore

pb.BookStore@localhost:9090> call CreateUser
name (TYPE_STRING) => rpctest1
password (TYPE_STRING) => 123456
email (TYPE_STRING) => rpctest1@bookstore.com
role (TYPE_STRING) => 
{
  "user": {
    "createdAt": "2025-02-07T00:55:04.167042Z",
    "email": "rpctest1@bookstore.com",
    "name": "rpctest1",
    "role": "buyer",
    "updatedAt": "2025-02-07T00:55:04.167042Z"
  }
}
```

#### gRPC Gateway
[gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway)是Google协议缓冲器编译器ProtoC的插件。它读取Protobuf服务定义并生成反向代理服务器，该服务器将静止的HTTP API转换为GRPC。该服务器是根据您的服务定义中的Google.api.http注释生成的。
1. [安装](https://github.com/grpc-ecosystem/grpc-gateway?tab=readme-ov-file#installation)gRPC-Gateway。
2. clone [googleapis](https://github.com/googleapis/googleapis)，将[google/api](https://github.com/googleapis/googleapis/tree/master/google/api)下需要的文件复制到[google/api](https://github.com/logosmjt/bookstore-go/tree/main/proto/google/api)。
3. 添加`--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \`到`Makefile`中的`proto`, 创建`runGatewayServer`。
4. clone [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway)，复制[options](https://github.com/grpc-ecosystem/grpc-gateway/tree/main/protoc-gen-openapiv2/options)下需要的文件到[protoc-gen-openapiv2/options](https://github.com/logosmjt/bookstore-go/tree/main/proto/protoc-gen-openapiv2/options)。更新[service_book_store.proto](https://github.com/logosmjt/bookstore-go/blob/main/proto/service_book_store.proto)，添加options。
```
rpc CreateUser (CreateUserRequest) returns (CreateUserResponse){
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
5. clone [swagger-ui](https://github.com/swagger-api/swagger-ui)，将[dist](https://github.com/swagger-api/swagger-ui/tree/master/dist)下的内容复制到[doc/swagger](https://github.com/logosmjt/bookstore-go/tree/main/doc/swagger)，修改swagger-initializer中的url，添加`--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bookstore \`到`Makefile`中的`proto`。
6. 安装[statik](https://github.com/rakyll/statik),并添加`statik -src=./doc/swagger -dest=./doc`到`Makefile`中的`proto`。
7. 打开`http://localhost:8080/swagger/`进行验证。
![swagger-ui](/doc/img/screenshot.jpg)

### 最后
gRPC中还没有完成用户登录，token验证和书籍相关的API，大家可尝试练习。