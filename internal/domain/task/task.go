package task

import (
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID             int64           `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Status         Status          `json:"status"`
	DueDate        time.Time       `json:"due_date"`
	RecurrenceRule *RecurrenceRule `json:"recurrence_rule,omitempty" gorm:"type:jsonb"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (t *Task) IsRecurring() bool {
	return t.RecurrenceRule != nil && t.RecurrenceRule.Type != ""
}
