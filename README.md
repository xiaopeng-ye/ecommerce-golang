# Go E-commerce Backend Service

[![12-Factor App](https://img.shields.io/badge/12--Factor-compliant-brightgreen)](https://12factor.net/)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker)](https://www.docker.com/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-ready-326CE5?logo=kubernetes)](https://kubernetes.io/)

This is a production-ready e-commerce backend service developed using Go, implementing [12-Factor App](https://12factor.net/) best practices. It provides a complete set of core e-commerce APIs, including user authentication, product management, shopping cart, and order system.

## ✨ Features

### Business Features
- **User System**
  - User registration and authentication
  - JWT-based authentication
  - User profile management
- **Product System**
  - Product catalog and search
  - Product details and inventory management
  - Product CRUD operations
- **Shopping Cart System**
  - Add/remove items from cart
  - Cart checkout process
- **Order System**
  - Order creation and tracking
  - Order management and history

### Technical Features
- **12-Factor App Compliant**: Following cloud-native best practices
- **Structured Logging**: JSON-formatted logs for easy parsing and monitoring
- **Health Checks**: `/health` and `/ready` endpoints for orchestration
- **Graceful Shutdown**: Proper signal handling and resource cleanup
- **Horizontal Scaling**: Stateless design, ready for container orchestration
- **Configuration Management**: Environment-based configuration
- **Security**: Non-root container user, minimal attack surface
- **Observability**: Structured logs with contextual information

## 🛠 Tech Stack

- **Language**: Go 1.22+
- **Database**: MySQL 8.0+
- **Authentication**: JWT (JSON Web Tokens)
- **HTTP Router**: Gorilla Mux
- **Migration**: golang-migrate
- **Containerization**: Docker & Docker Compose
- **Orchestration**: Kubernetes with HPA
- **Logging**: Custom structured logger (JSON format)

## Project Structure

```bash
.
├── cmd/                    # Main application entry points
│   ├── api/               # API server configuration
│   ├── main.go            # Main program entry
│   └── migrate/           # Database migration tool
├── config/                # Configuration files
├── db/                    # Database connection and management
├── service/               # Business logic layer
│   ├── auth/             # Authentication service
│   ├── cart/             # Shopping cart service
│   ├── order/            # Order service
│   ├── product/          # Product service
│   └── user/             # User service
├── types/                 # Type definitions
└── utils/                 # Utility functions

```

## 📋 Environment Requirements

- Go 1.22+
- MySQL 8.0+
- Make
- Docker & Docker Compose (for containerized deployment)
- Kubernetes (optional, for production deployment)

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended for Development)

1. **Clone the project**
   ```bash
   git clone https://github.com/xiaopeng-ye/ecommerce-golang.git
   cd ecommerce-golang
   ```

2. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env file - IMPORTANT: Change DB_PASSWORD and JWT_SECRET
   ```

3. **Start the services**
   ```bash
   docker-compose up --build
   ```

   The API will be available at `http://localhost:8080`

### Option 2: Local Development

1. **Clone and configure**
   ```bash
   git clone https://github.com/xiaopeng-ye/ecommerce-golang.git
   cd ecommerce-golang
   cp .env.example .env
   # Edit .env with your database credentials
   ```

2. **Start MySQL** (or use Docker)
   ```bash
   docker run -d \
     --name ecommerce-mysql \
     -e MYSQL_ROOT_PASSWORD=mypassword \
     -e MYSQL_DATABASE=ecom \
     -p 3306:3306 \
     mysql:8.0
   ```

3. **Run database migrations**
   ```bash
   make migrate-up
   ```

4. **Build and run**
   ```bash
   make run
   ```

### Option 3: Kubernetes (Production)

See [k8s/README.md](k8s/README.md) for detailed Kubernetes deployment instructions.

```bash
# Quick deploy (after configuring secrets)
kubectl apply -f k8s/
```

## 🔧 Available Make Commands

Run `make help` to see all available commands:

**Building**
- `make build` - Build the application with version information
- `make build-linux` - Build for Linux (cross-compile)
- `make docker-build` - Build Docker image with version tags

**Testing**
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report

**Running**
- `make run` - Build and run locally
- `make docker-run` - Run with Docker Compose

**Database**
- `make migration name=<name>` - Create a new migration file
- `make migrate-up` - Apply database migrations
- `make migrate-down` - Rollback database migrations

**Code Quality**
- `make fmt` - Format code
- `make vet` - Run go vet
- `make lint` - Run linter (requires golangci-lint)

**Kubernetes**
- `make k8s-apply` - Apply all Kubernetes manifests
- `make k8s-delete` - Delete Kubernetes resources
- `make k8s-logs` - Stream logs from pods
- `make k8s-status` - Show deployment status

**Utility**
- `make version` - Show version information
- `make clean` - Clean build artifacts

## 📡 API Endpoints

All API endpoints are prefixed with `/api/v1`

### Health & Monitoring
- `GET /health` - Liveness probe (returns 200 if service is alive)
- `GET /ready` - Readiness probe (checks database connectivity)

### User Management
- `POST /api/v1/register` - User registration
- `POST /api/v1/login` - User login (returns JWT token)
- `GET /api/v1/users/{userID}` - Get user information (requires authentication)

### Product Management
- `GET /api/v1/products` - List all products
- `GET /api/v1/products/{productID}` - Get product details
- `POST /api/v1/products` - Create a new product (requires authentication)

### Shopping & Orders
- `POST /api/v1/cart/checkout` - Checkout cart (requires authentication)
- `GET /api/v1/orders` - List user orders (requires authentication)

## 🔐 Authentication

The API uses JWT (JSON Web Token) for authentication.

### Getting a Token
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

### Using the Token
Include the JWT token in the Authorization header:

```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer <your_jwt_token>"
```

## 📊 Database Schema

The project uses MySQL database with the following main tables:

- `users` - User accounts and profiles
- `products` - Product catalog
- `orders` - Order records
- `order_items` - Individual items in orders

Migrations are managed using [golang-migrate](https://github.com/golang-migrate/migrate).

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage
```

## 🌐 Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ENVIRONMENT` | No | development | Application environment (development/staging/production) |
| `LOG_LEVEL` | No | info | Log level (debug/info/warn/error/fatal) |
| `LOG_FORMAT` | No | json | Log format (json/text) |
| `PORT` | No | 8080 | HTTP server port |
| `PUBLIC_HOST` | No | http://localhost | Public host URL |
| `DB_HOST` | No | localhost | Database host |
| `DB_PORT` | No | 3306 | Database port |
| `DB_USER` | No | root | Database user |
| `DB_PASSWORD` | **Yes** | - | Database password |
| `DB_NAME` | No | ecom | Database name |
| `JWT_SECRET` | **Yes** | - | JWT signing secret |
| `JWT_EXPIRATION_IN_SECONDS` | No | 604800 | JWT token expiration (7 days) |

**Note**: In production, `DB_PASSWORD` and `JWT_SECRET` must be set and should use strong, randomly generated values.

## 📖 Documentation

- [12-Factor App Implementation](docs/12-FACTOR.md) - How this project implements 12-factor methodology
- [Kubernetes Deployment Guide](k8s/README.md) - Kubernetes deployment instructions

## 🏗 Project Structure

```
.
├── cmd/
│   ├── api/              # API server implementation
│   ├── main.go           # Application entry point
│   └── migrate/          # Database migration tool
├── config/               # Configuration management
├── db/                   # Database connection
├── docs/                 # Documentation
├── k8s/                  # Kubernetes manifests
├── service/              # Business logic
│   ├── auth/            # Authentication
│   ├── cart/            # Shopping cart
│   ├── order/           # Order management
│   ├── product/         # Product management
│   └── user/            # User management
├── types/                # Type definitions
├── utils/                # Utility functions (logger, etc.)
├── Dockerfile            # Multi-stage Docker build
├── docker-compose.yml    # Local development environment
├── Makefile             # Build automation
└── README.md            # This file
```

## 🎯 12-Factor App Compliance

This project follows the [12-Factor App](https://12factor.net/) methodology:

- ✅ **I. Codebase**: One codebase in Git, many deploys
- ✅ **II. Dependencies**: Explicitly declared in go.mod
- ✅ **III. Config**: Configuration via environment variables
- ✅ **IV. Backing Services**: Database as attached resource
- ✅ **V. Build, Release, Run**: Strict separation of stages
- ✅ **VI. Processes**: Stateless, share-nothing processes
- ✅ **VII. Port Binding**: Self-contained HTTP service
- ✅ **VIII. Concurrency**: Scale out via process model
- ✅ **IX. Disposability**: Fast startup and graceful shutdown
- ✅ **X. Dev/Prod Parity**: Keep environments similar
- ✅ **XI. Logs**: Treat logs as event streams (JSON to stdout)
- ✅ **XII. Admin Processes**: Run admin tasks as one-off processes

See [docs/12-FACTOR.md](docs/12-FACTOR.md) for detailed implementation.

## 🚀 Production Deployment

### Building for Production

```bash
# Build with version tagging
VERSION=v1.0.0 make docker-build

# Push to registry
docker tag ecommerce-api:v1.0.0 your-registry/ecommerce-api:v1.0.0
docker push your-registry/ecommerce-api:v1.0.0
```

### Kubernetes Deployment

```bash
# Update k8s/secret.yaml with production secrets
# Update k8s/configmap.yaml with production config

# Deploy
kubectl apply -f k8s/

# Monitor deployment
kubectl get pods -n ecommerce -w
```

### Health Checks

The application provides health check endpoints:

- **Liveness**: `GET /health` - Returns 200 if app is running
- **Readiness**: `GET /ready` - Returns 200 if app is ready to serve traffic (DB connected)

## 🔒 Security Considerations

- Run as non-root user in containers (UID 1000)
- Secrets managed via environment variables or Kubernetes Secrets
- JWT tokens for API authentication
- Input validation on all endpoints
- SQL injection prevention via parameterized queries
- HTTPS recommended for production (configure via reverse proxy/ingress)

## 📈 Monitoring and Logging

### Structured Logging

All logs are output in JSON format to stdout:

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "message": "Starting HTTP server",
  "service": "ecommerce-api",
  "fields": {
    "address": ":8080",
    "version": "v1.0.0"
  }
}
```

### Viewing Logs

```bash
# Docker Compose
docker-compose logs -f api

# Kubernetes
kubectl logs -f -n ecommerce -l app=ecommerce-api

# With Makefile
make k8s-logs
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

MIT License

## 📧 Support

For issues and questions:
- Open an issue on GitHub
- Check the [documentation](docs/)
- Review [12-Factor implementation guide](docs/12-FACTOR.md)
