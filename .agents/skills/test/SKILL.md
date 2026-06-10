---
name: test
description: >
  Guide for writing effective tests using best practices: test pyramid, Arrange-Act-Assert pattern,
  mocking strategies, and coverage guidelines. Use when the user asks to write tests, improve test
  coverage, refactor tests, or understand testing best practices.
  Do NOT use when the user wants to run existing tests without modification, or when the focus is
  on test infrastructure setup (CI/CD, test frameworks installation).
---

# Testing Best Practices

Write reliable, maintainable tests that give confidence in code correctness.

## When to Use

- Writing new tests for features or bug fixes
- Improving test coverage in under-tested areas
- Refactoring flaky or brittle tests
- Reviewing test quality during code review
- Adding regression tests for discovered bugs

## When NOT to Use

- Running tests without modifications
- Setting up test infrastructure (CI/CD, framework installation)
- Performance/load testing (use dedicated load testing tools)

## Test Pyramid

Structure your test suite using the test pyramid model:

```
        /\
       /  \         Unit Tests (Fast, many)
      /    \
     /------\
    /        \      Integration Tests (Medium, fewer)
   /----------\
  /            \    E2E Tests (Slow, few)
 /--------------\
```

### Unit Tests (Base Layer)

- **Purpose**: Test individual functions/methods in isolation
- **Speed**: Milliseconds
- **Quantity**: Most tests should be unit tests
- **Isolation**: No external dependencies (network, database, filesystem)

```go
// Good: Pure unit test, no dependencies
func TestCalculateTotal(t *testing.T) {
    items := []Item{
        {Price: 10.00, Qty: 2},
        {Price: 5.00, Qty: 1},
    }
    
    got := CalculateTotal(items)
    want := 25.00
    
    if got != want {
        t.Errorf("CalculateTotal() = %v, want %v", got, want)
    }
}
```

### Integration Tests (Middle Layer)

- **Purpose**: Test component interactions and external integrations
- **Speed**: Seconds
- **Quantity**: Moderate
- **Isolation**: May use test databases, mock external APIs

```go
// Integration test with database
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewUserRepository(db)
    user := &User{Name: "Alice", Email: "alice@example.com"}
    
    err := repo.Create(user)
    if err != nil {
        t.Fatalf("Create() error = %v", err)
    }
    
    if user.ID == 0 {
        t.Error("expected user ID to be set")
    }
}
```

### End-to-End Tests (Top Layer)

- **Purpose**: Test complete user workflows through the full stack
- **Speed**: Seconds to minutes
- **Quantity**: Few, focused on critical paths
- **Isolation**: Full environment required

```go
// E2E test hitting real HTTP endpoint
func TestLoginEndpoint(t *testing.T) {
    server := startTestServer(t)
    defer server.Close()
    
    resp, err := http.Post(
        server.URL+"/api/login",
        "application/json",
        strings.NewReader(`{"email":"alice@example.com","password":"secret"}`),
    )
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
    }
}
```

## Arrange-Act-Assert Pattern

Structure every test with three clear phases:

```go
func TestDiscountCalculation(t *testing.T) {
    // ARRANGE: Set up test data and dependencies
    product := &Product{Price: 100.00}
    discount := 0.20 // 20% off
    
    // ACT: Execute the function under test
    result := CalculateDiscount(product.Price, discount)
    
    // ASSERT: Verify the expected outcome
    want := 80.00
    if result != want {
        t.Errorf("CalculateDiscount() = %v, want %v", result, want)
    }
}
```

### Why AAA?

- **Readability**: Clear separation of concerns
- **Maintainability**: Easy to identify which phase breaks
- **Documentation**: Test intent is obvious

### AAA Variants

```go
// Given-When-Then (BDD style)
func TestLoginValidation(t *testing.T) {
    // Given: valid credentials
    credentials := Credentials{Email: "alice@example.com", Password: "valid123"}
    
    // When: user attempts login
    result, err := AuthService.Login(credentials)
    
    // Then: login succeeds
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if result.Token == "" {
        t.Error("expected token to be set")
    }
}

// Arrange-Act-Assert-Cleanup (for resources)
func TestFileProcessing(t *testing.T) {
    // Arrange
    tmpFile := createTempFile(t, "test data")
    defer os.Remove(tmpFile)
    
    // Act
    result, err := ProcessFile(tmpFile)
    
    // Assert
    if err != nil {
        t.Fatalf("ProcessFile() error = %v", err)
    }
    if result.Processed != true {
        t.Error("expected file to be marked as processed")
    }
}
```

## Mocking and Stubbing

### When to Mock

- External services (APIs, databases, file systems)
- Time-dependent code
- Random/non-deterministic values
- Slow operations

### Interface-Based Mocking (Go)

Define interfaces for dependencies:

```go
// Production code depends on interface, not concrete type
type UserStore interface {
    GetByID(id int64) (*User, error)
    Save(user *User) error
}

// Mock implementation
type MockUserStore struct {
    Users  map[int64]*User
    SaveFn func(user *User) error
}

func (m *MockUserStore) GetByID(id int64) (*User, error) {
    user, ok := m.Users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return user, nil
}

func (m *MockUserStore) Save(user *User) error {
    if m.SaveFn != nil {
        return m.SaveFn(user)
    }
    m.Users[user.ID] = user
    return nil
}
```

### Using Mocks in Tests

```go
func TestGetUserProfile(t *testing.T) {
    // Arrange
    mockStore := &MockUserStore{
        Users: map[int64]*User{
            1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
        },
    }
    service := NewUserService(mockStore)
    
    // Act
    profile, err := service.GetProfile(1)
    
    // Assert
    if err != nil {
        t.Fatalf("GetProfile() error = %v", err)
    }
    if profile.Name != "Alice" {
        t.Errorf("Name = %q, want %q", profile.Name, "Alice")
    }
}

func TestSaveUser_RejectsInvalidEmail(t *testing.T) {
    // Arrange
    mockStore := &MockUserStore{
        Users: make(map[int64]*User),
    }
    service := NewUserService(mockStore)
    user := &User{Name: "Bob", Email: "invalid"}
    
    // Act
    err := service.Save(user)
    
    // Assert
    if err == nil {
        t.Error("expected error for invalid email, got nil")
    }
}
```

### Table-Driven Tests

For testing multiple cases efficiently:

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {name: "valid email", email: "alice@example.com", wantErr: false},
        {name: "missing @", email: "aliceexample.com", wantErr: true},
        {name: "missing domain", email: "alice@", wantErr: true},
        {name: "empty string", email: "", wantErr: true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail(%q) error = %v, wantErr %v",
                    tt.email, err, tt.wantErr)
            }
        })
    }
}
```

## Coverage Guidelines

### Coverage Targets

| Component Type | Minimum Coverage | Rationale |
|----------------|------------------|-----------|
| Business logic | 80-90% | Core correctness |
| Utilities/helpers | 90%+ | Reused everywhere |
| API handlers | 70-80% | HTTP concerns |
| UI components | 50-70% | Harder to test, visual |
| Configuration | 30-50% | Often trivial |

### What to Test

**Always test:**
- Happy paths (expected inputs)
- Error cases (invalid inputs, edge cases)
- Boundary conditions (empty, nil, max values)
- State transitions

**Don't obsess over:**
- Getters/setters
- Framework boilerplate
- Third-party library internals

### Coverage in Go

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML report for detailed view
go tool cover -html=coverage.out -o coverage.html
```

### Reading Coverage Output

```
ok  github.com/example/pkg  0.5s  coverage: 85.2% of statements
```

Focus on increasing coverage for files with lowest percentages:

```bash
# Show coverage per function
go tool cover -func=coverage.out | sort -k3 -n
```

## Test Organization

### File Structure

```
src/
├── user.go              # Production code
├── user_test.go         # Unit tests (same package)
├── user_integration_test.go  # Integration tests
└── testdata/
    └── fixtures/        # Test data files
```

### Naming Conventions

```go
// Function under test: CalculateDiscount

// Standard test (success case)
func TestCalculateDiscount(t *testing.T) { ... }

// Failure/edge case
func TestCalculateDiscount_ZeroDiscount(t *testing.T) { ... }
func TestCalculateDiscount_NegativeAmount(t *testing.T) { ... }

// Benchmark
func BenchmarkCalculateDiscount(b *testing.B) { ... }

// Example (documentation)
func ExampleCalculateDiscount() { ... }
```

## Edge Cases to Cover

### Nil/Empty Inputs

```go
func TestProcessItems_EmptySlice(t *testing.T) {
    result := ProcessItems([]Item{})
    if len(result) != 0 {
        t.Error("expected empty result for empty input")
    }
}

func TestProcessItems_NilSlice(t *testing.T) {
    result := ProcessItems(nil)
    if result != nil {
        t.Error("expected nil result for nil input")
    }
}
```

### Boundary Values

```go
func TestAgeValidation(t *testing.T) {
    tests := []struct {
        age     int
        wantErr bool
    }{
        {age: 0, wantErr: true},    // Minimum boundary
        {age: 1, wantErr: false},   // Just above minimum
        {age: 17, wantErr: false},  // Below adult age
        {age: 18, wantErr: false},  // Adult threshold
        {age: 150, wantErr: false}, // Reasonable maximum
        {age: 151, wantErr: true},  // Above reasonable max
    }
    // ... test loop
}
```

### Concurrent Access

```go
func TestConcurrentAccess(t *testing.T) {
    cache := NewCache()
    var wg sync.WaitGroup
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            cache.Set(fmt.Sprintf("key%d", n), n)
            _ = cache.Get(fmt.Sprintf("key%d", n))
        }(i)
    }
    
    wg.Wait()
    // Cache should not panic or corrupt data
}
```

## Common Anti-Patterns

### Avoid

- **Tests that depend on execution order** — each test should be independent
- **Mocking everything** — prefer real implementations when practical
- **Testing implementation details** — test behavior, not internals
- **Sleep-based synchronization** — use channels or callbacks
- **Ignored errors** — always check error returns

### Good vs Bad Test

```go
// BAD: Tests implementation detail, brittle
func TestUserRepository_InternalCache(t *testing.T) {
    repo := NewUserRepository(db)
    repo.cache = map[int64]*User{}  // Testing internal state
    
    user := repo.GetUser(1)
    if _, exists := repo.cache[1]; !exists {
        t.Error("expected cache to be populated")
    }
}

// GOOD: Tests observable behavior
func TestUserRepository_GetUser_CachesResult(t *testing.T) {
    db := &MockDB{QueryCount: 0}
    repo := NewUserRepository(db)
    
    // First call hits database
    _, _ = repo.GetUser(1)
    if db.QueryCount != 1 {
        t.Errorf("expected 1 query, got %d", db.QueryCount)
    }
    
    // Second call uses cache (no additional query)
    _, _ = repo.GetUser(1)
    if db.QueryCount != 1 {
        t.Errorf("expected still 1 query, got %d", db.QueryCount)
    }
}
```

## Quick Reference

### Test Command

```bash
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -run TestFunctionName    # Run specific test
go test -count=3 ./...           # Run 3 times (find flaky tests)
go test -race ./...              # Detect race conditions
```

### Test Flags

```bash
go test -cover                   # Show coverage percentage
go test -coverprofile=out.out    # Save coverage data
go test -bench=.                 # Run benchmarks
go test -short                   # Skip long-running tests
go test -timeout 30s             # Fail if test exceeds timeout
```

## Best Practices

1. **One assertion per concept** — each test should verify one behavior
2. **Descriptive names** — test names should explain what's being tested
3. **Independent tests** — no test should depend on another test's state
4. **Fast feedback** — keep unit tests under 100ms each
5. **Test behavior, not implementation** — verify outcomes, not methods called
6. **Keep tests readable** — a test is documentation
7. **Refactor tests** — extract helpers for repeated setup
8. **Delete dead tests** — if code is removed, remove its tests too
