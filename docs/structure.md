## 📁 **Project Structure**

```
otp-auth/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── domain/                     # Domain Layer (Core Business Logic)
│   │   ├── entities/
│   │   │   ├── user.go
│   │   │   ├── otp.go
│   │   │   └── refresh_token.go
│   │   └── valueobjects/
│   │       ├── phone_number.go
│   │       └── session_id.go
│   ├── application/                # Application Layer (Use Cases)
│   │   ├── usecases/
│   │   │   ├── send_otp.go
│   │   │   ├── login_register.go
│   │   │   ├── refresh_token.go
│   │   │   ├── logout.go
│   │   │   └── admin_users.go
│   │   ├── dto/                    # Data Transfer Objects
│   │   │   ├── requests.go
│   │   │   └── responses.go
│   │   └── ports/                  # Application Interfaces
│   │       ├── repositories/       # Repository Interfaces
│   │       │   ├── user_repository.go
│   │       │   ├── otp_repository.go
│   │       │   ├── token_repository.go
│   │       │   └── rate_limit_repository.go
│   │       └── services/           # Service Interfaces
│   │           ├── otp_sender.go
│   │           ├── jwt_service.go
│   │           └── hash_service.go
│   └── infrastructure/             # Infrastructure Layer (External Concerns)
│       ├── http/                   # HTTP Handlers
│       │   ├── handlers/
│       │   │   ├── auth_handler.go
│       │   │   └── admin_handler.go
│       │   ├── middleware/
│       │   │   ├── cors.go
│       │   │   ├── auth.go
│       │   │   └── rate_limit.go
│       │   └── router.go
│       ├── persistence/            # Database Implementations
│       │   ├── postgres/
│       │   │   ├── user_repository.go
│       │   │   ├── token_repository.go
│       │   │   └── migrations/
│       │   └── redis/
│       │       ├── otp_repository.go
│       │       └── rate_limiter.go
│       ├── services/               # External Service Implementations
│       │   ├── console_otp_sender.go
│       │   ├── ecdsa_jwt_service.go
│       │   └── bcrypt_hash_service.go
│       └── config/
│           ├── config.go
│           └── database.go
├── pkg/                            # Shared Utilities
│   ├── utils/
│   │   ├── hash.go
│   │   ├── random.go
│   │   └── validator.go
│   └── errors/
│       └── custom_errors.go
├── configs/
│   ├── config.yaml
│   └── config.prod.yaml
├── keys/                           # ECDSA Keys (auto-generated)
│   ├── private.pem
│   └── public.pem
├── docker-compose.yml
├── Dockerfile
└── docs/
    ├── blueprint.md
    ├── tech.md
    └── architecture.md
```