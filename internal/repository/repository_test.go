package repository_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/0xc0d/TodoService/internal/repository"
	"github.com/0xc0d/TodoService/internal/service"

	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
)

var ctx = context.Background()

type Suite struct {
	suite.Suite
	repo *repository.TaskRepository
}

func (s *Suite) Test01CreateFirstTask() {
	task := &service.Task{
		Title:    "Fist Task",
		Comments: "Inspiration is a guest that does not willingly visit the lazy.",
		DueDate:  time.Date(1840, time.May, 6, 0, 0, 0, 0, time.UTC),
		Done:     false,
	}

	createdTask, err := s.repo.Create(ctx, task)
	s.Assert().NoError(err)
	s.Assert().Equal(service.TaskID(0), task.ID)
	s.Assert().Equal(task, createdTask)
}

func (s *Suite) Test02CreateNextTask() {
	task := &service.Task{
		Title:    "Learn how to play an instrument",
		Comments: "Truly there wuold be reason to go mad were it not for music",
		Labels:   []string{"high_priority"},
		DueDate:  time.Date(1893, time.November, 6, 0, 0, 0, 0, time.UTC),
		Done:     false,
	}

	createdTask, err := s.repo.Create(ctx, task)
	s.Assert().NoError(err)
	s.Assert().Equal(service.TaskID(1), task.ID)
	s.Assert().Equal(task, createdTask)
}

func (s *Suite) Test03UpdateTask() {
	task := &service.Task{
		ID:       1,
		Title:    "Learn how to play an instrument",
		Comments: "Truly there would be reason to go mad were it not for music",
		Labels:   []string{"high_priority"},
		DueDate:  time.Date(1893, time.November, 6, 0, 0, 0, 0, time.UTC),
		Done:     false,
	}

	err := s.repo.Update(ctx, task)
	s.Assert().NoError(err)

	dummyTask := &service.Task{
		ID: 10000,
	}
	err = s.repo.Update(ctx, dummyTask)
	s.Assert().ErrorIs(err, repository.ErrNotFound)

	err = s.repo.Update(ctx, nil)
	s.Assert().ErrorIs(err, repository.ErrInvalidTask)
}

func (s *Suite) Test03DeleteTask() {
	s.repo.Delete(ctx, service.TaskID(0))
	_, err := s.repo.Get(ctx, service.TaskID(0))
	s.Assert().ErrorIs(err, repository.ErrNotFound)
}

func (s *Suite) Test04GetAllTasks() {
	tasks := s.repo.GetAll(ctx)
	s.Assert().Len(tasks, 1)
	s.Assert().Equal(service.TaskID(1), tasks[0].ID)
}

func (s *Suite) Test05CreateUpdateRace() {
	typoTask := &service.Task{
		Title:    "Do your things",
		Comments: "I pay no attention whatever to anybody's priase or blame.",
		Labels:   []string{"lessons"},
		DueDate:  time.Date(1791, time.December, 5, 0, 0, 0, 0, time.UTC),
		Done:     false,
	}

	task := &service.Task{
		ID:       2,
		Title:    "Do your things",
		Comments: "I pay no attention whatever to anybody's praise or blame.",
		Labels:   []string{"lessons"},
		DueDate:  time.Date(1791, time.December, 5, 0, 0, 0, 0, time.UTC),
		Done:     false,
	}

	c := make(chan struct{})
	go func() {
		_ = s.repo.Update(ctx, task)
		c <- struct{}{}
	}()
	_, _ = s.repo.Create(ctx, typoTask)
	<-c

	resultTask, err := s.repo.Get(ctx, service.TaskID(2))
	s.Assert().NoError(err)
	s.Assert().Equal(task.Comments, resultTask.Comments)
}

func (s *Suite) Test06CreateConcurrent() {
	task := service.Task{
		Title:    "Truth to be told",
		Comments: "No one feels another's grief, no one understands another's joy.",
		Labels:   []string{"fact"},
		DueDate:  time.Date(1828, time.November, 19, 0, 0, 0, 0, time.UTC),
		Done:     false,
	}

	for i := 0; i < 1000; i++ {
		go func(i int) {
			task := task
			task.Title += "#" + strconv.Itoa(i)
			_, _ = s.repo.Create(ctx, &task)
		}(i)
	}
}

func (s *Suite) Test07CreateInvalidTask() {
	_, err := s.repo.Create(ctx, nil)
	s.Assert().ErrorIs(err, repository.ErrInvalidTask)
}

func TestSuite(t *testing.T) {
	defer goleak.VerifyNone(t)

	testSuite := new(Suite)
	testSuite.repo = repository.New()
	suite.Run(t, testSuite)
}
