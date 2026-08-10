# 🛠️API Gateway
## _High-load microservices gateway on Go_
![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)
![Coverage](https://img.shields.io/badge/Coverage-91.4%25-brightgreen)

## 🗺️ Roadmap

### 📈 Infrastructure and Architecture
- [x] Monorepository architecture
- [x] Layer isolation (Handler -> Service -> Repository)
- [ ] Gateway integration for routing
- [ ] Inner interaction through gRPC

### 🛡️ Security and Fault Tolerance
- [x] Separation of rights in the database (separate users for migrations and production)
- [x] Graceful shutdown
- [x] Middleware for panic recovery
- [ ] Limit of requests to server (Rate-limiting) with Redis
- [ ] JWT authentication (Access/Refresh)

### 📊 Testing and Monitoring
- [ ] Tests coverage (Unit + Integration)
    - [x] `Product Service` (100% coverage)
    - [x] `User Service` (100% coverage)
    - [ ] `Order Service` (In progress)
    - [ ] `Gateway` (In progress)
- [ ] Metrics in Prometheus and visualization in Grafana

## 🚀 Tech Stack
- **Go** 1.26 🐹
- **chi**⚡️
- **PostgreSQL** 🐘
- **Docker & Docker Compose** 🐳

## Will be used
- **Redis** 📡 
- **JWT** 🔐
- **Prometheus** 📊

## 🏗️ Current Architecture and Development Progress
Client (REST) → [ User Service ] → PostgreSQL (CRUD DB User)

Client (REST) → [ Product Service ] → PostgreSQL (CRUD DB User)

## 💻 How to Run

The application is **fully containerized**. You do not need to have Go or PostgreSQL installed locally on your machine.

### Prerequisites
* [Docker](https://docker.com) installed and running
* [Docker Compose](https://docker.com)

### Step-by-Step Setup

1. **Configure environment variables**  
   Create a `.env` file in the root directory by copying the example file:
   ```bash
   cp .env.example .env
   ```

2. **Start the application**  
   Build the images and start all services (including databases and migrations) in detached mode:
   ```bash
   docker compose up --build -d
   ```

3. **Verify the services**  
   Once the containers are up, the services will be available at:
   * **Product-Service:** [http://localhost:8082](http://localhost:8082)
   * **User-Service:** [http://localhost:8083](http://localhost:8083)

### Useful Commands
* To view application logs: `docker compose logs -f`
* To stop all services: `docker compose down`
* To completely reset databases: `docker compose down -v`


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
## 🌐Endpoints (ORDER SERVICE)
- [ ] POST /orders
- [ ] PUT /orders/{id}
- [ ] GET /orders/{id}
- [ ] GET /orders
- [ ] DELETE /orders/{id}
- [ ] GET /health 
### Admin
- [ ] POST /admin/login
- [ ] PUT /admin/users/profile/{id}
- [ ] GET /admin/users/profile/{id}
- [ ] GET /admin/users
- [ ] POST /admin/promote
### Metrics
- [ ] GET /metrics

## 🌐Endpoints (GATEWAY)
- [ ] POST /api/v1/auth/register
- [ ] POST /api/v1/auth/login
- [ ] GET /api/v1/health
- [ ] GET /api/v1/metrics
