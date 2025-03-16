package memory

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/lowc1012/exercise-api-server/internal/cache"
	"github.com/lowc1012/exercise-api-server/internal/domain"
)

var (
	cachePrefix  = "task_"
	defaultTTL   = 3600 * time.Second
	ErrInvalidID = errors.New("invalid task ID")
)

type autoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	id         int
}

func (a *autoInc) ID() int {
	a.Lock()
	defer a.Unlock()

	a.id++
	return a.id
}

type TaskRepository struct {
	ai    *autoInc
	Cache *cache.Cache[domain.Task]
}

func NewTaskRepository() *TaskRepository {
	taskCache := cache.NewCache[domain.Task]()
	return &TaskRepository{
		ai:    &autoInc{},
		Cache: taskCache,
	}
}

func (t *TaskRepository) FetchAll(ctx context.Context) ([]domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	res := make([]domain.Task, 0)
	for _, task := range t.Cache.FetchAll() {
		res = append(res, task)
	}
	return res, nil
}

func (t *TaskRepository) GetByID(ctx context.Context, id string) (domain.Task, error) {
	if id == "" {
		return domain.Task{}, ErrInvalidID
	}
	key := generateKey(id)
	if task, found := t.Cache.Get(key); found {
		return task, nil
	}
	return domain.Task{}, nil
}

func (t *TaskRepository) Store(ctx context.Context, task domain.Task) error {
	id := strconv.Itoa(t.ai.ID())
	key := generateKey(id)
	task.ID = id
	t.Cache.Set(key, task, defaultTTL)
	return nil
}

func (t *TaskRepository) Update(ctx context.Context, task domain.Task) error {
	if task.ID == "" {
		return ErrInvalidID
	}
	key := generateKey(task.ID)
	if _, found := t.Cache.Get(key); !found {
		return ErrInvalidID
	}
	t.Cache.Set(key, task, defaultTTL)
	return nil
}

func (t *TaskRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidID
	}
	key := generateKey(id)
	t.Cache.Delete(key)
	return nil
}

func generateKey(id string) string {
	// TODO: needs to improve
	if id == "" {
		return ""
	}
	return cachePrefix + id
}
