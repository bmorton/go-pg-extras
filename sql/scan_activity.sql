/* Combined scan activity: index scans and sequential scans per table */

SELECT
  relname AS name,
  idx_scan AS index_scans,
  seq_scan AS sequential_scans,
  idx_scan + seq_scan AS total_scans,
  CASE WHEN (idx_scan + seq_scan) > 0
    THEN round(100.0 * idx_scan / (idx_scan + seq_scan), 1)
    ELSE 0
  END AS index_scan_pct
FROM
  pg_stat_user_tables
WHERE
  schemaname = '%{schema}'
ORDER BY (idx_scan + seq_scan) DESC;
