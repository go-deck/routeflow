# 🚀 RouteFlow - Declarative API Routing for Go

## Overview
RouteFlow simplifies API development in Go by allowing you to define routes, middleware, and server configurations in a **YAML file**. It provides a **framework-agnostic** approach, currently supporting **Gin**, with plans for **Fiber, Echo, and Chi**.

## Features
- 📝 **Declarative Routing**: Define routes easily in YAML.
- 🌍 **Multi-Framework Support**: Supports **Gin** (more coming soon).
- 🔌 **Middleware Configuration**: Logging, CORS, and security settings.
- 🔄 **Dynamic Handler Resolution**: No need to manually map functions.
- 🔧 **Flexible Context Management**: Extracts path, query, and body data.

## Installation
```sh
go get github.com/go-deck/routeflow
```

## Usage
### 1️⃣ Define Your `lib.yaml`
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

database:
  type: sqlite3            # postgres, mysql, sqlite3
  host:                    # localhost
  port:                    # 5432
  username:                # root
  password:                # root
  database: ./userdb       # file path for sqlite3 or database name for postgres and mysql
  sslmode:                 # disable
  max_idle_connections: 10 
  max_open_connections: 100 
  conn_max_lifetime: 1h    # 1 hour
  options: {}              # additional options for database connection

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
              props:
                min_length: 2
                max_length: 10
                required: true
                format: username
            - name: email
              type: string
              props:
                min_length: 5
                max_length: 12
                required: true
                format: email
            - name: phonenumber
              type: string
              props:
                min_length: 10
                max_length: 12
                required: true
                format: phone
```

### 2️⃣ Create Your Go Application
```go
package main

import (
    "log"
    "github.com/go-deck/routeflow/module/sample"
    routeflow "github.com/go-deck/routeflow/routeflow"
    ctx "github.com/go-deck/routeflow/routeflow/frameworks/ginframework/ctx"
)

func main() {

    handlerMap := map[string]func(*ctx.Context) (interface{}, int){
        "getUserData":     sample.ListUsers,
        "getUserDataById": sample.GetUserDataById,
        "createUser":      sample.CreateUser,
    }

    log.Println("Starting API Server with declarative routing...")
    routeflow.StartServer("lib.yaml", handlerMap)
}
```

### 3️⃣ Run Your Application
```sh
go run main.go
```

## Future Enhancements
- ✅ Support for **Fiber, Echo, Chi**
- 🔐 Authentication & Role-Based Access Control
- 📈 Rate Limiting & API Analytics
- 🔄 Hot Reload for Config Updates

## Contributing
Contributions are welcome! Feel free to open an issue on [GitHub](https://github.com/go-deck/routeflow).

## License
This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

🚀 Happy coding! 🎯