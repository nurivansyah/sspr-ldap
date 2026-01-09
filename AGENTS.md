# AGENTS.md

## Project Overview
SSPR-LDAP is a Go web application providing LDAP-backed self-service password reset and authentication. It authenticates users against LDAP, enables password changes with AD-compatible encoding, and offers a session-based web UI.

Purpose: Enable users to reset passwords securely via LDAP integration, reducing IT support load.

## Tech Stack
- Go 1.25.4 (main.go:1)
- LDAP: go-ldap/ldap/v3 (go.mod:6)
- Sessions: gorilla/sessions (go.mod:7)
- Environment: godotenv (go.mod:8)
- Templates: Go's html/template (infra/template/engine.go:1)

## Key Directories and Purposes
- config/: Configuration loading from environment variables (config/config.go:1)
- domain/: Domain models (e.g., User, Credentials) (domain/user.go:3)
- handlers/: HTTP request handlers (handlers/auth-handler.go:13, handlers/user-handler.go:1)
- services/: Business logic (services/auth-service.go:9, services/user-service.go:1)
- infra/: Infrastructure adapters (infra/ldap/repository.go:13, infra/session/store.go:1, infra/template/engine.go:1)
- ports/: Interface definitions (ports/auth.go:5)
- templates/: HTML templates (templates/login.html:1)

## Essential Build/Test Commands
- Build: `go build -o sspr-ldap ./` (README.md:50)
- Run locally: `go run ./` (README.md:47)
- Docker build: `docker build -t sspr-ldap:local .` (README.md:61)
- Docker run: `docker compose up --build` (README.md:91)
- Test: `go test ./...` (inferred; no explicit tests found, but standard Go testing)

## Adding New Features or Fixing Bugs

**IMPORTANT**: When you work on a new feature or bug, create a git branch first.
Then work on changes in that branch for the remainder of the session.

## Additional Documentation
For specialized topics, check these files:
- docs/architectural_patterns.md: Architectural patterns, design decisions, and conventions
- docs/ldap-integration.md: LDAP configuration and usage
- docs/session-management.md: Session handling details
- docs/deployment.md: Docker and production deployment