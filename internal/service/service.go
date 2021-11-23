package service

import (
	"context"
	"errors"
	"time"

	"todoService/internal/pb"

	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const serviceName = "api"

var (
	ErrEmptyTitle = errors.New("title is required")
	ErrEmptyTask  = errors.New("valid task required but received empty")
	ErrEmptyTime  = errors.New("valid due date required but received empty")
)

type TaskID uint64

type Task struct {
	ID       TaskID    `json:"id"`
	Title    string    `json:"title"`
	Comments string    `json:"comments"`
	Labels   []string  `json:"labels"`
	DueDate  time.Time `json:"due_date"`
	Done     bool      `json:"done"`
}

type Repository interface {
	Create(context.Context, *Task) (*Task, error)
	Get(context.Context, TaskID) (*Task, error)
	GetAll(context.Context) []*Task
	Update(context.Context, *Task) error
	Delete(context.Context, TaskID)
}

type Service struct {
	repo Repository
}

// New returns a new TODO service with the given repository as storage.
func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateTask creates a new task and returns the task with the associated ID. A task
// must and Title and DueDate.
func (s *Service) CreateTask(ctx context.Context, task *pb.CreateTaskRequest) (*pb.Task, error) {
	tr := otel.Tracer(serviceName)
	ctx, span := tr.Start(ctx, "CreateTask")
	defer span.End()

	if task == nil {
		return nil, ErrEmptyTask
	}

	if task.Title == "" {
		return nil, ErrEmptyTitle
	}

	if task.DueDate == nil {
		return nil, ErrEmptyTime
	}

	t := time.Unix(task.DueDate.Seconds, 0)

	createdTask, err := s.repo.Create(ctx, &Task{
		Title:    task.Title,
		Comments: task.Comments,
		Labels:   task.Labels,
		DueDate:  t,
		Done:     task.Done,
	})
	if err != nil {
		return nil, err
	}

	resultTask := &pb.Task{
		Id:       uint64(createdTask.ID),
		Title:    task.Title,
		Comments: task.Comments,
		Labels:   task.Labels,
		DueDate:  task.DueDate,
		Done:     task.Done,
	}

	return resultTask, nil
}

// GetAllTasks returns all tasks with no specified order.
func (s *Service) GetAllTasks(ctx context.Context, _ *pb.Empty) (*pb.Tasks, error) {
	tr := otel.Tracer(serviceName)
	ctx, span := tr.Start(ctx, "GetAllTasks")
	defer span.End()

	tasks := s.repo.GetAll(ctx)

	allTasks := make([]*pb.Task, 0, len(tasks))
	for _, task := range tasks {
		allTasks = append(allTasks, &pb.Task{
			Id:       uint64(task.ID),
			Title:    task.Title,
			Comments: task.Comments,
			Labels:   task.Labels,
			DueDate:  timestamppb.New(task.DueDate),
			Done:     task.Done,
		})
	}

	return &pb.Tasks{Tasks: allTasks}, nil
}

// GetTask returns the task associated to the given taskID. It returns repository.ErrNotFound if
// there is no task associated to the taskID.
func (s *Service) GetTask(ctx context.Context, in *pb.GetTaskRequest) (*pb.Task, error) {
	tr := otel.Tracer(serviceName)
	ctx, span := tr.Start(ctx, "GetTask")
	defer span.End()

	task, err := s.repo.Get(ctx, TaskID(in.TaskID))
	if err != nil {
		return nil, err
	}

	return &pb.Task{
		Id:       uint64(task.ID),
		Title:    task.Title,
		Comments: task.Comments,
		Labels:   task.Labels,
		DueDate:  timestamppb.New(task.DueDate),
		Done:     task.Done,
	}, nil
}

// UpdateTask updates a Task associated to the taskID. A task must and Title and DueDate.
func (s *Service) UpdateTask(ctx context.Context, in *pb.UpdateTaskRequest) (*pb.Task, error) {
	tr := otel.Tracer(serviceName)
	ctx, span := tr.Start(ctx, "UpdateTask")
	defer span.End()

	if in == nil || in.Task == nil {
		return nil, ErrEmptyTask
	}

	if in.Task.Title == "" {
		return nil, ErrEmptyTitle
	}

	if in.Task.DueDate == nil {
		return nil, ErrEmptyTime
	}

	t := time.Unix(in.Task.DueDate.Seconds, 0)

	err := s.repo.Update(ctx, &Task{
		ID:       TaskID(in.Task.Id),
		Title:    in.Task.Title,
		Comments: in.Task.Comments,
		Labels:   in.Task.Labels,
		DueDate:  t,
		Done:     in.Task.Done,
	})
	if err != nil {
		return nil, err
	}

	return in.Task, nil
}

// DeleteTask deletes task associated with the given id. it does not return any error
// even if the task does not exist.
func (s *Service) DeleteTask(ctx context.Context, in *pb.DeleteTaskRequest) (*pb.Empty, error) {
	tr := otel.Tracer(serviceName)
	ctx, span := tr.Start(ctx, "DeleteTask")
	defer span.End()

	s.repo.Delete(ctx, TaskID(in.TaskID))
	return &pb.Empty{}, nil
}
