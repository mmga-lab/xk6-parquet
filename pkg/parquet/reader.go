package parquet

import (
	"fmt"
	"os"

	"github.com/parquet-go/parquet-go"
)

// ReadOptions defines options for reading Parquet files.
type ReadOptions struct {
	Columns    []string `json:"columns"`    // Specific columns to read
	RowLimit   int      `json:"rowLimit"`   // Maximum number of rows to read (-1 for all)
	SkipRows   int      `json:"skipRows"`   // Number of rows to skip
	BufferSize int      `json:"bufferSize"` // Buffer size for reading
}

// Read reads an entire Parquet file and returns the data as a slice of maps.
// It supports optional filtering by columns, limiting rows, and skipping rows.
func (p *Parquet) Read(filename string, options ...map[string]interface{}) ([]map[string]interface{}, error) {
	// Check cache first
	if cached, ok := p.cache.Get(filename); ok {
		return cached, nil
	}

	// Parse options
	opts := ReadOptions{
		RowLimit:   -1, // Default: read all rows
		BufferSize: 1000,
	}

	if len(options) > 0 {
		if columns, ok := options[0]["columns"].([]interface{}); ok {
			opts.Columns = make([]string, len(columns))
			for i, col := range columns {
				if colStr, ok := col.(string); ok {
					opts.Columns[i] = colStr
				}
			}
		}
		if rowLimit, ok := options[0]["rowLimit"].(int); ok {
			opts.RowLimit = rowLimit
		} else if rowLimit, ok := options[0]["rowLimit"].(float64); ok {
			opts.RowLimit = int(rowLimit)
		}
		if skipRows, ok := options[0]["skipRows"].(int); ok {
			opts.SkipRows = skipRows
		} else if skipRows, ok := options[0]["skipRows"].(float64); ok {
			opts.SkipRows = int(skipRows)
		}
	}

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Create Parquet file reader
	pf, err := parquet.OpenFile(file, stat.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to open parquet file: %w", err)
	}

	// Read data using generic reader
	results := make([]map[string]interface{}, 0)
	rowsRead := 0
	rowsSkipped := 0

	// Create reader for the file
	reader := parquet.NewGenericReader[map[string]any](pf)
	defer reader.Close()

	// Read rows
	for {
		if opts.RowLimit > 0 && rowsRead >= opts.RowLimit {
			break
		}

		rows := make([]map[string]any, 1)
		n, err := reader.Read(rows)
		if n == 0 || err != nil {
			break
		}

		// Skip specified rows
		if rowsSkipped < opts.SkipRows {
			rowsSkipped++
			continue
		}

		row := rows[0]

		// Filter columns if specified
		if len(opts.Columns) > 0 {
			filtered := make(map[string]interface{})
			for _, col := range opts.Columns {
				if val, ok := row[col]; ok {
					filtered[col] = val
				}
			}
			results = append(results, filtered)
		} else {
			// Convert map[string]any to map[string]interface{}
			converted := make(map[string]interface{}, len(row))
			for k, v := range row {
				converted[k] = v
			}
			results = append(results, converted)
		}

		rowsRead++
	}

	// Cache results
	p.cache.Set(filename, results)

	return results, nil
}

// ReadChunked reads a Parquet file in chunks, calling the provided callback for each chunk.
// This is useful for processing large files without loading them entirely into memory.
func (p *Parquet) ReadChunked(filename string, chunkSize int, callback func([]map[string]interface{}) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	pf, err := parquet.OpenFile(file, stat.Size())
	if err != nil {
		return fmt.Errorf("failed to open parquet file: %w", err)
	}

	reader := parquet.NewGenericReader[map[string]any](pf)
	defer reader.Close()

	chunk := make([]map[string]interface{}, 0, chunkSize)

	for {
		rows := make([]map[string]any, chunkSize)
		n, err := reader.Read(rows)
		if n == 0 || err != nil {
			break
		}

		// Convert and add to chunk
		for i := 0; i < n; i++ {
			row := rows[i]
			converted := make(map[string]interface{}, len(row))
			for k, v := range row {
				converted[k] = v
			}
			chunk = append(chunk, converted)
		}

		// Call callback when chunk size is reached
		if len(chunk) >= chunkSize {
			if err := callback(chunk); err != nil {
				return err
			}
			chunk = make([]map[string]interface{}, 0, chunkSize)
		}
	}

	// Process remaining data
	if len(chunk) > 0 {
		if err := callback(chunk); err != nil {
			return err
		}
	}

	return nil
}
