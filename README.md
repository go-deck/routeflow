# Declarative Routing Library for Go

## Overview
The **Declarative Routing Library** simplifies API development in Go by allowing you to define routes, middleware, and server configurations in a **YAML file**. This eliminates the need for complex route handling and provides a **framework-agnostic** approach that currently supports **Gin**, with plans for **Fiber, Echo, Chi**, and more.

## Features
- 🚀 **Declarative Routing**: Define routes in YAML—no need to write route logic manually.
- 🌍 **Multi-Framework Support**: Supports **Gin** (more frameworks coming soon).
- 🛠 **Middleware Configuration**: Enable logging, CORS, cookies, and security settings easily.
- 🔒 **Private & Public Routes**: Manage authentication and role-based access control.
- 📄 **Easy Extensibility**: Modular structure makes it simple to add new frameworks.
- 🔄 **Auto-Generated Handlers**: Automatically maps handlers to routes, reducing boilerplate code.

---

## Installation
To install the library, run:
```sh
 go get github.com/go-deck/routeflow
```

---

## Usage
### 1️⃣ Define Your `lib.yaml`
Create a YAML file to configure routes, middleware, and server settings:
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
            - name: phonenumber
              type: string
```

---

### 2️⃣ Create Your Go Application
```go
package main

import (
    "log"
    "net/http"

    "github.com/go-deck/routeflow"
)

// Define API Handlers
func listUsers(pathParams map[string]string, queryParams map[string]string, bodyData map[string]interface{}) (interface{}, int) {
    return map[string]string{"message": "List of users"}, http.StatusOK
}

func getUserDataById(pathParams map[string]string, queryParams map[string]string, bodyData map[string]interface{}) (interface{}, int) {
    id, exists := pathParams["id"]
    if !exists {
        return map[string]string{"error": "Invalid user ID"}, http.StatusBadRequest
    }
    return map[string]string{"id": id, "username": "test_user"}, http.StatusOK
}

func createUser(pathParams map[string]string, queryParams map[string]string, bodyData map[string]interface{}) (interface{}, int) {
    username, _ := bodyData["username"].(string)
    phonenumber, _ := bodyData["phonenumber"].(string)
    return map[string]string{"message": "User created", "username": username, "phonenumber": phonenumber}, http.StatusCreated
}

func main() {
    handlerMap := map[string]func(map[string]string, map[string]string, map[string]interface{}) (interface{}, int){
        "getUserData":     listUsers,
        "getUserDataById": getUserDataById,
        "createUser":      createUser,
    }

    log.Println("Starting API Server with declarative routing...")
    declarative.StartServer("lib.yaml", handlerMap)
}
```

---

### 3️⃣ Run Your Application
```sh
go run main.go
```
Your API server will be up and running! 🎉


---

## Future Enhancements
- ✅ Support for **Fiber, Echo, Chi**
- 🔐 Authentication & Role-Based Access Control
- 📈 Rate Limiting & API Analytics
- 🔄 Hot Reload for Config Updates

## Contributing
Contributions are welcome! Feel free to open an issue or PR on [GitHub](https://github.com/go-deck/routeflow).

## License
This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

🚀 Happy coding! 🎯

