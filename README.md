# 🛠️API Gateway
## _High-load microservices gateway on Go_

## Roadmap 🗺️
- [x] Graceful shutdown
- [ ] Rate limiting (100 req/min per IP)
- [ ] JWT auth (access + refresh)
- [ ] Admin-only endpoints
- [ ] Prometheus metrics
- [ ] Unit tests (90% coverage)
- [ ] Integration tests

## Tech Stack🚀
- **Go** 1.26 🐹
- **chi**⚡️
- **PostgreSQL** 🐘
- **Docker & Docker Compose** 🐳

## Will be used
- **Redis** 📡 
- **JWT** 🔐
- **Prometheus** 📊

## Architecture🏗️
Client → Gateway → [User Service, Product Service] → [Redis, PostgreSQL]

## How to run
1. go run services/product/cmd/main.go
2. go run services/user/cmd/main.go
3. go run services/gateway/cmd/main.go

## 🌐Endpoints (USER SERVICE)
### User
- [x] POST /users/register
- [ ] POST /users/login
- [x] PUT /users/profile/{id}
- [x] GET /users/profile/{id}
- [ ] DELETE /users/profile/{id} (only logged user ID allowed)
- [x] GET /health
## 🌐Endpoints (PRODUCT SERVICE)
- [x] POST /products
- [x] PUT /products/{id}
- [x] GET /products/{id}
- [x] GET /products
- [x] DELETE /products/{id}
- [x] GET /health
### Admin
- [ ] POST /admin/login
- [ ] PUT /admin/users/profile/{id}
- [ ] GET /admin/users/profile/{id}
- [ ] GET /admin/users
- [ ] POST /admin/promote
### Metrics
- [ ] GET /metrics

## 🌐Endpoints (GATEWAY)
### public
- [ ] POST /api/v1/auth/register
- [ ] POST /api/v1/auth/login
- [ ] GET /api/v1/health
- [ ] GET /api/v1/metrics
