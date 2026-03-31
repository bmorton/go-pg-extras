package pgextras

import "context"

// IndexInfoResult holds aggregated index information.
type IndexInfoResult struct {
	IndexName  string `json:"index_name"`
	TableName  string `json:"table_name"`
	Columns    string `json:"columns"`
	IndexSize  string `json:"index_size"`
	IndexScans int64  `json:"index_scans"`
	NullFrac   string `json:"null_frac"`
}

// IndexInfo returns aggregated information about indexes.
func (c *Client) IndexInfo(ctx context.Context, tableName string) ([]IndexInfoResult, error) {
	indexes, err := c.Indexes(ctx)
	if err != nil {
		return nil, err
	}
	indexSizes, err := c.IndexSize(ctx)
	if err != nil {
		return nil, err
	}
	scans, err := c.IndexScans(ctx)
	if err != nil {
		return nil, err
	}
	nullIdxs, err := c.NullIndexes(ctx, NullIndexesParams{MinRelationSizeMB: 0})
	if err != nil {
		// NullIndexes may fail if no indexes match; treat as empty.
		nullIdxs = nil
	}

	// Build lookup maps.
	sizeMap := make(map[string]string)
	for _, s := range indexSizes {
		sizeMap[s.Name] = s.Size
	}
	scanMap := make(map[string]int64)
	for _, s := range scans {
		scanMap[s.Index] = s.IndexScans
	}
	nullFracMap := make(map[string]string)
	for _, n := range nullIdxs {
		nullFracMap[n.Index] = n.NullFrac
	}

	var results []IndexInfoResult
	for _, idx := range indexes {
		if tableName != "" && idx.TableName != tableName {
			continue
		}
		results = append(results, IndexInfoResult{
			IndexName:  idx.IndexName,
			TableName:  idx.TableName,
			Columns:    idx.Columns,
			IndexSize:  sizeMap[idx.IndexName],
			IndexScans: scanMap[idx.IndexName],
			NullFrac:   nullFracMap[idx.IndexName],
		})
	}
	return results, nil
}
