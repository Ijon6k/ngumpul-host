-- Ngumpul Host: Project State Architecture Decoupling
-- Decouples project lifecycle, visibility, and availability into separate concepts.

-- 1. Relax or drop legacy projects_status_check
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_status_check;

-- 2. Add decoupled lifecycle and availability columns
ALTER TABLE projects ADD COLUMN IF NOT EXISTS lifecycle_status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE' CHECK (lifecycle_status IN ('SETUP', 'ACTIVE', 'SUSPENDED', 'ARCHIVED'));
ALTER TABLE projects ADD COLUMN IF NOT EXISTS availability VARCHAR(32) NOT NULL DEFAULT 'UNKNOWN' CHECK (availability IN ('UNKNOWN', 'REACHABLE', 'UNREACHABLE'));
ALTER TABLE projects ADD COLUMN IF NOT EXISTS availability_reason VARCHAR(64) NOT NULL DEFAULT '';

-- 3. Backfill existing projects data based on previous status
UPDATE projects
SET lifecycle_status = CASE
        WHEN status = 'SETUP' OR status = 'PENDING' THEN 'SETUP'
        WHEN status = 'ARCHIVED' THEN 'ARCHIVED'
        ELSE 'ACTIVE'
    END,
    availability = CASE
        WHEN status = 'ONLINE' THEN 'REACHABLE'
        WHEN status = 'OFFLINE' THEN 'UNREACHABLE'
        ELSE 'UNKNOWN'
    END
WHERE lifecycle_status = 'ACTIVE' AND availability = 'UNKNOWN';

-- 4. Create performance indexes
CREATE INDEX IF NOT EXISTS idx_projects_lifecycle_status ON projects(lifecycle_status);
CREATE INDEX IF NOT EXISTS idx_projects_availability ON projects(availability);
