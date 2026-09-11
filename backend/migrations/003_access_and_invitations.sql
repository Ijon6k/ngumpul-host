-- 003_access_and_invitations.sql: Instance settings & invitation-based access control

CREATE TABLE IF NOT EXISTS instance_settings (
    key VARCHAR(64) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Default registration mode: INVITE_ONLY for private community homelab security
INSERT INTO instance_settings (key, value)
VALUES ('registration_mode', 'INVITE_ONLY')
ON CONFLICT (key) DO NOTHING;

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
