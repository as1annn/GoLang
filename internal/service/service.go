package service

import (
	"sis/internal/storage"
	"sis/internal/task"
)

type Service struct {
	repo storage.Repository
}

func New(repo storage.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(title, payload string) (task.Task, error) {
	t, err := task.New(title, payload)
	if err != nil {
		return task.Task{}, err
	}
	return s.repo.Create(t)
}

func (s *Service) Get(id int) (task.Task, error) { return s.repo.Get(id) }
func (s *Service) List() []task.Task             { return s.repo.List() }
func (s *Service) Delete(id int) error           { return s.repo.Delete(id) }

func (s *Service) Start(id int) (task.Task, error) {
	return s.apply(id, func(t *task.Task) error { return t.Start() })
}

func (s *Service) Complete(id int, result string) (task.Task, error) {
	return s.apply(id, func(t *task.Task) error { return t.Complete(result) })
}

func (s *Service) Fail(id int, reason string) (task.Task, error) {
	return s.apply(id, func(t *task.Task) error { return t.Fail(reason) })
}

func (s *Service) Cancel(id int) (task.Task, error) {
	return s.apply(id, func(t *task.Task) error { return t.Cancel() })
}

func (s *Service) apply(id int, mutate func(*task.Task) error) (task.Task, error) {
	t, err := s.repo.Get(id)
	if err != nil {
		return task.Task{}, err
	}
	if err := mutate(&t); err != nil {
		return task.Task{}, err
	}
	if err := s.repo.Update(t); err != nil {
		return task.Task{}, err
	}
	return t, nil
}
