package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo       Repository
	recService *RecurrenceService
	now        func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:       repo,
		recService: NewRecurrenceService(),
		now:        func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	if normalized.RecurrenceRule != nil {
		if err := normalized.RecurrenceRule.Validate(); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		}
	}

	model := &taskdomain.Task{
		Title:          normalized.Title,
		Description:    normalized.Description,
		Status:         normalized.Status,
		RecurrenceRule: normalized.RecurrenceRule,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	if normalized.RecurrenceRule != nil {
		if err := normalized.RecurrenceRule.Validate(); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		}
	}

	model := &taskdomain.Task{
		ID:             id,
		Title:          normalized.Title,
		Description:    normalized.Description,
		Status:         normalized.Status,
		RecurrenceRule: normalized.RecurrenceRule,
		UpdatedAt:      s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListInRange(ctx context.Context, startDate, endDate time.Time) ([]taskdomain.Task, error) {
	templates, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []taskdomain.Task

	for _, t := range templates {
		if t.RecurrenceRule != nil && t.RecurrenceRule.Type != "" {
			dates, err := s.recService.GenerateDatesInRange(t.RecurrenceRule, startDate, endDate)
			if err != nil {
				continue
			}

			for _, date := range dates {
				taskCopy := t
				taskCopy.DueDate = date
				result = append(result, taskCopy)
			}
		} else {
			if !t.DueDate.Before(startDate) && !t.DueDate.After(endDate) {
				result = append(result, t)
			}
		}
	}

	return result, nil
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}
