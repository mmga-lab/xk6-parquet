# Sample Parquet Data

This directory contains utilities to generate sample Parquet files for testing the xk6-parquet extension.

## Generating Sample Data

### Using Go

```bash
go run generate_sample.go
```

### Using Python

Requirements:
```bash
pip install pyarrow pandas faker
```

Generate data:
```bash
python3 generate_sample.py
```

## Generated Files

After running the generation script, you'll have:

| File | Records | Size | Description |
|------|---------|------|-------------|
| `sample.parquet` | 1,000 | ~100 KB | Small sample file for quick tests |
| `medium.parquet` | 10,000 | ~1 MB | Medium-sized file for testing |
| `large.parquet` | 100,000 | ~10 MB | Large file for performance testing |

## Data Schema

All generated files contain user data with the following schema:

| Column | Type | Description |
|--------|------|-------------|
| `id` | INT64 | Unique user ID |
| `username` | STRING | Username |
| `email` | STRING | Email address |
| `name` | STRING | Full name |
| `age` | INT32 | User age (18-80) |
| `subscription` | STRING | Subscription tier (free/premium/enterprise) |
| `active` | BOOLEAN | Account active status |
| `created_at` | STRING | Account creation timestamp (ISO 8601) |
| `balance` | DOUBLE | Account balance |
| `country` | STRING | Country code (ISO 2-letter) |

## Using Sample Data

### In k6 Scripts

```javascript
import parquet from 'k6/x/parquet';

export function setup() {
  const data = parquet.read('./examples/data/sample.parquet');
  return { users: data };
}

export default function(data) {
  const user = data.users[__VU % data.users.length];
  // Use user data in your tests...
}
```

### Inspecting Data

```javascript
import parquet from 'k6/x/parquet';

export default function() {
  const metadata = parquet.getMetadata('./examples/data/sample.parquet');
  console.log(JSON.stringify(metadata, null, 2));

  const schema = parquet.getSchema('./examples/data/sample.parquet');
  console.log(JSON.stringify(schema, null, 2));

  const sample = parquet.read('./examples/data/sample.parquet', {
    rowLimit: 5
  });
  console.log(JSON.stringify(sample, null, 2));
}
```

## Custom Data

You can modify the generation scripts to create custom data structures that match your testing needs.
