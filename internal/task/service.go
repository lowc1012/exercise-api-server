package task

import (
	"context"
	"time"

	"github.com/lowc1012/exercise-api-server/internal/domain"
)

type TaskRepository interface {
	FetchAll(ctx context.Context) ([]domain.Task, error)
	GetByID(ctx context.Context, id string) (domain.Task, error)
	Store(ctx context.Context, t domain.Task) error
	Update(ctx context.Context, t domain.Task) error
	Delete(ctx context.Context, id string) error
}

type Service struct {
	taskRepo TaskRepository
}

func NewService(t TaskRepository) *Service {
	return &Service{
		taskRepo: t,
	}
}

func (s *Service) FetchAll(ctx context.Context) ([]domain.Task, error) {
	tasks, err := s.taskRepo.FetchAll(ctx)
	return tasks, err
}

func (s *Service) GetByID(ctx context.Context, id string) (domain.Task, error) {
	res, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}
	return res, err
}

func (s *Service) Update(ctx context.Context, t domain.Task) error {
	t.UpdatedAt = time.Now()
	return s.taskRepo.Update(ctx, t)
}

func (s *Service) Store(ctx context.Context, task domain.Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	return s.taskRepo.Store(ctx, task)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	_, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.taskRepo.Delete(ctx, id)
}
