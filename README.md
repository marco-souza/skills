# Skills CLI

A Go CLI tool for managing AI agent skills — reusable skill definitions stored as `SKILL.md` files in `.agents/skills/` directories.

## Installation

```bash
go install github.com/marco-souza/skills@latest
```

Requires Go 1.21+. The binary installs to `$GOPATH/bin` or `$HOME/go/bin`.

## Quick Start

```bash
# List available skills
skills list

# Validate all skills
skills validate

# Search for a skill
skills s git

# Install a skill to a project
skills i git-commit-formatter -t ~/my-project

# Create a new skill
skills a my-new-skill
```

## Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `skills list [path]` | `ls` | List all available skills |
| `skills add <name>` | `a` | Create a new skill from template |
| `skills validate [path]` | — | Validate skill format and structure |
| `skills install <skill>...` | `i` | Install skill(s) to a target project |
| `skills init [path]` | — | Initialize a project with `.agents` structure |
| `skills search <query>` | `s` | Search skills by name, description, or content |

### Global Flags

| Flag | Description |
|------|-------------|
| `-r, --root string` | Root directory for skills operations (default: `.`) |
| `-h, --help` | Help for any command |
| `-v, --version` | Show version |

---

## Command Reference

### `skills list` / `skills ls`

List all skills in the current project's `.agents/skills/` directory.

```bash
skills list                    # List local skills
skills ls --format json        # Output as JSON
skills ls --category git       # Filter by category
skills ls -r /other/project    # List from a different project
```

### `skills add` / `skills a`

Create a new skill directory with a `SKILL.md` scaffold.

```bash
skills a my-skill              # Create in current project
skills a my-skill -t ~/project # Create in a specific project
skills a my-skill --from ./template  # Use a custom template
```

Creates:
```
.agents/skills/my-skill/
└── SKILL.md
```

### `skills validate`

Validate all skills for correct YAML frontmatter, required fields, and naming conventions.

```bash
skills validate                # Validate all skills
skills validate --fix          # Auto-fix issues (name mismatches)
skills validate .agents/skills/my-skill/SKILL.md  # Validate one file
```

Checks:
- Valid YAML frontmatter
- Required `name` and `description` fields
- Description under 1024 characters
- Name matches directory name
- Name uses lowercase + hyphens only
- Description includes "Use when" and "Do NOT use when" triggers

### `skills install` / `skills i`

Install one or more skills to a target project.

```bash
# Single skill
skills i git-commit-formatter -t ~/my-project

# Multiple skills
skills i git-commit-formatter pr-review -t ~/my-project

# All skills
skills i --all -t ~/my-project

# By category
skills i --category git -t ~/my-project

# From a GitHub repo (clones via SSH)
skills i marco/skills -t ~/my-project

# Dry run
skills i git-commit-formatter --dry-run -t ~/my-project
```

| Flag | Description |
|------|-------------|
| `-t, --target string` | Target project directory |
| `--all` | Install all available skills |
| `-c, --category string` | Install skills matching category |
| `--dry-run` | Show what would be copied |

### `skills init`

Initialize a project with the `.agents/skills/` directory structure and a starter `AGENTS.md`.

```bash
skills init                    # Init in current directory
skills init ~/new-project      # Init in a specific path
```

Creates:
```
.agents/
├── AGENTS.md
└── skills/
```

### `skills search` / `skills s`

Search skills by name, description, or content.

```bash
skills s git                   # Search for "git"
skills s "commit message"      # Multi-word search
skills s git --tag devops      # Filter by tag
skills s git --format json     # JSON output
```

---

## Skill Format

Every skill is a directory under `.agents/skills/` containing a single `SKILL.md` file.

### SKILL.md Structure

```markdown
---
name: my-skill
description: >
  What this skill does in one sentence.
  Use when: the user does X or Y.
  Do NOT use when: Z applies.
---

# My Skill

## When to Use

- Trigger scenario one
- Trigger scenario two

## When NOT to Use

- Anti-pattern one
- Anti-pattern two

## Instructions

[Core instructions here]

## Examples

[Concrete examples]

## Edge Cases

[Common issues and solutions]
```

### Frontmatter Rules

| Field | Required | Rules |
|-------|----------|-------|
| `name` | Yes | 1-64 chars, lowercase + hyphens only, must match directory |
| `description` | Yes | Max 1024 chars, must include "Use when" and "Do NOT use when" |
| `tags` | No | Array of strings |
| `category` | No | Single string (e.g., `git`, `review`, `terminal`) |
| `author` | No | Author name |
| `version` | No | Semver string |

---

## Project Structure

```
skills/
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go
│   ├── list.go
│   ├── add.go
│   ├── validate.go
│   ├── install.go
│   ├── init.go
│   └── search.go
├── internal/skills/        # Core logic
│   ├── skill.go
│   ├── loader.go
│   ├── validator.go
│   ├── installer.go
│   ├── resolver.go
│   └── frontmatter.go
├── main.go
├── .agents/skills/         # This repo's skill definitions
├── testdata/               # Test fixtures
└── go.mod
```

---

## Workflow

### Adding a skill to this repo

```bash
skills a my-new-skill          # Scaffold
# Edit .agents/skills/my-new-skill/SKILL.md
skills validate                # Check format
git add .agents/skills/my-new-skill/
git commit -m "fea: add my-new-skill"
```

### Installing skills to another project

```bash
cd ~/my-project
skills init                    # Create .agents/skills/
skills i git-commit-formatter pr-review -t .  # Install skills
```

### Validating before committing

```bash
skills validate                # Catch issues early
skills validate --fix          # Auto-fix name mismatches
```

---

## Configuration

Optional `.skills.yaml` in the project root:

```yaml
skills_dir: .agents/skills
categories:
  - git
  - review
  - terminal
```

Override the config path with `SKILLS_CONFIG` env var.

---

## Development

```bash
go build -o skills .           # Build
go test ./...                  # Run tests
go vet ./...                   # Lint
go install .                   # Install locally
```
