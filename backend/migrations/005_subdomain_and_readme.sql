-- Migration: 005_subdomain_and_readme.sql
-- Description: Add subdomain and readme to hosting_requests, and readme to projects.

ALTER TABLE hosting_requests ADD COLUMN IF NOT EXISTS subdomain VARCHAR(64) NOT NULL DEFAULT '';
ALTER TABLE hosting_requests ADD COLUMN IF NOT EXISTS readme TEXT NOT NULL DEFAULT '';
ALTER TABLE projects ADD COLUMN IF NOT EXISTS readme TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_hosting_requests_subdomain ON hosting_requests(subdomain);
