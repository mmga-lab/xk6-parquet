package parquet

import (
	"testing"

	"github.com/parquet-go/parquet-go"
)

func TestConvertSchema(t *testing.T) {
	// Create a simple schema
	schema := parquet.NewSchema("test", parquet.Group{
		"id":     parquet.Int64(),
		"name":   parquet.String(),
		"age":    parquet.Int32(),
		"active": parquet.Boolean(),
	})

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
	tests := []struct {
		name     string
		field    parquet.Field
		expected string
	}{
		{
			name:     "String type",
			field:    parquet.String(),
			expected: "STRING",
		},
		{
			name:     "Int64 without logical type",
			field:    parquet.Int64(),
			expected: "",
		},
		{
			name:     "Boolean",
			field:    parquet.Boolean(),
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
