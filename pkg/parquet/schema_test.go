package parquet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/parquet-go/parquet-go"
)

func TestGetSchema(t *testing.T) {
	// Create a test Parquet file with known schema
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_schema.parquet")

	type TestRow struct {
		ID       int64   `parquet:"id"`
		Name     string  `parquet:"name"`
		Active   bool    `parquet:"active"`
		Score    float64 `parquet:"score"`
		Age      int32   `parquet:"age,optional"`
		Metadata []byte  `parquet:"metadata,optional"`
	}

	// Write a sample file
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	writer := parquet.NewGenericWriter[TestRow](file)
	_, err = writer.Write([]TestRow{
		{ID: 1, Name: "Test", Active: true, Score: 95.5, Age: 30, Metadata: []byte("data")},
	})
	if err != nil {
		t.Fatalf("failed to write test data: %v", err)
	}
	writer.Close()
	file.Close()

	p := &Parquet{
		cache: NewReaderCache(),
	}

	t.Run("GetSchema success", func(t *testing.T) {
		schema, err := p.GetSchema(filename)
		if err != nil {
			t.Fatalf("GetSchema() error = %v", err)
		}

		if schema == nil {
			t.Fatal("expected non-nil schema")
		}

		if len(schema) == 0 {
			t.Error("expected at least one field in schema")
		}

		// Verify some expected fields exist
		expectedFields := []string{"id", "name", "active", "score", "age", "metadata"}
		for _, expectedField := range expectedFields {
			if _, ok := schema[expectedField]; !ok {
				t.Errorf("expected field %s not found in schema", expectedField)
			}
		}

		// Verify field has expected properties
		if idField, ok := schema["id"].(map[string]interface{}); ok {
			if _, ok := idField["type"]; !ok {
				t.Error("expected id field to have 'type' property")
			}
			if _, ok := idField["optional"]; !ok {
				t.Error("expected id field to have 'optional' property")
			}
		} else {
			t.Error("expected id field to be a map")
		}
	})

	t.Run("GetSchema non-existent file", func(t *testing.T) {
		_, err := p.GetSchema("/non/existent/file.parquet")
		if err == nil {
			t.Error("expected error when getting schema of non-existent file")
		}
	})

	t.Run("GetSchema invalid file", func(t *testing.T) {
		// Create a non-Parquet file
		invalidFile := filepath.Join(tmpDir, "invalid.parquet")
		if err := os.WriteFile(invalidFile, []byte("not a parquet file"), 0644); err != nil {
			t.Fatalf("failed to create invalid file: %v", err)
		}

		_, err := p.GetSchema(invalidFile)
		if err == nil {
			t.Error("expected error when getting schema of invalid file")
		}
	})
}

func TestGetMetadata(t *testing.T) {
	// Create a test Parquet file with known metadata
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_metadata.parquet")

	type TestRow struct {
		ID    int64  `parquet:"id"`
		Name  string `parquet:"name"`
		Value int32  `parquet:"value"`
	}

	// Write multiple rows to ensure we have meaningful metadata
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	rows := make([]TestRow, 100)
	for i := 0; i < 100; i++ {
		rows[i] = TestRow{
			ID:    int64(i + 1),
			Name:  "Test",
			Value: int32(i * 10),
		}
	}

	writer := parquet.NewGenericWriter[TestRow](file)
	_, err = writer.Write(rows)
	if err != nil {
		t.Fatalf("failed to write test data: %v", err)
	}
	writer.Close()
	file.Close()

	p := &Parquet{
		cache: NewReaderCache(),
	}

	t.Run("GetMetadata success", func(t *testing.T) {
		metadata, err := p.GetMetadata(filename)
		if err != nil {
			t.Fatalf("GetMetadata() error = %v", err)
		}

		if metadata == nil {
			t.Fatal("expected non-nil metadata")
		}

		// Check numRows
		numRows, ok := metadata["numRows"].(int64)
		if !ok {
			t.Fatal("expected numRows to be int64")
		}
		if numRows != 100 {
			t.Errorf("expected 100 rows, got %d", numRows)
		}

		// Check numColumns
		numColumns, ok := metadata["numColumns"].(int)
		if !ok {
			t.Fatal("expected numColumns to be int")
		}
		if numColumns != 3 {
			t.Errorf("expected 3 columns, got %d", numColumns)
		}

		// Check numRowGroups
		numRowGroups, ok := metadata["numRowGroups"].(int)
		if !ok {
			t.Fatal("expected numRowGroups to be int")
		}
		if numRowGroups < 1 {
			t.Error("expected at least 1 row group")
		}

		// Check size
		size, ok := metadata["size"].(int64)
		if !ok {
			t.Fatal("expected size to be int64")
		}
		if size <= 0 {
			t.Error("expected positive file size")
		}

		// Check rowGroups array
		rowGroups, ok := metadata["rowGroups"].([]map[string]interface{})
		if !ok {
			t.Fatal("expected rowGroups to be []map[string]interface{}")
		}
		if len(rowGroups) != numRowGroups {
			t.Errorf("expected %d row groups, got %d", numRowGroups, len(rowGroups))
		}

		// Check first row group has expected fields
		if len(rowGroups) > 0 {
			rg := rowGroups[0]
			if _, ok := rg["index"]; !ok {
				t.Error("expected row group to have 'index' field")
			}
			if _, ok := rg["numRows"]; !ok {
				t.Error("expected row group to have 'numRows' field")
			}
			if _, ok := rg["numColumns"]; !ok {
				t.Error("expected row group to have 'numColumns' field")
			}
		}

		// Check schema
		schema, ok := metadata["schema"].(map[string]interface{})
		if !ok {
			t.Fatal("expected schema to be map[string]interface{}")
		}
		if schema == nil {
			t.Error("expected non-nil schema in metadata")
		}
	})

	t.Run("GetMetadata non-existent file", func(t *testing.T) {
		_, err := p.GetMetadata("/non/existent/file.parquet")
		if err == nil {
			t.Error("expected error when getting metadata of non-existent file")
		}
	})

	t.Run("GetMetadata invalid file", func(t *testing.T) {
		// Create a non-Parquet file
		invalidFile := filepath.Join(tmpDir, "invalid_meta.parquet")
		if err := os.WriteFile(invalidFile, []byte("not a parquet file"), 0644); err != nil {
			t.Fatalf("failed to create invalid file: %v", err)
		}

		_, err := p.GetMetadata(invalidFile)
		if err == nil {
			t.Error("expected error when getting metadata of invalid file")
		}
	})
}

func TestClose(t *testing.T) {
	p := &Parquet{
		cache: NewReaderCache(),
	}

	// Add some data to cache
	testData := []map[string]interface{}{
		{"id": 1, "name": "test"},
	}
	p.cache.Set("key1", testData)
	p.cache.Set("key2", testData)

	// Verify data is in cache
	_, found := p.cache.Get("key1")
	if !found {
		t.Error("expected data to be in cache before Close")
	}

	// Call Close
	err := p.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Verify cache is cleared
	_, found = p.cache.Get("key1")
	if found {
		t.Error("expected cache to be cleared after Close")
	}

	_, found = p.cache.Get("key2")
	if found {
		t.Error("expected cache to be cleared after Close")
	}
}
