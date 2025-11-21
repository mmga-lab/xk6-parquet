package parquet

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/parquet-go/parquet-go"
)

// createTestParquetFile creates a temporary Parquet file for testing
func createTestParquetFile(t *testing.T) string {
	t.Helper()

	// Create a temporary file
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.parquet")

	// Create test data
	type TestRow struct {
		ID       int64   `parquet:"id"`
		Name     string  `parquet:"name"`
		Active   bool    `parquet:"active"`
		Score    float64 `parquet:"score"`
		Age      int32   `parquet:"age,optional"`
		Metadata []byte  `parquet:"metadata,optional"`
	}

	rows := []TestRow{
		{ID: 1, Name: "Alice", Active: true, Score: 95.5, Age: 30, Metadata: []byte("test1")},
		{ID: 2, Name: "Bob", Active: false, Score: 87.3, Age: 25, Metadata: []byte("test2")},
		{ID: 3, Name: "Charlie", Active: true, Score: 92.1, Age: 35, Metadata: []byte("test3")},
		{ID: 4, Name: "David", Active: false, Score: 78.9, Age: 28, Metadata: []byte("test4")},
		{ID: 5, Name: "Eve", Active: true, Score: 88.7, Age: 32, Metadata: []byte("test5")},
	}

	// Write the Parquet file
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	writer := parquet.NewGenericWriter[TestRow](file)
	_, err = writer.Write(rows)
	if err != nil {
		file.Close()
		t.Fatalf("failed to write test data: %v", err)
	}

	if err := writer.Close(); err != nil {
		file.Close()
		t.Fatalf("failed to close writer: %v", err)
	}

	file.Close()

	return filename
}

func TestValueToInterface(t *testing.T) {
	tests := []struct {
		name     string
		value    parquet.Value
		expected interface{}
	}{
		{
			name:     "Null value",
			value:    parquet.NullValue(),
			expected: nil,
		},
		{
			name:     "Boolean value",
			value:    parquet.BooleanValue(true),
			expected: true,
		},
		{
			name:     "Int32 value",
			value:    parquet.Int32Value(42),
			expected: int32(42),
		},
		{
			name:     "Int64 value",
			value:    parquet.Int64Value(9876543210),
			expected: int64(9876543210),
		},
		{
			name:     "Float value",
			value:    parquet.FloatValue(3.14),
			expected: float32(3.14),
		},
		{
			name:     "Double value",
			value:    parquet.DoubleValue(2.71828),
			expected: float64(2.71828),
		},
		{
			name:     "ByteArray value",
			value:    parquet.ByteArrayValue([]byte("test")),
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := valueToInterface(tt.value)
			if result != tt.expected {
				t.Errorf("valueToInterface() = %v (%T), want %v (%T)", result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestRowToMap(t *testing.T) {
	// Create a simple schema
	type SimpleRow struct {
		ID   int64  `parquet:"id"`
		Name string `parquet:"name"`
	}

	schema := parquet.SchemaOf(new(SimpleRow))

	// Create a row using the parquet library
	var buf bytes.Buffer
	writer := parquet.NewGenericWriter[SimpleRow](&buf)
	_, err := writer.Write([]SimpleRow{{ID: 123, Name: "Test"}})
	if err != nil {
		t.Fatalf("failed to write test row: %v", err)
	}
	writer.Close()

	// Now read as Row to test rowToMap
	rowReader := parquet.NewReader(bytes.NewReader(buf.Bytes()))
	defer rowReader.Close()

	rowBuf := make([]parquet.Row, 1)
	n, err := rowReader.ReadRows(rowBuf)
	if n == 0 {
		t.Fatalf("failed to read row: n=0, err=%v", err)
	}

	result := rowToMap(rowBuf[0], schema)

	if result["id"] != int64(123) {
		t.Errorf("expected id=123, got %v", result["id"])
	}

	if result["name"] != "Test" {
		t.Errorf("expected name=Test, got %v", result["name"])
	}
}

func TestRead(t *testing.T) {
	filename := createTestParquetFile(t)

	p := &Parquet{
		cache: NewReaderCache(),
	}

	t.Run("Read all rows", func(t *testing.T) {
		results, err := p.Read(filename)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}

		if len(results) != 5 {
			t.Errorf("expected 5 rows, got %d", len(results))
		}

		// Check first row
		if results[0]["id"] != int64(1) {
			t.Errorf("expected id=1, got %v", results[0]["id"])
		}
		if results[0]["name"] != "Alice" {
			t.Errorf("expected name=Alice, got %v", results[0]["name"])
		}
		if results[0]["active"] != true {
			t.Errorf("expected active=true, got %v", results[0]["active"])
		}
	})

	t.Run("Read with row limit", func(t *testing.T) {
		p.cache.Clear() // Clear cache to avoid getting cached full results
		options := map[string]interface{}{
			"rowLimit": 2,
		}
		results, err := p.Read(filename, options)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}

		if len(results) != 2 {
			t.Errorf("expected 2 rows, got %d", len(results))
		}
	})

	t.Run("Read with row limit as float64", func(t *testing.T) {
		p.cache.Clear()
		options := map[string]interface{}{
			"rowLimit": float64(3),
		}
		results, err := p.Read(filename, options)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}

		if len(results) != 3 {
			t.Errorf("expected 3 rows, got %d", len(results))
		}
	})

	t.Run("Read with skip rows", func(t *testing.T) {
		p.cache.Clear()
		options := map[string]interface{}{
			"skipRows": 2,
		}
		results, err := p.Read(filename, options)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}

		if len(results) != 3 {
			t.Errorf("expected 3 rows (5-2), got %d", len(results))
		}

		// First result should be the third row (ID=3)
		if results[0]["id"] != int64(3) {
			t.Errorf("expected first row id=3, got %v", results[0]["id"])
		}
	})

	t.Run("Read with skip rows as float64", func(t *testing.T) {
		p.cache.Clear()
		options := map[string]interface{}{
			"skipRows": float64(1),
		}
		results, err := p.Read(filename, options)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}

		if len(results) != 4 {
			t.Errorf("expected 4 rows, got %d", len(results))
		}

		if results[0]["id"] != int64(2) {
			t.Errorf("expected first row id=2, got %v", results[0]["id"])
		}
	})

	t.Run("Read with column filter", func(t *testing.T) {
		p.cache.Clear()
		options := map[string]interface{}{
			"columns": []interface{}{"id", "name"},
		}
		results, err := p.Read(filename, options)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}

		if len(results) != 5 {
			t.Errorf("expected 5 rows, got %d", len(results))
		}

		// Check that only requested columns are present
		if len(results[0]) != 2 {
			t.Errorf("expected 2 columns, got %d", len(results[0]))
		}

		if _, ok := results[0]["id"]; !ok {
			t.Error("expected 'id' column to be present")
		}
		if _, ok := results[0]["name"]; !ok {
			t.Error("expected 'name' column to be present")
		}
		if _, ok := results[0]["active"]; ok {
			t.Error("expected 'active' column to be filtered out")
		}
	})

	t.Run("Read with combined options", func(t *testing.T) {
		p.cache.Clear()
		options := map[string]interface{}{
			"columns":  []interface{}{"id", "name"},
			"rowLimit": 2,
			"skipRows": 1,
		}
		results, err := p.Read(filename, options)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}

		if len(results) != 2 {
			t.Errorf("expected 2 rows, got %d", len(results))
		}

		if len(results[0]) != 2 {
			t.Errorf("expected 2 columns, got %d", len(results[0]))
		}

		// Should start from second row (ID=2)
		if results[0]["id"] != int64(2) {
			t.Errorf("expected first row id=2, got %v", results[0]["id"])
		}
	})

	t.Run("Cached read", func(t *testing.T) {
		// Clear cache first
		p.cache.Clear()

		// First read should cache
		results1, err := p.Read(filename)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}

		// Second read should return from cache
		results2, err := p.Read(filename)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}

		if len(results1) != len(results2) {
			t.Errorf("cached results length mismatch: %d vs %d", len(results1), len(results2))
		}
	})

	t.Run("Read non-existent file", func(t *testing.T) {
		_, err := p.Read("/non/existent/file.parquet")
		if err == nil {
			t.Error("expected error when reading non-existent file")
		}
	})
}

func TestReadChunked(t *testing.T) {
	filename := createTestParquetFile(t)

	p := &Parquet{
		cache: NewReaderCache(),
	}

	t.Run("ReadChunked with chunk size 2", func(t *testing.T) {
		chunks := 0
		totalRows := 0

		err := p.ReadChunked(filename, 2, func(chunk []map[string]interface{}) error {
			chunks++
			totalRows += len(chunk)

			// Each chunk should have at most 2 rows
			if len(chunk) > 2 {
				t.Errorf("chunk size exceeded: got %d rows", len(chunk))
			}

			return nil
		})

		if err != nil {
			t.Fatalf("ReadChunked() error = %v", err)
		}

		// We have 5 rows, so we should get 3 chunks (2+2+1)
		if chunks != 3 {
			t.Errorf("expected 3 chunks, got %d", chunks)
		}

		if totalRows != 5 {
			t.Errorf("expected 5 total rows, got %d", totalRows)
		}
	})

	t.Run("ReadChunked with chunk size 10", func(t *testing.T) {
		chunks := 0

		err := p.ReadChunked(filename, 10, func(chunk []map[string]interface{}) error {
			chunks++
			return nil
		})

		if err != nil {
			t.Fatalf("ReadChunked() error = %v", err)
		}

		// All 5 rows should fit in 1 chunk
		if chunks != 1 {
			t.Errorf("expected 1 chunk, got %d", chunks)
		}
	})

	t.Run("ReadChunked with callback error", func(t *testing.T) {
		expectedErr := errors.New("callback error")

		err := p.ReadChunked(filename, 2, func(chunk []map[string]interface{}) error {
			return expectedErr
		})

		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("ReadChunked non-existent file", func(t *testing.T) {
		err := p.ReadChunked("/non/existent/file.parquet", 10, func(chunk []map[string]interface{}) error {
			return nil
		})

		if err == nil {
			t.Error("expected error when reading non-existent file")
		}
	})

	t.Run("ReadChunked verify data", func(t *testing.T) {
		allData := make([]map[string]interface{}, 0)

		err := p.ReadChunked(filename, 2, func(chunk []map[string]interface{}) error {
			allData = append(allData, chunk...)
			return nil
		})

		if err != nil {
			t.Fatalf("ReadChunked() error = %v", err)
		}

		// Verify the data is correct
		if len(allData) != 5 {
			t.Errorf("expected 5 rows, got %d", len(allData))
		}

		if allData[0]["name"] != "Alice" {
			t.Errorf("expected first row name=Alice, got %v", allData[0]["name"])
		}

		if allData[4]["name"] != "Eve" {
			t.Errorf("expected last row name=Eve, got %v", allData[4]["name"])
		}
	})
}
