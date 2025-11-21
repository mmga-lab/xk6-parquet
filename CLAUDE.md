# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

xk6-parquet is a k6 extension that enables reading Apache Parquet files directly in k6 load testing scripts. It's built using the xk6 extension framework and integrates with the parquet-go library.

## Architecture

### Module System

This is a k6 extension using the k6 modules system. The architecture follows k6's VU (Virtual User) model:

- **RootModule** (`pkg/parquet/module.go:13`): Global module instance registered as `k6/x/parquet`
- **Parquet** (`pkg/parquet/module.go:16`): Per-VU instance created for each virtual user, containing:
  - VU reference for accessing k6 runtime
  - ReaderCache for caching Parquet file reads per VU

Each VU gets its own module instance with an isolated cache to prevent data races in concurrent k6 scenarios.

### Core Components

1. **Reader** (`pkg/parquet/reader.go`): Main file reading logic
   - Uses `parquet.NewGenericReader[map[string]any]` for reading Parquet files
   - Supports column filtering, row limiting, and row skipping
   - Implements both full reads (`Read`) and chunked reads (`ReadChunked`)
   - Chunked reading processes data in batches to avoid memory issues with large files

2. **Cache** (`pkg/parquet/cache.go`): Thread-safe cache with TTL
   - Default TTL: 5 minutes
   - Uses RWMutex for concurrent access safety
   - Caches complete file reads by filename
   - Note: Cache is per-VU, not global

3. **Schema** (`pkg/parquet/schema.go`): Schema inspection and metadata
   - `GetSchema`: Returns field types, optionality, logical types
   - `GetMetadata`: Returns file-level stats (rows, row groups, size)
   - `Close`: Cleanup method to clear cache

4. **Converter** (`pkg/parquet/converter.go`): Schema type conversion
   - Converts parquet-go schema types to JavaScript-friendly representations

## Development Commands

### Building

```bash
# Install xk6 (required first time)
make install-xk6

# Build k6 with the extension (creates ./k6 binary)
make build
# OR using xk6 directly:
xk6 build --with github.com/mmga-lab/xk6-parquet=.
```

### Testing

```bash
# Run all tests with race detector
make test

# Run tests with full coverage details
make test-verbose

# View coverage report in browser
make coverage

# Run specific test
go test -v ./pkg/parquet -run TestRead
```

### Linting & Formatting

```bash
# Run linters (go vet + golangci-lint)
make lint

# Format all Go code
make fmt

# Tidy dependencies
make mod-tidy
```

### Running Examples

```bash
# Build and verify
make examples

# Run specific example
./k6 run examples/basic-read.js
./k6 run examples/chunked-read.js
./k6 run examples/schema-inspection.js
./k6 run examples/advanced-usage.js
```

### Generating Test Data

```bash
# Generate sample Parquet files in examples/data/
make generate-data

# Or manually:
cd examples/data && go run generate_sample.go
```

## Key Technical Details

### parquet-go API Usage

The extension uses the Row-based Reader API (`parquet.NewReader` + `ReadRows`):
- `parquet.NewGenericReader[T]` requires a concrete struct type, not `map[string]any`
- The Reader API works with `parquet.Row` (which is `[]parquet.Value`)
- Must convert `parquet.Value` to Go interface{} based on `Kind()` (see `valueToInterface` in reader.go:10)
- Type conversion handles: Boolean, Int32, Int64, Int96, Float, Double, ByteArray, FixedLenByteArray
- ByteArray values are converted to strings for better JavaScript compatibility

### JavaScript Interop

When adding new methods to the module:
1. Add method to `Parquet` struct in `pkg/parquet/`
2. Export it in `Exports()` method (`pkg/parquet/module.go:38`)
3. Handle JavaScript type conversions (arrays come as `[]interface{}`, numbers as `float64`)
4. Return Go error as second return value - k6 automatically converts to JavaScript exception

### Testing Strategy

- Unit tests use in-memory Parquet data
- Tests must handle race conditions (use `go test -race`)
- Cache tests verify TTL expiration and concurrent access
- Converter tests verify type mappings between parquet-go and JavaScript

## CI/CD

GitHub Actions runs on pushes to `main`, `develop`, and `claude/**` branches:
- **Test job**: Go 1.24 and 1.25 with race detector and coverage
- **Lint job**: golangci-lint with 5m timeout
- **Build job**: Full xk6 build to verify extension compiles

## Project Structure Requirements

### register.go

The root directory MUST contain a `register.go` file that imports the actual module implementation:
```go
package parquet

import (
	_ "github.com/mmga-lab/xk6-parquet/pkg/parquet"
)
```
This file is required for xk6 to find and register the extension. Without it, `xk6 build` will fail with "does not contain package" error.

## Project Constraints

- Go 1.24+ required (set in go.mod)
- Depends on k6 v0.49.0 and parquet-go v0.20.1
- Must maintain compatibility with k6's module system and VU isolation model
- The register.go file in the root is mandatory for xk6 builds
