-- migrations/001_add_recurrence_fields.sql
ALTER TABLE tasks 
ADD COLUMN IF NOT EXISTS is_recurring BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS recurrence_type VARCHAR(20),
ADD COLUMN IF NOT EXISTS recurrence_rule JSONB;

CREATE INDEX IF NOT EXISTS idx_tasks_is_recurring ON tasks(is_recurring);
CREATE INDEX IF NOT EXISTS idx_tasks_recurrence_type ON tasks(recurrence_type);