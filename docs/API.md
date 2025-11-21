# API Documentation

Complete API reference for xk6-parquet extension.

## Module Import

```javascript
import parquet from 'k6/x/parquet';
```

## Functions

### read()

Reads an entire Parquet file into memory.

#### Signature

```javascript
read(filename: string, options?: ReadOptions): Array<Object>
```

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `filename` | string | Yes | Path to the Parquet file |
| `options` | ReadOptions | No | Reading options |

#### ReadOptions

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `columns` | string[] | undefined | Array of column names to read. If not specified, all columns are read. |
| `rowLimit` | number | -1 | Maximum number of rows to read. -1 means read all rows. |
| `skipRows` | number | 0 | Number of rows to skip from the beginning. |

#### Returns

Array of objects where each object represents a row with column names as keys.

#### Examples

```javascript
// Read entire file
const data = parquet.read('./data.parquet');

// Read specific columns
const data = parquet.read('./data.parquet', {
  columns: ['id', 'name', 'email']
});

// Read limited rows
const data = parquet.read('./data.parquet', {
  rowLimit: 1000
});

// Pagination
const page1 = parquet.read('./data.parquet', {
  skipRows: 0,
  rowLimit: 100
});
const page2 = parquet.read('./data.parquet', {
  skipRows: 100,
  rowLimit: 100
});

// Combine options
const data = parquet.read('./data.parquet', {
  columns: ['id', 'username'],
  skipRows: 500,
  rowLimit: 100
});
```

#### Errors

Throws an error if:
- File doesn't exist or cannot be opened
- File is not a valid Parquet file
- Read operation fails

---

### readChunked()

Reads a Parquet file in chunks, calling a callback for each chunk. Useful for processing large files without loading everything into memory.

#### Signature

```javascript
readChunked(
  filename: string,
  chunkSize: number,
  callback: (chunk: Array<Object>) => Error | null
): Error | null
```

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `filename` | string | Yes | Path to the Parquet file |
| `chunkSize` | number | Yes | Number of rows per chunk |
| `callback` | function | Yes | Function to process each chunk |

#### Callback Function

The callback receives an array of objects (the chunk) and should return:
- `null` to continue reading
- An `Error` object to stop reading

#### Returns

- `null` if reading completed successfully
- `Error` if an error occurred

#### Examples

```javascript
// Basic chunked reading
parquet.readChunked('./large.parquet', 5000, (chunk) => {
  console.log(`Processing ${chunk.length} rows`);
  // Process chunk...
  return null; // Continue
});

// Filter data while reading
const filteredData = [];
parquet.readChunked('./data.parquet', 1000, (chunk) => {
  const filtered = chunk.filter(row => row.active === true);
  filteredData.push(...filtered);
  return null;
});

// Stop reading early
let collected = 0;
const maxRows = 10000;

parquet.readChunked('./data.parquet', 1000, (chunk) => {
  collected += chunk.length;

  if (collected >= maxRows) {
    return new Error('Collected enough data');
  }

  return null;
});

// Error handling
const err = parquet.readChunked('./data.parquet', 1000, (chunk) => {
  try {
    // Process chunk
    processData(chunk);
    return null;
  } catch (e) {
    return new Error(`Processing failed: ${e.message}`);
  }
});

if (err) {
  console.error('Reading failed:', err);
}
```

---

### getSchema()

Retrieves the schema definition of a Parquet file.

#### Signature

```javascript
getSchema(filename: string): Object
```

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `filename` | string | Yes | Path to the Parquet file |

#### Returns

Object with field definitions. Each field contains:
- `type`: Parquet physical type (e.g., "INT64", "BYTE_ARRAY")
- `optional`: Boolean indicating if field can be null
- `repeated`: Boolean indicating if field is repeated (array)
- `logical`: Logical type if defined (e.g., "UTF8", "TIMESTAMP_MILLIS")

#### Example

```javascript
const schema = parquet.getSchema('./data.parquet');

// Output example:
// {
//   "id": {
//     "type": "INT64",
//     "optional": false,
//     "repeated": false,
//     "logical": "none"
//   },
//   "username": {
//     "type": "BYTE_ARRAY",
//     "optional": true,
//     "repeated": false,
//     "logical": "UTF8"
//   },
//   "created_at": {
//     "type": "INT64",
//     "optional": false,
//     "repeated": false,
//     "logical": "TIMESTAMP_MILLIS"
//   },
//   "tags": {
//     "type": "BYTE_ARRAY",
//     "optional": true,
//     "repeated": true,
//     "logical": "UTF8"
//   }
// }

// Validate schema
const hasIdField = schema.id !== undefined;
const idIsInt = schema.id.type === 'INT64';

// Check field properties
if (schema.email && schema.email.optional) {
  console.log('Email field can be null');
}
```

---

### getMetadata()

Retrieves metadata and statistics about a Parquet file.

#### Signature

```javascript
getMetadata(filename: string): Object
```

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `filename` | string | Yes | Path to the Parquet file |

#### Returns

Object containing:
- `numRows`: Total number of rows
- `numColumns`: Total number of columns
- `numRowGroups`: Number of row groups
- `size`: File size in bytes
- `rowGroups`: Array of row group information
- `schema`: Schema definition (same as getSchema())

#### Example

```javascript
const metadata = parquet.getMetadata('./data.parquet');

console.log(`Total rows: ${metadata.numRows}`);
console.log(`Total columns: ${metadata.numColumns}`);
console.log(`Row groups: ${metadata.numRowGroups}`);
console.log(`File size: ${(metadata.size / 1024 / 1024).toFixed(2)} MB`);

// Row group details
metadata.rowGroups.forEach((rg, index) => {
  console.log(`Row group ${rg.index}:`);
  console.log(`  Rows: ${rg.numRows}`);
  console.log(`  Columns: ${rg.numColumns}`);
});

// Use metadata to make decisions
if (metadata.numRows > 100000) {
  console.log('Large file detected, using chunked reading');
  parquet.readChunked('./data.parquet', 10000, processChunk);
} else {
  console.log('Small file, reading all at once');
  const data = parquet.read('./data.parquet');
}
```

---

### close()

Cleans up resources and clears the internal cache.

#### Signature

```javascript
close(): Error | null
```

#### Returns

- `null` if successful
- `Error` if cleanup failed

#### Example

```javascript
export function setup() {
  const data = parquet.read('./data.parquet');
  return { data };
}

export default function(data) {
  // Use data...
}

export function teardown(data) {
  // Clean up
  parquet.close();
  console.log('Resources cleaned up');
}
```

---

## Type Conversions

The extension automatically converts Parquet types to JavaScript types:

| Parquet Type | JavaScript Type | Notes |
|--------------|-----------------|-------|
| BOOLEAN | boolean | |
| INT32 | number | |
| INT64 | number | May lose precision for very large values |
| INT96 | number | Converted to INT64 |
| FLOAT | number | |
| DOUBLE | number | |
| BYTE_ARRAY | string | Assumes UTF-8 encoding |
| FIXED_LEN_BYTE_ARRAY | Uint8Array | Raw bytes |

### Logical Types

| Logical Type | JavaScript Type | Notes |
|--------------|-----------------|-------|
| UTF8 | string | |
| STRING | string | |
| TIMESTAMP_MILLIS | number | Unix timestamp in milliseconds |
| TIMESTAMP_MICROS | number | Unix timestamp in microseconds |
| DATE | number | Days since epoch |
| TIME_MILLIS | number | Milliseconds since midnight |
| DECIMAL | number | May lose precision |

---

## Error Handling

All functions may throw errors. Use try-catch or check return values:

```javascript
// Method 1: Try-catch
try {
  const data = parquet.read('./data.parquet');
} catch (error) {
  console.error('Failed to read file:', error);
}

// Method 2: Check return value (for functions that return errors)
const err = parquet.close();
if (err) {
  console.error('Failed to close:', err);
}

// Method 3: Callback error handling
parquet.readChunked('./data.parquet', 1000, (chunk) => {
  try {
    processChunk(chunk);
    return null;
  } catch (e) {
    return new Error(`Processing failed: ${e.message}`);
  }
});
```

---

## Best Practices

### 1. Use Column Selection

Only read columns you need:

```javascript
const data = parquet.read('./data.parquet', {
  columns: ['id', 'name'] // Only these columns
});
```

### 2. Limit Row Reads

Don't read more than necessary:

```javascript
const sample = parquet.read('./data.parquet', {
  rowLimit: 100 // Only 100 rows
});
```

### 3. Use Chunked Reading for Large Files

```javascript
parquet.readChunked('./large.parquet', 5000, processChunk);
```

### 4. Share Data with SharedArray

```javascript
import { SharedArray } from 'k6/data';

const data = new SharedArray('my-data', function() {
  return parquet.read('./data.parquet');
});
```

### 5. Clean Up Resources

```javascript
export function teardown() {
  parquet.close();
}
```

### 6. Check Metadata First

```javascript
const metadata = parquet.getMetadata('./data.parquet');
if (metadata.numRows > 1000000) {
  // Use chunked reading for large files
}
```

---

## Performance Tips

1. **Caching**: The extension caches file reads automatically
2. **Column projection**: Use `columns` option to reduce memory
3. **Row limiting**: Use `rowLimit` for sampling
4. **Chunked reading**: For files larger than available memory
5. **SharedArray**: Share data across VUs to reduce memory usage

---

## Common Patterns

### Loading Test Data

```javascript
export function setup() {
  return {
    users: parquet.read('./users.parquet'),
    products: parquet.read('./products.parquet')
  };
}
```

### Data Sampling

```javascript
const sample = parquet.read('./large.parquet', {
  rowLimit: 1000
});
```

### Filtering Data

```javascript
const activeUsers = [];
parquet.readChunked('./users.parquet', 5000, (chunk) => {
  activeUsers.push(...chunk.filter(u => u.active));
  return null;
});
```

### Pagination

```javascript
function getPage(page, pageSize) {
  return parquet.read('./data.parquet', {
    skipRows: page * pageSize,
    rowLimit: pageSize
  });
}
```

---

For more examples, see the [examples](../examples/) directory.
