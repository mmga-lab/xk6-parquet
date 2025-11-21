package parquet

import (
	"github.com/parquet-go/parquet-go"
)

// ConvertSchema converts a Parquet schema to a readable map format.
func ConvertSchema(schema *parquet.Schema) map[string]interface{} {
	result := make(map[string]interface{})

	for _, field := range schema.Fields() {
		fieldInfo := map[string]interface{}{
			"type":     field.Type().String(),
			"optional": field.Optional(),
			"repeated": field.Repeated(),
		}

		// Add logical type if available
		if logicalType := getLogicalType(field); logicalType != "" {
			fieldInfo["logical"] = logicalType
		} else {
			fieldInfo["logical"] = "none"
		}

		result[field.Name()] = fieldInfo
	}

	return result
}

// getLogicalType extracts the logical type from a Parquet field.
func getLogicalType(field parquet.Field) string {
	if field.Type().LogicalType() != nil {
		return field.Type().LogicalType().String()
	}
	return ""
}
