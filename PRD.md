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

## Skill Collection Improvements

### Consistency & Standards

1. **Normalize YAML Frontmatter** — ✅ DONE. All 20 skills now use `description: >` consistently.
2. **Document `metadata` in AGENTS.md** — ✅ DONE. Added Metadata Field section documenting `scripts`, `runtime`, and `dependencies.skills`.
3. **Stick to Agent Skills Standard** — Do NOT add `tags`, `version`, or `scope` to frontmatter. Keep only `name` and `description` as the standard requires.

### New Skills

6. **`debug`** — Systematic debugging workflow (REPRO → GATHER → HYPOTHESIZE → TEST → FIX → VERIFY).
7. **`test`** — Testing best practices (test pyramid, Arrange-Act-Assert, mocking, coverage).
8. **`security-audit`** — Standalone security review (OWASP Top 10, dependency scanning, secret detection).
9. **`code-review`** — Local code review before pushing (readability, patterns, error handling).
10. **`docs`** — Documentation writing (README, API docs, ADRs, CHANGELOG).
11. **`refactor`** — Safe refactoring workflows (rename, extract, inline, move).

### CLI Improvements

12. **`skills validate`** — ✅ DONE. Validates frontmatter, name format, description length, and required fields. Use `--source` to validate skills in any directory.
13. **`skills graph`** — Visualize skill dependency graph (deferred).
14. **`skills available` vs `skills installed`** — Distinguish available vs installed skills (deferred).

### Content Improvements

15. **Expand `project-files` Quick-Start** — Add a simplified quick-start section at the top covering the 3 core files (PLAN.md, TODO.md, SESSION.md).
16. **Add Modern Tooling to `explore`** — Include `ripgrep`, `fd`, `bat`, `delta` alongside `grep`, `find`.
17. **Simplify `mixture-of-experts` Scripts** — Move inline bash to `.agents/scripts/` script files for maintainability.

## Out of Scope

- Skill versioning / lockfiles (install always gets the latest)
- Private registry or package index beyond GitHub
- Agent runtime integration (loading skills at inference time)
- Skill publishing / release workflow
- Splitting `mock-interview` into sub-skills (low priority)

## Implementation Phases

### Phase 1: Foundation (Consistency + Validation)

1. ~~Normalize YAML frontmatter across all 20 skills~~ ✅
2. ~~Document `metadata` field in AGENTS.md~~ ✅
3. ~~Add `skills validate` command to CLI~~ ✅
4. ~~Fix inconsistent `description:` usage in 3 skills~~ ✅

### Phase 2: New Skills + Content

5. ~~Add `debug` skill~~ ✅
6. ~~Add `test` skill~~ ✅
7. ~~Add `security-audit` skill~~ ✅
8. ~~Add `code-review` skill~~ ✅
9. ~~Add `docs` skill~~ ✅
10. ~~Add `refactor` skill~~ ✅
11. ~~Expand `project-files` with Quick-Start section~~ ✅
12. ~~Add modern tooling to `explore`~~ ✅
13. ~~Simplify `mixture-of-experts` scripts~~ ✅

## Success Criteria

- ✅ `skills install git-commit-formatter --source marco-souza/skills` works end-to-end
- ✅ All commands are covered by unit tests (`go test ./cmd/...` passes)
- ✅ The CLI binary is installable via `go install github.com/marco-souza/skills@latest`
- ✅ Skill validation (name format, description length, required fields) enforced on load
- ✅ All 20 skills in `.agents/skills/` have consistent YAML frontmatter
- ✅ 6 new skills added: `debug`, `test`, `security-audit`, `code-review`, `docs`, `refactor`
- ✅ AGENTS.md documents the `metadata` field
- ✅ `skills validate` command validates frontmatter, description, and required sections
