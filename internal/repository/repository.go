package repository

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"sync"

	"todoService/internal/service"
)

const serviceName = "repository"

var (
	ErrNotFound    = errors.New("task does not exist")
	ErrInvalidTask = errors.New("valid task expected but nil received")
)

type TaskRepository struct {
	tasks    map[service.TaskID]*service.Task
	mut      sync.RWMutex
	idCursor service.TaskID
}

func New() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[service.TaskID]*service.Task),
	}
}

// Create adds a task to the repository and returns the task with the given id.
func (s *TaskRepository) Create(ctx context.Context, task *service.Task) (*service.Task, error) {
	tr := otel.Tracer(serviceName)
	_, span := tr.Start(ctx, "Create")
	defer span.End()

	s.mut.Lock()
	defer s.mut.Unlock()

	if task == nil {
		return nil, ErrInvalidTask
	}

	task.ID = s.idCursor
	s.tasks[s.idCursor] = task
	s.idCursor++

	return task, nil
}

func (s *TaskRepository) Get(ctx context.Context, taskId service.TaskID) (*service.Task, error) {
	tr := otel.Tracer(serviceName)
	_, span := tr.Start(ctx, "Get")
	defer span.End()

	s.mut.RLock()
	defer s.mut.RUnlock()

	if task, ok := s.tasks[taskId]; ok {
		return task, nil

	}
	return nil, ErrNotFound
}

func (s *TaskRepository) GetAll(ctx context.Context) []*service.Task {
	tr := otel.Tracer(serviceName)
	_, span := tr.Start(ctx, "GetAll")
	defer span.End()

	s.mut.RLock()
	defer s.mut.RUnlock()

	all := make([]*service.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		all = append(all, task)
	}

	return all
}

func (s *TaskRepository) Update(ctx context.Context, task *service.Task) error {
	tr := otel.Tracer(serviceName)
	_, span := tr.Start(ctx, "Update")
	defer span.End()

	s.mut.Lock()
	defer s.mut.Unlock()

	if task == nil {
		return ErrInvalidTask
	}

	if _, ok := s.tasks[task.ID]; !ok {
		return ErrNotFound
	}

	s.tasks[task.ID] = task

	return nil
}

func (s *TaskRepository) Delete(ctx context.Context, taskId service.TaskID) {
	tr := otel.Tracer(serviceName)
	_, span := tr.Start(ctx, "Delete")
	defer span.End()

	s.mut.Lock()
	defer s.mut.Unlock()

	delete(s.tasks, taskId)
}
