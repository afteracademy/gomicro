[![Framework](https://img.shields.io/badge/Framework-blue?label=View&logo=go)](https://github.com/afteracademy/goserve)
[![Microservices](https://img.shields.io/badge/Architecture-Microservices-orange?logo=go)](https://github.com/afteracademy/gomicro)

<div align="center">

# GoMicro - Go Microservices Architecture

### Production-Ready Microservices with Kong Gateway & NATS

![Banner](.docs/gomicro-banner.png)

**A complete microservices implementation using GoServe framework, Kong API Gateway, NATS messaging, PostgreSQL, MongoDB, and Redis.**

[![Documentation](https://img.shields.io/badge/📚_Read_Documentation-goserve.afteracademy.com-blue?style=for-the-badge)](http://goserve.afteracademy.com)

---
[![GoServe Framework](https://img.shields.io/badge/🚀_Framework-GoServe-blue?style=for-the-badge)](https://github.com/afteracademy/goserve)
[![API Documentation](https://img.shields.io/badge/📚_API_Docs-View_Here-blue?style=for-the-badge)](https://documenter.getpostman.com/view/1552895/2sA3dxCWsa)
[![Docker Guide](https://img.shields.io/badge/🐳_Docker_Guide-README--DOCKER.md-green?style=for-the-badge)](README-DOCKER.md)
---
</div>

## Overview

This project demonstrates a production-ready microservices architecture built with the [GoServe framework](https://github.com/afteracademy/goserve). It breaks down a monolithic blogging platform into independent services using Kong as API Gateway and NATS for inter-service communication. Each service maintains its own database and cache, showcasing true microservices best practices with service isolation, independent scaling, and fault tolerance.

The architecture implements authentication, authorization, and API key validation across distributed services while maintaining clean separation of concerns and independent deployability.

## Features

- **GoServe Micro Framework** - Built on production-ready [GoServe v2](https://github.com/afteracademy/goserve) with microservices support
- **Kong API Gateway** - Single entry point with custom Go plugin for API key validation
- **NATS Messaging** - Asynchronous inter-service communication with request/reply patterns
- **Service Isolation** - Each service with dedicated database and Redis instance
- **PostgreSQL & MongoDB** - Auth service with PostgreSQL, Blog service with MongoDB
- **JWT Authentication** - Token-based authentication with refresh token support
- **Role-Based Authorization** - Fine-grained access control across services
- **Custom Kong Plugin** - Go-based API key validation plugin
- **Docker Compose Ready** - Multiple configurations for development, testing, and production
- **Load Balancing** - Pre-configured setup for horizontal scaling
- **Health Checks** - Service health monitoring and dependency management
- **Auto Migrations** - Database schema migrations on startup
- **Development Tools** - pgAdmin, Mongo Express, Redis Commander included

## Technology Stack

| Component | Technology |
|-----------|------------|
| **Language** | Go 1.21+ |
| **Framework** | [GoServe v2](https://github.com/afteracademy/goserve) |
| **API Gateway** | Kong |
| **Message Broker** | NATS |
| **Web Framework** | [Gin](https://github.com/gin-gonic/gin) |
| **Auth Database** | PostgreSQL ([pgx](https://github.com/jackc/pgx)) |
| **Blog Database** | MongoDB ([mongo-driver](https://github.com/mongodb/mongo-go-driver)) |
| **Cache** | Redis ([go-redis](https://github.com/redis/go-redis)) |
| **Migrations** | [golang-migrate](https://github.com/golang-migrate/migrate) |
| **Authentication** | JWT tokens |
| **Validation** | [validator](https://github.com/go-playground/validator) |
| **Configuration** | Environment variables |

## Quick Start

### Prerequisites

- Docker & Docker Compose ([Installation Guide](https://docs.docker.com/install/))
- Go 1.21+ (for local development)

### Installation

**1. Clone the Repository**

```bash
git clone https://github.com/afteracademy/gomicro.git
cd gomicro
```

**2. Generate RSA Keys**
```bash
go run .tools/rsa/keygen.go
```

**3. Create Environment Files**
```bash
go run .tools/copy/envs.go 
```

**4. Start with Docker Compose**
```bash
# Full stack with all services
docker compose up --build
```

The API will be available at: **http://localhost:8000** (via Kong Gateway)

**5. Health Check**
```bash
# Check Kong Gateway
curl http://localhost:8001/status

# Check NATS
curl http://localhost:8222/varz
```

### Access Points

| Service | URL | Description |
|---------|-----|-------------|
| **Kong Gateway** | http://localhost:8000 | API Entry Point |
| **Kong Admin** | http://localhost:8001 | Kong Configuration |
| **NATS Monitoring** | http://localhost:8222 | NATS Dashboard |
| **PostgreSQL** | localhost:5432 | Auth Database |
| **MongoDB** | localhost:27017 | Blog Database |
| **Redis (Auth)** | localhost:6379 | Auth Cache |
| **Redis (Blog)** | localhost:6380 | Blog Cache |

**Development Tools** (with docker-compose.dev.yml):

| Tool | URL | Purpose |
|------|-----|---------|
| **pgAdmin** | http://localhost:5050 | PostgreSQL Management |
| **Mongo Express** | http://localhost:8082 | MongoDB Management |
| **Redis Commander** | http://localhost:8083 | Redis Management |

### Troubleshooting

If you encounter issues:
- Ensure ports 8000, 8001, 5432, 27017, 6379, 6380, 4222 are available
- Check service logs: `docker compose logs -f [service_name]`
- Clean slate: `docker compose down -v && docker compose up --build`

For detailed setup, usage, and troubleshooting: **[README-DOCKER.md](README-DOCKER.md)**

## Deployment Scenarios

### 1. Full Stack Development (Recommended)
```bash
docker compose up --build
```
Starts all services with Kong, NATS, and shared databases.

Adds pgAdmin, Mongo Express, and Redis Commander for database management.

### 2. Individual Service Development
```bash
# Auth service only
cd auth_service && docker compose up --build

# Blog service only  
cd blog_service && docker compose up --build
```
Runs a single service in isolation for fast iteration and debugging.

### 3. Load Balanced Production
```bash
docker compose -f docker-compose-load-balanced.yml up --build
```
Runs 2 instances of each service behind Kong for production-like setup.

## Architecture

### Microservices Design Principles

This project follows microservices best practices:

- **Service Isolation** - Each service has its own database and cache
- **Independent Deployment** - Services can be deployed independently
- **API Gateway Pattern** - Single entry point via Kong
- **Event-Driven Communication** - NATS for async messaging
- **Database per Service** - No shared databases
- **Distributed Authentication** - Auth service validates tokens via NATS
- **Health Checks** - Service health monitoring and dependency management

### System Architecture

**1. Without Load Balancing**
![System Architecture](.docs/system.png)

**2. With Load Balancing**
![Load Balanced Architecture](.docs/system-load-balanced.png)

### Request Flow

```
Client → Kong Gateway → API Key Validation → Service → NATS → Response
```

1. **Client Request** → Kong Gateway (port 8000)
2. **API Key Validation** → Custom Kong plugin calls `auth:8080/verify/apikey`
3. **Route to Service** → Kong forwards to auth or blog service
4. **Service Processing** → Business logic execution
5. **NATS Communication** → Inter-service messaging (if needed)
6. **Response** → Kong → Client

### Authentication Flow

- **Users & Credentials** → Stored in auth_service PostgreSQL database
- **JWT Token Generation** → Auth service creates access/refresh tokens
- **Token Validation** → Auth service middleware validates JWT
- **Cross-Service Auth** → Blog service requests token validation via NATS
- **Distributed Security** → Each service can enforce its own authentication

### Authorization Flow

- **Roles & Permissions** → Stored in auth_service PostgreSQL database
- **Role Assignment** → Users can have multiple roles
- **Role Validation** → Auth service middleware checks permissions
- **Cross-Service Authorization** → Blog service requests role validation via NATS
- **Fine-Grained Control** → Each service decides which endpoints require which roles

> **Design Philosophy**: This distributed authentication/authorization gives each service autonomy to define public, protected, and restricted APIs independently while maintaining centralized user management.

### Service Communication

**Synchronous (HTTP)**
- Client ↔ Kong Gateway
- Kong ↔ Services (routing)
- Kong Plugin ↔ Auth Service (API key validation)

**Asynchronous (NATS)**
- Blog Service → Auth Service (token validation)
- Blog Service → Auth Service (role verification)
- Event-driven messaging between services

### Network Architecture

- **Custom Bridge Network** (`gomicro-network`) for service discovery
- **Container Names** as DNS (postgres, mongo, redis-auth, redis-blog, nats)
- **Internal Communication** via container names (no external IPs)
- **External Access** only via Kong Gateway
- **Database Access** exposed for development (can be restricted in production)

## Project Structure

```
gomicro/
├── auth_service/              # Authentication & Authorization Service
│   ├── api/
│   │   ├── auth/             # Auth endpoints (signup, signin, refresh)
│   │   │   ├── dto/          # Request/response DTOs
│   │   │   ├── message/      # NATS message definitions
│   │   │   ├── middleware/   # Auth & authorization middleware
│   │   │   └── model/        # PostgreSQL models
│   │   └── user/             # User management endpoints
│   ├── cmd/main.go           # Service entry point
│   ├── migrations/           # PostgreSQL migrations
│   ├── startup/              # Server initialization
│   └── docker-compose.yml    # Standalone development
│
├── blog_service/              # Blog Management Service
│   ├── api/
│   │   ├── author/           # Author-specific blog operations
│   │   ├── blog/             # Blog CRUD operations
│   │   ├── blogs/            # Blog listing & search
│   │   └── editor/           # Editorial operations
│   ├── cmd/main.go           # Service entry point
│   ├── startup/              # Server initialization & indexes
│   └── docker-compose.yml    # Standalone development
│
├── kong/                      # API Gateway Configuration
│   ├── kong.yml              # Kong declarative config
│   ├── kong-load-balanced.yml # Load balanced config
│   └── apikey_auth_plugin/   # Custom Go plugin for API key validation
│
├── docker-compose.yml         # Full stack orchestration
├── docker-compose.dev.yml     # Development tools
├── docker-compose-load-balanced.yml # Production setup
└── README-DOCKER.md          # Detailed Docker documentation
```

### Service Directories

| Directory | Purpose |
|-----------|---------|
| **api/** | Feature-based API implementations |
| **cmd/** | Application entry point (main.go) |
| **common/** | Shared code across APIs |
| **config/** | Environment configuration |
| **keys/** | RSA keys for JWT signing |
| **migrations/** | Database migration files (auth service) |
| **startup/** | Server, DB, Redis, NATS initialization |
| **utils/** | Utility functions |

### Helper Directories

| Directory | Purpose |
|-----------|---------|
| **.extra/** | Database init scripts, documentation, assets |
| **.tools/** | RSA key generator, env file copier |
| **.vscode/** | Editor configuration, debug settings |
| **.docs/** | Architecture diagrams, banners |

## Code Examples

### NATS Message Definition

Create message types for inter-service communication:

```go
package message

type TokenValidation struct {
  AccessToken string `json:"accessToken,omitempty"`
}

func NewTokenValidation(token string) *TokenValidation {
  return &TokenValidation{
    AccessToken: token,
  }
}
```

### Controller with NATS Endpoints

Controllers implement `micro.Controller` to handle both HTTP and NATS requests:

```go
package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/afteracademy/gomicro/auth_service/api/auth/message"
  "github.com/afteracademy/goserve/v2/micro"
  "github.com/afteracademy/goserve/v2/network"
)

type controller struct {
  micro.Controller
  service Service
}

func NewController(
  authMFunc network.AuthenticationProvider,
  authorizeMFunc network.AuthorizationProvider,
  service Service,
) micro.Controller {
  return &controller{
    Controller: micro.NewController("/auth", authMFunc, authorizeMFunc),
    service:    service,
  }
}

// MountNats - Endpoints other services can call via NATS
func (c *controller) MountNats(group micro.NatsGroup) {
  group.AddEndpoint("validate.token", micro.NatsHandlerFunc(c.validateTokenHandler))
  group.AddEndpoint("validate.role", micro.NatsHandlerFunc(c.validateRoleHandler))
}

// MountRoutes - HTTP endpoints for clients
func (c *controller) MountRoutes(group *gin.RouterGroup) {
  group.POST("/signup/basic", c.signupBasicHandler)
  group.POST("/signin/basic", c.signinBasicHandler)
  group.POST("/token/refresh", c.Authentication(), c.tokenRefreshHandler)
  group.DELETE("/logout", c.Authentication(), c.logoutHandler)
}

func (c *controller) validateTokenHandler(req micro.NatsRequest) {
  var msg message.TokenValidation
  if err := req.DecodeData(&msg); err != nil {
    micro.SendNatsErrorMessage(err)
    return
  }
  
  user, err := c.service.ValidateToken(msg.AccessToken)
  if err != nil {
    micro.SendNatsErrorMessage(err)
    return
  }
  
  micro.SendNatsMessage(user)
}

// HTTP handlers...
func (c *controller) signupBasicHandler(ctx *gin.Context) {
  // Implementation...
}
```

**Key Components:**
- `MountNats()` - Defines NATS endpoints for inter-service calls
- `MountRoutes()` - Defines HTTP endpoints for client requests
- `micro.Controller` - Interface for microservice controllers

### Service with NATS Communication

Services use `micro.RequestNats` to call other services via NATS messaging:

```go
package blog

import (
  authmsg "github.com/afteracademy/gomicro/auth_service/api/auth/message"
  "github.com/afteracademy/goserve/v2/micro"
  "github.com/afteracademy/goserve/v2/mongo"
  "github.com/afteracademy/goserve/v2/redis"
)

const NATS_AUTH_VALIDATE_TOKEN = "auth.validate.token"
const NATS_AUTH_VALIDATE_ROLE = "auth.validate.role"

type Service interface {
  ValidateUserToken(token string) (*authmsg.User, error)
  ValidateUserRole(userId, roleCode string) (bool, error)
  // ... other blog operations
}

type service struct {
  natsClient       micro.NatsClient
  blogQueryBuilder mongo.QueryBuilder[model.Blog]
  blogCache        redis.Cache[dto.BlogInfo]
}

func NewService(
  db mongo.Database, 
  store redis.Store, 
  natsClient micro.NatsClient,
) Service {
  return &service{
    natsClient:       natsClient,
    blogQueryBuilder: mongo.NewQueryBuilder[model.Blog](db, model.CollectionName),
    blogCache:        redis.NewCache[dto.BlogInfo](store),
  }
}

// Call auth service via NATS to validate token
func (s *service) ValidateUserToken(token string) (*authmsg.User, error) {
  request := authmsg.NewTokenValidation(token)
  
  // Send request to auth service via NATS and wait for response
  user, err := micro.RequestNats[authmsg.TokenValidation, authmsg.User](
    s.natsClient, 
    NATS_AUTH_VALIDATE_TOKEN, 
    request,
  )
  
  if err != nil {
    return nil, err
  }
  
  return user, nil
}

// Call auth service via NATS to validate role
func (s *service) ValidateUserRole(userId, roleCode string) (bool, error) {
  request := authmsg.NewRoleValidation(userId, roleCode)
  
  result, err := micro.RequestNats[authmsg.RoleValidation, authmsg.RoleResult](
    s.natsClient,
    NATS_AUTH_VALIDATE_ROLE,
    request,
  )
  
  if err != nil {
    return false, err
  }
  
  return result.HasRole, nil
}
```

**Key Features:**
- **Type-Safe NATS Calls** - Generic `RequestNats[Request, Response]`
- **Async Communication** - Non-blocking inter-service messaging
- **Error Handling** - Proper error propagation across services
- **Distributed Auth** - Services don't need auth logic, just call auth service

### NATS Client Setup

Initialize NATS client in service startup:

```go
package startup

import (
  "time"
  "github.com/afteracademy/goserve/v2/micro"
)

func SetupNats(env *config.Env) micro.NatsClient {
  natsConfig := micro.Config{
    NatsUrl:            env.NatsUrl,            // "nats://nats:4222"
    NatsServiceName:    env.NatsServiceName,    // "auth" or "blog"
    NatsServiceVersion: env.NatsServiceVersion, // "1.0.0"
    Timeout:            time.Second * 10,
  }

  return micro.NewNatsClient(&natsConfig)
}
```

> **NATS Documentation**: GoServe wraps [nats-io/nats.go](https://github.com/nats-io/nats.go/blob/main/micro/README.md) for simplified microservice patterns.

## Migration Guide

### From GoServe Monolith to Microservices

If you're coming from the [GoServe](https://github.com/afteracademy/goserve) monolithic framework:

| Monolithic | Microservices | Change |
|-----------|---------------|---------|
| `network.Module[T]` | `micro.Module[T]` | Module initialization |
| `network.NewRouter()` | `micro.NewRouter()` | Router creation |
| `network.BaseController` | `micro.BaseController` | Base controller interface |
| `network.Controller` | `micro.Controller` | Controller interface |
| N/A | `MountNats(group)` | NATS endpoint registration |

**New Capabilities:**
- `micro.NatsClient` - NATS client for inter-service communication
- `micro.RequestNats[Req, Res]()` - Type-safe NATS request/response
- `micro.NatsHandlerFunc` - NATS message handlers
- `MountNats()` - Register NATS endpoints alongside HTTP routes

## API Documentation

<div align="center">

[![API Documentation](https://img.shields.io/badge/📚_View_Full_API_Documentation-blue?style=for-the-badge)](https://documenter.getpostman.com/view/1552895/2sA3dxCWsa)

Complete API documentation with authentication flows, request/response examples, and microservices patterns

</div>

## Related Projects

Explore the GoServe ecosystem:

1. **[GoServe Framework](https://github.com/afteracademy/goserve)**  
   Core framework with PostgreSQL, MongoDB, Redis, and NATS support

2. **[PostgreSQL API Server](https://github.com/afteracademy/goserve-example-api-server-postgres)**  
   Monolithic REST API with PostgreSQL and clean architecture

3. **[MongoDB API Server](https://github.com/afteracademy/goserve-example-api-server-mongo)**  
   Complete REST API with MongoDB implementation

4. **[GoServeGen CLI](https://github.com/afteracademy/goservegen)**  
   Code generator for scaffolding new projects and APIs

## Articles & Tutorials

Learn the concepts behind this project:

- [How to Create Microservices — A Practical Guide Using Go](https://afteracademy.com/article/how-to-create-microservices-a-practical-guide-using-go)
- [How to Architect Good Go Backend REST API Services](https://afteracademy.com/article/how-to-architect-good-go-backend-rest-api-services)
- [Implement JSON Web Token (JWT) Authentication using AccessToken and RefreshToken](https://afteracademy.com/article/implement-json-web-token-jwt-authentication-using-access-token-and-refresh-token)

## Contributing

We welcome contributions! Please feel free to:

- **Fork** the repository
- **Open** issues for bugs or feature requests
- **Submit** pull requests with improvements
- **Share** your feedback and suggestions

## Learn More

Subscribe to **AfterAcademy** on YouTube for in-depth tutorials and microservices concepts:

[![YouTube](https://img.shields.io/badge/YouTube-Subscribe-red?style=for-the-badge&logo=youtube&logoColor=white)](https://www.youtube.com/@afteracad)

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support This Project

If you find this project useful, please consider:

- **Starring** ⭐ this repository
- **Sharing** with the Go community
- **Contributing** improvements
- **Reporting** bugs and issues
- **Writing** articles about your experience

---

<div align="center">

**Built with ❤️ by [AfterAcademy](https://github.com/afteracademy)**

[GoServe Framework](https://github.com/afteracademy/goserve) • [API Documentation](https://documenter.getpostman.com/view/1552895/2sA3dxCWsa) • [Docker Guide](README-DOCKER.md) • [Articles](https://afteracademy.com) • [YouTube](https://www.youtube.com/@afteracad)

</div>
