/* Index and table hit rate */

SELECT
  'index hit rate' AS name,
  COALESCE((sum(idx_blks_hit)) / nullif(sum(idx_blks_hit + idx_blks_read),0), 0) AS ratio
FROM pg_statio_user_indexes
WHERE schemaname = '%{schema}'
UNION ALL
SELECT
 'table hit rate' AS name,
  COALESCE(sum(heap_blks_hit) / nullif(sum(heap_blks_hit) + sum(heap_blks_read),0), 0) AS ratio
FROM pg_statio_user_tables
WHERE schemaname = '%{schema}';
