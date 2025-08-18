# Architecture Documentation

## Overview
This document describes the architectural patterns, design principles, and implementation details for the OTP-based authentication backend service.

## 🏗️ **Architectural Pattern: Clean Architecture (Hexagonal)**

### **Core Principles**

Our system follows **Clean Architecture** principles with **Hexagonal Architecture** implementation, ensuring:

- **Independence of Frameworks**: Business logic doesn't depend on external libraries
- **Testability**: Business rules can be tested without UI, database, or external services
- **Independence of UI**: Easy to change UI without changing business rules
- **Independence of Database**: Business rules not bound to specific database
- **Independence of External Services**: Business rules don't know about outside world

### **Dependency Rule**

```
External Services → Infrastructure → Application → Domain
                                                    ↑
                                            Dependencies point inward
```

**Key Rule**: Dependencies can only point inward. Inner layers cannot know about outer layers.

---

## 🎯 **Layer Responsibilities**

### **1. Domain Layer (Core)**

**Purpose**: Contains business logic, entities, and business rules.
#### **Entities**
#### **Value Objects**
#### **Repository Interfaces**
#### **Domain Services**

### **2. Application Layer (Use Cases)**

**Purpose**: Orchestrates domain objects to fulfill business use cases.
#### **Use Case Example**

### **3. Infrastructure Layer (External Concerns)**

**Purpose**: Implements interfaces defined in inner layers and handles external dependencies.
#### **HTTP Handlers**
#### **Repository Implementations**
#### **Service Implementations**

---

## 🔄 **Design Patterns Used**

### **1. Repository Pattern**

**Purpose**: Abstracts data access logic and provides a uniform interface for accessing domain objects.
**Implementation**:
- **Interface** defined in domain layer
- **Implementation** in infrastructure layer
- **Separation** of read and write operations

### **2. Dependency Injection Pattern**
**Purpose**: Inverts control of dependencies, making the system more testable and flexible.

### **3. Factory Pattern**
**Purpose**: Creates objects without specifying their concrete classes.

### **4. Strategy Pattern**

**Purpose**: Defines a family of algorithms and makes them interchangeable.

### **5. Middleware Pattern**
**Purpose**: Provides a way to filter and process HTTP requests.

---

## 🧪 **Testing Architecture**

### **Testing Strategy by Layer**
#### **1. Domain Layer Testing**
#### **2. Application Layer Testing (Use Cases)**
#### **3. Infrastructure Layer Testing**

---

## 🔒 **Security Architecture**

### **Authentication Flow Security**

1. **OTP Generation & Storage**:
   - OTPs are **hashed** before storage in Redis
   - **Session ID** prevents brute force attacks
   - **Rate limiting** prevents abuse

2. **JWT Token Security**:
   - **ECDSA signing** with P-256 curve
   - **Short-lived access tokens** (15 minutes)
   - **Refresh token rotation** on each use
   - **Token revocation** support

3. **Session Management**:
   - **HttpOnly cookies** prevent XSS
   - **CSRF protection** for cookie-based auth
   - **Session binding** with refresh tokens

### **Data Protection**

```go
// Security layers in data flow
Phone Number → Validation → Normalization → Storage
OTP → Hashing → Redis Storage (TTL)
Refresh Token → Hashing → Database Storage
Passwords → Bcrypt → Database Storage (if added)
```

---

## 📊 **Performance Architecture**

### **Caching Strategy**

1. **Redis for Hot Data**:
   - OTP storage (2-minute TTL)
   - Rate limiting counters (10-minute TTL)
   - Session validation cache

2. **Database Optimization**:
   - Connection pooling (25 connections)
   - Prepared statements
   - Proper indexing strategy

### **Scalability Considerations**

1. **Horizontal Scaling**:
   - Stateless application design
   - Database read replicas
   - Redis clustering

2. **Performance Monitoring**:
   - Request duration metrics
   - Database query performance
   - Redis operation latency

---

## 🎯 **Summary**

This architecture provides:

- **Clean separation** of concerns across layers
- **High testability** through dependency injection
- **Flexibility** to change external dependencies
- **Scalability** through stateless design
- **Security** through multiple defense layers
- **Maintainability** through clear patterns and structure
- **Observability** through comprehensive logging and metrics

The architecture follows **SOLID principles** and **Clean Architecture** patterns, ensuring the codebase remains maintainable and extensible as the system grows.