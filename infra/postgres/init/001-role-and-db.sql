-- Initialization script for Pontoon Postgres
-- Creates role and database if they do not already exist.

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'game') THEN
        CREATE ROLE game WITH LOGIN PASSWORD 'game';
        -- Grant superuser for dev convenience; remove later if not needed.
        ALTER ROLE game WITH SUPERUSER;
    END IF;
END
$$;

-- Create database if missing and set owner
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'pontoon') THEN
        CREATE DATABASE pontoon OWNER game;
    END IF;
END
$$;

-- Future placeholder: create schemas/tables here.
-- Example:
-- CREATE SCHEMA IF NOT EXISTS pontoon AUTHORIZATION game;
