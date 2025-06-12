# SSHTrust Codebase Review

## Current State Summary

**Branch:** ACLs (3 commits ahead of main)  
**Build Status:** ✅ Builds successfully  
**Test Status:** ✅ All tests pass (8 test files, 33 total Go files)

## Architecture Assessment

**Structure:** Well-organized Go project with clear separation of concerns:
- CLI layer (cmd/) using Cobra framework
- HTTP server (internal/server/) using Echo framework  
- Core business logic (pkg/) with interfaces for extensibility
- Good test coverage across key components

**Code Quality:** Generally solid with some areas for improvement:
- Clean interfaces (CAStore, UserList) enabling dependency injection
- Proper concurrency handling with sync.RWMutex in memory store
- Consistent error handling patterns
- Some minor typos in struct names (`InMemortCaStore`)

## Test Coverage Analysis

**Coverage:** 24% (8/33 files have tests)
- ✅ Core packages well-tested: cert, certStore, handlers
- ✅ Integration tests via launch-server.sh
- ❌ Missing: cmd/ package tests, auth package tests

## Security Review

**Resolved Issues:**
- ✅ JWT secret: Now loads from JWT_SECRET environment variable with secure random fallback
- ✅ JWT secret validation: Minimum 32-character length requirement with base64 support

**Remaining Critical Issues:**
- 🔴 Token file permissions: JWT tokens written with insecure permissions (internal/client/auth.go:85-97)
- 🔴 Missing input validation: URL parameters not validated in handlers (pkg/handlers/ca.go:34, sign.go:26)
- 🔴 Hardcoded server URLs: `http://localhost:8080` throughout client code

**Medium Priority Security Issues:**
- 🟡 JWT tokens valid for 72 hours (pkg/auth/user.go:71) - should be configurable
- 🟡 Hardcoded server port `:8080` (internal/server/main.go:20)
- 🟡 No rate limiting on authentication endpoints
- 🟡 Missing request size limits

**Good Practices:**
- ✅ bcrypt password hashing
- ✅ JWT-based authentication with secure secret management
- ✅ Optional no-auth mode for development
- ✅ Proper mutex usage for concurrent access
- ✅ Cryptographically secure random JWT secret generation

## ACL Implementation Status

**Current Progress:** Very early stage
- ✅ Basic ACL struct defined (pkg/auth/acl.go)
- ✅ Example ACL rules with wildcard and specific CA permissions
- ❌ No integration with actual authorization logic
- ❌ No user/group resolution system
- ❌ No enforcement in handlers or middleware

## Key Issues Identified

**High Priority:**
1. **Security:** Fix token file permissions (0600) and input validation
2. **Configuration:** Remove hardcoded URLs and make server configurable
3. **ACL Implementation:** Placeholder code not connected to authorization flow
4. **Dockerfile:** Broken - switched from Ubuntu to Rocky Linux but kept apt-get commands

**Medium Priority:**
1. Add rate limiting and request size limits for security
2. Make JWT expiration configurable via environment variables
3. Typo in `InMemortCaStore` (should be `InMemoryCAStore`)
4. Missing tests for auth and cmd packages
5. No configuration management system

**Low Priority:**
1. Clean up stale TODO comments
2. Standardize logging approach (mix of log types)
3. Add graceful shutdown handling in server

## Recommendations

1. **Immediate:** Fix remaining security issues (file permissions, input validation, configuration)
2. **Next:** Complete ACL implementation with proper middleware integration
3. **Future:** Add comprehensive configuration management, improve test coverage

## Recent Improvements

**JWT Security Enhancement (Fixed):**
- Implemented environment variable support for JWT_SECRET
- Added cryptographically secure random secret generation as fallback
- Added input validation for secret length and base64 decoding
- Updated documentation with configuration instructions

## Conclusion

The codebase has a solid foundational architecture with recent security improvements to JWT handling. The remaining high-priority security issues (file permissions, input validation, configuration hardcoding) should be addressed before production deployment. The ACL feature remains incomplete and needs integration with the authorization flow.