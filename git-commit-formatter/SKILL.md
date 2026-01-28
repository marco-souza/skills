---
name: git-commit-formatter
description: Formats git commit messages according to Conventional Commits specification. Use this when the user asks to commit changes or write a commit message.
---

# Git Commit Formatter Skill

When writing a git commit message, you MUST follow the Conventional Commits specification.

Format

`<type>(optional scope): <description>`

Allowed Types

- **fea**: A new feature
- **fix**: A bug fix
- **ref**: A code change that neither fixes a bug nor adds a feature
- **ai**: AI-related changes
- **test**: Adding missing tests or correcting existing tests
- **oth**: Anything that doesn't fit into the other types

Instructions

1. Analyze the changes to determine the primary `type`.
2. Identify the `scope` if applicable (e.g., specific component or file).
3. Write a concise `description` in an imperative mood (e.g., "add feature" not "added feature").
4. If there are breaking changes, add a footer starting with `BREAKING CHANGE:`.

Example
`feat(auth): implement login with google`
