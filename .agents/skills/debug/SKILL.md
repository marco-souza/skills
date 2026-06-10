---
name: debug
description: >
  Systematically debug issues using a structured REPRO → GATHER → HYPOTHESIZE →
  TEST → FIX → VERIFY workflow. Use when diagnosing bugs, investigating failures,
  tracing errors, or troubleshooting unexpected behavior in code.
  Do NOT use when the issue is already identified and a fix is obvious,
  for code exploration without a specific problem, or for proactive code review.
---

# Debug

Apply a systematic, repeatable workflow to diagnose and resolve bugs.
Avoid guessing — every hypothesis must be grounded in evidence.

## When to Use

- A test is failing and the root cause is unclear
- An error occurs in production or development
- Behavior doesn't match expected output
- A regression was introduced after a recent change
- Performance degrades without obvious cause
- An integration or API call returns unexpected results

## When NOT to Use

- The bug is already identified and the fix is trivial
- You are exploring code without a specific problem to solve
- The task is feature implementation, not debugging
- The issue is a known limitation documented in project docs
- The user is asking for code review, not debugging

## Workflow

Follow these six stages in order. Do not skip ahead — each stage
builds on the previous one.

### Stage 1: REPRO — Reproduce the Bug

Before investigating, confirm you can trigger the bug reliably.

**Actions:**

1. Identify the exact steps that trigger the issue
2. Run the failing command, test, or code path
3. Capture the exact error message, stack trace, or unexpected output
4. Note the environment: OS, runtime version, dependencies, recent changes

**Commands:**

```bash
# Run the specific failing test
go test ./path/to/pkg -run TestFunctionName -v

# Re-run the command that triggers the error
<exact command the user reported>

# Check recent commits that may have introduced the bug
git log --oneline -10 -- <affected-file>

# Record environment details
go version
git log --oneline -3
```

**Checkpoint:** You can reproduce the bug consistently before moving on.

### Stage 2: GATHER — Collect Evidence

Read the relevant code, logs, and context surrounding the failure.

**Actions:**

1. Read the file(s) where the error originates
2. Trace the call stack upward to understand the flow
3. Check logs, error messages, and stack traces for clues
4. Review related tests for expected behavior
5. Look at git blame/log for recent changes to the area

**Commands:**

```bash
# Read the file where the error occurs
cat src/affected-file.ts

# Trace imports and callers
grep -rn "functionName" --include="*.ts" .

# Check git blame for recent changes
git log --oneline -5 -- src/affected-file.ts
git blame src/affected-file.ts | head -40

# Find related tests
grep -rn "TestFunctionName\|functionName" --include="*_test.go" .

# Search for error handling around the area
grep -B 5 -A 5 "errorMessage" src/affected-file.ts
```

**Checkpoint:** You understand what the code is doing, what it should do,
and where the gap is.

### Stage 3: HYPOTHESIZE — Form Theories

Based on gathered evidence, propose possible root causes ranked by likelihood.

**Actions:**

1. List 2–4 plausible hypotheses
2. For each, identify what evidence would confirm or rule it out
3. Prioritize: which hypothesis is easiest to test first?

**Template:**

```markdown
## Hypothesis 1 (Most Likely)
- **Theory:** [description of suspected root cause]
- **Evidence for:** [what you've observed that supports this]
- **Evidence against:** [what contradicts this]
- **How to test:** [specific action to confirm or rule out]

## Hypothesis 2
- **Theory:** [description]
- **Evidence for:** [supporting observations]
- **Evidence against:** [contradicting observations]
- **How to test:** [specific action]
```

**Checkpoint:** You have at least one testable hypothesis with a clear
verification step.

### Stage 4: TEST — Validate Hypotheses

Run targeted experiments to confirm or eliminate each hypothesis.

**Actions:**

1. Start with the most likely hypothesis
2. Add debugging output, assertions, or breakpoints as needed
3. Run minimal experiments — change one thing at a time
4. Record results for each hypothesis

**Commands:**

```bash
# Add debug logging temporarily
echo "DEBUG: variableName = $variableName" >&2

# Run with verbose output
go test ./path -v -run TestName

# Use a debugger or print statements to inspect state
# Check intermediate values at key points in the code

# If hypothesis is about a specific commit, test before/after
git stash          # revert changes
go test ./...      # does it pass?
git stash pop      # restore changes
go test ./...      # does it fail?
```

**Checkpoint:** You have confirmed the root cause through targeted testing.

### Stage 5: FIX — Implement the Solution

Apply the minimal fix that addresses the confirmed root cause.

**Actions:**

1. Write the smallest change that fixes the bug
2. Avoid unrelated refactors — keep the fix focused
3. Add a regression test if one doesn't exist
4. Check for similar issues in nearby code

**Guidelines:**

- Fix the root cause, not just the symptom
- Prefer explicit error handling over silent failures
- If the fix is complex, break it into smaller changes
- Document why the bug occurred if the reason is non-obvious

### Stage 6: VERIFY — Confirm the Fix

Validate that the bug is resolved and nothing else broke.

**Actions:**

1. Re-run the exact reproduction steps from Stage 1
2. Run the full test suite to catch regressions
3. Verify edge cases related to the fix
4. Clean up any temporary debugging code

**Commands:**

```bash
# Re-run the original failing test
go test ./path/to/pkg -run TestFunctionName -v

# Run the full test suite
go test ./...

# Run vet/lint to catch additional issues
go vet ./...

# Check for leftover debug code
grep -rn "DEBUG\|TODO.*fix\|HACK" --include="*.go" .
```

**Checkpoint:** The bug no longer reproduces, all tests pass, and no
debugging artifacts remain.

## Examples

### Example 1: Failing Unit Test

**Scenario:** `TestUserCreate` returns a nil pointer error.

**REPRO:**

```bash
go test ./internal/user -run TestUserCreate -v
# Output: panic: runtime error: invalid memory address
#   at user.Create(): user.go:42
```

**GATHER:**

```bash
cat internal/user/user.go | head -50
# Line 42: result.Email = input.Email
# result is the return value of repository.FindByEmail()

grep -rn "FindByEmail" --include="*.go" .
# Found in internal/user/repository.go:28

cat internal/user/repository.go | sed -n '25,35p'
# FindByEmail returns nil, nil when no user is found (correct behavior)
# but caller doesn't check for nil before accessing .Email
```

**HYPOTHESIS:** The `Create` function calls `FindByEmail` to check for
duplicates, but doesn't handle the nil-return case when no existing user
is found.

**TEST:** Add a nil check before line 42 and re-run the test.

**FIX:** Add `if result != nil { return nil, ErrAlreadyExists }` before
accessing `result.Email`.

**VERIFY:**

```bash
go test ./internal/user -run TestUserCreate -v  # passes
go test ./internal/user/... -v                   # all user tests pass
```

### Example 2: Production API Error

**Scenario:** API returns 500 on `POST /api/orders` intermittently.

**REPRO:**

```bash
curl -X POST localhost:3000/api/orders \
  -H "Content-Type: application/json" \
  -d '{"items": [{"id": 1, "qty": 1}]}'
# {"error": "internal server error"} — fails ~30% of the time
```

**GATHER:**

```bash
# Check server logs
tail -f logs/app.log | grep "order"
# Intermittent: "ERROR: connection refused to database"

# Check database connection pool config
cat config/database.go | grep -A 5 "Pool"
# MaxOpenConns: 5 (too low for production traffic)

# Check for long-running queries
grep -rn "SELECT" --include="*.go" src/order/ | head -10
```

**HYPOTHESIS:** The connection pool is exhausted under load, causing
intermittent connection failures. The 30% failure rate matches peak
traffic patterns.

**TEST:** Increase pool to 20 and load-test locally.

**FIX:** Update `MaxOpenConns` to 20 and add connection timeout config.

**VERIFY:**

```bash
# Load test
ab -n 1000 -c 50 http://localhost:3000/api/orders
# All requests return 200, no connection errors in logs
```

## Edge Cases

### Cannot Reproduce

If the bug cannot be reproduced:

1. Ask the user for exact steps, environment details, and error output
2. Check for environment-specific issues (OS, versions, config)
3. Look for race conditions or timing-dependent bugs
4. Search for similar issues in project issue tracker or upstream docs

### Bug Is in a Dependency

If the root cause is in a third-party library:

1. Document the library version and the specific issue
2. Check if a newer version fixes it
3. Consider workarounds (monkey-patching, alternative libraries)
4. File an upstream issue if no fix exists

### Multiple Root Causes

If several independent issues contribute to the bug:

1. Fix them one at a time, verifying after each fix
2. Prioritize by impact — fix the most severe first
3. Create separate commits for each fix if they are independent

## Best Practices

- **Reproduce before fixing** — never guess at a fix without confirming the bug
- **One change at a time** — isolate variables to avoid confounding results
- **Write regression tests** — ensure the bug stays fixed
- **Keep a decision log** — note hypotheses tested and why they were ruled out
- **Clean up after yourself** — remove debug code before committing
- **Ask for help early** — if stuck after 30 minutes, consult a teammate or spawn
  a subagent for a second perspective
