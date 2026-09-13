-- Ngumpul Host Initial Database Schema

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(64) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(128) NOT NULL,
    avatar_url TEXT NOT NULL DEFAULT '',
    bio TEXT NOT NULL DEFAULT '',
    role VARCHAR(32) NOT NULL DEFAULT 'USER' CHECK (role IN ('USER', 'ADMIN')),
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('PENDING_VERIFICATION', 'ACTIVE', 'SUSPENDED', 'DELETED')),
    email_verified BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id VARCHAR(128) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);

-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    name VARCHAR(128) NOT NULL,
    slug VARCHAR(128) UNIQUE NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    cover_image_url TEXT NOT NULL DEFAULT '',
    cover_image_key TEXT NOT NULL DEFAULT '',
    repository_url TEXT NOT NULL DEFAULT '',
    documentation_url TEXT NOT NULL DEFAULT '',
    demo_url TEXT NOT NULL DEFAULT '',
    readme TEXT NOT NULL DEFAULT '',
    technology_stack TEXT[] NOT NULL DEFAULT '{}',
    hosting_type VARCHAR(32) NOT NULL DEFAULT 'HOSTED_HERE' CHECK (hosting_type IN ('HOSTED_HERE', 'EXTERNAL')),
    public_url TEXT NOT NULL DEFAULT '',
    status VARCHAR(32) NOT NULL DEFAULT 'ONLINE' CHECK (status IN ('PENDING', 'SETUP', 'ONLINE', 'OFFLINE', 'ARCHIVED')),
    visibility VARCHAR(32) NOT NULL DEFAULT 'PUBLIC' CHECK (visibility IN ('PUBLIC', 'UNPUBLISHED', 'ARCHIVED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_projects_slug ON projects(slug);
CREATE INDEX IF NOT EXISTS idx_projects_owner_id ON projects(owner_id);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_visibility ON projects(visibility);

-- Hosting Requests table
CREATE TABLE IF NOT EXISTS hosting_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requester_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_name VARCHAR(128) NOT NULL,
    subdomain VARCHAR(64) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    cover_image_url TEXT NOT NULL DEFAULT '',
    readme TEXT NOT NULL DEFAULT '',
    repository_url TEXT NOT NULL DEFAULT '',
    documentation_url TEXT NOT NULL DEFAULT '',
    deployment_notes TEXT NOT NULL DEFAULT '',
    technology_stack TEXT[] NOT NULL DEFAULT '{}',
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'REVIEWING', 'APPROVED', 'REJECTED', 'SETUP', 'COMPLETED', 'CANCELLED')),
    admin_notes TEXT NOT NULL DEFAULT '',
    reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_hosting_requests_requester_id ON hosting_requests(requester_id);
CREATE INDEX IF NOT EXISTS idx_hosting_requests_status ON hosting_requests(status);

-- Activities table
CREATE TABLE IF NOT EXISTS activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    type VARCHAR(64) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    visibility VARCHAR(32) NOT NULL DEFAULT 'PUBLIC' CHECK (visibility IN ('PUBLIC', 'ADMIN')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_activities_visibility_created ON activities(visibility, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activities_project_id ON activities(project_id);

-- Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'::jsonb,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id, created_at DESC);

-- Audit logs table (Admin only)
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(64) NOT NULL,
    target_type VARCHAR(64) NOT NULL,
    target_id VARCHAR(128) NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- System status items
CREATE TABLE IF NOT EXISTS system_status (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'OPERATIONAL',
    last_checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    response_time_ms INT NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb
);

-- Availability state transitions (persists only status changes, not every tick)
CREATE TABLE IF NOT EXISTS availability_events (
    id BIGSERIAL PRIMARY KEY,
    status VARCHAR(16) NOT NULL, -- 'UP' or 'DOWN'
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ,
    duration_seconds BIGINT DEFAULT 0,
    latency_ms INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_avail_events_status ON availability_events(status);
CREATE INDEX IF NOT EXISTS idx_avail_events_started ON availability_events(started_at);

-- Availability persistent single-row state (tracks session and heartbeat in-place, zero row spamming)
CREATE TABLE IF NOT EXISTS availability_state (
    id INT PRIMARY KEY DEFAULT 1,
    boot_id VARCHAR(64) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    last_seen_at TIMESTAMPTZ NOT NULL,
    first_monitored_at TIMESTAMPTZ NOT NULL,
    heartbeat_count BIGINT NOT NULL DEFAULT 1,
    current_status VARCHAR(16) NOT NULL DEFAULT 'OPERATIONAL',
    CONSTRAINT single_row_state CHECK (id = 1)
);

-- Meaningful downtime gaps and outage incidents only
CREATE TABLE IF NOT EXISTS availability_incidents (
    id BIGSERIAL PRIMARY KEY,
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ NOT NULL,
    duration_seconds BIGINT NOT NULL,
    cause VARCHAR(64) NOT NULL,
    details TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_avail_incidents_started ON availability_incidents(started_at DESC);

-- Media objects storage metadata (tracks binaries persisted in SeaweedFS / S3 / Local)
CREATE TABLE IF NOT EXISTS media_objects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    purpose VARCHAR(32) NOT NULL DEFAULT 'UPLOAD',
    object_key TEXT NOT NULL UNIQUE,
    content_type VARCHAR(64) NOT NULL,
    byte_size BIGINT NOT NULL,
    width INT NOT NULL DEFAULT 0,
    height INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_media_objects_owner_id ON media_objects(owner_id);
CREATE INDEX IF NOT EXISTS idx_media_objects_object_key ON media_objects(object_key);

-- Instance Settings (Key-Value configuration for self-hosted community node)
CREATE TABLE IF NOT EXISTS instance_settings (
    key VARCHAR(64) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO instance_settings (key, value)
VALUES ('registration_mode', 'INVITE_ONLY')
ON CONFLICT (key) DO NOTHING;

-- Invitations for controlled community entry
CREATE TABLE IF NOT EXISTS invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    token VARCHAR(128),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    invited_email VARCHAR(255),
    max_uses INT NOT NULL DEFAULT 1 CHECK (max_uses > 0),
    used_count INT NOT NULL DEFAULT 0 CHECK (used_count >= 0),
    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE invitations ADD COLUMN IF NOT EXISTS token VARCHAR(128);

CREATE INDEX IF NOT EXISTS idx_invitations_token_hash ON invitations(token_hash);
CREATE INDEX IF NOT EXISTS idx_invitations_expires_at ON invitations(expires_at);
CREATE INDEX IF NOT EXISTS idx_invitations_created_by ON invitations(created_by);

-- 1. Comments table
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

ALTER TABLE comments ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_comments_project_created ON comments(project_id, created_at ASC);
CREATE INDEX IF NOT EXISTS idx_comments_author_id ON comments(author_id);
CREATE INDEX IF NOT EXISTS idx_comments_deleted_at ON comments(deleted_at);

-- 2. Reports table (Projects and Comments)
CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_type VARCHAR(32) NOT NULL CHECK (target_type IN ('PROJECT', 'COMMENT')),
    target_id UUID NOT NULL,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    reason VARCHAR(64) NOT NULL CHECK (reason IN ('SPAM', 'ABUSE_HARASSMENT', 'INAPPROPRIATE', 'MALICIOUS_SUSPICIOUS', 'OTHER')),
    details TEXT NOT NULL DEFAULT '',
    status VARCHAR(32) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'REVIEWED', 'DISMISSED', 'RESOLVED')),
    resolved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    resolved_at TIMESTAMPTZ,
    action_notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reports_status_created ON reports(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_reports_target ON reports(target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_reports_project ON reports(project_id);

-- 3. Outbound Project Visits ("Visits from Ngumpul" via /go/:slug)
CREATE TABLE IF NOT EXISTS project_visits (
    id BIGSERIAL PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    visited_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    date_bucket DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE INDEX IF NOT EXISTS idx_project_visits_project_date ON project_visits(project_id, date_bucket);

-- 4. Daily Project Page Views (Privacy-conscious daily aggregates)
CREATE TABLE IF NOT EXISTS project_page_views (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    views_count INT NOT NULL DEFAULT 1,
    PRIMARY KEY (project_id, date)
);

CREATE INDEX IF NOT EXISTS idx_project_views_date ON project_page_views(date);

-- 5. Project Health & Availability Checks (HTTP probes)
CREATE TABLE IF NOT EXISTS project_availability_checks (
    id BIGSERIAL PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status_code INT NOT NULL DEFAULT 0,
    response_time_ms INT NOT NULL DEFAULT 0,
    is_successful BOOLEAN NOT NULL DEFAULT true,
    consecutive_failures INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_proj_avail_project_checked ON project_availability_checks(project_id, checked_at DESC);

-- Progressive migrations
ALTER TABLE hosting_requests ADD COLUMN IF NOT EXISTS subdomain VARCHAR(64) NOT NULL DEFAULT '';
ALTER TABLE hosting_requests ADD COLUMN IF NOT EXISTS cover_image_url TEXT NOT NULL DEFAULT '';
ALTER TABLE hosting_requests ADD COLUMN IF NOT EXISTS readme TEXT NOT NULL DEFAULT '';
ALTER TABLE projects ADD COLUMN IF NOT EXISTS readme TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_hosting_requests_subdomain ON hosting_requests(subdomain);

ALTER TABLE hosting_requests ADD COLUMN IF NOT EXISTS project_id UUID REFERENCES projects(id) ON DELETE SET NULL;
ALTER TABLE hosting_requests ADD COLUMN IF NOT EXISTS request_type VARCHAR(32) NOT NULL DEFAULT 'NEW_PROJECT';

CREATE INDEX IF NOT EXISTS idx_hosting_requests_project_id ON hosting_requests(project_id);
CREATE INDEX IF NOT EXISTS idx_hosting_requests_request_type ON hosting_requests(request_type);

-- Project state decoupling
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_status_check;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS lifecycle_status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE' CHECK (lifecycle_status IN ('SETUP', 'ACTIVE', 'SUSPENDED', 'ARCHIVED'));
ALTER TABLE projects ADD COLUMN IF NOT EXISTS availability VARCHAR(32) NOT NULL DEFAULT 'UNKNOWN' CHECK (availability IN ('UNKNOWN', 'REACHABLE', 'UNREACHABLE'));
ALTER TABLE projects ADD COLUMN IF NOT EXISTS availability_reason VARCHAR(64) NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_projects_lifecycle_status ON projects(lifecycle_status);
CREATE INDEX IF NOT EXISTS idx_projects_availability ON projects(availability);



