-- Migration: 006_subdomain_change_request.sql
-- Description: Add project_id and request_type to hosting_requests for subdomain change requests.

ALTER TABLE hosting_requests ADD COLUMN IF NOT EXISTS project_id UUID REFERENCES projects(id) ON DELETE SET NULL;
ALTER TABLE hosting_requests ADD COLUMN IF NOT EXISTS request_type VARCHAR(32) NOT NULL DEFAULT 'NEW_PROJECT';

CREATE INDEX IF NOT EXISTS idx_hosting_requests_project_id ON hosting_requests(project_id);
CREATE INDEX IF NOT EXISTS idx_hosting_requests_request_type ON hosting_requests(request_type);
