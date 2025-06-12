# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SSHTrust is an SSH Certificate Authority (CA) management system that provides:
- HTTP REST API for creating and managing SSH CAs
- SSH public key signing using in-memory CAs  
- CLI client for interacting with the server
- JWT-based authentication with optional no-auth mode
- Swagger API documentation at `/swagger/index.html`

## Development Commands

### Build and Test
```bash
make                    # Run full pipeline: test, gen, build, launchtest
make build             # Build the sshtrust binary
make test              # Run all tests with go test -fullpath ./...
make gen               # Generate Swagger documentation
make launchtest        # Run integration test script
```

### Server Operations
```bash
./sshtrust serve                    # Start server with authentication
./sshtrust serve --no-auth          # Start server without authentication
```

### CA Management
```bash
./sshtrust ca new -n <name> -p <principal> -t <type>    # Create new CA
./sshtrust ca list                                       # List all CAs
./sshtrust ca get <name>                                # Get CA details
./sshtrust sign -n <name> --ttl <minutes> -p <principal> -k "<pubkey>"  # Sign public key
```

### Authentication
```bash
./sshtrust register -u <username>      # Register new user
./sshtrust login -u <username>         # Login and get JWT token
```

## Environment Variables

### JWT Configuration
```bash
export JWT_SECRET="your-secret-key-here"    # JWT signing secret (minimum 32 characters)
```

**JWT_SECRET**: Used for signing and validating JWT tokens in authenticated mode.
- Can be a plain string (minimum 32 characters) or base64-encoded value
- If not provided, a cryptographically secure random secret is generated for each server session
- For production deployments, always set this to a strong, persistent secret

## Architecture

### Core Components

**CLI Layer (`cmd/`)**
- Cobra-based CLI with subcommands for server, CA operations, auth
- Entry point is `main.go` → `cmd.Execute()`

**HTTP Server (`internal/server/`)**
- Echo framework with JWT middleware (optional)
- Swagger documentation integration
- Routes: `/CA/*`, `/login`, `/register`

**Handlers (`pkg/handlers/`)**
- REST API controllers for CA operations
- App struct holds CAStore interface for dependency injection

**Certificate Management (`pkg/cert/`)**
- SSH CA creation and key signing logic
- Support for RSA, ECDSA, and Ed25519 key types
- Certificate validity and principal management

**Storage Layer (`pkg/certStore/`)**
- CAStore interface for pluggable storage backends
- Current implementation: InMemoryCAStore
- Planned: SQL, PostgreSQL, etcd backends

**Authentication (`pkg/auth/`)**
- JWT token generation and validation
- User registration/login with InMemoryUserList
- ACL system (in development)

### Data Flow
1. CLI commands → HTTP client → Server handlers
2. Handlers → CAStore interface → In-memory storage
3. Certificate operations use `golang.org/x/crypto/ssh` package
4. JWT tokens required for authenticated endpoints (unless `--no-auth`)

### Testing
- Unit tests alongside source files (`*_test.go`)
- Integration testing via `launch-server.sh` script
- Test SSH server in Docker for end-to-end validation

### Current Development Focus
Working on ACL branch - implementing access control lists for CA operations.