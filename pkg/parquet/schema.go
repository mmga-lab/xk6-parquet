package parquet

import (
	"fmt"
	"os"

	"github.com/parquet-go/parquet-go"
)

// GetSchema retrieves and returns the schema of a Parquet file.
func (p *Parquet) GetSchema(filename string) (map[string]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	pf, err := parquet.OpenFile(file, stat.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to open parquet file: %w", err)
	}

	return ConvertSchema(pf.Schema()), nil
}

// GetMetadata retrieves and returns metadata about a Parquet file.
func (p *Parquet) GetMetadata(filename string) (map[string]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	pf, err := parquet.OpenFile(file, stat.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to open parquet file: %w", err)
	}

	metadata := make(map[string]interface{})

	// File-level information
	metadata["numRows"] = pf.NumRows()
	metadata["numRowGroups"] = len(pf.RowGroups())
	metadata["numColumns"] = len(pf.Schema().Fields())
	metadata["size"] = stat.Size()

	// Row group information
	rowGroups := make([]map[string]interface{}, 0)
	for i, rg := range pf.RowGroups() {
		rgInfo := map[string]interface{}{
			"index":      i,
			"numRows":    rg.NumRows(),
			"numColumns": len(rg.ColumnChunks()),
		}
		rowGroups = append(rowGroups, rgInfo)
	}
	metadata["rowGroups"] = rowGroups

	// Schema information
	metadata["schema"] = ConvertSchema(pf.Schema())

	return metadata, nil
}

// Close cleans up resources and clears the cache.
func (p *Parquet) Close() error {
	p.cache.Clear()
	return nil
}
