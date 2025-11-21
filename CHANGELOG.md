# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial implementation of xk6-parquet extension
- `read()` function for reading entire Parquet files
- `readChunked()` function for memory-efficient chunked reading
- `getSchema()` function for schema inspection
- `getMetadata()` function for file metadata retrieval
- `close()` function for resource cleanup
- Automatic type conversion from Parquet to JavaScript types
- Built-in caching for improved performance
- Support for column selection and row limiting
- Comprehensive examples and documentation
- GitHub Actions workflows for CI/CD
- Sample data generation utilities

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

## [0.1.0] - TBD

### Added
- Initial release
- Basic Parquet file reading capabilities
- Schema and metadata inspection
- Chunked reading for large files
- Type conversion support
- Caching mechanism
- Documentation and examples

[Unreleased]: https://github.com/mmga-lab/xk6-parquet/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/mmga-lab/xk6-parquet/releases/tag/v0.1.0
