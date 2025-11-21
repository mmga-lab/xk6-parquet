package parquet

import (
	"github.com/parquet-go/parquet-go"
)

// convertValue converts a Parquet value to a Go native type.
func convertValue(v parquet.Value) interface{} {
	if v.IsNull() {
		return nil
	}

	switch v.Kind() {
	case parquet.Boolean:
		return v.Boolean()

	case parquet.Int32:
		return v.Int32()

	case parquet.Int64:
		return v.Int64()

	case parquet.Int96:
		// Int96 is typically used for timestamps
		return v.Int64()

	case parquet.Float:
		return v.Float()

	case parquet.Double:
		return v.Double()

	case parquet.ByteArray:
		return string(v.ByteArray())

	case parquet.FixedLenByteArray:
		return v.ByteArray()

	default:
		// For unknown types, return string representation
		return v.String()
	}
}

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

// convertRow converts a Parquet row to a map with optional column filtering.
func convertRow(schema *parquet.Schema, row []parquet.Value, columns []string) map[string]interface{} {
	rowMap := make(map[string]interface{})

	// Create column filter set if columns are specified
	columnFilter := make(map[string]bool)
	if len(columns) > 0 {
		for _, col := range columns {
			columnFilter[col] = true
		}
	}

	fields := schema.Fields()
	for i, field := range fields {
		fieldName := field.Name()

		// Check column filter
		if len(columnFilter) > 0 && !columnFilter[fieldName] {
			continue
		}

		if i < len(row) {
			rowMap[fieldName] = convertValue(row[i])
		}
	}

	return rowMap
}
