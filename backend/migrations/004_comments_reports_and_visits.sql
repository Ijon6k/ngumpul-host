-- 004_comments_reports_and_visits.sql: Community comments, content moderation reports, outbound visits, and project availability checks

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
