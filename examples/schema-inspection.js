// Example: Parquet Schema and Metadata Inspection
// This example shows how to inspect Parquet file structure before reading data

import { check } from 'k6';
import parquet from 'k6/x/parquet';

export const options = {
  vus: 1,
  iterations: 1,
};

export default function() {
  const filename = './examples/data/sample.parquet';

  console.log('=== Parquet File Inspection ===\n');

  // Get and display schema
  console.log('--- Schema ---');
  const schema = parquet.getSchema(filename);
  console.log(JSON.stringify(schema, null, 2));

  // Validate schema structure
  const schemaChecks = check(schema, {
    'schema is not empty': (s) => Object.keys(s).length > 0,
    'has expected fields': (s) => {
      // Add your expected fields here
      return s.id !== undefined;
    },
  });

  console.log('\n--- Metadata ---');
  const metadata = parquet.getMetadata(filename);
  console.log(JSON.stringify(metadata, null, 2));

  // Display summary information
  console.log('\n--- Summary ---');
  console.log(`Total Rows: ${metadata.numRows}`);
  console.log(`Total Columns: ${metadata.numColumns}`);
  console.log(`Row Groups: ${metadata.numRowGroups}`);
  console.log(`File Size: ${(metadata.size / 1024 / 1024).toFixed(2)} MB`);

  // Validate metadata
  const metadataChecks = check(metadata, {
    'has rows': (m) => m.numRows > 0,
    'has columns': (m) => m.numColumns > 0,
    'has row groups': (m) => m.numRowGroups > 0,
  });

  // Read a small sample to verify data structure
  console.log('\n--- Sample Data ---');
  const sampleData = parquet.read(filename, {
    rowLimit: 5
  });

  console.log(`Sample of ${sampleData.length} records:`);
  sampleData.forEach((record, index) => {
    console.log(`Record ${index + 1}:`, JSON.stringify(record, null, 2));
  });

  // Validate data
  const dataChecks = check(sampleData, {
    'sample data retrieved': (d) => d.length > 0,
    'records have expected structure': (d) => {
      return d.every(record => typeof record === 'object');
    },
  });

  // Read specific columns only
  console.log('\n--- Column Selection ---');
  const columnData = parquet.read(filename, {
    columns: ['id', 'name'],
    rowLimit: 3
  });

  console.log('Selected columns data:');
  columnData.forEach((record, index) => {
    console.log(`Record ${index + 1}:`, JSON.stringify(record));
  });

  console.log('\n=== Inspection Complete ===');
  console.log(`Schema validation: ${schemaChecks ? 'PASSED' : 'FAILED'}`);
  console.log(`Metadata validation: ${metadataChecks ? 'PASSED' : 'FAILED'}`);
  console.log(`Data validation: ${dataChecks ? 'PASSED' : 'FAILED'}`);

  // Clean up
  parquet.close();
}
