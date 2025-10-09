# Testing Guide

This document describes how to run tests for the Stellar RPC Indexer.

## Quick Start

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# View coverage in browser
make test-coverage-html
```

## Test Structure

```
indexer/
├── pkg/
│   ├── client/
│   │   ├── rpc.go
│   │   └── rpc_test.go          # Client tests
│   ├── parser/
│   │   ├── parser.go
│   │   ├── parser_test.go        # Parser tests
│   │   ├── token_operations.go
│   │   └── token_operations_test.go
│   └── models/
│       ├── event.go
│       ├── transaction.go
│       └── models_test.go        # Database model tests
└── Makefile
```

## Available Test Commands

### Basic Testing

```bash
# Run all tests
make test

# Run tests with race detector
make test-verbose

# Run short tests only (skip slow tests)
make test-short

# Run tests for specific package
make test-package PKG=./pkg/parser
```

### Coverage

```bash
# Generate coverage report
make test-coverage

# View coverage in browser
make test-coverage-html

# Check coverage percentage
go test -cover ./...
```

### Benchmarks

```bash
# Run benchmarks
make benchmark

# Run benchmarks for specific package
go test -bench=. ./pkg/parser
```

## Test Categories

### 1. Parser Tests (`pkg/parser/*_test.go`)

Tests for XDR parsing and data transformation:

- ✅ `TestParseEvent` - Event parsing from RPC
- ✅ `TestScValToInterface` - ScVal conversion
- ✅ `TestParseTransaction` - Transaction parsing
- ✅ `TestParseTokenOperation_*` - Token operation extraction
- ✅ `TestParseTokenMetadata` - Token metadata extraction
- ✅ `TestParseTokenBalance` - Balance parsing

**Run parser tests:**
```bash
make test-package PKG=./pkg/parser
```

### 2. Client Tests (`pkg/client/rpc_test.go`)

Tests for RPC client and JSON-RPC communication:

- ✅ `TestClient_GetLatestLedger` - Fetch latest ledger
- ✅ `TestClient_GetEvents` - Event fetching
- ✅ `TestClient_GetTransaction` - Transaction fetching
- ✅ `TestClient_Health` - Health check endpoint
- ✅ `TestClient_GetLedgerEntries` - Ledger entry fetching
- ✅ Error handling tests

**Run client tests:**
```bash
make test-package PKG=./pkg/client
```

### 3. Model Tests (`pkg/models/models_test.go`)

Tests for database models and GORM operations:

- ✅ `TestUpsertEvent` - Event database operations
- ✅ `TestUpsertTransaction` - Transaction storage
- ✅ `TestCursor` - Cursor management
- ✅ `TestUpsertTokenMetadata` - Token metadata storage
- ✅ `TestUpsertTokenOperation` - Token operation storage
- ✅ `TestUpsertTokenBalance` - Balance storage
- ✅ `TestInt128` - Custom Int128 type

**Run model tests:**
```bash
make test-package PKG=./pkg/models
```

## Writing New Tests

### Test File Naming

- Test files must end with `_test.go`
- Place test files next to the code they test
- Example: `parser.go` → `parser_test.go`

### Test Function Naming

```go
func TestFunctionName(t *testing.T) {
    // Test single behavior
}

func TestFunctionName_SpecificCase(t *testing.T) {
    // Test specific edge case
}
```

### Table-Driven Tests

Use table-driven tests for multiple scenarios:

```go
func TestParse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Result
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    Result{Value: "test"},
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Parse(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Parse() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Helpers

```go
// setupTestDB creates in-memory database for tests
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to open test database: %v", err)
    }
    return db
}
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      - name: Run tests
        run: make ci
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

### Pre-commit Hook

Add to `.git/hooks/pre-commit`:

```bash
#!/bin/bash
cd indexer
make pre-commit
```

## Test Data

### Mock RPC Server

Tests use `httptest.NewServer` for mocking RPC responses:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Mock response
    resp := jsonRPCResponse{
        Result: json.RawMessage(`{"status":"healthy"}`),
    }
    json.NewEncoder(w).Encode(resp)
}))
defer server.Close()

client := NewClient(server.URL)
```

### In-Memory Database

Model tests use SQLite in-memory database:

```go
db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
```

## Debugging Tests

### Verbose Output

```bash
# See all test output
go test -v ./pkg/parser

# See only failed tests
go test ./pkg/parser
```

### Run Single Test

```bash
# Run specific test function
go test -v -run TestParseSingleTest ./pkg/parser

# Run tests matching pattern
go test -v -run TestParse.* ./pkg/parser
```

### Debug with Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug specific test
dlv test ./pkg/parser -- -test.run TestParseEvent
```

## Best Practices

### ✅ Do

- Write tests for all public functions
- Use table-driven tests for multiple scenarios
- Test error cases
- Keep tests independent (no shared state)
- Use meaningful test names
- Mock external dependencies

### ❌ Don't

- Skip cleanup (`defer server.Close()`)
- Use hardcoded ports or files
- Test implementation details
- Write flaky tests
- Ignore race conditions

## Coverage Goals

- **Overall:** 80%+ coverage
- **Critical paths:** 90%+ coverage
  - Parser functions
  - Database operations
  - RPC client

## Performance Testing

### Benchmarks

```bash
# Run all benchmarks
make benchmark

# Run specific benchmark
go test -bench=BenchmarkParseEvent ./pkg/parser

# With memory stats
go test -bench=. -benchmem ./pkg/parser
```

### Example Benchmark

```go
func BenchmarkParseEvent(b *testing.B) {
    event := client.Event{/* test data */}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ParseEvent(event)
    }
}
```

## Troubleshooting

### Tests fail with import errors

```bash
make deps
make tidy
```

### Database connection errors in tests

Tests use in-memory SQLite, no external DB needed.

### Timeout errors

Increase test timeout:
```bash
go test -timeout 30s ./...
```

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Test Coverage](https://go.dev/blog/cover)
- [GORM Testing](https://gorm.io/docs/testing.html)
