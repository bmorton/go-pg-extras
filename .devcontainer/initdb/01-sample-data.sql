-- =============================================================================
-- Sample dataset for go-pg-extras testing
--
-- Creates tables, indexes, and data that exercise pg-extras diagnostic queries:
--   - Table/index sizes and bloat
--   - Index usage and unused indexes
--   - Sequential scans vs index scans
--   - Cache hit ratios
--   - Vacuum and analyze stats
--   - Bloat estimation
--   - Null/duplicate indexes
-- =============================================================================

-- Enable pg_stat_statements if available (useful for calls/outliers queries)
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Enable sslinfo for ssl_used query
CREATE EXTENSION IF NOT EXISTS sslinfo;

-- =============================================================================
-- 1. orders: a well-indexed, heavily-used table
-- =============================================================================
CREATE TABLE orders (
    id          BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    product_id  BIGINT NOT NULL,
    quantity    INT NOT NULL DEFAULT 1,
    total_cents BIGINT NOT NULL DEFAULT 0,
    status      TEXT NOT NULL DEFAULT 'pending',
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_customer_id ON orders (customer_id);
CREATE INDEX idx_orders_product_id  ON orders (product_id);
CREATE INDEX idx_orders_status      ON orders (status);
CREATE INDEX idx_orders_created_at  ON orders (created_at);

-- =============================================================================
-- 2. customers: table with some unused indexes to detect
-- =============================================================================
CREATE TABLE customers (
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    region     TEXT,
    tier       TEXT DEFAULT 'free',
    signed_up  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- These indexes are intentionally unused (will show up in unused_indexes)
CREATE INDEX idx_customers_region     ON customers (region);
CREATE INDEX idx_customers_tier       ON customers (tier);
CREATE INDEX idx_customers_signed_up  ON customers (signed_up);

-- =============================================================================
-- 3. products: table with a duplicate index scenario
-- =============================================================================
CREATE TABLE products (
    id          BIGSERIAL PRIMARY KEY,
    sku         TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    category    TEXT,
    price_cents BIGINT NOT NULL DEFAULT 0,
    weight_kg   NUMERIC(10,2),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Duplicate index on (category) — useful for duplicate_indexes query
CREATE INDEX idx_products_category   ON products (category);
CREATE INDEX idx_products_category_2 ON products (category);

-- =============================================================================
-- 4. events: an append-heavy table for bloat/vacuum testing
-- =============================================================================
CREATE TABLE events (
    id         BIGSERIAL PRIMARY KEY,
    event_type TEXT NOT NULL,
    payload    JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_events_type       ON events (event_type);
CREATE INDEX idx_events_created_at ON events (created_at);

-- =============================================================================
-- 5. audit_log: a large text-heavy table for size queries
-- =============================================================================
CREATE TABLE audit_log (
    id         BIGSERIAL PRIMARY KEY,
    actor      TEXT NOT NULL,
    action     TEXT NOT NULL,
    target     TEXT,
    detail     TEXT,
    recorded   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Partial index (exercises null_indexes detection differently)
CREATE INDEX idx_audit_log_target ON audit_log (target) WHERE target IS NOT NULL;

-- =============================================================================
-- 6. sessions: table with null-frac index potential
-- =============================================================================
CREATE TABLE sessions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    BIGINT NOT NULL,
    token      TEXT NOT NULL,
    ip_address INET,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_user_id    ON sessions (user_id);
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);

-- =============================================================================
-- Seed data
-- =============================================================================

-- Customers (1,000 rows)
INSERT INTO customers (email, name, region, tier, signed_up)
SELECT
    'user' || g || '@example.com',
    'Customer ' || g,
    (ARRAY['us-east','us-west','eu-west','eu-central','ap-south'])[1 + (g % 5)],
    (ARRAY['free','starter','pro','enterprise'])[1 + (g % 4)],
    now() - (random() * interval '730 days')
FROM generate_series(1, 1000) AS g;

-- Products (500 rows)
INSERT INTO products (sku, name, category, price_cents, weight_kg, created_at)
SELECT
    'SKU-' || lpad(g::text, 6, '0'),
    'Product ' || g,
    (ARRAY['electronics','clothing','home','sports','books','toys'])[1 + (g % 6)],
    (random() * 50000)::bigint + 100,
    round((random() * 25)::numeric, 2),
    now() - (random() * interval '365 days')
FROM generate_series(1, 500) AS g;

-- Orders (50,000 rows — provides meaningful stats)
INSERT INTO orders (customer_id, product_id, quantity, total_cents, status, notes, created_at, updated_at)
SELECT
    1 + (random() * 999)::bigint,
    1 + (random() * 499)::bigint,
    1 + (random() * 9)::int,
    (random() * 100000)::bigint + 500,
    (ARRAY['pending','processing','shipped','delivered','cancelled'])[1 + (g % 5)],
    CASE WHEN random() < 0.3 THEN 'Note for order ' || g ELSE NULL END,
    now() - (random() * interval '365 days'),
    now() - (random() * interval '30 days')
FROM generate_series(1, 50000) AS g;

-- Events (20,000 rows)
INSERT INTO events (event_type, payload, created_at)
SELECT
    (ARRAY['page_view','click','signup','purchase','logout','error'])[1 + (g % 6)],
    jsonb_build_object(
        'request_id', gen_random_uuid(),
        'value', (random() * 1000)::int
    ),
    now() - (random() * interval '90 days')
FROM generate_series(1, 20000) AS g;

-- Audit log (10,000 rows)
INSERT INTO audit_log (actor, action, target, detail, recorded)
SELECT
    'user-' || (1 + (random() * 99)::int),
    (ARRAY['create','update','delete','login','export'])[1 + (g % 5)],
    CASE WHEN random() < 0.8 THEN 'resource-' || (random() * 500)::int ELSE NULL END,
    'Detail line ' || g,
    now() - (random() * interval '180 days')
FROM generate_series(1, 10000) AS g;

-- Sessions (5,000 rows — some with NULL expires_at for null_indexes query)
INSERT INTO sessions (user_id, token, ip_address, expires_at, created_at)
SELECT
    1 + (random() * 999)::bigint,
    md5(random()::text),
    ('10.' || (random()*255)::int || '.' || (random()*255)::int || '.' || (random()*255)::int)::inet,
    CASE WHEN random() < 0.7 THEN now() + (random() * interval '30 days') ELSE NULL END,
    now() - (random() * interval '60 days')
FROM generate_series(1, 5000) AS g;

-- =============================================================================
-- Force stats collection so pg-extras queries return data immediately
-- =============================================================================
ANALYZE;

-- Run a few queries to populate pg_stat_user_tables / pg_stat_user_indexes
SELECT count(*) FROM orders WHERE customer_id = 42;
SELECT count(*) FROM orders WHERE status = 'shipped';
SELECT * FROM orders ORDER BY created_at DESC LIMIT 10;
SELECT count(*) FROM events WHERE event_type = 'purchase';
