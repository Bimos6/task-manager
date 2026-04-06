package task

import (
	"errors"
	"fmt"
	"time"
)

type RecurrenceRule struct {
	Type       string   `json:"type"`
	Interval   int      `json:"interval,omitempty"`
	DayOfMonth int      `json:"day_of_month,omitempty"`
	Dates      []string `json:"dates,omitempty"`
	Parity     string   `json:"parity,omitempty"`
}

func (r *RecurrenceRule) Validate() error {
	if r == nil {
		return nil
	}

	switch r.Type {
	case "daily":
		if r.Interval < 1 {
			return errors.New("daily interval must be >= 1")
		}
	case "monthly":
		if r.DayOfMonth < 1 || r.DayOfMonth > 31 {
			return errors.New("day_of_month must be between 1 and 31")
		}
	case "specific_dates":
		if len(r.Dates) == 0 {
			return errors.New("specific_dates requires at least one date")
		}
		for _, date := range r.Dates {
			if _, err := time.Parse("2006-01-02", date); err != nil {
				return fmt.Errorf("invalid date format: %s", date)
			}
		}
	case "parity":
		if r.Parity != "even" && r.Parity != "odd" {
			return errors.New("parity must be 'even' or 'odd'")
		}
	default:
		return fmt.Errorf("unknown recurrence type: %s", r.Type)
	}

	return nil
}
