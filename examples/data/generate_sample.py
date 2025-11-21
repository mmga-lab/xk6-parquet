#!/usr/bin/env python3
"""
Generate sample Parquet files for testing xk6-parquet extension.
Requires: pip install pyarrow pandas faker
"""

import pyarrow as pa
import pyarrow.parquet as pq
import pandas as pd
from faker import Faker
from datetime import datetime, timedelta
import random

fake = Faker()
Faker.seed(42)
random.seed(42)


def generate_sample_data(num_records=1000):
    """Generate sample user data."""
    data = {
        'id': range(1, num_records + 1),
        'username': [fake.user_name() for _ in range(num_records)],
        'email': [fake.email() for _ in range(num_records)],
        'name': [fake.name() for _ in range(num_records)],
        'age': [random.randint(18, 80) for _ in range(num_records)],
        'subscription': [random.choice(['free', 'premium', 'enterprise']) for _ in range(num_records)],
        'active': [random.choice([True, False]) for _ in range(num_records)],
        'created_at': [(datetime.now() - timedelta(days=random.randint(0, 365))).isoformat() for _ in range(num_records)],
        'balance': [round(random.uniform(0, 10000), 2) for _ in range(num_records)],
        'country': [fake.country_code() for _ in range(num_records)],
    }
    return pd.DataFrame(data)


def main():
    print("Generating sample Parquet files...")

    # Generate small sample file (1,000 records)
    print("Creating sample.parquet (1,000 records)...")
    df_small = generate_sample_data(1000)
    df_small.to_parquet('sample.parquet', compression='snappy', index=False)
    print(f"✓ Created sample.parquet: {len(df_small)} records")

    # Generate medium file (10,000 records)
    print("Creating medium.parquet (10,000 records)...")
    df_medium = generate_sample_data(10000)
    df_medium.to_parquet('medium.parquet', compression='snappy', index=False)
    print(f"✓ Created medium.parquet: {len(df_medium)} records")

    # Generate large file (100,000 records)
    print("Creating large.parquet (100,000 records)...")
    df_large = generate_sample_data(100000)
    df_large.to_parquet('large.parquet', compression='snappy', index=False)
    print(f"✓ Created large.parquet: {len(df_large)} records")

    # Display schema and sample data
    print("\n=== Sample Schema ===")
    parquet_file = pq.read_table('sample.parquet')
    print(parquet_file.schema)

    print("\n=== Sample Data (first 5 records) ===")
    print(df_small.head())

    print("\n=== File Statistics ===")
    import os
    for filename in ['sample.parquet', 'medium.parquet', 'large.parquet']:
        if os.path.exists(filename):
            size_mb = os.path.getsize(filename) / (1024 * 1024)
            print(f"{filename}: {size_mb:.2f} MB")

    print("\n✓ All sample files generated successfully!")


if __name__ == '__main__':
    main()
