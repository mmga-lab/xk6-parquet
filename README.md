# xk6-parquet

[![Go Report Card](https://goreportcard.com/badge/github.com/mmga-lab/xk6-parquet)](https://goreportcard.com/report/github.com/mmga-lab/xk6-parquet)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A [k6 extension](https://k6.io/docs/extensions/) that enables reading Apache Parquet files directly in k6 load testing scripts.

## Features

- ✨ **Native Parquet Support** - Read Parquet files directly without conversion
- 🚀 **High Performance** - Efficient data reading with minimal memory overhead
- 📊 **Schema Inspection** - Query file schema and metadata before reading
- 🔄 **Chunked Reading** - Process large files in manageable chunks
- 🎯 **Column Selection** - Read only the columns you need
- 💾 **Built-in Caching** - Automatic caching for improved performance
- 🔧 **Type Conversion** - Automatic conversion to JavaScript-friendly types

## Use Cases

- **Large-Scale Data Testing** - Read test data from Parquet files without loading everything into memory
- **Data-Driven Testing** - Use real production data for load testing
- **Performance Benchmarking** - Compare different data formats and access patterns
- **API Testing** - Use structured data for comprehensive API endpoint testing

## Installation

### Prerequisites

- Go 1.21 or later
- [xk6](https://github.com/grafana/xk6) - k6 extension builder

### Build from Source

```bash
# Install xk6
go install go.k6.io/xk6/cmd/xk6@latest

# Build k6 with the parquet extension
xk6 build --with github.com/mmga-lab/xk6-parquet@latest

# Verify the build
./k6 version
```

### Build with Local Development Version

```bash
# Clone the repository
git clone https://github.com/mmga-lab/xk6-parquet.git
cd xk6-parquet

# Build with local version
xk6 build --with github.com/mmga-lab/xk6-parquet=.

# Run examples
./k6 run examples/basic-read.js
```

## Quick Start

### Basic Usage

```javascript
import parquet from 'k6/x/parquet';
import { check } from 'k6';

export default function() {
  // Read entire file
  const data = parquet.read('./data/users.parquet');

  console.log(`Loaded ${data.length} records`);

  // Use the data in your tests
  const user = data[0];
  console.log('First user:', user);
}
```

### Reading with Options

```javascript
import parquet from 'k6/x/parquet';

export function setup() {
  // Read only specific columns
  const data = parquet.read('./data/users.parquet', {
    columns: ['id', 'username', 'email'],
    rowLimit: 1000,
    skipRows: 100
  });

  return { users: data };
}

export default function(data) {
  const user = data.users[Math.floor(Math.random() * data.users.length)];
  // Use user in your test...
}
```

### Chunked Reading for Large Files

```javascript
import parquet from 'k6/x/parquet';

export function setup() {
  const processedData = [];

  // Read file in chunks to avoid memory issues
  parquet.readChunked('./data/large-file.parquet', 5000, (chunk) => {
    // Process each chunk
    const filtered = chunk.filter(row => row.active === true);
    processedData.push(...filtered);

    // Return null to continue, or an Error to stop
    return null;
  });

  return { data: processedData };
}
```

### Schema Inspection

```javascript
import parquet from 'k6/x/parquet';

export default function() {
  // Get schema information
  const schema = parquet.getSchema('./data/users.parquet');
  console.log('Schema:', JSON.stringify(schema, null, 2));

  // Get file metadata
  const metadata = parquet.getMetadata('./data/users.parquet');
  console.log('Total rows:', metadata.numRows);
  console.log('Total columns:', metadata.numColumns);
  console.log('Row groups:', metadata.numRowGroups);
}
```

## API Reference

### `read(filename, options?)`

Reads a Parquet file and returns all data as an array of objects.

**Parameters:**
- `filename` (string): Path to the Parquet file
- `options` (object, optional):
  - `columns` (string[]): Specific columns to read
  - `rowLimit` (number): Maximum rows to read (-1 for all)
  - `skipRows` (number): Number of rows to skip

**Returns:** Array of objects representing the data rows

**Example:**
```javascript
const data = parquet.read('./data.parquet', {
  columns: ['id', 'name'],
  rowLimit: 100,
  skipRows: 50
});
```

### `readChunked(filename, chunkSize, callback)`

Reads a Parquet file in chunks, calling the callback for each chunk.

**Parameters:**
- `filename` (string): Path to the Parquet file
- `chunkSize` (number): Number of rows per chunk
- `callback` (function): Function to process each chunk
  - Receives: array of objects (the chunk)
  - Returns: null to continue, Error to stop

**Example:**
```javascript
parquet.readChunked('./large.parquet', 1000, (chunk) => {
  console.log(`Processing ${chunk.length} rows`);
  // Process chunk...
  return null; // Continue reading
});
```

### `getSchema(filename)`

Retrieves the schema of a Parquet file.

**Parameters:**
- `filename` (string): Path to the Parquet file

**Returns:** Object with field definitions

**Example:**
```javascript
const schema = parquet.getSchema('./data.parquet');
// {
//   "id": { "type": "INT64", "optional": false, ... },
//   "name": { "type": "BYTE_ARRAY", "optional": true, ... }
// }
```

### `getMetadata(filename)`

Retrieves metadata about a Parquet file.

**Parameters:**
- `filename` (string): Path to the Parquet file

**Returns:** Object with file metadata

**Example:**
```javascript
const metadata = parquet.getMetadata('./data.parquet');
// {
//   "numRows": 10000,
//   "numColumns": 15,
//   "numRowGroups": 10,
//   "size": 52428800,
//   ...
// }
```

### `close()`

Cleans up resources and clears the cache.

**Example:**
```javascript
export function teardown() {
  parquet.close();
}
```

## Examples

Check out the [examples](examples/) directory for complete working examples:

- [basic-read.js](examples/basic-read.js) - Basic file reading
- [chunked-read.js](examples/chunked-read.js) - Chunked reading for large files
- [schema-inspection.js](examples/schema-inspection.js) - Schema and metadata inspection
- [advanced-usage.js](examples/advanced-usage.js) - Advanced features and patterns

## Generating Test Data

The repository includes utilities to generate sample Parquet files for testing:

### Using Go

```bash
cd examples/data
go run generate_sample.go
```

### Using Python

```bash
cd examples/data
pip install pyarrow pandas faker
python3 generate_sample.py
```

This will create three sample files:
- `sample.parquet` - 1,000 records (~100 KB)
- `medium.parquet` - 10,000 records (~1 MB)
- `large.parquet` - 100,000 records (~10 MB)

## Performance Tips

1. **Use Column Selection** - Only read the columns you need:
   ```javascript
   const data = parquet.read(file, { columns: ['id', 'name'] });
   ```

2. **Limit Row Reads** - Don't read more data than necessary:
   ```javascript
   const data = parquet.read(file, { rowLimit: 10000 });
   ```

3. **Use Chunked Reading** - For large files, use chunked reading:
   ```javascript
   parquet.readChunked(file, 5000, processChunk);
   ```

4. **Use SharedArray** - Share data across VUs to save memory:
   ```javascript
   import { SharedArray } from 'k6/data';

   const data = new SharedArray('parquet-data', function() {
     return parquet.read('./data.parquet');
   });
   ```

5. **Cache Benefits** - The extension automatically caches reads. Same file reads are served from cache.

## Development

### Building

```bash
# Clone repository
git clone https://github.com/mmga-lab/xk6-parquet.git
cd xk6-parquet

# Download dependencies
go mod download

# Run tests
go test ./...

# Build with xk6
xk6 build --with github.com/mmga-lab/xk6-parquet=.
```

### Project Structure

```
xk6-parquet/
├── pkg/parquet/          # Core implementation
│   ├── module.go         # k6 module registration
│   ├── reader.go         # File reading logic
│   ├── converter.go      # Type conversion
│   ├── cache.go          # Caching layer
│   └── schema.go         # Schema & metadata
├── examples/             # Usage examples
├── tests/                # Test files
└── docs/                 # Documentation
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Troubleshooting

### Out of Memory Errors

If you encounter memory issues:
- Use `readChunked()` instead of `read()`
- Limit the number of rows with `rowLimit`
- Use `SharedArray` to share data across VUs
- Select only needed columns with `columns` option

### File Not Found

Ensure file paths are correct:
- Use relative paths from your k6 script location
- Or use absolute paths
- Check file permissions

### Type Conversion Issues

Check the schema first:
```javascript
const schema = parquet.getSchema('./file.parquet');
console.log(schema);
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Credits

- [k6](https://k6.io) - Modern load testing tool
- [parquet-go](https://github.com/parquet-go/parquet-go) - Go implementation of Apache Parquet
- [xk6](https://github.com/grafana/xk6) - k6 extension builder

## Support

- **Issues**: [GitHub Issues](https://github.com/mmga-lab/xk6-parquet/issues)
- **Discussions**: [GitHub Discussions](https://github.com/mmga-lab/xk6-parquet/discussions)
- **k6 Community**: [k6 Community Forum](https://community.grafana.com/c/grafana-k6/)

---

Made with ❤️ for the k6 community
