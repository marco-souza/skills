---
name: docs
description: >
  Write and maintain project documentation including READMEs, API docs, ADRs, CHANGELOGs, and inline code comments.
  Use when: the user asks to write documentation, update a README, create an API reference, draft an ADR, add a CHANGELOG entry, or improve inline comments.
  Do NOT use when: the user wants to format markdown (use markdown-format), review code (use pr-review), or create a PRD (use create-prd).
---

# Documentation Writing

Write clear, consistent, and maintainable documentation for software projects.

## Prerequisites

- Access to the project source code
- Understanding of the project's purpose and structure
- Familiarity with the project's conventions

## Documentation Types

### 1. README Structure

Every project should have a `README.md` at the root with these sections:

```markdown
# Project Name

One-sentence description of what the project does.

## Features

- Feature 1: brief description
- Feature 2: brief description

## Quick Start

### Prerequisites

- Requirement 1
- Requirement 2

### Installation

```bash
# Installation command
```

### Usage

```bash
# Basic usage example
```

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `VAR_NAME` | What it does | `default` |

## API Reference

Link to detailed API docs or inline summary.

## Contributing

Link to CONTRIBUTING.md or brief guidelines.

## License

License type with link.
```

**Rules:**
- Start with a clear, concise project name and tagline
- Provide a "Quick Start" that works in under 30 seconds
- Use real commands, not pseudocode
- Include links to detailed docs for complex topics
- Keep README under 500 lines; move details to separate files

### 2. API Documentation

Document APIs using these patterns:

#### REST APIs

```markdown
## Endpoints

### POST /api/v1/users

Create a new user account.

**Request:**

```json
{
  "email": "user@example.com",
  "name": "Jane Doe",
  "role": "member"
}
```

**Response (201 Created):**

```json
{
  "id": "usr_abc123",
  "email": "user@example.com",
  "name": "Jane Doe",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Errors:**

| Status | Code | Description |
|--------|------|-------------|
| 400 | `invalid_email` | Email format is invalid |
| 409 | `email_exists` | Email already registered |
```

#### CLI Tools

```markdown
## Commands

### `tool init`

Initialize a new project.

**Usage:**

```bash
tool init [options] <project-name>
```

**Options:**

| Flag | Description | Default |
|------|-------------|---------|
| `--template, -t` | Template to use | `default` |
| `--yes, -y` | Skip prompts | `false` |

**Examples:**

```bash
# Initialize with default template
tool init my-project

# Initialize with specific template
tool init my-project --template typescript
```
```

#### Function/Method Documentation

```markdown
### `createUser(options: CreateUserOptions): Promise<User>`

Creates a new user in the system.

**Parameters:**

- `options.email` (string, required): User's email address
- `options.name` (string, required): User's display name
- `options.role` (enum, optional): User role. One of: `admin`, `member`, `viewer`. Default: `member`

**Returns:** `Promise<User>` - The created user object

**Throws:**

- `ValidationError` - If email is invalid
- `ConflictError` - If email already exists

**Example:**

```typescript
const user = await createUser({
  email: 'jane@example.com',
  name: 'Jane Doe'
});
```
```

### 3. Architecture Decision Records (ADRs)

Use ADRs to document significant architectural decisions. Store in `docs/adr/` or `docs/decisions/`.

#### ADR Template

```markdown
# ADR-{NUMBER}: {TITLE}

## Status

{Proposed | Accepted | Deprecated | Superseded by ADR-XXX}

## Date

{YYYY-MM-DD}

## Context

What is the issue that we're seeing that is motivating this decision or change?

## Decision

What is the change that we're proposing and/or doing?

## Consequences

### Positive

- {benefit 1}
- {benefit 2}

### Negative

- {tradeoff 1}
- {tradeoff 2}

### Risks

- {risk 1}

## Alternatives Considered

### {Alternative 1}

{Description of alternative}

**Pros:** {benefits}
**Cons:** {drawbacks}
**Why rejected:** {reason}
```

#### ADR Examples

```markdown
# ADR-001: Use PostgreSQL as Primary Database

## Status

Accepted

## Date

2024-01-15

## Context

We need a relational database that supports complex queries, JSON storage, and has strong ecosystem support.

## Decision

We will use PostgreSQL as our primary database.

## Consequences

### Positive

- ACID compliance for transactional data
- Native JSON support for flexible schemas
- Excellent tooling and community support

### Negative

- Heavier operational overhead than SQLite
- Requires dedicated database server

### Risks

- Team needs to learn PostgreSQL-specific features

## Alternatives Considered

### SQLite

**Pros:** Zero config, embedded, fast for read-heavy workloads
**Cons:** Limited concurrency, no network access
**Why rejected:** Doesn't support concurrent writes from multiple services
```

**Rules:**
- Number ADRs sequentially: `001`, `002`, etc.
- Write in clear, factual language
- Record the decision, not just the options
- Document rejected alternatives and why

### 4. CHANGELOG Format

Follow [Keep a Changelog](https://keepachangelog.com/) format. File: `CHANGELOG.md`.

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
- New feature description

### Changed
- Changes to existing functionality

### Deprecated
- Features that will be removed

### Removed
- Features that have been removed

### Fixed
- Bug fixes

### Security
- Vulnerability fixes

## [1.2.0] - 2024-01-15

### Added
- User authentication with JWT tokens
- Rate limiting on API endpoints

### Fixed
- Memory leak in connection pool
- Incorrect date formatting in exports

## [1.1.0] - 2024-01-01

### Added
- CSV export functionality

### Changed
- Improved query performance by 50%

## [1.0.0] - 2023-12-01

Initial stable release.
```

**Rules:**
- Add entries to `[Unreleased]` section
- Move to versioned section on release
- Use past tense for Fixed/Changed, imperative for Added
- Group by type: Added, Changed, Deprecated, Removed, Fixed, Security
- Include issue/PR references: `- Fix login timeout (#123)`

### 5. Inline Documentation

#### General Principles

- Document **why**, not **what**
- Keep comments close to the code they describe
- Update comments when updating code
- Use complete sentences with proper punctuation

#### When to Comment

```go
// GOOD: Explains why
// Retry 3 times because the API intermittently fails under load
for i := 0; i < 3; i++ {
    result, err := callAPI()
    if err == nil {
        return result
    }
}

// BAD: Explains what (code is self-explanatory)
// Loop 3 times
for i := 0; i < 3; i++ {
    // ...
}
```

#### Documentation Comments

```go
// ProcessOrder handles the complete order lifecycle including payment
// processing, inventory allocation, and shipping notification.
//
// It returns an error if any step fails, but partial side effects
// may have occurred (e.g., payment charged but not shipped).
//
// See: https://docs.internal.com/orders/lifecycle
func ProcessOrder(ctx context.Context, orderID string) error {
    // ...
}
```

#### TODO/FIXME Format

```typescript
// TODO(username): description [ticket-id]
// FIXME(username): description [ticket-id]
// HACK(username): temporary workaround for #123

// Examples:
// TODO(jane): add retry logic for transient failures [PROJ-456]
// FIXME(bob): remove after API v1 deprecation [PROJ-789]
```

## Style Guidelines

### Writing Style

- Use **active voice**: "The function returns..." not "The value is returned..."
- Use **present tense**: "This method creates..." not "This method will create..."
- Be **concise**: Aim for 1-2 sentences per comment
- Use **consistent terminology**: Pick one term and stick with it

### Formatting

- Use **Markdown** for all documentation files
- Use `code blocks` for commands, paths, and code
- Use **bold** for emphasis on key terms
- Use numbered lists for sequential steps
- Use bullet lists for non-sequential items

### Structure

- Start with the most important information
- Use headings to break up long sections
- Include a table of contents for docs over 200 lines
- Link to related documentation

## Process

1. **Identify documentation type** needed (README, API, ADR, etc.)
2. **Review existing docs** for patterns and conventions
3. **Write draft** following the appropriate template
4. **Review for accuracy** - test all code examples
5. **Proofread** - check for typos and clarity
6. **Commit** with appropriate message type (`oth`)

## Edge Cases

### Outdated Documentation

When documentation is outdated:
1. Update the documentation
2. Add an entry to CHANGELOG under `[Unreleased] > Fixed`
3. Consider if the old documentation should be archived

### Conflicting Sources

When documentation conflicts with code:
1. Trust the code as the source of truth
2. Update documentation to match
3. Add tests if the behavior is intentional

### Large Documentation Updates

For major documentation overhauls:
1. Break into multiple commits
2. Use ADR if it involves architectural decisions
3. Update CHANGELOG with summary

## Examples

### Example: Writing a README

**Input:** User asks "Write a README for my project"

**Process:**
1. Read `package.json`, `go.mod`, or similar to understand the project
2. Review existing source for key features
3. Draft README following the template
4. Include real, tested examples

**Output:**

```markdown
# myproject

A fast, lightweight JSON parser for Go.

## Quick Start

```bash
go get github.com/user/myproject
```

```go
import "github.com/user/myproject"

data, err := myproject.Parse(jsonString)
```

## Features

- Streaming parser for large files
- Zero allocations in hot path
- Full RFC 8259 compliance
```

### Example: Creating an ADR

**Input:** User asks "Document why we chose Redis"

**Process:**
1. Review meeting notes or Slack discussions
2. Identify the context, decision, and tradeoffs
3. Draft ADR using the template
4. Save as `docs/adr/003-use-redis-for-caching.md`

## Best Practices

- **Write for your audience**: Developers vs users vs operators
- **Keep docs close to code**: Store with the relevant source
- **Version documentation**: Tag docs with releases
- **Automate what you can**: Generate API docs from code annotations
- **Review docs in PRs**: Treat documentation like code
- **Measure usage**: Track which docs are accessed

## Troubleshooting

### Documentation Feels Outdated

- Add a "Last updated: YYYY-MM-DD" header
- Set up reminders to review quarterly
- Consider automating doc generation from code

### Too Much Documentation

- Focus on what users need to get started
- Move detailed reference to separate files
- Use links instead of duplicating content

### Inconsistent Style

- Create a style guide
- Use linting tools (markdownlint, Vale)
- Review docs in pull requests
