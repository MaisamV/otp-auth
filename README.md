# OTP Authentication Service

## Project Summary

A robust, production-ready OTP (One-Time Password) authentication service built with Go, following clean architecture principles. This service provides secure phone number-based authentication using JWT tokens and Redis for OTP storage.

The system implements a complete authentication flow with phone number validation, OTP generation and verification, JWT token management, and user session handling. Built with Clean Architecture, it ensures maintainability, testability, and scalability.

### Key Features

- 📱 **Phone Number Authentication**: Secure OTP-based authentication
- 🔐 **JWT Tokens**: ECDSA-signed access and refresh tokens
- 🚀 **Clean Architecture**: Domain-driven design with clear separation of concerns
- 📊 **Rate Limiting**: Configurable rate limiting for API endpoints and OTP requests
- 🔒 **Security**: Bcrypt password hashing, secure session management
- 🐳 **Docker Support**: Complete containerization with Docker Compose
- 📈 **Monitoring**: Health checks and metrics endpoints
- ⚡ **High Performance**: Redis caching and PostgreSQL persistence
- 🛡️ **CORS Support**: Configurable cross-origin resource sharing
- 📝 **Comprehensive Logging**: Structured logging with configurable levels

## Quick Start

### Prerequisites

- Docker and Docker Compose

## How to run with Docker compose

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd otp-auth
   ```

2. **Start all services** (PostgreSQL, Redis, and OTP Auth Service):
   ```bash
   docker-compose up --build -d
   ```

3. **Verify services are running**:
   ```bash
   docker-compose ps
   ```

4. **Check service health**:
   ```bash
   curl http://localhost:8080/health
   ```

5. **View logs** (optional):
   ```bash
   docker-compose logs -f otp-auth-service
   ```

6. **Open the swagger to test the service**:
   http://localhost:8080/swagger/index.html

## How to run locally

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### Setup Steps

1. **Install Go dependencies**:
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL database**:
   ```bash
   # Create database
   createdb otp_auth
   
   # Or using psql
   psql -c "CREATE DATABASE otp_auth;"
   ```

3. **Start Redis server**:
   ```bash
   # Using default configuration
   redis-server
   
   # Or with custom config
   redis-server /path/to/redis.conf
   ```

4. **Configure application**:
   ```bash
   # Copy and edit configuration
   cp configs/config.yaml configs/config.local.yaml
   # Edit configs/config.local.yaml with your database and Redis settings
   ```

5. **Run database migrations** (automatic on startup):
   ```bash
   # Migrations will run automatically when starting the application
   ```

6. **Start the application**:
   ```bash
   # Development mode
   go run cmd/server/main.go
   
   # Or build and run
   go build -o bin/otp-auth cmd/server/main.go
   ./bin/otp-auth
   ```

7. **Verify the application is running**:
   ```bash
   curl http://localhost:8080/health
   ```

8. **Open the swagger to test the service**:
   http://localhost:8080/swagger/index.html

## API Endpoints

### Authentication

- `POST /api/v1/auth/send-otp` - Send OTP to phone number
- `POST /api/v1/auth/login` - Login with OTP
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout user

### User Management

- `GET /api/v1/users` - Get users (admin only)
- `GET /api/v1/users/profile` - Get current user profile
- `PUT /api/v1/users/:id/scope` - Update user scope (admin only)

### Health & Monitoring

- `GET /health` - Health check
- `GET /ready` - Readiness check
- `GET /live` - Liveness check


## Development

### Project Structure

- **Domain Layer**: Contains business entities, repository interfaces, and domain services
- **Application Layer**: Contains use cases, DTOs, and application services
- **Infrastructure Layer**: Contains concrete implementations of repositories, external services, and HTTP handlers

### Adding New Features

1. **Define domain entities** in `internal/domain/entities/`
2. **Create repository interfaces** in `internal/domain/repositories/`
3. **Implement use cases** in `internal/application/usecases/`
4. **Add HTTP handlers** in `internal/infrastructure/http/handlers/`
5. **Update routing** in `internal/infrastructure/http/router/`

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/domain/entities/
```

## Deployment

### Docker Production Deployment

1. **Build production image**:
   ```bash
   docker build -t otp-auth:latest .
   ```

2. **Use production compose**:
   ```bash
   docker-compose -f docker-compose.yml --profile production up -d
   ```

### Environment-Specific Configurations

- `configs/config.yaml` - Development configuration
- `configs/config.prod.yaml` - Production configuration

### Health Checks

The service includes comprehensive health checks:

- **Health**: Overall service health
- **Ready**: Service readiness (dependencies available)
- **Live**: Service liveness (service is running)

## Security Considerations

- **JWT Tokens**: Use ECDSA signing with secure key management
- **Rate Limiting**: Prevent brute force attacks
- **OTP Security**: Short-lived, hashed OTP codes
- **CORS**: Configurable cross-origin policies
- **Input Validation**: Comprehensive request validation
- **Secure Headers**: Security headers in responses

## Monitoring and Observability

### Metrics

- Request/response metrics
- Database connection metrics
- Redis connection metrics
- Custom business metrics

### Logging

- Structured JSON logging
- Configurable log levels
- Request/response logging
- Error tracking

### Health Monitoring

```bash
# Check all health endpoints
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:8080/live
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:

- Create an issue in the repository
- Check the documentation
- Review the configuration examples

---

**Built with ❤️ using Go and Clean Architecture principles**