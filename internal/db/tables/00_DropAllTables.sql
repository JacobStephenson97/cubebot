-- Drop all tables in the correct order to handle foreign key dependencies
-- This file should be executed first before any table creation scripts

-- First drop tables that depend on other tables
DROP TABLE IF EXISTS user_stats;
DROP TABLE IF EXISTS draft_matches;
DROP TABLE IF EXISTS draft_rounds;
DROP TABLE IF EXISTS draft_participants;
DROP TABLE IF EXISTS draft_sessions;

-- Then drop tables that are referenced by other tables
DROP TABLE IF EXISTS draft_formats;
DROP TABLE IF EXISTS guilds;
DROP TABLE IF EXISTS users;

-- Note: This ensures that tables with foreign key constraints are dropped before
-- the tables they reference, avoiding foreign key constraint violations during drops. 