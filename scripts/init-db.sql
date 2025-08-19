-- Database initialization script for OTP Auth Service
-- This script is executed when the PostgreSQL container starts for the first time

-- Create database if it doesn't exist (this is handled by POSTGRES_DB env var)
-- But we can set up additional configurations here

-- Set timezone
SET timezone = 'UTC';

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create a dedicated user for the application (optional, for better security)
-- DO $$
-- BEGIN
--     IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'otp_auth_user') THEN
--         CREATE ROLE otp_auth_user WITH LOGIN PASSWORD 'secure_password';
--     END IF;
-- END
-- $$;

-- Grant necessary permissions
-- GRANT CONNECT ON DATABASE otp_auth TO otp_auth_user;
-- GRANT USAGE ON SCHEMA public TO otp_auth_user;
-- GRANT CREATE ON SCHEMA public TO otp_auth_user;

-- Create schema_migrations table to track migrations
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial migration records (these will be handled by the application)
-- But we can pre-populate if needed

-- Set up logging (optional)
-- ALTER SYSTEM SET log_statement = 'all';
-- ALTER SYSTEM SET log_duration = on;
-- ALTER SYSTEM SET log_min_duration_statement = 1000; -- Log queries taking more than 1 second

-- Performance tuning (optional, adjust based on your needs)
-- ALTER SYSTEM SET shared_buffers = '256MB';
-- ALTER SYSTEM SET effective_cache_size = '1GB';
-- ALTER SYSTEM SET maintenance_work_mem = '64MB';
-- ALTER SYSTEM SET checkpoint_completion_target = 0.9;
-- ALTER SYSTEM SET wal_buffers = '16MB';
-- ALTER SYSTEM SET default_statistics_target = 100;
-- ALTER SYSTEM SET random_page_cost = 1.1;
-- ALTER SYSTEM SET effective_io_concurrency = 200;

-- Reload configuration
-- SELECT pg_reload_conf();

-- Create indexes for better performance (these will be created by migrations)
-- But we can add some general ones here if needed

-- Log the completion
INSERT INTO schema_migrations (version) VALUES ('init_db_script') 
ON CONFLICT (version) DO NOTHING;

-- Display completion message
DO $$
BEGIN
    RAISE NOTICE 'Database initialization completed successfully';
END
$$;