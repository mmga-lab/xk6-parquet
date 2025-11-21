package parquet

import (
	"testing"

	"github.com/parquet-go/parquet-go"
)

func TestConvertValue(t *testing.T) {
	tests := []struct {
		name     string
		input    parquet.Value
		expected interface{}
	}{
		{
			name:     "Boolean true",
			input:    parquet.BooleanValue(true),
			expected: true,
		},
		{
			name:     "Boolean false",
			input:    parquet.BooleanValue(false),
			expected: false,
		},
		{
			name:     "Int32",
			input:    parquet.Int32Value(42),
			expected: int32(42),
		},
		{
			name:     "Int64",
			input:    parquet.Int64Value(12345678901),
			expected: int64(12345678901),
		},
		{
			name:     "Float",
			input:    parquet.FloatValue(3.14),
			expected: float32(3.14),
		},
		{
			name:     "Double",
			input:    parquet.DoubleValue(3.14159265359),
			expected: float64(3.14159265359),
		},
		{
			name:     "ByteArray string",
			input:    parquet.ByteArrayValue([]byte("hello world")),
			expected: "hello world",
		},
		{
			name:     "Null value",
			input:    parquet.NullValue(),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertValue(tt.input)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
				return
			}

			// Type-specific comparisons
			switch expected := tt.expected.(type) {
			case bool:
				if result != expected {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case int32:
				if result != expected {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case int64:
				if result != expected {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case float32:
				// Use approximate comparison for floats
				if r, ok := result.(float32); !ok || !approximatelyEqual(float64(r), float64(expected), 0.0001) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case float64:
				// Use approximate comparison for doubles
				if r, ok := result.(float64); !ok || !approximatelyEqual(r, expected, 0.0000001) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case string:
				if result != expected {
					t.Errorf("expected %v, got %v", expected, result)
				}
			default:
				t.Errorf("unhandled type in test: %T", expected)
			}
		})
	}
}

func TestConvertRow(t *testing.T) {
	// Create a simple schema
	schema := parquet.NewSchema("test", parquet.Group{
		"id":   parquet.Int64(),
		"name": parquet.String(),
		"age":  parquet.Int32(),
	})

	// Create test row values
	row := []parquet.Value{
		parquet.Int64Value(1),
		parquet.ByteArrayValue([]byte("Alice")),
		parquet.Int32Value(30),
	}

	t.Run("convert all columns", func(t *testing.T) {
		result := convertRow(schema, row, nil)

		if len(result) != 3 {
			t.Errorf("expected 3 fields, got %d", len(result))
		}

		if result["id"] != int64(1) {
			t.Errorf("expected id=1, got %v", result["id"])
		}

		if result["name"] != "Alice" {
			t.Errorf("expected name=Alice, got %v", result["name"])
		}

		if result["age"] != int32(30) {
			t.Errorf("expected age=30, got %v", result["age"])
		}
	})

	t.Run("convert selected columns", func(t *testing.T) {
		columns := []string{"id", "name"}
		result := convertRow(schema, row, columns)

		if len(result) != 2 {
			t.Errorf("expected 2 fields, got %d", len(result))
		}

		if _, exists := result["id"]; !exists {
			t.Error("expected id field to exist")
		}

		if _, exists := result["name"]; !exists {
			t.Error("expected name field to exist")
		}

		if _, exists := result["age"]; exists {
			t.Error("expected age field to not exist")
		}
	})

	t.Run("empty columns filter", func(t *testing.T) {
		columns := []string{}
		result := convertRow(schema, row, columns)

		// Empty filter should return all columns
		if len(result) != 3 {
			t.Errorf("expected 3 fields, got %d", len(result))
		}
	})
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLogicalType(tt.field)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Helper function for approximate float comparison
func approximatelyEqual(a, b, epsilon float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}
