/* Combined table overview: size, total size (with indexes), index size, and estimated row count */

SELECT
  c.relname AS name,
  pg_size_pretty(pg_table_size(c.oid)) AS table_size,
  pg_table_size(c.oid) AS table_size_bytes,
  pg_size_pretty(pg_total_relation_size(c.oid)) AS total_size,
  pg_total_relation_size(c.oid) AS total_size_bytes,
  pg_size_pretty(pg_indexes_size(c.oid)) AS index_size,
  pg_indexes_size(c.oid) AS index_size_bytes,
  COALESCE(s.n_live_tup, 0) AS estimated_rows
FROM pg_class c
LEFT JOIN pg_namespace n ON (n.oid = c.relnamespace)
LEFT JOIN pg_stat_user_tables s ON (s.relid = c.oid)
WHERE n.nspname = '%{schema}'
AND n.nspname !~ '^pg_toast'
AND c.relkind IN ('r', 'm')
ORDER BY pg_total_relation_size(c.oid) DESC;
