-- +goose Up
-- +goose StatementBegin

-- Create the partman schema first
CREATE SCHEMA IF NOT EXISTS partman;

-- Create extensions (pg_partman will use the partman schema)
CREATE EXTENSION IF NOT EXISTS pg_partman WITH SCHEMA partman;
CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Create the audit log table
CREATE TABLE IF NOT EXISTS "audit_log"(
  id SERIAL,
  uuid UUID NOT NULL,
  event_type VARCHAR(64) NOT NULL,
  actor_id UUID,
  payload JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at); 

-- Configure monthly partitioning
SELECT partman.create_parent(
    p_parent_table := 'public.audit_log',
    p_control := 'created_at',
    p_interval := '1 month',
    p_premake := 2
);

-- Set retention policy
UPDATE partman.part_config
SET retention = '12 months',
    retention_keep_table = false
WHERE parent_table = 'public.audit_log';

-- Schedule automatic maintenance (runs daily at 1 AM)
SELECT cron.schedule(
    'partman-maintenance', 
    '0 1 * * *', 
    'SELECT partman.run_maintenance();'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove cron job
SELECT cron.unschedule('partman-maintenance');

-- Drop the table (this will also drop all partitions)
DROP TABLE IF EXISTS "audit_log" CASCADE;

-- Drop extensions
DROP EXTENSION IF EXISTS pg_partman CASCADE;
DROP EXTENSION IF EXISTS pg_cron;

-- Drop the partman schema
DROP SCHEMA IF EXISTS partman CASCADE;

-- +goose StatementEnd
