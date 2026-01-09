# docs/architectural_patterns.md

## Architectural Patterns

### Hexagonal Architecture (Ports and Adapters)
The codebase follows hexagonal architecture to decouple core logic from external dependencies.
- Ports (interfaces) define contracts (ports/auth.go:5-13)
- Adapters (infra) implement ports (infra/ldap/repository.go:13-21, infra/session/store.go:1)
- Domain contains business models (domain/user.go:3-21)
- Services handle application logic (services/auth-service.go:9-30)
- Handlers manage presentation (handlers/auth-handler.go:13-25)

### Dependency Injection
Dependencies are injected via constructors to enable testability and flexibility.
- Services receive repositories (services/auth-service.go:13-17)
- Handlers receive services and infra (handlers/auth-handler.go:19-25)
- Infra components receive config (infra/ldap/repository.go:17-21)

### Repository Pattern
Data access is abstracted through repository interfaces.
- Auth operations (ports/auth.go:6-8, infra/ldap/repository.go:23-...)
- User operations (ports/auth.go:11-13, services/user-service.go:1)

### Service Layer Pattern
Business logic is encapsulated in service structs.
- Auth service validates and delegates (services/auth-service.go:19-30)
- User service handles password changes (services/user-service.go:1)

### Handler Pattern for HTTP
HTTP handlers process requests and responses.
- Auth handler manages login/logout (handlers/auth-handler.go:27-...)
- User handler manages dashboard/password change (handlers/user-handler.go:1)

### Configuration Management
Environment variables loaded via config package.
- Centralized loading (config/config.go:1)
- Used throughout main.go (main.go:31-39)

### Session Management
Gorilla sessions for user state.
- Store initialization (infra/session/store.go:1, main.go:34)
- Authentication checks (handlers/auth-handler.go:28-33)

### Error Handling
Errors propagated up the call stack.
- Service validation (services/auth-service.go:20-27)
- Handler responses (handlers/auth-handler.go:50-...)

### Graceful Shutdown
Server shuts down cleanly on signals (main.go:74-86)