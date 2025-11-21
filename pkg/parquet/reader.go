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

	// Read data
	results := make([]map[string]interface{}, 0)
	rowsRead := 0

	for _, rowGroup := range pf.RowGroups() {
		if opts.RowLimit > 0 && rowsRead >= opts.RowLimit {
			break
		}

		rows := rowGroup.Rows()
		defer rows.Close()

		for {
			if opts.RowLimit > 0 && rowsRead >= opts.RowLimit {
				break
			}

			row := make([]parquet.Value, len(pf.Schema().Fields()))
			n, err := rows.ReadValues(row)
			if n == 0 || err != nil {
				break
			}

			// Skip specified rows
			if rowsRead < opts.SkipRows {
				rowsRead++
				continue
			}

			// Convert row to map
			rowMap := convertRow(pf.Schema(), row, opts.Columns)
			results = append(results, rowMap)
			rowsRead++
		}
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

	chunk := make([]map[string]interface{}, 0, chunkSize)

	for _, rowGroup := range pf.RowGroups() {
		rows := rowGroup.Rows()
		defer rows.Close()

		for {
			row := make([]parquet.Value, len(pf.Schema().Fields()))
			n, err := rows.ReadValues(row)
			if n == 0 || err != nil {
				break
			}

			rowMap := convertRow(pf.Schema(), row, nil)
			chunk = append(chunk, rowMap)

			// Call callback when chunk size is reached
			if len(chunk) >= chunkSize {
				if err := callback(chunk); err != nil {
					return err
				}
				chunk = make([]map[string]interface{}, 0, chunkSize)
			}
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
