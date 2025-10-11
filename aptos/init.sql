-- Initialize Aptos Indexer Database
-- This script sets up the necessary database structure for the Aptos indexer

-- Create schema if not exists
CREATE SCHEMA IF NOT EXISTS public;

-- Grant privileges
GRANT ALL ON SCHEMA public TO postgres;
GRANT CREATE ON SCHEMA public TO postgres;

-- Create required extensions
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Set default settings for better performance
ALTER DATABASE aptos_indexer SET maintenance_work_mem = '2GB';
ALTER DATABASE aptos_indexer SET work_mem = '256MB';
ALTER DATABASE aptos_indexer SET shared_buffers = '4GB';
ALTER DATABASE aptos_indexer SET effective_cache_size = '12GB';
ALTER DATABASE aptos_indexer SET random_page_cost = 1.1;

-- Create indexer user if needed (for better security in production)
-- CREATE USER indexer_user WITH PASSWORD 'indexer_password';
-- GRANT ALL PRIVILEGES ON DATABASE aptos_indexer TO indexer_user;
-- GRANT ALL ON SCHEMA public TO indexer_user;