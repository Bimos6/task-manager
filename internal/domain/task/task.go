package task

import (
	"fmt"
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
	IsRecurring    bool            `json:"is_recurring" gorm:"default:false"`
	RecurrenceType string          `json:"recurrence_type,omitempty"`
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

func (t *Task) Normalize() {
	if t.RecurrenceRule != nil && t.RecurrenceRule.Type != "" {
		t.IsRecurring = true
		t.RecurrenceType = t.RecurrenceRule.Type
	} else {
		t.IsRecurring = false
		t.RecurrenceType = ""
		t.RecurrenceRule = nil
	}
}

func (t *Task) BeforeSave() error {
	t.Normalize()

	if t.IsRecurring && t.RecurrenceRule != nil {
		if err := t.RecurrenceRule.Validate(); err != nil {
			return fmt.Errorf("invalid recurrence rule: %w", err)
		}
	}

	return nil
}
