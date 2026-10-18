package task

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusDone      Status = "done"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// Допустимые переходы между статусами (конечный автомат).
var transitions = map[Status][]Status{
	StatusPending: {StatusRunning, StatusCancelled},
	StatusRunning: {StatusDone, StatusFailed, StatusCancelled},
}

var (
	ErrEmptyTitle        = errors.New("task title must not be empty")
	ErrInvalidTransition = errors.New("invalid status transition")
)

type TransitionError struct {
	From, To Status
}

func (e *TransitionError) Error() string {
	return fmt.Sprintf("%v: %s -> %s", ErrInvalidTransition, e.From, e.To)
}

func (e *TransitionError) Unwrap() error { return ErrInvalidTransition }

func (s Status) IsTerminal() bool {
	_, hasNext := transitions[s]
	return !hasNext
}

func (s Status) CanTransitionTo(next Status) bool {
	for _, allowed := range transitions[s] {
		if allowed == next {
			return true
		}
	}
	return false
}

type Task struct {
	ID        int
	Title     string
	Payload   string
	Status    Status
	Result    string
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(title, payload string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, ErrEmptyTitle
	}
	now := time.Now()
	return Task{
		Title:     title,
		Payload:   payload,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (t *Task) transition(next Status) error {
	if !t.Status.CanTransitionTo(next) {
		return &TransitionError{From: t.Status, To: next}
	}
	t.Status = next
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) Start() error  { return t.transition(StatusRunning) }
func (t *Task) Cancel() error { return t.transition(StatusCancelled) }

func (t *Task) Complete(result string) error {
	if err := t.transition(StatusDone); err != nil {
		return err
	}
	t.Result = result
	return nil
}

func (t *Task) Fail(reason string) error {
	if err := t.transition(StatusFailed); err != nil {
		return err
	}
	t.Error = reason
	return nil
}
