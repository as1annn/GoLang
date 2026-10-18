package storage

import (
	"errors"
	"fmt"
	"sort"

	"sis/internal/task"
)

var ErrNotFound = errors.New("task not found")

type Repository interface {
	Create(t task.Task) (task.Task, error)
	Get(id int) (task.Task, error)
	List() []task.Task
	Update(t task.Task) error
	Delete(id int) error
}

type MemoryStore struct {
	tasks  map[int]task.Task
	nextID int
}

var _ Repository = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{tasks: make(map[int]task.Task), nextID: 1}
}

func (s *MemoryStore) Create(t task.Task) (task.Task, error) {
	t.ID = s.nextID
	s.nextID++
	s.tasks[t.ID] = t
	return t, nil
}

func (s *MemoryStore) Get(id int) (task.Task, error) {
	t, ok := s.tasks[id]
	if !ok {
		return task.Task{}, fmt.Errorf("get task %d: %w", id, ErrNotFound)
	}
	return t, nil
}

func (s *MemoryStore) List() []task.Task {
	result := make([]task.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func (s *MemoryStore) Update(t task.Task) error {
	if _, ok := s.tasks[t.ID]; !ok {
		return fmt.Errorf("update task %d: %w", t.ID, ErrNotFound)
	}
	s.tasks[t.ID] = t
	return nil
}

func (s *MemoryStore) Delete(id int) error {
	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("delete task %d: %w", id, ErrNotFound)
	}
	delete(s.tasks, id)
	return nil
}
