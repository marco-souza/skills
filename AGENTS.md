# AGENTS.md - AI Agent Guidelines for Skills Repository

This document provides guidelines for AI agents working in this repository.

## Repository Overview

This is the **Skills CLI** repository — a Go CLI tool (`skills`) for managing AI agent skill definitions stored as `SKILL.md` files in `.agents/skills/` directories. The repo also ships the canonical collection of reusable skills that agents can install into any project.

**Key Characteristics:**

- Go CLI (`github.com/marco-souza/skills`) built with Cobra
- Skill definitions live in `.agents/skills/` as Markdown files with YAML frontmatter
- Follows the [Agent Skills standard](https://agentskills.io/specification)
- Tests, linting, and builds are all available via standard Go tooling

## Repository Structure

```
skills/
├── AGENTS.md                        # This file — agent guidelines
├── README.md                        # User-facing documentation
├── PRD.md                           # Product requirements document
├── main.go                          # CLI entry point
├── go.mod / go.sum                  # Go module files
├── cmd/                             # Cobra CLI commands
│   ├── root.go
│   ├── list.go
│   ├── add.go
│   ├── install.go
│   ├── uninstall.go
│   ├── init.go
│   └── config.go
├── internal/
│   ├── skills/                      # Core skill logic
│   │   ├── skill.go                 # Skill type, validation, frontmatter parsing
│   │   ├── loader.go                # Discovers/loads skills from disk
│   │   ├── installer.go             # Copies skills to target projects
│   │   ├── resolver.go              # Resolves source directories
│   │   └── remote.go               # GitHub clone support
│   └── config/                      # Persistent CLI config (~/.config/skills/)
│       └── config.go
├── testdata/                        # Test fixtures
└── .agents/skills/                  # Canonical skill definitions shipped with this repo
    ├── brainstorm/
    ├── create-prd/
    ├── explore/
    ├── git-commit-formatter/
    ├── grill-me/
    ├── implement-tasks/
    ├── markdown-format/
    ├── mixture-of-experts/
    ├── pr-review/
    ├── prd-to-tasks/
    ├── project-files/
    ├── spawn-subagents/
    └── terminal-multiplexer/
```

## Build/Lint/Test Commands

```bash
go build -o skills .     # Build the CLI binary
go test ./...            # Run all tests
go vet ./...             # Static analysis / lint
go install .             # Install the binary locally
```

## File Format: SKILL.md

Each skill is defined in a `SKILL.md` file with the following structure:

### YAML Frontmatter (Required)

```yaml
---
name: skill-name-here
description: >
  Clear description of what the skill does.
  Use when: specific triggers for when to use this skill.
  Do NOT use when: negative triggers to avoid misuse.
---
```

**Frontmatter Fields:**

- `name`: Lowercase, hyphenated identifier (e.g., `git-commit-formatter`, `pr-review`)
  - Must match parent directory name
  - 1-64 characters
  - Only lowercase letters, numbers, and hyphens
- `description`: Multi-line description with explicit triggers
  - Include positive triggers ("Use when...")
  - Include negative triggers ("Do NOT use when...")
  - Max 1024 characters

### Description Best Practices

Good descriptions include:

1. **What** the skill does (one sentence)
2. **When** to use it (positive triggers)
3. **When NOT** to use it (negative triggers)
4. **Prerequisites** if any (optional)

**Example:**

```yaml
description: >
  Format and write git commit messages using a structured type-based format.
  Use when the user asks to commit, stage, create a commit message, or summarize code changes.
  Do NOT use for merge commits, revert operations, or work-in-progress commits.
```

### Markdown Content Structure

All skills should follow this standard section order:

```markdown
# Skill Title

[Brief introduction paragraph]

## Prerequisites
[If any requirements before using the skill]

## When to Use
[Positive triggers - when this skill applies]

## When NOT to Use
[Negative triggers - when to avoid this skill]

## [Main Content]
[Core instructions, commands, workflows]

## Examples
[Concrete, copy-pasteable examples]

## Edge Cases / Troubleshooting
[Common issues and solutions]

## Best Practices
[Tips for optimal usage]
```

## Code Style Guidelines

### Directory Naming

- Use lowercase letters only
- Use hyphens to separate words (kebab-case)
- Be descriptive but concise
- Must match the `name` field in SKILL.md frontmatter

**Good:** `git-commit-formatter`, `pr-review`, `terminal-multiplexer`
**Bad:** `GitCommitFormatter`, `pr_review`, `tmux-skill`

### File Naming

- Skill definition files MUST be named `SKILL.md` (uppercase)
- One `SKILL.md` per directory

### Markdown Formatting

- Use `#` for main title (H1) - one per file
- Use `##` for major sections (H2)
- Use `###` for subsections (H3)
- Use `-` (dash) for bullet points
- Use numbered lists (`1.`, `2.`) for sequential steps
- Use **bold** for emphasis on key terms
- Use code blocks with language specifiers: ` ```bash `, ` ```yaml `, etc.
- Use inline code for commands, file names, and technical terms
- Use `>` for notes and callouts

### Writing Style

- Use imperative mood for instructions ("Do this", not "You should do this")
- Be concise and direct
- Provide examples for complex concepts
- Use consistent terminology throughout
- Include edge cases and error handling

## Git Commit Convention

Follow the format defined in `git-commit-formatter/SKILL.md`:

```
<type>(optional scope): <description>
```

### Commit Types

| Type   | Description                          |
|--------|--------------------------------------|
| `fea`  | A new feature or skill               |
| `fix`  | A bug fix or correction              |
| `ref`  | Refactoring or restructuring         |
| `ai`   | AI-related changes                   |
| `test` | Test-related changes                 |
| `oth`  | Other changes (docs, config, etc.)   |

### Commit Message Guidelines

- Use lowercase for the description
- Keep the first line under 72 characters
- Do not end with a period
- Use present tense ("add feature" not "added feature")

**Examples:**

```
fea: add kubernetes deployment skill
fix: correct yaml syntax in pr-review skill
ref: reorganize sections in terminal-multiplexer skill
oth: update AGENTS.md with new guidelines
```

## Creating a New Skill

1. Create a new directory with a descriptive, hyphenated name:

   ```bash
   mkdir .agents/skills/my-new-skill
   ```

2. Create the `SKILL.md` file with required frontmatter:

   ```bash
   touch .agents/skills/my-new-skill/SKILL.md
   ```

3. Add content following the structure guidelines above

4. Ensure description includes:
   - What the skill does
   - Positive triggers ("Use when...")
   - Negative triggers ("Do NOT use when...")

5. Commit with appropriate type:

   ```bash
   git add .agents/skills/my-new-skill/
   git commit -m "fea: add my-new-skill for X purpose"
   ```

## Modifying Existing Skills

- Preserve the existing structure and formatting style
- Update the `description` in frontmatter if the skill's purpose changes
- Keep changes focused and atomic
- Test that YAML frontmatter remains valid after edits

## Common Patterns in This Repository

### Description Template

Use this template for all skill descriptions:

```yaml
description: >
  [What the skill does in one sentence].
  Use when: [positive triggers].
  Do NOT use when: [negative triggers].
  Requires: [prerequisites if any].
```

### When to Use / When NOT to Use Sections

```markdown
## When to Use

- First use case
- Second use case
- Third use case

## When NOT to Use

- Anti-pattern one
- Anti-pattern two
- When user explicitly requests alternatives
```

### Prerequisites Section

```markdown
## Prerequisites

Verify requirements before use:

```bash
# Check if tool is available
command -v <tool> >/dev/null 2>&1 || {
    echo "Error: <tool> is not installed"
    echo "Install with: <install command>"
    exit 1
}
```

```

### Troubleshooting Section

```markdown
## Troubleshooting

### Common Issue 1

```bash
# Diagnosis command
# Fix command
```

### Common Issue 2

```bash
# Diagnosis command
# Fix command
```

```

### Code Examples

Always specify the language for syntax highlighting:

```markdown
## Example

Here's how to use this command:

` ` `bash
tmux new-session -s mysession
` ` `
```

### Key Concepts

Use bold for introducing important terms:

```markdown
The **session** is the top-level container in tmux...
```

## Validation Checklist

Before committing changes, verify:

- [ ] `SKILL.md` exists in the skill directory
- [ ] Directory name matches `name` field in frontmatter
- [ ] YAML frontmatter has valid syntax
- [ ] `name` field uses lowercase hyphenated format (1-64 chars)
- [ ] `description` field includes positive triggers
- [ ] `description` field includes negative triggers ("Do NOT use when...")
- [ ] `description` field is under 1024 characters
- [ ] Markdown headings follow hierarchy (H1 > H2 > H3)
- [ ] Code blocks have language specifiers
- [ ] Sections follow standard order (When to Use, Main Content, Examples, etc.)
- [ ] Edge cases and troubleshooting are documented
- [ ] Commit message follows convention

## References

- [Agent Skills Specification](https://agentskills.io/specification)
- [Agent Skills Integration](https://agentskills.io/integrate-skills)
- [Pi Documentation](https://github.com/badlogic/pi)
