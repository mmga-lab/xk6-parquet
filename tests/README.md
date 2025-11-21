# Tests

This directory contains test files and test data for xk6-parquet.

## Directory Structure

```
tests/
├── unit/          # Unit tests (currently in pkg/parquet/*_test.go)
├── integration/   # Integration tests
└── testdata/      # Test Parquet files and fixtures
```

## Running Tests

### Unit Tests

Unit tests are located alongside the source code in `pkg/parquet/`:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.txt ./...

# View coverage report
go tool cover -html=coverage.txt

# Run specific test
go test -v -run TestCacheSetAndGet ./pkg/parquet
```

### Integration Tests

Integration tests will test the extension with actual k6:

```bash
# Build k6 with extension
xk6 build --with github.com/mmga-lab/xk6-parquet=.

# Run integration test scripts
./k6 run tests/integration/basic_test.js
```

## Test Data

Place test Parquet files in `testdata/` directory:

```bash
# Generate test data
cd examples/data
go run generate_sample.go
cp *.parquet ../../tests/testdata/
```

## Writing Tests

### Unit Tests

Follow Go testing conventions:

```go
func TestFeatureName(t *testing.T) {
    // Setup
    cache := NewReaderCache()

    // Execute
    result := cache.Get("key")

    // Assert
    if result != expected {
        t.Errorf("expected %v, got %v", expected, result)
    }
}
```

### Integration Tests

Use k6 test syntax:

```javascript
import { check } from 'k6';
import parquet from 'k6/x/parquet';

export default function() {
    const data = parquet.read('./testdata/sample.parquet');

    check(data, {
        'data loaded': (d) => d.length > 0,
    });
}
```

## Current Test Coverage

- ✅ Cache operations (cache_test.go)
- ✅ Type conversion (converter_test.go)
- 🔲 File reading
- 🔲 Schema inspection
- 🔲 Chunked reading
- 🔲 Error handling
- 🔲 Integration tests

## TODO

- [ ] Add reader_test.go for file operations
- [ ] Add schema_test.go for schema/metadata
- [ ] Add integration tests with k6
- [ ] Add benchmark tests
- [ ] Add test Parquet files to testdata/
