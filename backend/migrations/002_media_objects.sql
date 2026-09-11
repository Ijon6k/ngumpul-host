-- 002_media_objects.sql: Media storage metadata and project cover key support
ALTER TABLE projects ADD COLUMN IF NOT EXISTS cover_image_key TEXT NOT NULL DEFAULT '';

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
