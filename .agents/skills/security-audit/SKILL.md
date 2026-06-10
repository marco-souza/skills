---
name: security-audit
description: >
  Perform standalone security audits covering OWASP Top 10, dependency scanning,
  secret detection, input validation, and authentication review.
  Use when: the user asks for a security review, vulnerability assessment, audit,
  or wants to check for security issues in code or dependencies.
  Do NOT use when: the user wants a general code review (use pr-review),
  is asking about deploying to production, or wants penetration testing guidance.
---

# Security Audit

Conduct systematic security reviews of codebases to identify vulnerabilities, misconfigurations, and compliance gaps.

## Prerequisites

Ensure the following tools are available:

```bash
# Check for security tools (install if missing)
command -v gitleaks >/dev/null 2>&1 || echo "Install gitleaks: brew install gitleaks"
command -v trivy >/dev/null 2>&1 || echo "Install trivy: brew install trivy"
command -v semgrep >/dev/null 2>&1 || echo "Install semgrep: brew install semgrep"
```

**Optional but recommended:**
- `npm audit` / `yarn audit` (for Node.js projects)
- `pip-audit` (for Python projects)
- `bundler-audit` (for Ruby projects)

## When to Use

- Security review before production deployment
- Periodic security audits of existing codebases
- After adding new authentication or authorization logic
- When handling sensitive data (PII, financial, health)
- After a security incident or breach report
- When onboarding to a new codebase with security concerns

## When NOT to Use

- For general code quality reviews (use `pr-review`)
- When the user wants penetration testing guidance
- For infrastructure security (cloud config, network)
- When scope is unclear and needs brainstorming first

## Audit Workflow

### Phase 1: Reconnaissance (5-10 min)

Gather project context:

```bash
# Identify project type and stack
cat package.json 2>/dev/null | head -20
cat requirements.txt 2>/dev/null | head -20
cat go.mod 2>/dev/null | head -20

# Find configuration files
find . -maxdepth 3 -name "*.env*" -o -name "config.*" -o -name ".env*" | grep -v node_modules | head -20

# Check for security-related files
ls -la .github/dependabot.yml .snyk 2>/dev/null
ls -la SECURITY.md 2>/dev/null
```

### Phase 2: OWASP Top 10 Checklist

Review each category systematically:

#### A01:2021 - Broken Access Control

```bash
# Find authentication/authorization middleware
grep -r "auth\|permission\|role\|admin\|middleware" --include="*.ts" --include="*.js" --include="*.py" | grep -v node_modules | head -20

# Check for unprotected routes
grep -r "router\.\(get\|post\|put\|delete\)" --include="*.ts" --include="*.js" | grep -v "auth\|protect\|secure" | head -20

# Look for IDOR vulnerabilities
grep -r "req\.params\|req\.query\|params\[" --include="*.ts" --include="*.js" | head -15
```

**Checklist:**
- [ ] All endpoints have proper authorization checks
- [ ] Users can only access their own resources (no IDOR)
- [ ] Admin routes are protected and audited
- [ ] CORS is properly configured (not `*` in production)
- [ ] Directory traversal is prevented

#### A02:2021 - Cryptographic Failures

```bash
# Find hardcoded secrets or weak crypto
grep -rn "secret\|password\|api.key\|token" --include="*.ts" --include="*.js" --include="*.py" | grep -v "node_modules\|test\|spec" | head -20

# Check for deprecated algorithms
grep -rn "md5\|sha1\|des\|rc4" --include="*.ts" --include="*.js" --include="*.py" | grep -v node_modules | head -10

# Find encryption usage
grep -rn "encrypt\|decrypt\|hash\|cipher" --include="*.ts" --include="*.js" --include="*.py" | grep -v node_modules | head -15
```

**Checklist:**
- [ ] No hardcoded secrets in source code
- [ ] TLS/HTTPS enforced for all connections
- [ ] Strong algorithms used (AES-256, bcrypt, Argon2)
- [ ] Sensitive data encrypted at rest
- [ ] Proper key management (not in code)

#### A03:2021 - Injection

```bash
# Find SQL queries (look for string concatenation)
grep -rn "SELECT\|INSERT\|UPDATE\|DELETE" --include="*.ts" --include="*.js" --include="*.py" | grep -E "\+|`|\"|\$" | head -20

# Find command execution
grep -rn "exec\|spawn\|system\|eval\|Function(" --include="*.ts" --include="*.js" --include="*.py" | grep -v node_modules | head -15

# Find template literal usage in queries
grep -rn "\`\s*SELECT\|\`\s*INSERT" --include="*.ts" --include="*.js" | head -10
```

**Checklist:**
- [ ] Parameterized queries used (no string concatenation)
- [ ] Input validation before processing
- [ ] ORM used instead of raw queries
- [ ] No `eval()` or dynamic code execution
- [ ] Command injection prevented

#### A04:2021 - Insecure Design

```bash
# Look for business logic patterns
grep -rn "TODO.*security\|FIXME.*security\|HACK\|XXX" --include="*.ts" --include="*.js" --include="*.py" | head -10

# Find rate limiting configuration
grep -rn "rate.limit\|throttle\|rateLimit" --include="*.ts" --include="*.js" | head -10

# Check for security headers
grep -rn "helmet\|csp\|Content-Security-Policy\|X-Frame" --include="*.ts" --include="*.js" | head -10
```

**Checklist:**
- [ ] Rate limiting implemented on sensitive endpoints
- [ ] Business logic cannot be abused
- [ ] Proper error handling (no info leakage)
- [ ] Security headers configured

#### A05:2021 - Security Misconfiguration

```bash
# Find configuration files
find . -maxdepth 3 -name "*.config.*" -o -name "docker-compose*" -o -name "Dockerfile" | grep -v node_modules | head -15

# Check for default credentials
grep -rn "admin\|password\|default" --include="*.yml" --include="*.yaml" --include="*.json" | grep -v node_modules | head -15

# Find debug mode flags
grep -rn "debug.*true\|DEBUG\|verbose" --include="*.ts" --include="*.js" --include="*.py" | grep -v node_modules | head -10
```

**Checklist:**
- [ ] Debug mode disabled in production
- [ ] Default credentials changed
- [ ] Unnecessary features/services disabled
- [ ] Error messages don't expose internals
- [ ] Security headers configured

#### A06:2021 - Vulnerable and Outdated Components

```bash
# Check for outdated dependencies
npm audit 2>/dev/null || echo "Not a Node.js project"
yarn audit 2>/dev/null || echo "Not a Yarn project"
pip-audit 2>/dev/null || echo "Not a Python project"

# Look for known vulnerable packages
grep -rn "lodash.*4\.17\.1[0-5]\|express.*4\.17\.\|minimist.*1\.2" package.json 2>/dev/null | head -10
```

**Checklist:**
- [ ] All dependencies updated to latest patch versions
- [ ] No known vulnerabilities (check audit output)
- [ ] Dependencies pinned to specific versions
- [ ] No unused/abandoned dependencies

#### A07:2021 - Identification and Authentication Failures

```bash
# Find authentication logic
grep -rn "login\|authenticate\|password\|bcrypt\|jwt" --include="*.ts" --include="*.js" --include="*.py" | grep -v node_modules | head -20

# Check for session management
grep -rn "session\|cookie\|jwt\|token" --include="*.ts" --include="*.js" | grep -v node_modules | head -15

# Look for brute force protection
grep -rn "attempt\|lockout\|brute\|rate.limit" --include="*.ts" --include="*.js" | head -10
```

**Checklist:**
- [ ] Passwords hashed with bcrypt/Argon2 (not MD5/SHA1)
- [ ] Multi-factor authentication available
- [ ] Session timeout configured
- [ ] Account lockout after failed attempts
- [ ] JWT tokens properly validated and expired

#### A08:2021 - Software and Data Integrity Failures

```bash
# Find dependency lock files
ls -la package-lock.json yarn.lock poetry.lock go.sum 2>/dev/null

# Check for CI/CD integrity
ls -la .github/workflows/ 2>/dev/null | head -10

# Find auto-update mechanisms
grep -rn "auto.update\|self.update" --include="*.ts" --include="*.js" --include="*.py" | head -10
```

**Checklist:**
- [ ] Lock files committed and used
- [ ] CI/CD pipeline integrity verified
- [ ] Code signing for releases
- [ ] No unsigned/untrusted sources

#### A09:2021 - Security Logging and Monitoring Failures

```bash
# Find logging patterns
grep -rn "console.log\|logger\.\|log\.\|winston\|pino" --include="*.ts" --include="*.js" --include="*.py" | grep -v node_modules | head -20

# Check for security event logging
grep -rn "audit\|security\|unauthorized\|forbidden" --include="*.ts" --include="*.js" | grep -i "log" | head -15
```

**Checklist:**
- [ ] Security events logged (login failures, access denied)
- [ ] Logs not exposing sensitive data
- [ ] Centralized logging configured
- [ ] Alerting for suspicious activity

#### A10:2021 - Server-Side Request Forgery (SSRF)

```bash
# Find HTTP client usage
grep -rn "fetch\|axios\|request\|http\." --include="*.ts" --include="*.js" --include="*.py" | grep -v node_modules | head -20

# Check for URL validation
grep -rn "url\|href\|redirect\|proxy" --include="*.ts" --include="*.js" | grep -v node_modules | head -15
```

**Checklist:**
- [ ] User-supplied URLs validated against allowlist
- [ ] Internal network access restricted
- [ ] Response content not returned to user without validation
- [ ] DNS rebinding protection

## Secret Detection

### Run Automated Scanning

```bash
# Using gitleaks (recommended)
gitleaks detect --source . --verbose

# Or using trufflehog
trufflehog filesystem --directory .

# Manual pattern search
grep -rn --include="*.ts" --include="*.js" --include="*.py" --include="*.env*" \
  -E "(AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}" . | head -10
```

### Common Secret Patterns

Look for these patterns in source code:

```
# AWS Keys
AKIA[0-9A-Z]{16}

# GitHub/GitLab Tokens
ghp_[A-Za-z0-9]{36}
glpat-[A-Za-z0-9\-]{20,}

# Slack Tokens
xox[bpsar]-[A-Za-z0-9\-]+

# Private Keys
-----BEGIN (RSA |EC )?PRIVATE KEY-----

# Generic patterns
password\s*[:=]\s*['"][^'"]+['"]
secret\s*[:=]\s*['"][^'"]+['"]
api[_-]?key\s*[:=]\s*['"][^'"]+['"]
```

### Verify Secret Handling

```bash
# Check .gitignore for sensitive files
cat .gitignore | grep -E "\.env|secret|credential|key"

# Verify secrets are in environment variables, not code
grep -rn "process\.env\|os\.environ\|os\.getenv" --include="*.ts" --include="*.js" --include="*.py" | head -15
```

## Input Validation

### Check Validation Patterns

```bash
# Find validation libraries
grep -rn "zod\|joi\|yup\|ajv\|class-validator\|validator" --include="*.ts" --include="*.js" | head -10

# Find manual validation
grep -rn "\.trim\(\)\|\.length\|\.match\|\.test\|regex" --include="*.ts" --include="*.js" | grep -v node_modules | head -15

# Find sanitization
grep -rn "sanitize\|escape\|encode\|xss\|DOMPurify" --include="*.ts" --include="*.js" | head -10
```

### Validation Checklist

- [ ] All user inputs validated on server side
- [ ] Schema validation used (Zod, Joi, etc.)
- [ ] Output encoding applied (XSS prevention)
- [ ] File upload validation (type, size)
- [ ] URL validation for redirects

## Authentication Review

### Analyze Auth Implementation

```bash
# Find auth-related code
find . -type f \( -name "*auth*" -o -name "*login*" -o -name "*session*" \) | grep -v node_modules | head -15

# Check password handling
grep -rn "bcrypt\|argon2\|scrypt\|hash" --include="*.ts" --include="*.js" --include="*.py" | head -10

# Find JWT usage
grep -rn "jwt\|jsonwebtoken\|jose" --include="*.ts" --include="*.js" --include="*.py" | head -15
```

### Authentication Checklist

- [ ] Passwords stored with strong hashing (bcrypt/Argon2)
- [ ] JWT secrets are strong and rotated
- [ ] Token expiration properly configured
- [ ] Refresh token rotation implemented
- [ ] Password reset flow secure
- [ ] Account enumeration prevented

## Reporting

### Generate Audit Report

Create a structured report:

```markdown
# Security Audit Report

**Date:** YYYY-MM-DD
**Scope:** [application name and version]
**Auditor:** [AI Agent]

## Executive Summary

- Critical: X findings
- High: X findings
- Medium: X findings
- Low: X findings
- Informational: X findings

## Findings

### [CRITICAL] Finding Title

**Category:** OWASP A0X
**File(s):** path/to/file.ts:line
**Description:** ...
**Impact:** ...
**Remediation:** ...
**Reference:** [CWE-XXX] [link]

## OWASP Top 10 Compliance

| Category | Status | Notes |
|----------|--------|-------|
| A01: Broken Access Control | Pass/Fail | ... |
| A02: Cryptographic Failures | Pass/Fail | ... |
| ... | ... | ... |

## Recommendations

1. Priority-ordered remediation steps
2. Quick wins vs long-term fixes
3. Tools/processes to prevent recurrence
```

### Severity Classification

- **Critical:** Immediate exploitation risk, data breach potential
- **High:** Significant vulnerability, requires prompt fix
- **Medium:** Exploitable under certain conditions
- **Low:** Defense-in-depth improvement
- **Informational:** Best practice recommendation

## Example Audit

```markdown
# Security Audit Report

**Date:** 2024-01-15
**Scope:** api.example.com v2.1.0

## Executive Summary

- Critical: 1 findings
- High: 2 findings
- Medium: 3 findings
- Low: 2 findings

## Findings

### [CRITICAL] SQL Injection in User Search

**Category:** OWASP A03:2021 - Injection
**File:** src/api/users.ts:45
**Description:** User search endpoint concatenates input directly into SQL query.
**Impact:** Attacker can extract entire database, modify data, or execute commands.
**Remediation:** Use parameterized queries or ORM.
```typescript
// Before (vulnerable)
const query = `SELECT * FROM users WHERE name = '${name}'`;

// After (secure)
const user = await db.user.findUnique({ where: { name } });
```

### [HIGH] Hardcoded JWT Secret

**Category:** OWASP A02:2021 - Cryptographic Failures
**File:** src/config/auth.ts:12
**Description:** JWT secret hardcoded in source code.
**Impact:** Token forgery if source code is compromised.
**Remediation:** Move to environment variable, rotate immediately.
```

## Edge Cases

### Legacy Code

For older codebases without modern tooling:

1. Focus on injection and authentication first
2. Add manual validation before automated scanning
3. Create a remediation roadmap

### Third-Party Dependencies

When dependencies have vulnerabilities:

1. Check if vulnerability is exploitable in your usage
2. Look for patches or alternatives
3. Document accepted risk if not exploitable

### Microservices

For distributed systems:

1. Audit each service independently
2. Check inter-service authentication
3. Verify network policies

## Troubleshooting

### Tool Not Found

```bash
# Install missing tools
brew install gitleaks trivy semgrep

# Or use Docker
docker run --rm -v $(pwd):/src ghcr.io/gitleaks/gitleaks:latest detect -s /src
```

### Too Many Findings

If audit produces overwhelming results:

1. Filter by severity (Critical and High first)
2. Focus on user-facing endpoints
3. Prioritize authentication and data handling code

### False Positives

Common false positives to ignore:

- Test files with mock secrets
- Documentation examples
- Environment variable references (not actual values)

## Best Practices

1. **Automate early** — Integrate security scanning into CI/CD
2. **Shift left** — Review security during design, not just before release
3. **Layer defense** — Multiple security controls, not just one
4. **Document decisions** — Record accepted risks and mitigations
5. **Regular cadence** — Schedule periodic audits, not just pre-release
6. **Scope clearly** — Define what's in/out of audit scope
