# OTP Authentication Service

A robust, production-ready OTP (One-Time Password) authentication service built with Go, following clean architecture principles. This service provides secure phone number-based authentication using JWT tokens and Redis for OTP storage.

## Features

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

## Architecture

This project follows **Clean Architecture** principles with the following layers:

```
├── cmd/                    # Application entry points
├── internal/
│   ├── domain/            # Business logic and entities
│   │   ├── entities/      # Domain entities
│   │   ├── repositories/  # Repository interfaces
│   │   ├── services/      # Service interfaces
│   │   └── valueobjects/  # Value objects
│   ├── application/       # Application layer
│   │   ├── dto/          # Data transfer objects
│   │   ├── ports/        # Ports (interfaces)
│   │   └── usecases/     # Use cases
│   ├── infrastructure/    # Infrastructure layer
│   │   ├── http/         # HTTP handlers, middleware, routing
│   │   ├── persistence/  # Database implementations
│   │   └── services/     # Service implementations
│   └── config/           # Configuration management
├── pkg/                   # Shared utilities
├── configs/              # Configuration files
└── scripts/              # Deployment and utility scripts
```

## Quick Start

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL 15+
- Redis 7+

### Using Docker Compose (Recommended)

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd otp-auth
   ```

2. **Start the services**:
   ```bash
   docker-compose up -d
   ```

3. **Check service health**:
   ```bash
   curl http://localhost:8080/health
   ```

### Manual Setup

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL and Redis**:
   ```bash
   # PostgreSQL
   createdb otp_auth
   
   # Redis (default configuration)
   redis-server
   ```

3. **Configure environment**:
   ```bash
   cp configs/config.yaml configs/config.local.yaml
   # Edit configs/config.local.yaml with your settings
   ```

4. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```

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

## Configuration

The service uses YAML configuration files with environment variable overrides:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug" # debug, release, test

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "otp_auth"
  sslmode: "disable"

redis:
  addr: "localhost:6379"
  password: ""
  db: 0

jwt:
  access_token_ttl: "15m"
  refresh_token_ttl: "168h" # 7 days
  issuer: "otp-auth-service"

otp:
  length: 6
  ttl: "5m"
  sender_type: "console" # console, sms

security:
  rate_limit:
    enabled: true
    requests: 100
    window: "1m"
    otp_limit: 5
    otp_window: "1h"
```

### Environment Variables

All configuration can be overridden using environment variables with the `OTP_AUTH_` prefix:

```bash
export OTP_AUTH_DATABASE_HOST=localhost
export OTP_AUTH_DATABASE_PASSWORD=secret
export OTP_AUTH_REDIS_ADDR=redis:6379
```

## Usage Examples

### Send OTP

```bash
curl -X POST http://localhost:8080/api/v1/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+989123456789"}'
```

### Login with OTP

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "+989123456789",
    "otp_code": "123456",
    "session_id": "session-uuid"
  }'
```

### Access Protected Endpoint

```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <access_token>"
```

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