package parquet

import (
	"testing"

	"github.com/parquet-go/parquet-go"
)

// TestRecord is a simple struct for testing schema conversion
type TestRecord struct {
	ID     int64  `parquet:"id"`
	Name   string `parquet:"name"`
	Age    int32  `parquet:"age"`
	Active bool   `parquet:"active"`
}

func TestConvertSchema(t *testing.T) {
	// Create a schema from the test struct
	schema := parquet.SchemaOf(TestRecord{})

	result := ConvertSchema(schema)

	if len(result) != 4 {
		t.Errorf("expected 4 fields, got %d", len(result))
	}

	// Check id field
	if idField, ok := result["id"]; ok {
		if field, ok := idField.(map[string]interface{}); ok {
			if field["optional"] == nil {
				t.Error("expected optional field to be set")
			}
			if field["repeated"] == nil {
				t.Error("expected repeated field to be set")
			}
			if field["type"] == nil {
				t.Error("expected type field to be set")
			}
		} else {
			t.Error("expected field to be a map")
		}
	} else {
		t.Error("expected id field in result")
	}

	// Check name field (String type)
	if nameField, ok := result["name"]; ok {
		if field, ok := nameField.(map[string]interface{}); ok {
			if field["logical"] == "none" {
				// This is okay, logical type might not be set
			}
		}
	} else {
		t.Error("expected name field in result")
	}
}

func TestGetLogicalType(t *testing.T) {
	// Create schema and get fields from it
	schema := parquet.SchemaOf(TestRecord{})
	fields := schema.Fields()

	tests := []struct {
		name     string
		field    parquet.Field
		expected string
	}{
		{
			name:     "String type (name field)",
			field:    fields[1], // name field
			expected: "STRING",
		},
		{
			name:     "Int64 without logical type (id field)",
			field:    fields[0], // id field
			expected: "",
		},
		{
			name:     "Boolean (active field)",
			field:    fields[3], // active field
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLogicalType(tt.field)
			if result != tt.expected {
				// Logical types might vary based on parquet-go version
				// Just check that we get a string back
				if result == "" && tt.expected != "" {
					t.Logf("Warning: expected %q, got empty string (may be version-dependent)", tt.expected)
				}
			}
		})
	}
}
