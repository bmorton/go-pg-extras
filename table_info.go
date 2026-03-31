package pgextras

import (
	"context"
	"sort"
)

// TableSchemaDetail holds grouped schema and foreign key info for a single table.
type TableSchemaDetail struct {
	TableName   string               `json:"table_name"`
	Columns     []TableSchemasResult `json:"columns"`
	ForeignKeys []ForeignKeysResult  `json:"foreign_keys,omitempty"`
}

// TableSchemasGrouped returns table schema information grouped by table,
// including column definitions and foreign key relationships for each table.
func (c *Client) TableSchemasGrouped(ctx context.Context) ([]TableSchemaDetail, error) {
	schemas, err := c.TableSchemas(ctx)
	if err != nil {
		return nil, err
	}
	fks, err := c.ForeignKeys(ctx)
	if err != nil {
		return nil, err
	}

	// Group columns by table name, preserving order of first appearance.
	tableOrder := make([]string, 0)
	colsByTable := make(map[string][]TableSchemasResult)
	for _, col := range schemas {
		if _, exists := colsByTable[col.TableName]; !exists {
			tableOrder = append(tableOrder, col.TableName)
		}
		colsByTable[col.TableName] = append(colsByTable[col.TableName], col)
	}

	// Group foreign keys by table name.
	fksByTable := make(map[string][]ForeignKeysResult)
	for _, fk := range fks {
		fksByTable[fk.TableName] = append(fksByTable[fk.TableName], fk)
	}

	// Sort table names alphabetically for predictable output.
	sort.Strings(tableOrder)

	results := make([]TableSchemaDetail, 0, len(tableOrder))
	for _, name := range tableOrder {
		results = append(results, TableSchemaDetail{
			TableName:   name,
			Columns:     colsByTable[name],
			ForeignKeys: fksByTable[name],
		})
	}
	return results, nil
}

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
