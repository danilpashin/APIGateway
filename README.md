# 🛠️API Gateway
## _High-load microservices gateway on Go_

## Roadmap 🗺️
- [x] Graceful shutdown
- [ ] Rate limiting (100 req/min per IP)
- [ ] JWT auth (access + refresh)
- [ ] Admin-only endpoints
- [ ] Prometheus metrics
- [ ] Unit tests (90% coverage)

## Tech Stack🚀
- **Go** 1.26 🐹
- **chi**⚡️
- **PostgreSQL** 🐘
- **Redis** 📡
- **JWT** 🔐
- **Prometheus** 📊

## Architecture🏗️
Client → Gateway → [User Service, Product Service] → [Redis, PostgreSQL]

## How to run
1. docker-compose up -d
2. go run services/product/cmd/main.go
3. go run services/user/cmd/main.go
4. go run services/gateway/cmd/main.go

## 🌐Endpoints (USER SERVICE)
### User
- POST /users/register
- POST /users/login
- PUT /users/profile/{id}
- GET /users/profile/{id}
- DELETE /users/profile/{id} (only logged user ID allowed)
### Admin
- POST /admin/login
- PUT /admin/users/profile/{id}
- GET /admin/users/profile/{id}
- GET /admin/users
- POST /admin/promote
### Metrics
- GET /metrics
## 🌐Endpoints (PRODUCT SERVICE)
- POST /products
- PUT /products/{id}
- GET /products/{id}
- GET /products
- DELETE /products/{id}
## 🌐Endpoints (GATEWAY)
### public
- POST /api/v1/auth/register
- POST /api/v1/auth/login
- GET /api/v1/health
- GET /api/v1/metrics
