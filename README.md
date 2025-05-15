# 🚀 RouteFlow: Declarative API Routing for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/go-deck/routeflow.svg)](https://pkg.go.dev/github.com/go-deck/routeflow)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-deck/routeflow)](https://goreportcard.com/report/github.com/go-deck/routeflow)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![View Documentation](https://img.shields.io/badge/View-Documentation-2ea44f?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/go-deck/routeflow)

```go
import "github.com/go-deck/routeflow"  // [Documentation](https://pkg.go.dev/github.com/go-deck/routeflow)
```

## 🔍 Overview
RouteFlow is a high-performance, declarative API routing library for Go that simplifies RESTful API development. Define your entire API structure, including routes, middlewares, and database configurations, in a clean YAML file for better maintainability and scalability.

## 🌟 Key Features
- **Declarative Configuration** - Define your entire API in YAML
- **Middleware Support** - Built-in and custom middleware chaining
- **Request Validation** - Schema-based request validation
- **Database Integration** - Seamless ORM support with GORM
- **RESTful Routes** - Intuitive route definitions
- **High Performance** - Optimized for speed and low latency
- **Extensible** - Easy to extend with custom functionality

## 🚀 Quick Start

### Prerequisites
- Go 1.16+
- Database (PostgreSQL/MySQL/SQLite3)

### Installation

```bash
go get -u github.com/go-deck/routeflow
```

### Basic Project Structure

```
myapp/
├── main.go
├── go.mod
├── go.sum
├── lib.yaml
└── handlers/
    └── user_handler.go
```


## 🛠 Middleware System

### Built-in Middlewares

| Middleware | Description |
|------------|-------------|
| `logging`  | Request/response logging |
| `cors`     | CORS headers management |
| `auth`     | Basic authentication |
| `recovery` | Panic recovery |

### Custom Middleware Example

```go
// Define your middleware methods
type AuthMiddleware struct{}


// JWT Authentication Middleware
func (m *AuthMiddleware) JWTValidation(c *ctx.Context) (interface{}, int) {
    token := c.GetHeader("Authorization")
    if token == "" {
        return map[string]string{"error": "Authorization header required"}, 401
    }
    
    // Add your JWT validation logic here
    if !isValidToken(token) {
        return map[string]string{"error": "Invalid or expired token"}, 401
    }
    
    return nil, 0 // Continue to next middleware/handler
}

// Rate Limiting Middleware
func (m *AuthMiddleware) RateLimit(c *ctx.Context) (interface{}, int) {
    ip := c.ClientIP()
    if isRateLimited(ip) {
        return map[string]string{"error": "Too many requests"}, 429
    }
    return nil, 0
}
```

### Registering Middleware

```yaml
middlewares:
  built_in: [logging, recovery]
  custom: [JWTValidation, RateLimit]

## 💾 Database Integration

RouteFlow provides seamless integration with multiple databases through GORM:

### Supported Databases

| Database | Connection String Example |
|----------|--------------------------|
| PostgreSQL | `postgres://user:pass@localhost:5432/dbname` |
| MySQL      | `user:pass@tcp(127.0.0.1:3306)/dbname` |
| SQLite3    | `file:test.db` |

### `lib.yaml` Example:

```yaml
database:
  type: postgres  # postgres, mysql, sqlite3
  host: localhost
  port: 5432
  username: dbuser
  password: dbpass
  database: myapp_db
  sslmode: disable
  max_idle_connections: 10
  max_open_connections: 100
  conn_max_lifetime: 1h
  migrate: true  # Auto-migrate models
```

## 🚀 Getting Started

### Example Project Structure

```
myapp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   └── user_handler.go
│   ├── middleware/
│   │   └── auth.go
│   └── model/
│       └── user.go
├── migrations/
│   └── 001_initial_schema.sql
├── pkg/
│   └── utils/
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── lib.yaml
```

### Example Handler

```go
// handlers/user_handler.go
package handlers

import (
	"net/http"
	"github.com/go-deck/routeflow/context"
)

type UserHandler struct{}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *context.Context) (interface{}, int) {
	id := c.Param("id")
	var user User
	
	if err := c.DB.First(&user, id).Error; err != nil {
		return map[string]string{"error": "User not found"}, http.StatusNotFound
	}
	
	return user, http.StatusOK
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *context.Context) (interface{}, int) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		return map[string]string{"error": err.Error()}, http.StatusBadRequest
	}

	user := User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashPassword(input.Password),
	}

	if err := c.DB.Create(&user).Error; err != nil {
		return map[string]string{"error": "Failed to create user"}, http.StatusInternalServerError
	}

	return user, http.StatusCreated
}
```

## 🔄 Development

### Running Tests

```bash
go test -v ./...
```

### Building

```bash
go build -o app cmd/server/main.go
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Useful Links

- [Documentation](https://github.com/go-deck/routeflow/wiki)
- [Examples](https://github.com/go-deck/routeflow/tree/main/examples)
- [Report Bug](https://github.com/go-deck/routeflow/issues)
- [Request Feature](https://github.com/go-deck/routeflow/issues/new?template=feature_request.md)

## 🌟 Show Your Support

Give a ⭐️ if this project helped you!

---

<div align="center">
  Made with ❤️ by RouteFlow Team
</div>

🚀 **Happy coding!** 🎯