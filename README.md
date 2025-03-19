# 🚀 RouteFlow: Declarative API Routing for Go

## Overview
RouteFlow is a **declarative API routing library** for Go that simplifies API development by allowing users to define routes, middleware, and database configurations using a **YAML file**. It supports multiple frameworks, starting with **Gin**, and provides built-in database handling with **PostgreSQL, MySQL, and SQLite3**.

## 🌟 Features
- **Declarative Routing**: Define routes in YAML, eliminating manual route handling.
- **Multi-Framework Support**: Currently supports **Gin**, with plans for Fiber, Echo, and Chi.
- **Middleware Configuration**: Easily enable logging, CORS, security, and more via YAML.
- **Automatic Validation**: Supports field validation, including **required fields, email, phone number, min/max length**, etc.
- **Database Integration**: Built-in support for PostgreSQL, MySQL, and SQLite3.
- **Dynamic Context**: Auto-injects request parameters and database connection into handlers.

---

## 🛠 Installation
Install RouteFlow using:
```sh
go get github.com/go-deck/routeflow
```

---

## 🚀 Usage

### 1️⃣ Define Your `lib.yaml`
Create a YAML file to configure the server, middleware, routes, and database settings.
```yaml
server:
  port: 8080
  timeout: 30s
  allow_cors: true
  allowed_origins: ["*"]
  cookie:
    secure: true
    http_only: true
    same_site: Strict

framework: gin

database:
  type: sqlite3            # postgres, mysql, sqlite3
  host:                    # localhost
  port:                    # 5432
  username:                # root
  password:                # root
  database: ./userdb       # File path for SQLite3 or DB name for Postgres/MySQL
  sslmode:                 # disable
  max_idle_connections: 10
  max_open_connections: 100
  conn_max_lifetime: 1h    # 1 hour
  migrate: true            # Auto-migrate database

middlewares:
  global: [logging]

routes:
  groups:
    - base: /api/v1
      routes:
        - path: /get-user
          handler: getUserData
          method: GET
        - path: /get-user/:id
          handler: getUserDataById
          method: GET
        - path: /userpost
          handler: createUser
          method: POST
          body_params:
            - name: username
              type: string
              validation:
                min_length: 2
                max_length: 10
                required: true
                pattern: username
            - name: email
              type: string
              validation:
                min_length: 5
                max_length: 12
                required: true
                pattern: email
            - name: phonenumber
              type: string
              validation:
                min_length: 10
                max_length: 12
                required: true
                pattern: phone
```

---

### 2️⃣ Create Your Go Application

```go
package main

import (
    "github.com/apps/sample"
    ctx "github.com/go-deck/routeflow"
    log "github.com/sirupsen/logrus"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    app, err := ctx.New("lib.yaml", &sample.SampleHandler{})

    if err != nil {
        log.Error(err)
    }

    err = app.InitDB()

    if err != nil {
        log.Error(err)
    }

    err = app.DB.AutoMigrate(&sample.User{})

    if err != nil {
        log.Error(err)
    }

    log.Println("Starting API Server with declarative routing...")

    app.Serve()
}
```

---

### 3️⃣ Define Your Handlers

```go
package sample

import (
    "net/http"
    ctx "github.com/go-deck/routeflow"
    "gorm.io/gorm"
)

// User represents a user model
type User struct {
    ID          int    `json:"id"`
    Username    string `json:"username"`
    PhoneNumber string `json:"phonenumber"`
    Email       string `json:"email"`
}

type SampleHandler struct{}

// Get all users
func (b *SampleHandler) ListUsers(c *ctx.Context) (interface{}, int) {
    var users []User

    c.DB.AutoMigrate(&User{})

    // Use GORM to fetch all users
    if err := c.DB.Find(&users).Error; err != nil {
        return map[string]string{"error": "Database error"}, http.StatusInternalServerError
    }

    return users, http.StatusOK
}

// Get user data by ID
func (b *SampleHandler) GetUserDataById(c *ctx.Context) (interface{}, int) {
    c.DB.AutoMigrate(&User{})
    id, exists := c.PathParams["id"]
    if !exists {
        return map[string]string{"error": "Invalid user ID"}, http.StatusBadRequest
    }

    var user User
    // Use GORM to fetch user by ID
    if err := c.DB.Where("id = ?", id).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return map[string]string{"error": "User not found"}, http.StatusNotFound
        }
        return map[string]string{"error": "Database error"}, http.StatusInternalServerError
    }

    return user, http.StatusOK
}

// Create a new user
func (b *SampleHandler) CreateUser(c *ctx.Context) (interface{}, int) {
    c.DB.AutoMigrate(&User{})
    // Extract parameters from bodyData
    username, ok1 := c.BodyData["username"].(string)
    phonenumber, ok2 := c.BodyData["phonenumber"].(string)
    email, ok3 := c.BodyData["email"].(string)

    if !ok1 || !ok2 || !ok3 {
        return map[string]string{"error": "Invalid input"}, http.StatusBadRequest
    }

    // Insert into database using GORM
    user := User{
        Username:    username,
        PhoneNumber: phonenumber,
        Email:       email,
    }

    if err := c.DB.Create(&user).Error; err != nil {
        return map[string]string{"error": err.Error()}, http.StatusInternalServerError
    }

    return map[string]string{"message": "User created successfully"}, http.StatusCreated
}
```

---

### 4️⃣ Run Your Application

```sh
go run main.go
```

Your API server is now running! 🎉

---

## 🔥 Future Enhancements
- ✅ Support for **Fiber, Echo, Chi**
- 🔐 Authentication & Role-Based Access Control
- 📈 Rate Limiting & API Analytics
- 🔄 Hot Reload for Config Updates

## 🤝 Contributing
Contributions are welcome! Feel free to open an issue or PR on [GitHub](https://github.com/go-deck/routeflow).

## 📜 License
This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

🚀 **Happy coding!** 🎯