-- =============================================================================
-- Generate traffic to populate pg_stat_* views for the pg-extras dashboard.
-- Run with: psql "$DATABASE_URL" -q -f scripts/generate-traffic.sql
-- =============================================================================

-- Index scans: hit various indexes so they show usage
SELECT count(*) FROM orders WHERE customer_id BETWEEN 1 AND 100;
SELECT count(*) FROM orders WHERE customer_id BETWEEN 200 AND 400;
SELECT count(*) FROM orders WHERE product_id IN (1,5,10,50,100,200,499);
SELECT count(*) FROM orders WHERE status = 'shipped';
SELECT count(*) FROM orders WHERE status = 'delivered';
SELECT count(*) FROM orders WHERE status = 'pending';
SELECT count(*) FROM orders WHERE created_at > now() - interval '7 days';
SELECT count(*) FROM orders WHERE created_at > now() - interval '30 days';
SELECT * FROM orders WHERE customer_id = 42 ORDER BY created_at DESC LIMIT 10;
SELECT * FROM orders WHERE product_id = 7 ORDER BY created_at DESC LIMIT 10;

-- Sequential scans on customers (intentionally skip using indexes)
SELECT count(*) FROM customers;
SELECT * FROM customers WHERE name LIKE '%500%';
SELECT * FROM customers WHERE signed_up > now() - interval '90 days' ORDER BY name;

-- Events: mix of index and sequential access
SELECT count(*) FROM events WHERE event_type = 'purchase';
SELECT count(*) FROM events WHERE event_type = 'page_view';
SELECT count(*) FROM events WHERE created_at > now() - interval '7 days';
SELECT * FROM events ORDER BY created_at DESC LIMIT 50;

-- Audit log
SELECT count(*) FROM audit_log WHERE target = 'resource-42';
SELECT count(*) FROM audit_log WHERE action = 'create';
SELECT * FROM audit_log ORDER BY recorded DESC LIMIT 20;

-- Sessions
SELECT count(*) FROM sessions WHERE user_id = 42;
SELECT count(*) FROM sessions WHERE expires_at IS NOT NULL AND expires_at < now();

-- Products
SELECT * FROM products WHERE category = 'electronics' ORDER BY price_cents DESC;
SELECT * FROM products ORDER BY created_at DESC LIMIT 20;

-- Heavier queries for outlier/call stats
SELECT o.id, c.name, p.name, o.total_cents
  FROM orders o
  JOIN customers c ON c.id = o.customer_id
  JOIN products p ON p.id = o.product_id
  WHERE o.status = 'shipped' ORDER BY o.total_cents DESC LIMIT 100;

SELECT customer_id, count(*), sum(total_cents)
  FROM orders GROUP BY customer_id ORDER BY sum(total_cents) DESC LIMIT 50;

SELECT event_type, count(*) FROM events GROUP BY event_type ORDER BY count(*) DESC;

SELECT category, count(*), avg(price_cents) FROM products GROUP BY category;

-- Generate some UPDATE churn for vacuum/bloat stats
UPDATE events SET payload = payload || '{"touched": true}'::jsonb
  WHERE id IN (SELECT id FROM events ORDER BY random() LIMIT 2000);

UPDATE orders SET updated_at = now()
  WHERE id IN (SELECT id FROM orders ORDER BY random() LIMIT 5000);

DELETE FROM events WHERE id IN (SELECT id FROM events ORDER BY random() LIMIT 1000);

-- Re-analyze to refresh stats
ANALYZE;
