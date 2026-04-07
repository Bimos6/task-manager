package task

import (
	"errors"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type RecurrenceService struct{}

func NewRecurrenceService() *RecurrenceService {
	return &RecurrenceService{}
}

func (s *RecurrenceService) GenerateDatesInRange(rule *taskdomain.RecurrenceRule, startDate, endDate time.Time) ([]time.Time, error) {

	if rule == nil {
		return nil, errors.New("recurrence rule is nil")
	}

	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

	switch rule.Type {
	case "daily":
		return s.generateDaily(rule, startDate, endDate)
	case "monthly":
		return s.generateMonthly(rule, startDate, endDate)
	case "specific_dates":
		return s.generateSpecificDates(rule, startDate, endDate)
	case "parity":
		return s.generateParity(rule, startDate, endDate)
	default:
		return nil, errors.New("unknown recurrence type: " + rule.Type)
	}
}

func (s *RecurrenceService) generateDaily(rule *taskdomain.RecurrenceRule, startDate, endDate time.Time) ([]time.Time, error) {
	var dates []time.Time
	interval := rule.Interval
	if interval <= 0 {
		interval = 1
	}
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, interval) {
		dates = append(dates, d)
	}
	return dates, nil
}

func (s *RecurrenceService) generateMonthly(rule *taskdomain.RecurrenceRule, startDate, endDate time.Time) ([]time.Time, error) {
	var dates []time.Time
	day := rule.DayOfMonth
	if day < 1 || day > 31 {
		return dates, nil
	}
	current := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	for !current.After(endDate) {
		target := time.Date(current.Year(), current.Month(), day, 0, 0, 0, 0, current.Location())
		if target.Day() == day && !target.Before(startDate) && !target.After(endDate) {
			dates = append(dates, target)
		}
		current = current.AddDate(0, 1, 0)
	}
	return dates, nil
}

func (s *RecurrenceService) generateSpecificDates(rule *taskdomain.RecurrenceRule, startDate, endDate time.Time) ([]time.Time, error) {
	var dates []time.Time
	for _, dateStr := range rule.Dates {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		if !date.Before(startDate) && !date.After(endDate) {
			dates = append(dates, date)
		}
	}
	return dates, nil
}

func (s *RecurrenceService) generateParity(rule *taskdomain.RecurrenceRule, startDate, endDate time.Time) ([]time.Time, error) {
	var dates []time.Time
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if rule.Parity == "even" && d.Day()%2 == 0 {
			dates = append(dates, d)
		} else if rule.Parity == "odd" && d.Day()%2 == 1 {
			dates = append(dates, d)
		}
	}
	return dates, nil
}
