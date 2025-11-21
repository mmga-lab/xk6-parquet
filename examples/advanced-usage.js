// Example: Advanced Parquet Usage
// Demonstrates advanced features like column selection, pagination, and data filtering

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { SharedArray } from 'k6/data';
import parquet from 'k6/x/parquet';

export const options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '1m', target: 20 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.1'],
  },
};

// Use SharedArray to share data across VUs efficiently
const users = new SharedArray('users', function() {
  console.log('Loading user data from Parquet...');

  // Read only specific columns to reduce memory usage
  const data = parquet.read('./examples/data/sample.parquet', {
    columns: ['id', 'username', 'email', 'created_at'],
    rowLimit: 10000, // Limit to first 10,000 records
  });

  console.log(`Loaded ${data.length} user records`);

  // Pre-process and validate data
  return data.filter(user => {
    return user.email && user.email.includes('@');
  });
});

const premiumUsers = new SharedArray('premium-users', function() {
  console.log('Loading premium users from Parquet...');

  const allUsers = [];

  // Use chunked reading for large files
  parquet.readChunked('./examples/data/sample.parquet', 5000, (chunk) => {
    // Filter premium users only
    const premium = chunk.filter(user => user.subscription === 'premium');
    allUsers.push(...premium);

    console.log(`Processed chunk: found ${premium.length} premium users`);
    return null;
  });

  console.log(`Total premium users: ${allUsers.length}`);
  return allUsers;
});

export function setup() {
  // Get file metadata
  const metadata = parquet.getMetadata('./examples/data/sample.parquet');

  console.log('\n=== Dataset Information ===');
  console.log(`Total records: ${metadata.numRows}`);
  console.log(`Columns: ${metadata.numColumns}`);
  console.log(`Row groups: ${metadata.numRowGroups}`);

  return {
    metadata: metadata,
    startTime: Date.now(),
  };
}

export default function(data) {
  group('User API Tests', () => {
    // Test with regular users
    group('Regular User Flow', () => {
      const user = users[Math.floor(Math.random() * users.length)];

      const loginResponse = http.post(
        'https://httpbin.org/post',
        JSON.stringify({
          username: user.username,
          email: user.email,
        }),
        {
          headers: { 'Content-Type': 'application/json' },
          tags: { type: 'login' },
        }
      );

      check(loginResponse, {
        'login successful': (r) => r.status === 200,
      });
    });

    // Test with premium users
    if (premiumUsers.length > 0) {
      group('Premium User Flow', () => {
        const premiumUser = premiumUsers[Math.floor(Math.random() * premiumUsers.length)];

        const premiumResponse = http.get(
          `https://httpbin.org/get?user_id=${premiumUser.id}&tier=premium`,
          {
            tags: { type: 'premium' },
          }
        );

        check(premiumResponse, {
          'premium access granted': (r) => r.status === 200,
        });
      });
    }
  });

  sleep(1);
}

export function teardown(data) {
  const duration = (Date.now() - data.startTime) / 1000;

  console.log('\n=== Test Summary ===');
  console.log(`Test duration: ${duration.toFixed(2)}s`);
  console.log(`Regular users loaded: ${users.length}`);
  console.log(`Premium users loaded: ${premiumUsers.length}`);
  console.log(`Total dataset size: ${data.metadata.numRows} rows`);

  // Clean up resources
  parquet.close();
  console.log('Resources cleaned up');
}
