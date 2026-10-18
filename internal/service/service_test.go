package service

import (
	"errors"
	"testing"

	"sis/internal/storage"
	"sis/internal/task"
)

func TestCreateValidation(t *testing.T) {
	svc := New(storage.NewMemoryStore())
	if _, err := svc.Create("  ", "p"); !errors.Is(err, task.ErrEmptyTitle) {
		t.Fatalf("err = %v, want ErrEmptyTitle", err)
	}
	if got := svc.List(); len(got) != 0 {
		t.Errorf("invalid task was stored: %v", got)
	}
}

func TestFullLifecycle(t *testing.T) {
	svc := New(storage.NewMemoryStore())
	created, err := svc.Create("job", "data")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Start(created.ID); err != nil {
		t.Fatal(err)
	}
	done, err := svc.Complete(created.ID, "result")
	if err != nil {
		t.Fatal(err)
	}
	if done.Status != task.StatusDone || done.Result != "result" {
		t.Errorf("unexpected: %+v", done)
	}
	stored, _ := svc.Get(created.ID)
	if stored.Status != task.StatusDone {
		t.Errorf("not persisted: %+v", stored)
	}
}

func TestInvalidTransitionDoesNotChangeStore(t *testing.T) {
	svc := New(storage.NewMemoryStore())
	c, _ := svc.Create("job", "")
	if _, err := svc.Complete(c.ID, "x"); !errors.Is(err, task.ErrInvalidTransition) {
		t.Fatalf("err = %v", err)
	}
	got, _ := svc.Get(c.ID)
	if got.Status != task.StatusPending {
		t.Errorf("status = %s, want pending", got.Status)
	}
}

func TestOperationsOnMissingTask(t *testing.T) {
	svc := New(storage.NewMemoryStore())
	if _, err := svc.Start(7); !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("start: %v", err)
	}
	if _, err := svc.Get(7); !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("get: %v", err)
	}
}
