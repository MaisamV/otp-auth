# Technology Stack & Architecture Decisions

## Overview
This document explains the technology choices made for the OTP-based authentication backend service and the reasoning behind each decision compared to alternatives.

## Core Technologies

### 🔧 **Programming Language: Go 1.24.4**

**Why Go?**
- **High performance** with low memory footprint
- **Excellent concurrency** with goroutines for handling multiple requests
- **Fast compilation** and deployment
- **Strong standard library** for HTTP servers and cryptography
- **Static binary compilation** simplifies Docker deployment

**Alternatives Considered:**
- **Node.js**: Faster development but single-threaded, higher memory usage
- **Python**: Great ecosystem but slower performance
- **Java**: Enterprise-grade but heavier resource usage and slower startup
- **Rust**: Excellent performance but steeper learning curve

---

## 🏗️ **Architecture Pattern: Clean/Hexagonal Architecture**

**Why Clean Architecture?**
- **Separation of concerns** between domain, application, and infrastructure
- **Testability** through dependency inversion
- **Maintainability** with clear boundaries
- **Technology independence** - easy to swap databases or frameworks

**Key Components:**
- **Domain Layer**: Business logic, entities, value objects
- **Application Layer**: Use cases, interfaces
- **Infrastructure Layer**: Database, Redis, HTTP handlers

**Alternatives Considered:**
- **MVC**: Simpler but tends to create tightly coupled code
- **Layered Architecture**: Good but less flexible than hexagonal
- **Microservices**: Overkill for authentication service

---

## 🔐 **Cryptography: ECDSA (Elliptic Curve Digital Signature Algorithm)**

**Why ECDSA?**
- **Smaller key sizes** (256-bit ECDSA ≈ 3072-bit RSA security)
- **Faster signature verification** than RSA
- **Smaller JWT tokens** due to compact signatures
- **Modern standard** widely supported
- **Better performance** on mobile devices

**Implementation:**
- **P-256 curve** (secp256r1) for optimal security/performance balance
- **Automatic key generation** if keys don't exist
- **PEM format** for key storage

**Alternatives Considered:**
- **RSA**: Larger signatures, slower verification, but more universal support
- **HMAC**: Symmetric, faster but requires shared secrets across services
- **Ed25519**: Excellent but less ecosystem support

---

## 🗄️ **Database Strategy: PostgreSQL + Citus (Progressive Scaling)**

### **Primary Database: PostgreSQL**

**Why PostgreSQL?**
- **ACID compliance** for critical authentication data
- **JSON support** for flexible user metadata
- **Excellent performance** for millions of users
- **Rich ecosystem** and tooling
- **Strong consistency** for refresh token management

**Schema Design:**
- **Users table**: Basic user information with phone numbers
- **Refresh tokens table**: Token management with revocation support
- **Indexes**: Optimized for phone number and token lookups

### **Scaling Strategy: Citus Extension**

**Why Citus for Scaling?**
- **Horizontal scaling** when needed (>10M users)
- **Distributed SQL** maintains familiar PostgreSQL interface
- **Automatic sharding** by user_id or phone_number
- **Zero downtime migration** from single PostgreSQL

**Alternatives Considered:**
- **MongoDB**: Good for user profiles but lacks ACID for financial data
- **CockroachDB**: Great distributed SQL but more complex setup
- **MySQL**: Good performance but weaker JSON support
- **Cassandra**: Excellent scale but eventual consistency issues

---

## ⚡ **Caching & Session Storage: Redis**

**Why Redis?**
- **Sub-millisecond latency** for OTP validation
- **Built-in TTL** perfect for temporary OTP storage
- **Atomic operations** for rate limiting counters
- **Memory-based storage** with optional persistence
- **Simple data structures** (strings, counters)

**Use Cases:**
- **OTP storage**: `<phone_number>` → `<session_id>-<hashed_otp>` (2min TTL)
- **Rate limiting**: `rate_limit:<phone_number>` → counter (10min TTL)
- **Session validation**: Fast lookup for session_id validation

**Alternatives Considered:**
- **Memcached**: No TTL support, no atomic operations
- **Database**: Too slow for real-time OTP validation
- **In-memory maps**: No persistence, scaling issues

---

## 🌐 **HTTP Server & Framework**

### **Recommended: Gin Web Framework**

**Why Gin?**
- **High performance** (40x faster than Martini)
- **Minimal memory footprint**
- **Middleware support** for CORS, logging, authentication
- **JSON binding** and validation
- **Easy testing** with httptest package

**Alternatives:**
- **Echo**: Similar performance, not mature as Gin
- **Fiber V3**: Express.js-like API, very fast but in beta
- **Standard net/http**: More control but more boilerplate
- **Chi**: Lightweight router, good for simple APIs

---

## 🔌 **Database Connectivity**

### **Recommended: GORM + pgx Driver**

**Why GORM?**
- **ORM convenience** with raw SQL fallback
- **Migration support** for schema management
- **Connection pooling** built-in
- **Hooks and callbacks** for audit trails
- **Type safety** with Go structs

**Why pgx Driver?**
- **Best PostgreSQL performance** in Go
- **Native PostgreSQL features** support
- **Connection pooling** with pgxpool
- **Prepared statements** for security

**Alternatives:**
- **sqlx**: More control, less convenience
- **database/sql**: Standard but more boilerplate
- **Ent**: Type-safe but more complex

---

## 📡 **Redis Connectivity**

### **Recommended: go-redis/redis**

**Why go-redis?**
- **Most popular** Redis client for Go
- **Connection pooling** built-in
- **Pipeline support** for batch operations
- **Cluster support** for scaling
- **Context support** for timeouts

**Alternatives:**
- **redigo**: Older, less features
- **rueidis**: Newer, very fast but less mature

---

## 🔒 **Security Mechanisms**

### **Session ID Anti-Brute Force**

**Why Session ID?**
- **Prevents OTP enumeration** without valid session
- **Stateless validation** via Redis
- **Minimal performance impact**
- **Elegant security layer** beyond rate limiting

**Implementation:**
- Generated during "Send OTP" flow
- Required for "Login/Register" flow
- Stored with OTP in Redis: `<session_id>-<hashed_otp>`

### **Rate Limiting Strategy**

**Why Redis Counters?**
- **Atomic increment** operations
- **Automatic expiration** with TTL
- **Distributed** across multiple server instances
- **Simple implementation** with high performance

**Configuration:**
- **3 attempts per phone number** in 10 minutes
- **Sliding window** with Redis TTL

---

## 🧪 **Testing Strategy**

### **Recommended Testing Stack**

**Unit Testing:**
- **testify/suite**: Structured test suites
- **testify/mock**: Interface mocking
- **testify/assert**: Rich assertions

**Integration Testing:**
- **testcontainers-go**: Real PostgreSQL/Redis in tests
- **httptest**: HTTP endpoint testing
- **dockertest**: Alternative container testing

**Why This Approach?**
- **Real database testing** catches integration issues
- **Isolated test environments** prevent test interference
- **Fast feedback loop** with parallel test execution

---

## 📦 **Deployment & DevOps**

### **Containerization: Docker**

**Why Docker?**
- **Consistent environments** across dev/staging/prod
- **Easy dependency management**
- **Horizontal scaling** with container orchestration
- **Go static binaries** create minimal images

### **Orchestration: Docker Compose**

**Why Docker Compose?**
- **Simple multi-service setup** for development
- **Easy local testing** with real dependencies
- **Production-ready** for small to medium deployments
- **Clear service dependencies** and networking

**Services:**
- **app**: Go application
- **postgres**: Database with initialization scripts
- **redis**: Cache and session storage
- **nginx**: Reverse proxy (optional)

---

## 🔧 **Configuration Management**

### **Recommended: Viper + YAML**

**Why Viper?**
- **Multiple config sources** (files, env vars, flags)
- **Hot reloading** for development
- **Type-safe** configuration binding
- **Environment-specific** configs

---

## 📊 **Monitoring & Observability**

### **Recommended Stack**

**Logging:**
- **logrus** or **zap**: Structured logging
- **JSON format** for log aggregation
- **Different log levels** for environments

**Metrics:**
- **Prometheus**: Metrics collection
- **Custom metrics**: OTP success rate, login attempts
- **HTTP metrics**: Request duration, status codes

**Health Checks:**
- **Database connectivity** checks
- **Redis connectivity** checks
- **Kubernetes-ready** health endpoints

---

## 🚀 **Performance Considerations**

### **Optimization Strategies**

**Database:**
- **Connection pooling** (max 25 connections)
- **Prepared statements** for frequent queries
- **Indexes** on phone_number, refresh_token_hash
- **Read replicas** for user lookups (future)

**Redis:**
- **Pipeline operations** for batch requests
- **Connection pooling** with go-redis
- **Memory optimization** with appropriate TTLs

**Application:**
- **Goroutine pools** for concurrent processing
- **Context timeouts** for all external calls
- **Graceful shutdown** for zero-downtime deployments

---

## 📈 **Scaling Strategy**

### **Horizontal Scaling Path**

**Phase 1: Single Instance** (0-100K users)
- Single PostgreSQL instance
- Single Redis instance
- Single application instance

**Phase 2: Load Balancing** (100K-1M users)
- Multiple application instances
- PostgreSQL read replicas
- Redis clustering

**Phase 3: Database Sharding** (1M-10M users)
- Citus extension for PostgreSQL
- Distributed Redis cluster
- Microservice decomposition

**Phase 4: Multi-Region** (10M+ users)
- Geographic distribution
- Regional databases
- CDN for static assets

---

## 🔄 **Migration & Deployment Strategy**

### **Database Migrations**
- **GORM AutoMigrate** for development
- **Custom migration scripts** for production
- **Backward compatibility** for zero-downtime deployments

### **Blue-Green Deployment**
- **Health checks** before traffic switching
- **Database migration** compatibility
- **Rollback procedures** for failed deployments

---

## 📝 **Summary**

This technology stack provides:
- **High performance** with Go and optimized databases
- **Strong security** with ECDSA and session-based anti-brute force
- **Excellent scalability** with progressive database scaling
- **Developer productivity** with clean architecture and good tooling
- **Operational simplicity** with Docker and comprehensive monitoring

The choices balance **immediate development speed** with **long-term scalability** and **operational excellence**.