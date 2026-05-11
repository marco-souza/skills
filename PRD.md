# PRD: Skills CLI

## What We're Building

A lightweight Go CLI tool (`skills`) for managing AI agent skill definitions across projects. Skills are reusable Markdown files stored in `.agents/skills/` directories. The CLI lets developers list, install, uninstall, and create skills — either from a local directory or a remote GitHub repository — so that AI assistants can be augmented with consistent, project-specific capabilities.

## Why

AI coding agents are increasingly embedded in developer workflows, but there is no standard way to distribute or reuse agent instructions across projects. Each team ends up copy-pasting prompt snippets or maintaining bespoke instruction files with no versioning or tooling. The Skills CLI solves this by treating skills as first-class artifacts: versioned, installable, and shareable via GitHub.

## Key Features

- **List skills** from the local `.agents/skills/` directory, a local path, or a remote GitHub repo
- **Install skills** (single, multiple, or all) to any target project directory
- **Uninstall skills** from a project
- **Scaffold new skills** with a validated `SKILL.md` template
- **Initialize projects** with the `.agents/skills/` directory structure
- **Persistent configuration** (`~/.config/skills/config.yaml`) for default source and root
- **Script dependencies** — skills can declare helper scripts that are automatically copied on install
- **Skill dependencies** — skills can declare other skills they depend on, which are installed automatically
- **Remote GitHub support** — clone and install skills from any `owner/repo` reference

## Technical Constraints

- Written in Go 1.22+; single binary with no runtime dependencies
- Uses [Cobra](https://github.com/spf13/cobra) for CLI structure
- Skill format follows the [Agent Skills standard](https://agentskills.io/specification): YAML frontmatter + Markdown body
- GitHub clone is done via `git clone --depth 1` over SSH (requires `git` on PATH)
- Config stored as YAML at `~/.config/skills/config.yaml`

## Out of Scope

- Skill versioning / lockfiles (install always gets the latest)
- Private registry or package index beyond GitHub
- Agent runtime integration (loading skills at inference time)
- Skill publishing / release workflow

## Success Criteria

- `skills install git-commit-formatter --source marco-souza/skills` works end-to-end in a fresh project
- All commands are covered by unit tests (`go test ./...` passes)
- The CLI binary is installable via `go install github.com/marco-souza/skills@latest`
- Skill validation (name format, description length, required fields) is enforced on load
