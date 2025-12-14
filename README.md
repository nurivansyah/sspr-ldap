sspr-ldap
========

A small Go web app that provides LDAP-backed self-service password reset and authentication helpers.

Quick overview
--------------
- Language: Go
- Primary functions: authenticate users against LDAP, change passwords (AD-compatible unicodePwd encoding), simple session-based web UI.

Prerequisites
-------------
- Go 1.25+ installed to build/run locally
- Docker (optional) to build and run container image
- LDAP server accessible from runtime with a service account (if required)

Configuration
-------------
The app reads configuration from environment variables. Use a `.env` file for local development (do NOT commit real secrets).
Example variables (see `.env.example`):

- `PORT` – HTTP port (default `8080`)
- `SESSION_KEY` – session encryption key (rotate regularly, min 32 bytes)
- `LDAP_SERVER` – LDAP host
- `LDAP_PORT` – LDAP port (eg. 636 for LDAPS)
- `LDAP_BASE_DN` – base DN used when searching for users
- `LDAP_BIND_DN` – service account DN used for binds (optional)
- `LDAP_BIND_PASSWORD` – service account password
- `LDAP_USER_FILTER` – search filter, e.g. `(userPrincipalName=%s)`
- `LDAP_USE_TLS` – `true`/`false`
- `LDAP_TLS_SKIP_VERIFY` – `true`/`false` (avoid true in production)
- `SESSION_COOKIE_SECURE` - session secure option

Run locally with Go
-------------------
1. Copy `.env.example` to `.env` and populate secrets (do NOT commit):

```bash
cp .env.example .env
# edit .env
```

2. Run the app:

```bash
# using go run
go run ./

# or build and run binary
go build -o sspr-ldap ./
./sspr-ldap
```

3. Open http://localhost:8080

Run in Docker
-------------
Build image:

```bash
docker build -t sspr-ldap:local .
```

Run container (pass environment variables; do NOT mount a file with secrets into the image):

```bash
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e SESSION_KEY="your-session-key" \
  -e LDAP_SERVER="ldap.example.local" \
  -e LDAP_PORT=636 \
  -e LDAP_BIND_DN="cn=svc,dc=example,dc=com" \
  -e LDAP_BIND_PASSWORD="secret" \
  sspr-ldap:local
```

Run with Docker Compose
-----------------------
This repository includes a `docker-compose.yml` to build and run the app using an `.env` file for configuration.

1. Create `.env` from `.env.example` and fill values (do NOT commit):

```bash
cp .env.example .env
# edit .env
```

2. Build and start the app:

```bash
docker compose up --build
```

Or run detached:

```bash
docker compose up --build -d
```

3. Stop and remove containers:

```bash
docker compose down
```

Notes:
- `docker-compose.yml` maps the container port to the host using the `PORT` variable from `.env` (defaults to `8080`).
- The compose file references `.env` using `env_file: - .env`. Ensure the file exists locally and contains no production secrets.
- A healthcheck is commented out in `docker-compose.yml` because it requires a probe binary (e.g., `wget` or `curl`) in the runtime image. If you enable a healthcheck, ensure the runtime image contains the necessary tool or add one to the `Dockerfile`.
