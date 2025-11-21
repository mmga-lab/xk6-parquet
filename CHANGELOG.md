# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- N/A

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- N/A

## [0.1.0] - 2025-11-21

### Added
- Initial implementation of xk6-parquet extension
- `read()` function for reading entire Parquet files with column filtering and row limiting
- `readChunked()` function for memory-efficient chunked reading of large files
- `getSchema()` function for Parquet schema inspection
- `getMetadata()` function for file metadata retrieval (row count, size, etc.)
- `close()` function for cache cleanup and resource management
- Automatic type conversion from Parquet types to JavaScript types
- Built-in caching mechanism with 5-minute TTL for improved performance
- Support for column selection and row limiting/skipping
- k6 VU (Virtual User) isolation for thread-safe concurrent testing
- Comprehensive documentation (README, API docs, CONTRIBUTING guide)
- Working examples (basic read, chunked read, schema inspection, advanced usage)
- GitHub Actions CI/CD workflows (test, lint, build, release)
- Sample Parquet data generation utilities (Go and Python scripts)

### Technical Details
- Go 1.24+ required
- parquet-go v0.25.1 dependency
- k6 v0.49.0 compatibility
- Multi-platform binaries (Linux, macOS, Windows for amd64 and arm64)

[Unreleased]: https://github.com/mmga-lab/xk6-parquet/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/mmga-lab/xk6-parquet/releases/tag/v0.1.0
