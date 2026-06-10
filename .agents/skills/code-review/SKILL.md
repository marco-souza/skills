---
name: code-review
description: >
  Review your own local code changes before pushing or creating a PR.
  Use when: the user wants to self-review staged or unstaged changes, check code quality before committing,
  or get feedback on work-in-progress. Covers readability, patterns, error handling, performance, and security.
  Do NOT use when: reviewing an existing pull request (use pr-review), reviewing code you didn't write,
  or when changes are already pushed to remote.
---

# Code Review (Pre-Push)

Perform a thorough self-review of local code changes before pushing. Catch issues early when they're cheap to fix.

## When to Use

- Before committing significant changes
- Before pushing to a shared branch
- When you want a second opinion on your own code
- After completing a feature or bug fix
- Before opening a PR

## When NOT to Use

- When reviewing someone else's PR (use `pr-review`)
- When changes are already pushed and live
- For one-line typo fixes or trivial changes

## Prerequisites

No special tools required. Optionally, ensure you have a diff viewer available:

```bash
# Verify git is available
git --version

# Optional: check if diffstat is available
which diffstat
```

## Workflow

### 1. Identify Changes to Review

```bash
# Review all uncommitted changes (staged + unstaged)
git diff HEAD

# Or review only staged changes (ready to commit)
git diff --cached

# Review specific files
git diff HEAD -- path/to/file.go path/to/other.js

# See what files changed
git diff --stat HEAD
```

### 2. Review by File

Work through each changed file systematically:
1. Start with the most critical/business logic files
2. Then review utility/helper functions
3. Finally check tests and documentation

### 3. Apply the Checklist

Use the review checklist below to evaluate each file.

## Review Checklist

### Readability
- [ ] Code is self-documenting; comments explain **why**, not **what**
- [ ] Variable and function names are clear and descriptive
- [ ] No magic numbers or unexplained constants
- [ ] Consistent formatting and style
- [ ] Appropriate use of whitespace and grouping

### Patterns & Consistency
- [ ] Follows existing codebase patterns
- [ ] Consistent error handling approach
- [ ] Similar problems solved similarly (no reinventing)
- [ ] Naming conventions match the project style
- [ ] Architecture boundaries respected

### Error Handling
- [ ] All error paths handled explicitly
- [ ] Errors propagated with context
- [ ] No swallowed exceptions
- [ ] Resource cleanup on failure (defer, finally, etc.)
- [ ] User-facing errors are helpful, not cryptic

### Performance
- [ ] No unnecessary allocations or copies
- [ ] Appropriate data structures for use case
- [ ] No N+1 queries or repeated expensive operations
- [ ] Pagination for large datasets
- [ ] No blocking in async/hot paths

### Security
- [ ] No hardcoded secrets or credentials
- [ ] Input validation at trust boundaries
- [ ] No injection vulnerabilities (SQL, command, XSS)
- [ ] Sensitive data not logged
- [ ] Authentication/authorization checks in place

### Testing
- [ ] Changes covered by existing tests
- [ ] New tests added for new functionality
- [ ] Edge cases considered
- [ ] Tests are meaningful (not just coverage)

### Completeness
- [ ] No TODO/FIXME/HACK left unintentionally
- [ ] No debug code or console.log statements
- [ ] No commented-out code
- [ ] Documentation updated if needed
- [ ] No unrelated changes mixed in

## Review Process

For each file in the diff:

```bash
# View the full file for context (not just the diff)
git show HEAD:path/to/file

# Or open in your editor
git diff HEAD -- path/to/file
```

Ask yourself:
1. **Would I understand this in 6 months?**
2. **Is this the simplest way to solve the problem?**
3. **Did I handle failure cases?**
4. **Will this scale?**
5. **Is this secure?**

## Severity Levels

When flagging issues, categorize by severity:

| Level | Meaning | Action |
|-------|---------|--------|
| **[blocker]** | Must fix before push | Security flaw, data loss risk, broken logic |
| **[warning]** | Should fix now | Code smell, poor pattern, missing error handling |
| **[nit]** | Nice to fix | Style, naming, minor optimization |

## Output Format

Provide a structured review summary:

```markdown
## Code Review: [brief description]

### Files Reviewed
- file1.go (added)
- file2.ts (modified)

### Summary
[1-2 sentence overview of the changes]

### Findings
- [blocker] file:line — Description of issue and why it matters
- [warning] file:line — Description and suggested fix
- [nit] file:line — Minor improvement suggestion

### What Looks Good
- [Positive observations about the code]

### Suggestions (Optional)
- [Non-blocking improvements for consideration]
```

## Example

```markdown
## Code Review: User authentication flow

### Files Reviewed
- src/auth/login.ts (added)
- src/middleware/auth.ts (modified)
- src/types/user.ts (modified)

### Summary
Implements JWT-based login with refresh tokens. Generally well-structured,
but has a security concern and missing input validation.

### Findings
- [blocker] src/auth/login.ts:45 — JWT secret hardcoded as string literal.
  Move to environment variable: `process.env.JWT_SECRET`.
- [warning] src/auth/login.ts:67 — No rate limiting on login attempts.
  Consider adding rate limit middleware.
- [warning] src/middleware/auth.ts:23 — Missing null check on token payload.
  Add guard: `if (!decoded.userId) return res.status(401)`.
- [nit] src/types/user.ts:12 — `User` type exported but unused in this PR.

### What Looks Good
- Clean separation of auth concerns
- Proper token expiry handling
- Good TypeScript types

### Suggestions
- Consider adding refresh token rotation for better security
- Add integration test for full login → protected route flow
```

## Quick Self-Review (2-minute version)

For small changes, run through this quick check:

```bash
# 1. Any secrets exposed?
git diff HEAD | grep -iE '(password|secret|token|key).*=.*["\x27]'

# 2. Any console.log or debug statements?
git diff HEAD | grep -E '(console\.log|fmt\.Print|debugger|TODO|FIXME)'

# 3. Any large functions (>50 lines)?
git diff HEAD | grep -c '^+' | head -1
```

## Troubleshooting

### "I see too many changes to review"
- Use `git diff --stat` to identify the most-changed files
- Focus on business logic first, tests second, config last
- Consider splitting into smaller commits

### "I'm not sure if this is a blocker"
- Ask: "Could this cause data loss or security breach?" → If yes, it's a blocker
- Ask: "Will this cause bugs in production?" → If yes, it's at least a warning

### "I need context on what changed"
```bash
# See the commit history leading to this point
git log --oneline -10

# See the full file for context
git show HEAD:path/to/file
```
