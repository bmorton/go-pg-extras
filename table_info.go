package pgextras

import "context"

// TableInfoResult holds aggregated table information.
type TableInfoResult struct {
	TableName     string `json:"table_name"`
	TableSize     string `json:"table_size"`
	TableCacheHit string `json:"table_cache_hit"`
	IndexCacheHit string `json:"index_cache_hit"`
	EstimatedRows int64  `json:"estimated_rows"`
	SeqScans      int64  `json:"seq_scans"`
	IdxScans      int64  `json:"idx_scans"`
}

// TableInfo returns aggregated information about tables.
func (c *Client) TableInfo(ctx context.Context, tableName string) ([]TableInfoResult, error) {
	sizes, err := c.TableSize(ctx)
	if err != nil {
		return nil, err
	}
	tableCacheHits, err := c.TableCacheHit(ctx)
	if err != nil {
		return nil, err
	}
	indexCacheHits, err := c.IndexCacheHit(ctx)
	if err != nil {
		return nil, err
	}
	records, err := c.RecordsRank(ctx)
	if err != nil {
		return nil, err
	}
	seqScans, err := c.SeqScans(ctx)
	if err != nil {
		return nil, err
	}
	idxScans, err := c.TableIndexScans(ctx)
	if err != nil {
		return nil, err
	}

	// Build lookup maps.
	sizeMap := make(map[string]string)
	for _, s := range sizes {
		sizeMap[s.Name] = s.Size
	}
	tableCacheMap := make(map[string]string)
	for _, t := range tableCacheHits {
		tableCacheMap[t.Name] = t.Ratio
	}
	indexCacheMap := make(map[string]string)
	for _, i := range indexCacheHits {
		indexCacheMap[i.Name] = i.Ratio
	}
	recordsMap := make(map[string]int64)
	for _, r := range records {
		recordsMap[r.Name] = r.EstimatedCount
	}
	seqScanMap := make(map[string]int64)
	for _, s := range seqScans {
		seqScanMap[s.Name] = s.Count
	}
	idxScanMap := make(map[string]int64)
	for _, i := range idxScans {
		idxScanMap[i.Name] = i.Count
	}

	var results []TableInfoResult
	for _, s := range sizes {
		if tableName != "" && s.Name != tableName {
			continue
		}
		results = append(results, TableInfoResult{
			TableName:     s.Name,
			TableSize:     s.Size,
			TableCacheHit: tableCacheMap[s.Name],
			IndexCacheHit: indexCacheMap[s.Name],
			EstimatedRows: recordsMap[s.Name],
			SeqScans:      seqScanMap[s.Name],
			IdxScans:      idxScanMap[s.Name],
		})
	}
	return results, nil
}
