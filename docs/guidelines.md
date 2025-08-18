# Development Guidelines

## Overview
This document establishes the development standards and practices that must be followed by all contributors (human and AI) working on the OTP Authentication Backend Service.

---

## 🎯 **General Principles**

### **SOLID Principles**
- **Single Responsibility**: Each class/function should have one reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Liskov Substitution**: Derived classes must be substitutable for base classes
- **Interface Segregation**: Clients shouldn't depend on interfaces they don't use
- **Dependency Inversion**: Depend on abstractions, not concretions

### **Clean Code Principles**
- **Readability**: Code should be self-documenting
- **Simplicity**: Prefer simple solutions over complex ones
- **DRY (Don't Repeat Yourself)**: Eliminate code duplication
- **YAGNI (You Aren't Gonna Need It)**: Don't implement features until needed
- **Boy Scout Rule**: Leave code cleaner than you found it

---

## 📝 **Coding Standards**

### **Go Language Standards**

#### **Formatting**
- Use `gofmt` for automatic formatting
- Use `goimports` for import organization
- Line length: Maximum 120 characters
- Indentation: Use tabs, not spaces
- No trailing whitespace

#### **Naming Conventions**

**Variables:**
- Use camelCase for local variables: `userID`, `phoneNumber`
- Use short names for short-lived variables: `i`, `err`, `ctx`
- Use descriptive names for longer-lived variables: `refreshTokenRepository`

**Functions:**
- Use PascalCase for exported functions: `GenerateOTP()`, `ValidateToken()`
- Use camelCase for private functions: `hashPassword()`, `generateSessionID()`
- Use verb-noun pattern: `CreateUser()`, `SendOTP()`, `ValidatePhone()`

**Types:**
- Use PascalCase for all types: `User`, `PhoneNumber`, `OTPRequest`
- Use descriptive names: `UserRepository` not `UR`
- Interface names should describe behavior: `TokenGenerator`, `OTPSender`

**Constants:**
- Use PascalCase for exported constants: `DefaultOTPLength`
- Use camelCase for private constants: `maxRetryAttempts`
- Group related constants in blocks

**Packages:**
- Use lowercase, single word when possible: `domain`, `handlers`, `postgres`
- Use descriptive names: `usecases` not `uc`
- Avoid abbreviations: `repositories` not `repos`

#### **File Organization**
- One primary type per file
- File names should match primary type: `user.go` for `User` type
- Use snake_case for file names: `user_repository.go`
- Group related functionality in same package

#### **Import Organization**
- Standard library imports first
- Third-party imports second
- Local imports last
- Separate groups with blank lines
- Use blank imports only when necessary

#### **Error Handling**
- Always handle errors explicitly
- Use custom error types for domain errors
- Wrap errors with context using `fmt.Errorf`
- Return errors as last return value
- Use early returns to reduce nesting

#### **Comments and Documentation**
- All exported functions must have comments
- Comments should explain "why", not "what"
- Use complete sentences in comments
- Start comments with the function/type name
- Use `//` for single-line comments
- Use `/* */` for multi-line comments only when necessary

---

## 🏗️ **Architecture Guidelines**

### **Layer Separation**
- **Domain Layer**: No dependencies on external packages
- **Application Layer**: Can depend only on domain layer
- **Infrastructure Layer**: Can depend on application and domain layers
- Never import from outer layers to inner layers

### **Dependency Injection**
- Use constructor injection for dependencies
- Inject interfaces, not concrete types
- Keep constructors simple and focused
- Validate dependencies in constructors

### **Interface Design**
- Keep interfaces small and focused
- Define interfaces in the package that uses them
- Use composition over inheritance
- Prefer many small interfaces over few large ones

### **Repository Pattern**
- Separate read and write operations
- Use context for all repository methods
- Return domain entities, not database models
- Handle database-specific errors in repository layer

---

## 🧪 **Testing Guidelines**

### **Testing Strategy**
- **Unit Tests**: Test individual functions/methods in isolation
- **Integration Tests**: Test component interactions
- **End-to-End Tests**: Test complete user workflows
- Aim for 80%+ code coverage

### **Test Organization**
- Place tests in same package as code being tested
- Use `_test.go` suffix for test files
- One test file per source file when possible
- Group related tests in test suites

### **Test Naming**
- Use `TestFunctionName_Scenario_ExpectedResult` pattern
- Be descriptive: `TestSendOTP_ValidPhoneNumber_ReturnsSuccess`
- Use table-driven tests for multiple scenarios

### **Test Structure**
- Follow Arrange-Act-Assert pattern
- Use meaningful test data
- Clean up resources in test teardown
- Use test helpers to reduce duplication

### **Mocking Guidelines**
- Mock external dependencies only
- Use interface-based mocking
- Keep mocks simple and focused
- Verify mock expectations in tests

### **Integration Testing**
- Use test containers for real database testing
- Test with realistic data volumes
- Test error scenarios and edge cases
- Ensure tests are isolated and repeatable

---

## 🔒 **Security Guidelines**

### **Data Protection**
- Never log sensitive data (passwords, tokens, OTPs)
- Hash all passwords and OTPs before storage
- Use secure random number generation
- Validate all input data

### **Authentication & Authorization**
- Use strong cryptographic algorithms (ECDSA, bcrypt)
- Implement proper session management
- Use HTTPS for all communications
- Implement rate limiting for all endpoints

### **Error Handling**
- Don't expose internal system details in error messages
- Use generic error messages for security-sensitive operations
- Log detailed errors internally for debugging
- Implement proper error codes for client applications

### **Configuration**
- Never commit secrets to version control
- Use environment variables for sensitive configuration
- Implement proper key rotation procedures
- Use secure defaults for all configuration options

---

## 📊 **Performance Guidelines**

### **Database Operations**
- Use connection pooling
- Implement proper indexing strategy
- Use prepared statements for repeated queries
- Avoid N+1 query problems

### **Caching Strategy**
- Cache frequently accessed data
- Use appropriate TTL values
- Implement cache invalidation strategies
- Monitor cache hit rates

### **Memory Management**
- Avoid memory leaks in long-running processes
- Use appropriate data structures for the task
- Implement proper resource cleanup
- Monitor memory usage in production

### **Concurrency**
- Use goroutines for I/O-bound operations
- Implement proper synchronization mechanisms
- Avoid race conditions
- Use context for cancellation and timeouts

---

## 🔄 **Git Workflow**

### **Branch Strategy**
- **main**: Production-ready code only
- **develop**: Integration branch for features
- **feature/**: Individual feature development
- **hotfix/**: Critical production fixes
- **release/**: Release preparation

### **Branch Naming**
- Use descriptive names: `feature/otp-rate-limiting`
- Include ticket numbers: `feature/AUTH-123-jwt-refresh`
- Use kebab-case: `hotfix/security-vulnerability-fix`

### **Commit Guidelines**
- Use conventional commit format: `type(scope): description`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Keep commits atomic and focused
- Write clear, descriptive commit messages
- Reference issue numbers when applicable

### **Pull Request Process**
- Create PR from feature branch to develop
- Include clear description of changes
- Add relevant reviewers
- Ensure all tests pass
- Update documentation if needed
- Squash commits before merging

### **Code Review Standards**
- Review for correctness, security, and performance
- Check adherence to coding standards
- Verify test coverage
- Ensure documentation is updated
- Provide constructive feedback

---

## 📚 **Documentation Standards**

### **Code Documentation**
- Document all public APIs
- Include usage examples for complex functions
- Document configuration options
- Maintain up-to-date README files

### **API Documentation**
- Use OpenAPI/Swagger for REST APIs
- Include request/response examples
- Document error codes and messages
- Provide authentication requirements

### **Architecture Documentation**
- Keep architecture diagrams current
- Document design decisions and rationale
- Include deployment instructions
- Maintain troubleshooting guides

---

## 🚀 **Deployment Guidelines**

### **Environment Management**
- Use consistent environments: dev, staging, production
- Implement infrastructure as code
- Use environment-specific configurations
- Maintain environment parity

### **Release Process**
- Use semantic versioning (MAJOR.MINOR.PATCH)
- Tag releases in git
- Maintain changelog
- Implement rollback procedures

### **Monitoring and Logging**
- Implement structured logging
- Use appropriate log levels
- Monitor key performance metrics
- Set up alerting for critical issues

### **Database Migrations**
- Use versioned migration scripts
- Test migrations on staging first
- Implement rollback procedures
- Document breaking changes

---

## 🔧 **Development Tools**

### **Required Tools**
- **Go**: Version 1.24.4 or later
- **Docker**: For containerization
- **Git**: For version control
- **Make**: For build automation

### **Recommended IDE Setup**
- Use Go language server (gopls)
- Configure automatic formatting on save
- Enable linting and static analysis
- Set up debugging configuration

### **Code Quality Tools**
- **golangci-lint**: For comprehensive linting
- **gosec**: For security analysis
- **go vet**: For static analysis
- **gocyclo**: For complexity analysis

---

## 📋 **Code Review Checklist**

### **Functionality**
- [ ] Code solves the intended problem
- [ ] Edge cases are handled properly
- [ ] Error handling is comprehensive
- [ ] Business logic is correct

### **Code Quality**
- [ ] Code follows naming conventions
- [ ] Functions are appropriately sized
- [ ] Code is readable and well-structured
- [ ] No code duplication

### **Testing**
- [ ] Unit tests cover new functionality
- [ ] Integration tests are included where appropriate
- [ ] Tests are meaningful and not just for coverage
- [ ] Test names are descriptive

### **Security**
- [ ] Input validation is implemented
- [ ] Sensitive data is not logged
- [ ] Authentication/authorization is correct
- [ ] No security vulnerabilities introduced

### **Performance**
- [ ] No obvious performance issues
- [ ] Database queries are optimized
- [ ] Memory usage is reasonable
- [ ] Caching is used appropriately

### **Documentation**
- [ ] Public APIs are documented
- [ ] Complex logic is explained
- [ ] README is updated if needed
- [ ] API documentation is current

---

## 🚨 **Common Anti-Patterns to Avoid**

### **Code Smells**
- **God Objects**: Classes that do too much
- **Long Parameter Lists**: Functions with too many parameters
- **Deep Nesting**: Excessive if/else or loop nesting
- **Magic Numbers**: Unexplained numeric constants

### **Architecture Anti-Patterns**
- **Circular Dependencies**: Modules depending on each other
- **Tight Coupling**: Components too dependent on each other
- **Anemic Domain Model**: Domain objects with no behavior
- **Service Locator**: Using global service registry

### **Testing Anti-Patterns**
- **Testing Implementation Details**: Testing how instead of what
- **Fragile Tests**: Tests that break with minor changes
- **Slow Tests**: Tests that take too long to run
- **Test Interdependence**: Tests that depend on other tests

---

## 📈 **Continuous Improvement**

### **Regular Reviews**
- Conduct monthly architecture reviews
- Review and update guidelines quarterly
- Analyze code quality metrics regularly
- Gather feedback from development team

### **Learning and Development**
- Stay updated with Go best practices
- Learn from industry standards and patterns
- Participate in code review discussions
- Share knowledge through documentation

### **Metrics and Monitoring**
- Track code quality metrics
- Monitor test coverage trends
- Measure deployment frequency and success rate
- Analyze bug reports and root causes

---

## 🎯 **Enforcement**

### **Automated Checks**
- Pre-commit hooks for formatting and linting
- CI/CD pipeline checks for tests and quality
- Automated security scanning
- Code coverage reporting

### **Manual Reviews**
- Mandatory code reviews for all changes
- Architecture review for significant changes
- Security review for sensitive modifications
- Performance review for critical paths

### **Consequences**
- Failed checks block merge to main branches
- Repeated violations require additional training
- Significant violations require architecture review
- Security violations require immediate attention

---

## 📞 **Support and Questions**

### **Getting Help**
- Consult this document first
- Ask team members for clarification
- Create issues for guideline improvements
- Escalate architectural questions to tech leads

### **Guideline Updates**
- Propose changes through pull requests
- Discuss significant changes in team meetings
- Update documentation when practices evolve
- Communicate changes to all team members

---

**Remember**: These guidelines exist to ensure code quality, maintainability, and team productivity. When in doubt, prioritize clarity and simplicity over cleverness.