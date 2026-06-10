---
name: mixture-of-experts
description: >
  Solve complex problems by spawning multiple specialized expert agents
  that analyze from different angles, then synthesize their insights.
  Use for architecture decisions, code reviews, complex debugging,
  or when you need comprehensive analysis.
metadata:
  scripts:
    - ../../scripts/moe-spawn.sh
    - ../../scripts/moe-wait.sh
    - ../../scripts/moe-aggregate.sh
    - ../../scripts/moe-cleanup.sh
    - ../../scripts/moe-code-review.sh
  dependencies:
    skills:
      - spawn-subagents
      - terminal-multiplexer
---

# Mixture of Experts (MoE)

Spawn multiple specialized experts in parallel, each analyzing from a unique
angle. Aggregate their insights into a comprehensive, multi-dimensional answer.

> **Foundation:** This skill builds on `spawn-subagents` for the underlying
> tmux/pi subagent mechanics. Read `spawn-subagents` for the core patterns
> before using MoE.

## When to Use

- Architecture or design decisions
- Complex code reviews
- Security audits
- Performance optimization
- Debugging tricky issues
- Evaluating trade-offs

## Expert Specializations

### Core Experts

| Expert | Focus | System Prompt |
|--------|-------|---------------|
| `architect` | Design patterns, coupling, cohesion, long-term maintainability | "You are a software architect. Focus on design patterns, separation of concerns, and long-term maintainability." |
| `security` | Vulnerabilities, injection risks, auth flaws, data exposure | "You are a security engineer. Focus on vulnerabilities, injection risks, authentication, and data exposure." |
| `performance` | Algorithmic complexity, resource usage, bottlenecks | "You are a performance engineer. Focus on time/space complexity, resource usage, and optimization opportunities." |
| `maintainer` | Readability, documentation, testing, onboarding cost | "You are a senior maintainer. Focus on code readability, documentation, testing coverage, and onboarding new developers." |
| `minimalist` | Simplicity, YAGNI, removing unnecessary complexity | "You are a minimalist engineer. Focus on simplifying, removing unnecessary code, and YAGNI violations." |

### Domain Experts

| Expert | Focus |
|--------|-------|
| `api-designer` | REST/GraphQL conventions, versioning, backward compatibility |
| `data-modeler` | Schema design, normalization, query patterns, migrations |
| `dx-specialist` | Developer experience, tooling, error messages, debugging |
| `ops-engineer` | Deployment, monitoring, observability, rollback strategies |

## Scripts

This skill provides helper scripts in `.agents/scripts/` for automating the MoE workflow:

| Script | Purpose |
|--------|---------|
| `moe-spawn.sh` | Spawn expert agents in parallel |
| `moe-wait.sh` | Wait for all experts to complete |
| `moe-aggregate.sh` | Combine results and run aggregator |
| `moe-cleanup.sh` | Kill sessions and remove temp files |
| `moe-code-review.sh` | Complete code review workflow |

## Workflow

### 1. Define the Problem

Create a clear, specific prompt that all experts will analyze:

```bash
PROBLEM="Review the authentication flow in src/auth/ for issues and improvements"
TASK_ID=$(date +%s)
```

### 2. Spawn Experts in Parallel

Use `moe-spawn.sh` to spawn experts:

```bash
source .agents/scripts/moe-spawn.sh
spawn_moe_experts "architect security performance maintainer"
```

### 3. Wait for All Experts

Use `moe-wait.sh` to wait for completion:

```bash
source .agents/scripts/moe-wait.sh
wait_moe_experts "architect security performance maintainer"
```

### 4. Aggregate Results

Use `moe-aggregate.sh` to combine insights:

```bash
source .agents/scripts/moe-aggregate.sh
aggregate_moe_results "architect security performance maintainer"
```

### 5. Cleanup

Use `moe-cleanup.sh` to remove sessions and temp files:

```bash
source .agents/scripts/moe-cleanup.sh
cleanup_moe "architect security performance maintainer"
```

## Complete Example: Code Review

Use the `moe-code-review.sh` script for a complete review workflow:

```bash
# Run MoE code review on a specific file
./.agents/scripts/moe-code-review.sh src/auth/login.ts
```

The script will:
1. Read the file content
2. Spawn architect, security, performance, and maintainer experts
3. Wait for all experts to complete
4. Synthesize results into actionable recommendations
5. Clean up sessions and temp files

## Advanced: Weighted Aggregation

For weighted aggregation, modify the aggregator prompt in `moe-aggregate.sh`:

```bash
# Define weights
architect=3
security=3
performance=2
maintainer=2

# Build weighted prompt
WEIGHTED_PROMPT="Synthesize with these expert weights:\n"
for expert in $EXPERTS; do
  weight=$(eval echo \$$expert)
  WEIGHTED_PROMPT="$WEIGHTED_PROMPT\n- $expert (weight: $weight/10)"
done
WEIGHTED_PROMPT="$WEIGHTED_PROMPT\n\nHigher weight = more influence on final recommendation."
```

## Best Practices

- **Choose 3-5 experts** — Too few misses angles, too many adds noise
- **Make prompts specific** — Generic prompts yield generic answers
- **Include file content** — Don't make experts hunt for context
- **Define clear aggregation strategy** — Consensus-based, weighted, or hierarchical
- **Cache expert outputs** — Save to files for inspection if aggregation fails
- **Set timeouts** — Kill hung experts after 5 minutes

## Common Patterns

| Scenario | Expert Mix |
|----------|------------|
| API Design | `architect`, `api-designer`, `security`, `dx-specialist` |
| Database Schema | `data-modeler`, `performance`, `architect` |
| Frontend Component | `maintainer`, `minimalist`, `performance`, `dx-specialist` |
| DevOps Pipeline | `ops-engineer`, `security`, `maintainer` |
| Full Feature | `architect`, `security`, `performance`, `maintainer`, `minimalist` |

## Limitations

- Expert outputs may conflict (aggregator must resolve)
- Token cost scales linearly with expert count
- No cross-expert communication during analysis
- Synthesis quality depends on aggregator prompt quality
