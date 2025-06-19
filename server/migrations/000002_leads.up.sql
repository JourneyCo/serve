-- Remove the existing leads column
ALTER TABLE projects DROP COLUMN IF EXISTS leads;

-- Add the leads column back with proper JSONB type and default
ALTER TABLE projects ADD COLUMN leads JSONB NOT NULL DEFAULT '[]'::jsonb;
