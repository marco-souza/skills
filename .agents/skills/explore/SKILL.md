---
name: explore
description: >
  Research and explore codebases to build context before making changes.
  Use when starting work on an unfamiliar project, investigating a bug,
  planning a feature, or when you need to understand how something works.
---

# Explore

Systematically research codebases to build accurate mental models before acting.

## When to Explore

- Starting work on a new/unfamiliar codebase
- Investigating a bug or error
- Planning a new feature
- Reviewing code you didn't write
- Context feels incomplete or ambiguous
- About to make assumptions

## Exploration Principles

### 1. Outside-In (Macro → Micro)

Start broad, then drill down:

```
Project Structure → Module Organization → File Purpose → Implementation Details
```

### 2. Follow the Data

Trace how information flows:

```
Input → Validation → Transformation → Storage → Output
```

### 3. Question Everything

Before assuming:

- "What is this file's responsibility?"
- "Who calls this function?"
- "What would break if I changed this?"
- "Is this used in production or just tests?"

## Modern CLI Tools

Use modern alternatives for faster, more ergonomic exploration. Install them if available:

### Tool Comparison

| Task | Traditional | Modern Alternative | Why Modern is Better |
|------|-------------|-------------------|----------------------|
| Search text in files | `grep` | `ripgrep` (`rg`) | Faster, respects `.gitignore`, better output formatting |
| Find files | `find` | `fd` | Faster, simpler syntax, respects `.gitignore` |
| View files | `cat` | `bat` | Syntax highlighting, line numbers, Git integration |
| View diffs | `git diff` | `delta` | Syntax-highlighted diffs, side-by-side view, line numbers |

### Installation

```bash
# macOS (Homebrew)
brew install ripgrep fd bat delta

# Ubuntu/Debian
sudo apt install ripgrep fd-find bat
# Note: Ubuntu binary is 'fdfind', symlink: sudo ln -sf $(which fdfind) /usr/local/bin/fd

# Arch Linux
sudo pacman -S ripgrep fd bat delta

# Check if installed
command -v rg fd bat delta 2>/dev/null || echo 'Some tools not installed'
```

### Quick Usage

```bash
# ripgrep - faster grep
grep -r 'TODO' src/           # Traditional
rg 'TODO' src/                # Modern (faster, cleaner output)
rg -t ts 'interface'          # Search only TypeScript files

# fd - faster find
find . -name '*.ts'           # Traditional
fd '\.ts$'                    # Modern (regex support)
fd -e ts                      # By extension

# bat - better cat
cat README.md                 # Traditional
bat README.md                 # Modern (syntax highlighting)
bat -l yaml config.yml        # Force language

# delta - better diffs
git diff                      # Traditional
git diff | delta              # Modern (syntax highlighting)
delta --side-by-side          # Side-by-side view
```

---

## Phase 1: Project Reconnaissance (5 min)

### Read Entry Points

```bash
# Start here — every project has these
cat README.md                    # Project overview, setup, conventions
cat package.json 2>/dev/null     # Dependencies, scripts, metadata
cat AGENTS.md 2>/dev/null        # Project-specific agent instructions

# Or with bat (if available)
bat README.md                     # Syntax-highlighted view with line numbers
bat package.json 2>/dev/null     # Better readability for JSON
```

### Map Structure

```bash
# Understand directory layout
ls -la                           # Root contents
find . -maxdepth 2 -type d       # Top-level directories
tree -L 2 2>/dev/null || find . -maxdepth 2 -type d | head -20

# Or with fd (if available)
fd -d 2 -t d                    # Top-level directories (simpler syntax)
fd -d 3 -e ts -e go             # Find source files up to 3 levels deep
```

### Identify Technology Stack

```bash
# Detect frameworks and patterns
grep -l "react\|vue\|angular" package.json 2>/dev/null && echo "Frontend framework detected"
grep -l "express\|fastify\|hono" package.json 2>/dev/null && echo "Backend framework detected"
grep -l "prisma\|drizzle\|typeorm" package.json 2>/dev/null && echo "ORM detected"

# Or with ripgrep (if available)
rg -l 'react|vue|angular' package.json 2>/dev/null && echo "Frontend framework detected"
rg -l 'prisma|drizzle' package.json 2>/dev/null && echo "ORM detected"
```

### Output: Create SESSION.md Entry

```markdown
## Exploration: Project Overview

**Stack:** React + Node + Prisma
**Structure:** src/{components,pages,api}/
**Conventions:**

- Feature-based folders
- Tests co-located (\*.test.ts)
- API routes in src/api/
```

## Phase 2: Architecture Mapping (10 min)

### Find the Core Modules

```bash
# Identify high-traffic files
git log --pretty=format: --name-only | sort | uniq -c | sort -rg | head -20

# Or find most imported modules
grep -r "^import.*from" --include="*.ts" --include="*.js" | cut -d'"' -f2 | sort | uniq -c | sort -rg | head -20

# With ripgrep (faster)
rg -o 'from ["\'][^"\']+["\']' -t ts | sed 's/from //g' | sort | uniq -c | sort -rg | head -20
```

### Trace Data Flow

Pick a key entity (e.g., "User", "Order") and trace it:

```bash
# Find where it's defined
grep -r "interface User\|type User\|class User" --include="*.ts" | head -5

# With ripgrep
rg 'interface User|type User|class User' -t ts | head -5

# Find where it's used
grep -r "User" --include="*.ts" | grep -v node_modules | wc -l

# With ripgrep (auto-ignores node_modules)
rg -c 'User' -t ts | head -10

# Find API endpoints that handle it
grep -r "user\|User" src/api/ --include="*.ts" | head -10

# With ripgrep
cd src/api && rg -i 'user' -t ts | head -10
```

### Map Dependencies

```bash
# What does this module depend on?
cat src/auth/login.ts | grep "^import"

# Or with bat + rg
bat src/auth/login.ts | rg '^import'

# What depends on this module?
grep -r "from.*auth/login" --include="*.ts" | head -10

# With ripgrep
rg 'from.*auth/login' -t ts | head -10
```

### Document in STRUCTURE.md

```markdown
# Architecture

## Entry Points

- `src/main.ts` — Application bootstrap
- `src/api/index.ts` — API route registration

## Core Modules

- `auth/` — Authentication, session management
- `models/` — Database schemas (Prisma)
- `services/` — Business logic

## Data Flow
```

Request → Middleware → Handler → Service → Model → DB

```

## Key Files
| File | Purpose |
|------|---------|
| `src/auth/jwt.ts` | Token generation/validation |
| `src/models/user.ts` | User entity definition |
```

## Phase 3: Deep Dive (Targeted)

### Locate Relevant Code

Given a task (e.g., "fix login bug"):

```bash
# Search for keywords
grep -r "login\|signin\|authenticate" --include="*.ts" | grep -v test | head -10

# With ripgrep (cleaner, faster)
rg -g '!*.test.ts' 'login|signin|authenticate' -t ts | head -10

# Find related tests (show usage patterns)
grep -r "login" --include="*.test.ts" | head -5

# With ripgrep
fd -e test.ts && rg 'login' -g '*.test.ts' | head -5

# Check recent changes
git log --oneline --all --grep="login" | head -5
```

### Read Call Stacks

Start from entry point, trace down:

```bash
# API route
cat src/api/auth.ts
# → Calls authService.login()

# Implementation
cat src/services/auth.ts
# → Calls userRepository.findByEmail()

# Data layer
cat src/repositories/user.ts
# → Calls prisma.user.findUnique()

# With bat (better readability)
bat src/api/auth.ts src/services/auth.ts src/repositories/user.ts
```

### Understand Edge Cases

Look for:

- Error handling (`throw`, `catch`, `if (error)`)
- Validation logic (`zod`, `joi`, manual checks)
- Permissions/authorization (`canAccess`, `requireAuth`)
- Environment-specific code (`process.env`, `import.meta.env`)

## Phase 4: Context Clarification

### When You're Stuck

If something doesn't make sense:

1. **Find examples** — How is this used elsewhere?

   ```bash
   grep -r "similarFunctionName" --include="*.ts" -A 3 | head -20
   ```

2. **Check tests** — Tests show intended behavior

   ```bash
   cat src/auth/login.test.ts | grep -A 10 "should"
   ```

3. **Look for docs** — Comments, JSDoc, READMEs

   ```bash
   grep -B 5 "function login" src/auth/login.ts
   ```

4. **Trace git history** — Why was this added?

   ```bash
   git log -p --all -S "suspiciousCode" -- src/auth/login.ts | head -50
   ```

### Validate Assumptions

Before proceeding, verify:

```bash
# "This function is only called from X"
grep -r "functionName" --include="*.ts" | grep -v "def\|export" | wc -l

# With ripgrep
rg -c 'functionName' -t ts  # Shows count per file

# "This is always a string"
grep -r "variableName:" --include="*.ts" | head -5

# With ripgrep
rg 'variableName:' -t ts | head -5

# "This mutation updates the database"
grep -A 10 "mutationName" src/services/*.ts | grep -E "prisma|save|update"

# With ripgrep
rg -A 10 'mutationName' -t ts | rg 'prisma|save|update'
```

## Exploration Outputs

### Required: CONTEXT.md

After exploration, create/update:

```markdown
# Context: <Feature/Area>

## What I Learned

- X is handled by Y module
- Z is the source of truth for W data
- Authentication uses JWT with 24h expiry

## Open Questions

- [ ] Why is X implemented as Y instead of Z?
- [ ] How does the caching layer work?

## Relevant Files

| File                     | Why It Matters         |
| ------------------------ | ---------------------- |
| `src/auth/jwt.ts`        | Token generation logic |
| `src/middleware/auth.ts` | Route protection       |

## Risks/Watchouts

- Changing X requires updating Y and Z
- No tests for edge case A
```

### Optional: Update Project Files

Based on your exploration:

- Add to `DECISIONS.md` if you discovered why something is the way it is
- Update `TODO.md` with tasks that emerged
- Create `SPEC.md` for areas you now understand

## Common Exploration Patterns

### Pattern: Bug Investigation

```bash
# 1. Find error location
grep -r "errorMessage" --include="*.ts"

# 2. Trace backward
cat src/fileWithError.ts
# → Find caller
# → Find caller's caller

# 3. Check recent changes
git log --oneline --all -- src/fileWithError.ts | head -5

# 4. Reproduce
# Look for test that exercises this path
```

### Pattern: Feature Addition

```bash
# 1. Find similar features
grep -r "similarFeature" --include="*.ts" -l

# 2. Study the pattern
cat src/features/similar/index.ts

# 3. Identify all touchpoints
grep -r "similarFeature" --include="*.ts" | grep -v "def\|export" | cut -d: -f1 | sort -u

# 4. Note conventions
# - How are routes registered?
# - How are tests structured?
# - What validation is used?
```

### Pattern: Code Review Prep

```bash
# 1. Understand the change
git diff main...feature-branch --stat

# 2. Read modified files in dependency order
# (models → services → handlers → tests)

# 3. Check for missing pieces
# - Are there tests?
# - Is there error handling?
# - Are types defined?

# 4. Verify assumptions
# - Does this break existing code?
# - Are there migration concerns?

# 5. Review with delta (if available)
git diff main...feature-branch | delta          # Syntax-highlighted diff
delta --side-by-side main...feature-branch     # Side-by-side view
```

## Anti-Patterns (DON'Ts)

- **Don't skim** — Read code line by line when it matters
- **Don't guess types** — Check the actual definitions
- **Don't assume single usage** — Always grep for callers
- **Don't ignore tests** — They document expected behavior
- **Don't skip error paths** — Understand failure modes
- **Don't make changes during exploration** — Research first, act second

## Quick Reference: Grep Patterns

```bash
# Find function definitions
grep -r "function name\|const name\|async function name" --include="*.ts"

# Find all imports of a module
grep -r "from.*module-name" --include="*.ts"

# Find where a variable is used (excluding definition)
grep -r "varName" --include="*.ts" | grep -v "const\|let\|var\|import"

# Find exported items
grep -r "^export" --include="*.ts" src/some-module/

# Find TODO/FIXME comments
grep -r "TODO\|FIXME\|XXX\|HACK" --include="*.ts" src/
```

### With ripgrep (faster alternatives)

```bash
# Find function definitions
rg '(function|const|async function) name' -t ts

# Find all imports of a module
rg 'from.*module-name' -t ts

# Find exported items
rg '^export' -t ts src/some-module/

# Find TODO/FIXME comments
rg 'TODO|FIXME|XXX|HACK' -t ts src/
```
