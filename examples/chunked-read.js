// Example: Chunked Parquet File Reading
// This example demonstrates how to read large Parquet files in chunks
// to avoid loading the entire dataset into memory at once

import http from 'k6/http';
import { check, sleep } from 'k6';
import parquet from 'k6/x/parquet';

export const options = {
  vus: 5,
  duration: '1m',
};

export function setup() {
  console.log('Loading large Parquet file in chunks...');

  const allData = [];
  let chunkCount = 0;
  let totalRows = 0;

  // Read file in chunks of 1000 rows
  const chunkSize = 1000;

  parquet.readChunked('./examples/data/sample.parquet', chunkSize, (chunk) => {
    chunkCount++;
    totalRows += chunk.length;

    console.log(`Processing chunk ${chunkCount}: ${chunk.length} rows`);

    // Filter or process each chunk as needed
    const processedChunk = chunk.filter(row => row.active === true);

    // Only keep the data you need
    allData.push(...processedChunk);

    // Return null to continue, or an Error to stop reading
    // Example: Stop after collecting 5000 records
    if (allData.length >= 5000) {
      console.log('Collected enough data, stopping...');
      return new Error('Sufficient data collected');
    }

    return null;
  });

  console.log(`Total chunks processed: ${chunkCount}`);
  console.log(`Total rows read: ${totalRows}`);
  console.log(`Records kept after filtering: ${allData.length}`);

  return { testData: allData };
}

export default function(data) {
  if (data.testData.length === 0) {
    console.log('No test data available');
    return;
  }

  const record = data.testData[__VU % data.testData.length];

  const response = http.post(
    'https://httpbin.org/post',
    JSON.stringify(record),
    {
      headers: { 'Content-Type': 'application/json' },
    }
  );

  check(response, {
    'status is 200': (r) => r.status === 200,
    'has data': (r) => r.json('json') !== undefined,
  });

  sleep(1);
}

export function teardown(data) {
  console.log('Test completed.');
  console.log(`Final dataset size: ${data.testData.length} records`);
}
