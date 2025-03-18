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
    "log"

    "github.com/go-deck/routeflow/module/sample"
    routeflow "github.com/go-deck/routeflow/routeflow"
    "github.com/go-deck/routeflow/routeflow/ctx"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    handlerMap := map[string]func(*ctx.Context) (interface{}, int){
        "getUserData":     sample.ListUsers,
        "getUserDataById": sample.GetUserDataById,
        "createUser":      sample.CreateUser,
    }

    app, _ := routeflow.New("lib.yaml")

    log.Println("Starting API Server with declarative routing...")
    
    app.InitDB()
    app.DB.AutoMigrate(&sample.User{})  // Auto-migrate if enabled
    app.Serve(handlerMap)
}
```

---

### 3️⃣ Define Your Handlers

```go
package sample

import (
    "net/http"
    "github.com/go-deck/routeflow/routeflow/ctx"
)

type User struct {
    ID          int    `json:"id"`
    Username    string `json:"username"`
    PhoneNumber string `json:"phone_number"`
    Email       string `json:"email"`
}

// Get all users
func ListUsers(c *ctx.Context) (interface{}, int) {
    var users []User
    if err := c.DB.Find(&users).Error; err != nil {
        return map[string]string{"error": "Database error"}, http.StatusInternalServerError
    }
    return users, http.StatusOK
}

// Get user by ID
func GetUserDataById(c *ctx.Context) (interface{}, int) {
    id, exists := c.PathParams["id"]
    if !exists {
        return map[string]string{"error": "Invalid user ID"}, http.StatusBadRequest
    }

    var user User
    if err := c.DB.First(&user, id).Error; err != nil {
        return map[string]string{"error": "User not found"}, http.StatusNotFound
    }
    return user, http.StatusOK
}

// Create a new user
func CreateUser(c *ctx.Context) (interface{}, int) {
    user := User{
        Username:    c.BodyData["username"].(string),
        PhoneNumber: c.BodyData["phonenumber"].(string),
        Email:       c.BodyData["email"].(string),
    }
    
    if err := c.DB.Create(&user).Error; err != nil {
        return map[string]string{"error": "Failed to create user"}, http.StatusInternalServerError
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

