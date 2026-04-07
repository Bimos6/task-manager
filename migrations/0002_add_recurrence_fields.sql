ALTER TABLE tasks 
ADD COLUMN IF NOT EXISTS due_date TIMESTAMP,
ADD COLUMN IF NOT EXISTS recurrence_rule JSONB;

CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);