---
name: learning
description: >
  Capture and organize learnings from coding experiences, debugging sessions,
  and problem-solving. Use when: the user wants to save a learning, document
  a discovery, record a pattern, or create a reference for future work.
  Do NOT use when: creating project documentation (use docs), writing PRDs
  (use create-prd), or when the user wants to read existing learnings without
  creating new ones.
metadata:
  scripts:
    - ../../scripts/learning-create.ts
    - ../../scripts/learning-list.ts
    - ../../scripts/learning-search.ts
    - ../../scripts/learning-show.ts
  runtime: bun
---

# Learning — Capture and Organize Learnings

Capture valuable insights from coding experiences and save them for future reference. Learnings are stored in `.agents/learnings/` as Markdown files with YAML frontmatter.

## When to Use

- After debugging a tricky issue — capture the root cause and solution
- When discovering a new pattern or best practice
- After solving a problem that might come up again
- When learning how a library or framework works
- After code review feedback that reveals a new insight
- When making architecture decisions with lessons learned

## When NOT to Use

- Creating project-level documentation (use `docs` skill)
- Writing product requirements (use `create-prd` skill)
- Recording TODO items or task status (use `project-files` skill)
- The learning is trivial or easily searchable (just use web search)

## Learning Format

Each learning is a Markdown file in `.agents/learnings/` with this structure:

```markdown
---
title: Short descriptive title
date: YYYY-MM-DD
tags: [category1, category2]
context: Brief context (optional)
---

# Learning Title

## What I Learned

[Clear, concise statement of the learning]

## Context

[When and why this was discovered — what problem were you solving?]

## Solution / Pattern

[The actual solution, pattern, or insight — make it actionable]

## Application

[How to apply this in the future — include code examples if helpful]

## References

- Related files: `path/to/file.ts`
- Related commits: `abc1234`
- External docs: [link](url)
```

## Frontmatter Fields

| Field | Required | Description |
|-------|----------|-------------|
| `title` | Yes | Short, descriptive title (max 100 chars) |
| `date` | Yes | ISO date format: `YYYY-MM-DD` |
| `tags` | Yes | Array of lowercase tags for categorization |
| `context` | No | One-line context about when this applies |

### Tag Conventions

Use lowercase, hyphenated tags:

- **Language/tech:** `typescript`, `go`, `rust`, `sql`
- **Framework:** `react`, `hono`, `prisma`, `cloudflare-workers`
- **Concept:** `debugging`, `performance`, `security`, `testing`
- **Pattern:** `error-handling`, `caching`, `authentication`
- **Area:** `api`, `database`, `frontend`, `devops`

## Workflow

### Option A: Use Helper Scripts (Recommended)

```bash
# Create a new learning interactively
bun run scripts/learning-create.ts

# Create with title and tags
bun run scripts/learning-create.ts "Prisma connection pooling" --tag prisma,database

# List all learnings
bun run scripts/learning-list.ts

# List most recent 5
bun run scripts/learning-list.ts --recent 5

# Search learnings by tag or keyword
bun run scripts/learning-search.ts "debugging"
bun run scripts/learning-search.ts --tag typescript

# Show a specific learning (supports partial filename match)
bun run scripts/learning-show.ts 2024-01-15-react
bun run scripts/learning-show.ts --latest
```

### Option B: Manual Creation

1. Create a new file in `.agents/learnings/`:

   ```bash
   # Use a descriptive filename: YYYY-MM-DD-topic-slug.md
   touch .agents/learnings/2024-01-15-react-usecallback-memoization.md
   ```

2. Add the frontmatter and content following the format above

3. The filename should be:
   - Date-prefixed for chronological sorting
   - Lowercase, hyphenated
   - Descriptive of the learning topic

## Naming Convention

```
.agents/learnings/YYYY-MM-DD-descriptive-slug.md
```

Examples:
- `2024-01-15-react-usecallback-memoization.md`
- `2024-01-20-cloudflare-workers-cron-triggers.md`
- `2024-02-01-go-context-cancellation-patterns.md`
- `2024-02-10-prisma-connection-pooling-gotcha.md`

## Examples

### Example 1: Debugging Learning

```markdown
---
title: Prisma connection pool exhaustion under load
date: 2024-01-15
tags: [prisma, database, debugging, performance]
context: Production API returning 503 errors under high traffic
---

# Prisma Connection Pool Exhaustion Under Load

## What I Learned

Prisma's default connection pool size (10) is too small for high-traffic APIs.
Under load, requests queue waiting for available connections, causing timeouts.

## Context

Production API was returning 503 errors during peak traffic. Logs showed
"Timeout: Pool connection limit reached" errors from Prisma.

## Solution / Pattern

Increase the connection pool size in your Prisma schema:

```prisma
datasource db {
  provider = "postgresql"
  url      = env("DATABASE_URL")
  relationMode = "prisma"
}

generator client {
  provider = "prisma-client-js"
}
```

Then set pool size via environment variable:
```bash
DATABASE_URL="postgresql://...?connection_limit=20&pool_timeout=10"
```

## Application

- Always set `connection_limit` based on expected concurrent requests
- Monitor connection pool metrics in production
- Consider using PgBouncer for very high throughput

## References

- Related files: `prisma/schema.prisma`, `src/db.ts`
- External docs: [Prisma Connection Management](https://www.prisma.io/docs/orm/more/comparisons/prisma-and-typeorm/connection-management)
```

### Example 2: Pattern Learning

```markdown
---
title: TypeScript template literal types for API routes
date: 2024-02-01
tags: [typescript, api, patterns]
context: Building a type-safe API client
---

# TypeScript Template Literal Types for API Routes

## What I Learned

Template literal types can enforce correct HTTP method + path combinations at
compile time, preventing invalid API calls.

## Context

Building an API client and wanted to ensure callers only use valid endpoint
combinations (e.g., GET /users, POST /users, but not DELETE /users).

## Solution / Pattern

```typescript
type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';
type ApiPath = '/users' | '/posts' | '/comments';

type ValidEndpoint = `${HttpMethod} ${ApiPath}`;

function callApi<T extends ValidEndpoint>(endpoint: T): Promise<Response> {
  const [method, path] = endpoint.split(' ');
  return fetch(path, { method });
}

// ✅ Valid
callApi('GET /users');
callApi('POST /posts');

// ❌ Compile error
callApi('PATCH /users'); // Type error: not a ValidEndpoint
```

## Application

Use this pattern when building type-safe clients, routers, or any API
abstraction where method + path combinations matter.

## References

- Related files: `src/api/client.ts`
- External docs: [TypeScript Template Literal Types](https://www.typescriptlang.org/docs/handbook/2/template-literal-types.html)
```

### Example 3: Quick Gotcha

```markdown
---
title: Bun automatically kills long-running scripts
date: 2024-02-10
tags: [bun, devops, gotcha]
context: CI pipeline timing out
---

# Bun Automatically Kills Long-Running Scripts

## What I Learned

Bun has a default 30-second timeout for script execution. Long-running scripts
will be killed without warning.

## Context

CI pipeline was failing intermittently. Investigation revealed Bun was killing
scripts that ran longer than 30 seconds.

## Solution / Pattern

Use `--timeout` flag or `BUN_TIMEOUT` environment variable:

```bash
# Extend timeout to 5 minutes
bun run script.ts --timeout 300000

# Or via environment variable
BUN_TIMEOUT=300000 bun run script.ts
```

## Application

Always set explicit timeouts for long-running scripts in CI/CD pipelines.

## References

- External docs: [Bun CLI Reference](https://bun.sh/docs/cli/run)
```

## Managing Learnings

### Search by Tag

```bash
# Find all TypeScript learnings
grep -l "tags:.*typescript" .agents/learnings/*.md

# Or using the search script
./scripts/learning-search.sh --tag typescript
```

### Search by Keyword

```bash
# Search learning content
grep -rl "connection pool" .agents/learnings/

# Or using the search script
./scripts/learning-search.sh "connection pool"
```

### List by Date

```bash
# Learnings are naturally sorted by date via filename
ls .agents/learnings/

# Show most recent
ls -r .agents/learnings/ | head -10
```

### Archive Old Learnings

```bash
# Move learnings older than 1 year to archive
mkdir -p .agents/learnings/archive/2023
mv .agents/learnings/2023-*.md .agents/learnings/archive/2023/
```

## Best Practices

- **Be specific** — "Prisma connection pool size" is better than "Database tips"
- **Include code** — Working examples are more useful than descriptions
- **Add context** — What problem were you solving? What was the error?
- **Tag consistently** — Use the same tags for related topics
- **Keep it actionable** — Future you should be able to apply this without re-reading the original problem
- **Date everything** — Context decays; dates help prioritize newer learnings

## Anti-Patterns

- Don't create learnings for things easily found via web search
- Don't duplicate learnings — consolidate related insights into one file
- Don't use vague titles — "TypeScript trick" is useless; be specific
- Don't skip the "Application" section — that's the whole point

## Integration with Other Skills

- **After `debug`** — Capture the root cause and solution as a learning
- **After `code-review`** — Record patterns or anti-patterns discovered
- **With `explore`** — Document findings about unfamiliar codebases
- **With `mixture-of-experts`** — Save expert insights for future reference
