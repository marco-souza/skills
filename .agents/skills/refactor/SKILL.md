---
name: refactor
description: >
  Perform safe refactoring operations while preserving behavior and maintaining code integrity.
  Use when: renaming identifiers, extracting functions/modules, inlining code, moving files,
  or restructuring code without changing functionality.
  Do NOT use when: adding new features, fixing bugs (use standard editing), or making
  breaking API changes.
---

# Refactor

Execute safe refactoring operations with verification at each step. Refactoring changes code structure without altering external behavior.

## Core Principle

**Preserve behavior.** Every refactoring must pass tests before and after. If behavior changes, it's not refactoring—it's a rewrite.

## Refactoring Operations

### Rename

Rename variables, functions, classes, files, or directories while maintaining all references.

**When to use:**
- Identifier name is unclear or violates naming conventions
- Code review feedback on naming
- Aligning names with business domain

**Steps:**
1. Find all occurrences (definitions + usages)
2. Verify the scope (global vs local)
3. Rename atomically (all at once)
4. Run tests to confirm behavior unchanged

```bash
# Find all usages before renaming
grep -r "oldName" --include="*.ts" --include="*.js" --include="*.py"
```

### Extract

Pull out a code block into a new function, method, class, or module.

**When to use:**
- Function is too long (>30 lines)
- Code block is duplicated elsewhere
- Logic can be reused in other contexts
- Improving readability by naming complex operations

**Steps:**
1. Identify the code block to extract
2. Determine inputs (parameters) and outputs (return value)
3. Create the new function/module
4. Replace original code with a call to the new function
5. Verify tests still pass

**Extraction targets:**
- Function/method extraction
- Class extraction
- Module/file extraction
- Interface/type extraction

### Inline

Replace a function call with the function body.

**When to use:**
- Function is trivial (one-liner)
- Function is called only once
- Indirection adds unnecessary complexity
- Making code easier to follow

**Steps:**
1. Verify the function has no side effects
2. Copy function body to call site
3. Replace arguments with actual values
4. Remove the now-unused function
5. Run tests

### Move

Relocate code to a more appropriate location in the codebase.

**When to use:**
- Code belongs in a different module/package
- Organizing by feature instead of type
- Reducing coupling between modules
- Aligning with project architecture

**Steps:**
1. Identify the target location
2. Update all import paths
3. Ensure no circular dependencies
4. Move the file/code
5. Update tests and documentation

## Safety Checklist

Before any refactoring, verify:

```bash
# 1. Ensure clean git state
git status
git stash  # or commit work-in-progress

# 2. Run existing tests to establish baseline
# Pick the appropriate test command for your project
npm test 2>/dev/null || go test ./... 2>/dev/null || pytest 2>/dev/null || cargo test 2>/dev/null

# 3. Create a checkpoint
git checkout -b refactor/<description>
```

**Pre-refactoring checklist:**
- [ ] All tests pass on current code
- [ ] Git working directory is clean (or changes stashed)
- [ ] You understand the code's current behavior
- [ ] You've identified all usages of code being refactored
- [ ] You have a rollback plan (git)

## Step-by-Step Workflow

### Phase 1: Preparation

1. **Understand current behavior**
   - Read the code thoroughly
   - Run existing tests
   - Note edge cases and dependencies

2. **Plan the refactoring**
   - Choose the operation type (rename/extract/inline/move)
   - Identify all affected files
   - Estimate risk level

3. **Set up safety net**
   ```bash
   # Ensure tests exist
   # If no tests, consider adding characterization tests first
   
   # Create refactoring branch
   git checkout -b refactor/<description>
   ```

### Phase 2: Execution

4. **Make the change incrementally**
   - One logical change per commit
   - Small, reviewable steps
   - Never mix refactoring with feature changes

5. **Verify after each step**
   ```bash
   # Run tests
   npm test 2>/dev/null || go test ./...
   
   # Verify build succeeds
   npm run build 2>/dev/null || go build ./...
   ```

6. **Commit frequently**
   ```bash
   git add -A
   git commit -m "ref: <what changed>"
   ```

### Phase 3: Verification

7. **Run full test suite**
   ```bash
   # Comprehensive verification
   npm test 2>/dev/null && npm run lint 2>/dev/null
   ```

8. **Manual smoke test** (if applicable)
   - Start the application
   - Verify critical paths work
   - Check edge cases

9. **Clean up**
   ```bash
   # Remove any dead code
   # Update documentation if needed
   # Merge or create PR
   ```

## Common Patterns

### Pattern: Extract Function

```typescript
// Before
function processOrder(order: Order) {
  // 50 lines of validation
  // 30 lines of calculation
  // 20 lines of persistence
}

// After
function processOrder(order: Order) {
  const validated = validateOrder(order);
  const calculated = calculateTotals(validated);
  return persistOrder(calculated);
}
```

### Pattern: Move to Module

```bash
# Before: utils.js has 500 lines
# After: Split into organized modules

mkdir -p src/utils/validation
mkdir -p src/utils/formatting

# Move validation functions
mv src/utils/validateUser.js src/utils/validation/
mv src/utils/validateOrder.js src/utils/validation/

# Update imports
find . -name "*.ts" -exec sed -i 's|from.*utils/validateUser|from../utils/validation/validateUser|g' {} \;
```

### Pattern: Rename with IDE Support

```bash
# For large codebases, use language server
# TypeScript
npx tsserver --rename <file> <oldName> <newName>

# Or use find + sed for simpler cases
grep -rl "oldFunctionName" src/ | xargs sed -i 's/oldFunctionName/newFunctionName/g'
```

## Edge Cases

### Circular Dependencies

When moving code creates circular imports:
1. Extract shared types/interfaces to a separate module
2. Use dependency injection
3. Create a facade module

### Breaking Changes

If refactoring must change public API:
1. Create new API alongside old
2. Deprecate old API
3. Migrate callers gradually
4. Remove old API in future release

### Large-Scale Refactoring

For codebase-wide changes:
1. Script the change (sed, ast-grep, codemods)
2. Run on entire codebase at once
3. Commit as single atomic change
4. Verify everything still works

## Troubleshooting

### Tests Fail After Refactoring

```bash
# 1. Check what changed
git diff HEAD~1

# 2. Verify test expectations
# Look for tests that assert implementation details

# 3. Consider if test needs updating (rare)
# Only if refactoring exposed test as brittle
```

### Build Errors

```bash
# TypeScript: Check for missing imports
npx tsc --noEmit

# Go: Check for unused imports
goimports -l .

# General: Search for references to old names
grep -r "oldName" --include="*.ts" --include="*.js" --include="*.go"
```

### Performance Regression

```bash
# Profile before and after
# Compare key metrics
# Ensure refactoring didn't introduce N+1 queries or unnecessary copies
```

## Best Practices

1. **Small steps** — One logical change per commit
2. **Tests first** — Never refactor without passing tests
3. **Version control** — Use branches for risky refactoring
4. **No behavior changes** — Refactor only, then change behavior separately
5. **Document intent** — Commit messages explain WHY, not WHAT
6. **Review carefully** — Refactoring bugs are subtle
7. **Incremental rollout** — For large changes, merge progressively
