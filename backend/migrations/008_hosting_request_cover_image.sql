-- Migration 008: Add cover_image_url to hosting_requests table
-- Enables users to attach project cover artwork during initial hosting submission
-- and propagates cover artwork automatically to the project catalog upon admin approval.

ALTER TABLE hosting_requests ADD COLUMN IF NOT EXISTS cover_image_url TEXT NOT NULL DEFAULT '';
