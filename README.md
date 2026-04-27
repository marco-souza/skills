# Skills CLI

A lightweight Go CLI tool for managing AI agent skills — reusable skill definitions stored as `SKILL.md` files in `.agents/skills/` directories.

## Installation

```bash
go install github.com/marco-souza/skills@latest
```

Requires Go 1.22+. The binary installs to `$GOPATH/bin` or `$HOME/go/bin`.

## Quick Start

```bash
# List available skills
skills list

# Install a skill to a project
skills install git-commit-formatter -t ~/my-project

# Create a new skill
skills add my-new-skill

# Remove a skill
skills uninstall my-new-skill
```

## Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `skills list [path]` | `ls` | List skills from local or remote repo |
| `skills install <skill>...` | `i` | Install skill(s) to a target project |
| `skills uninstall <skill>...` | `rm`, `remove` | Remove skill(s) from a project |
| `skills add <name>` | `a` | Create a new skill from template |
| `skills init [path]` | — | Initialize a project with `.agents` structure |
| `skills config` | — | Manage persistent CLI configuration |

### Global Flags

| Flag | Description |
|------|-------------|
| `-r, --repo string` | GitHub repo (owner/repo) for remote operations |
| `-h, --help` | Help for any command |
| `-v, --version` | Show version |

---

## Command Reference

### `skills list` / `skills ls`

List skills from the local `.agents/skills` directory. Use `--repo` to list from a remote GitHub repository.

```bash
skills ls                      # List local skills
skills ls -r marco-souza/skills  # List from a remote repo
skills ls ~/other-project      # List from a different project
```

### `skills add` / `skills a`

Create a new skill directory with a `SKILL.md` scaffold.

```bash
skills a my-skill              # Create in current project
skills a my-skill -t ~/project # Create in a specific project
```

Creates:
```
.agents/skills/my-skill/
└── SKILL.md
```

### `skills install` / `skills i`

Install one or more skills to a target project.

```bash
# Single skill
skills i git-commit-formatter -t ~/my-project

# Multiple skills
skills i git-commit-formatter pr-review -t ~/my-project

# From a GitHub repo
skills i git-commit-formatter -r marco-souza/skills -t ~/my-project
```

| Flag | Description |
|------|-------------|
| `-t, --target string` | Target project directory (default: current) |
| `-r, --repo string` | GitHub repo to install from |

### `skills uninstall` / `skills rm` / `skills remove`

Remove one or more skills from a project.

```bash
# Single skill
skills uninstall git-commit-formatter

# Multiple skills
skills uninstall git-commit-formatter pr-review -t ~/my-project
```

| Flag | Description |
|------|-------------|
| `-t, --target string` | Target project directory (default: current) |

### `skills config`

Manage persistent CLI configuration stored at `~/.config/skills/config.yaml`.

```bash
skills config list             # Show all settings
skills config get default_repo # Get default repo
skills config set default_repo my-org/skills
skills config set default_root ~/projects
```

| Setting | Description | Default |
|---------|-------------|--------|
| `default_repo` | Fallback repo when `--repo` is not set | `marco-souza/skills` |
| `default_root` | Fallback root when `--root` is not set | `.` |

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
│   ├── install.go
│   ├── uninstall.go
│   ├── init.go
│   └── config.go
├── internal/
│   ├── skills/               # Core logic
│   │   ├── skill.go
│   │   ├── loader.go
│   │   ├── installer.go
│   │   ├── resolver.go
│   │   ├── remote.go
│   │   └── frontmatter.go
│   └── config/               # Persistent CLI config
│       └── config.go
├── main.go
├── .agents/skills/         # This repo's skill definitions
├── testdata/               # Test fixtures
└── go.mod
```

---

## Configuration

Persistent config lives at `~/.config/skills/config.yaml`:

```yaml
default_repo: marco-souza/skills
default_root: .
```

These values are used as fallbacks when `--repo` or `--root` flags are not provided.

Manage via CLI:

```bash
skills config list              # Show current settings
skills config set default_repo my-org/skills
skills config set default_root ~/workspace
```

---

## Development

```bash
go build -o skills .           # Build
go test ./...                  # Run tests
go vet ./...                   # Lint
go install .                   # Install locally
```
