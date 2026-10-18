package task

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{"valid", "Do work", nil},
		{"empty", "", ErrEmptyTitle},
		{"spaces only", "   ", ErrEmptyTitle},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.title, "p")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && got.Status != StatusPending {
				t.Errorf("status = %s, want pending", got.Status)
			}
		})
	}
}

func TestTransitions(t *testing.T) {
	tests := []struct {
		name  string
		from  Status
		act   func(*Task) error
		want  Status
		valid bool
	}{
		{"pending->running", StatusPending, (*Task).Start, StatusRunning, true},
		{"pending->cancelled", StatusPending, (*Task).Cancel, StatusCancelled, true},
		{"pending->done", StatusPending, func(t *Task) error { return t.Complete("x") }, StatusPending, false},
		{"running->done", StatusRunning, func(t *Task) error { return t.Complete("x") }, StatusDone, true},
		{"running->failed", StatusRunning, func(t *Task) error { return t.Fail("boom") }, StatusFailed, true},
		{"done->running", StatusDone, (*Task).Start, StatusDone, false},
		{"failed->cancelled", StatusFailed, (*Task).Cancel, StatusFailed, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{Status: tt.from}
			err := tt.act(&task)
			if tt.valid && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.valid {
				if !errors.Is(err, ErrInvalidTransition) {
					t.Fatalf("err = %v, want ErrInvalidTransition", err)
				}
				var te *TransitionError
				if !errors.As(err, &te) || te.From != tt.from {
					t.Fatalf("expected *TransitionError from %s, got %v", tt.from, err)
				}
			}
			if task.Status != tt.want {
				t.Errorf("status = %s, want %s", task.Status, tt.want)
			}
		})
	}
}

func TestResultAndError(t *testing.T) {
	a := Task{Status: StatusRunning}
	_ = a.Complete("ok")
	if a.Result != "ok" {
		t.Errorf("result = %q", a.Result)
	}
	b := Task{Status: StatusRunning}
	_ = b.Fail("boom")
	if b.Error != "boom" {
		t.Errorf("error = %q", b.Error)
	}
}

func TestIsTerminal(t *testing.T) {
	for _, s := range []Status{StatusDone, StatusFailed, StatusCancelled} {
		if !s.IsTerminal() {
			t.Errorf("%s should be terminal", s)
		}
	}
	for _, s := range []Status{StatusPending, StatusRunning} {
		if s.IsTerminal() {
			t.Errorf("%s should not be terminal", s)
		}
	}
}
