package storage

import (
	"errors"
	"testing"

	"sis/internal/task"
)

func newTask(t *testing.T, title string) task.Task {
	t.Helper()
	tk, err := task.New(title, "")
	if err != nil {
		t.Fatal(err)
	}
	return tk
}

func TestCreateAssignsIncreasingIDs(t *testing.T) {
	s := NewMemoryStore()
	a, _ := s.Create(newTask(t, "a"))
	b, _ := s.Create(newTask(t, "b"))
	if a.ID != 1 || b.ID != 2 {
		t.Errorf("ids = %d, %d; want 1, 2", a.ID, b.ID)
	}
}

func TestGetNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.Get(42)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestGetReturnsCopy(t *testing.T) {
	s := NewMemoryStore()
	created, _ := s.Create(newTask(t, "orig"))
	got, _ := s.Get(created.ID)
	got.Title = "changed"
	again, _ := s.Get(created.ID)
	if again.Title != "orig" {
		t.Errorf("store state leaked: title = %q", again.Title)
	}
}

func TestListSortedByID(t *testing.T) {
	s := NewMemoryStore()
	for _, n := range []string{"a", "b", "c", "d"} {
		_, _ = s.Create(newTask(t, n))
	}
	list := s.List()
	if len(list) != 4 {
		t.Fatalf("len = %d", len(list))
	}
	for i := 1; i < len(list); i++ {
		if list[i-1].ID >= list[i].ID {
			t.Errorf("not sorted: %v", list)
		}
	}
}

func TestUpdateAndDelete(t *testing.T) {
	s := NewMemoryStore()
	c, _ := s.Create(newTask(t, "x"))
	c.Title = "y"
	if err := s.Update(c); err != nil {
		t.Fatal(err)
	}
	got, _ := s.Get(c.ID)
	if got.Title != "y" {
		t.Errorf("title = %q", got.Title)
	}
	if err := s.Update(task.Task{ID: 99}); !errors.Is(err, ErrNotFound) {
		t.Errorf("update missing: %v", err)
	}
	if err := s.Delete(c.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(c.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("delete twice: %v", err)
	}
}
