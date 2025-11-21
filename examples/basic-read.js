// Example: Basic Parquet File Reading
// This example demonstrates how to read a Parquet file and use the data in k6 tests

import http from 'k6/http';
import { check, sleep } from 'k6';
import parquet from 'k6/x/parquet';

// Test configuration
export const options = {
  vus: 10,
  duration: '30s',
};

// Setup function - runs once before the test
export function setup() {
  console.log('Loading test data from Parquet file...');

  // Read the entire Parquet file
  const data = parquet.read('./examples/data/sample.parquet');

  console.log(`Loaded ${data.length} records from Parquet file`);

  // Display first record as example
  if (data.length > 0) {
    console.log('Sample record:', JSON.stringify(data[0], null, 2));
  }

  return { testData: data };
}

// Main test function - runs for each VU
export default function(data) {
  // Select a random record from the loaded data
  const record = data.testData[Math.floor(Math.random() * data.testData.length)];

  // Use the record data in your test
  // Example: Making an HTTP request with data from Parquet file
  const response = http.get(`https://httpbin.org/get?id=${record.id}`);

  // Check response
  check(response, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(1);
}

// Teardown function - runs once after the test
export function teardown(data) {
  console.log(`Test completed. Processed ${data.testData.length} records.`);
}
