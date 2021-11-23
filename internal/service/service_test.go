package service_test

import (
	"context"
	"go.uber.org/goleak"
	"testing"
	"time"
	"todoService/internal/repository"
	"todoService/internal/service"

	"todoService/internal/pb"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ctx = context.TODO()

type Suite struct {
	suite.Suite
	service *service.Service
}

func (s *Suite) Test01CreateTask() {
	task := &pb.CreateTaskRequest{
		Title:    "Last task before release",
		Comments: "Sorrowful and great is the artist's destiny.",
		DueDate:  timestamppb.New(time.Date(1827, time.March, 26, 0, 0, 0, 0, time.UTC)),
	}

	res, err := s.service.CreateTask(ctx, task)
	s.Assert().NoError(err)
	s.Assert().Equal(task.Title, res.Title)
	s.Assert().Equal(task.Comments, res.Comments)
	s.Assert().Equal(task.DueDate, res.DueDate)
}

func (s *Suite) Test02CreateDummyTask() {
	res, err := s.service.CreateTask(ctx, nil)
	s.Assert().ErrorIs(err, service.ErrEmptyTask)
	s.Assert().Nil(res)
}

func (s *Suite) Test03CreateTaskNoTitle() {
	task := &pb.CreateTaskRequest{
		Comments: "There is no such thing as happy music.",
		DueDate:  timestamppb.Now(),
	}

	res, err := s.service.CreateTask(ctx, task)
	s.Assert().Nil(res)
	s.Assert().ErrorIs(err, service.ErrEmptyTitle)
}

func (s *Suite) Test04CreateTaskNoTime() {
	task := &pb.CreateTaskRequest{
		Title:    "On-going task for life",
		Comments: "To regret the past, to hope in the future, and never to be satisfied with the present",
	}

	res, err := s.service.CreateTask(ctx, task)
	s.Assert().Nil(res)
	s.Assert().ErrorIs(err, service.ErrEmptyTime)
}

func (s *Suite) Test05UpdateDummyTask() {
	res, err := s.service.UpdateTask(ctx, &pb.UpdateTaskRequest{Task: nil})
	s.Assert().ErrorIs(err, service.ErrEmptyTask)
	s.Assert().Nil(res)
}

func (s *Suite) Test06UpdateTaskNoTitle() {
	task := &pb.Task{
		Id:       0,
		Comments: "Inspiration is a guest that does not willingly visit the lazy.",
		DueDate:  timestamppb.Now(),
	}

	res, err := s.service.UpdateTask(
		ctx,
		&pb.UpdateTaskRequest{Task: task},
	)
	s.Assert().Nil(res)
	s.Assert().ErrorIs(err, service.ErrEmptyTitle)
}

func (s *Suite) Test07UpdateTaskNoTime() {
	task := &pb.Task{
		Id:       0,
		Title:    "On-going task for life",
		Comments: "I shall seize fate by the throat.",
	}

	res, err := s.service.UpdateTask(
		ctx,
		&pb.UpdateTaskRequest{Task: task},
	)
	s.Assert().Nil(res)
	s.Assert().ErrorIs(err, service.ErrEmptyTime)
}

func (s *Suite) Test08GetAllTask() {
	res, err := s.service.GetAllTasks(ctx, &pb.Empty{})
	s.Assert().NoError(err)
	s.Assert().Len(res.Tasks, 1)

	theTask := res.Tasks[0]
	s.Assert().EqualValues(service.TaskID(0), theTask.Id)
	s.Assert().Equal("Last task before release", theTask.Title)
}

func (s *Suite) Test09GetTask() {
	expectedTask := &pb.Task{
		Title:    "Last task before release",
		Comments: "Sorrowful and great is the artist's destiny.",
		DueDate:  timestamppb.New(time.Date(1827, time.March, 26, 0, 0, 0, 0, time.UTC)),
	}

	task, err := s.service.GetTask(ctx, &pb.GetTaskRequest{TaskID: 0})
	s.Assert().NoError(err)
	s.Assert().Equal(expectedTask, task)

}

func (s *Suite) Test10UpdateTask() {
	task := &pb.Task{
		Id:       0,
		Title:    "Last task before release",
		Comments: "Sorrowful and great is the artist's destiny.",
		DueDate:  timestamppb.Now(),
		Done:     true,
	}

	res, err := s.service.UpdateTask(
		ctx,
		&pb.UpdateTaskRequest{Task: task},
	)
	s.Assert().NoError(err)
	s.Assert().NotNil(res)
	s.Assert().Equal(res.Id, task.Id)
	s.Assert().Equal(res.Done, task.Done)
}

func (s *Suite) Test11UpdateNotExistTask() {
	task := &pb.Task{
		Id:       1000,
		Title:    "The truth",
		Comments: "To play a wrong note is insignificant.",
		DueDate:  timestamppb.Now(),
		Done:     true,
	}

	res, err := s.service.UpdateTask(
		ctx,
		&pb.UpdateTaskRequest{Task: task},
	)
	s.Assert().Nil(res)
	s.Assert().ErrorIs(err, repository.ErrNotFound)
}

func (s *Suite) Test12DeleteTask() {
	res, err := s.service.DeleteTask(ctx, &pb.DeleteTaskRequest{TaskID: 0})
	s.Assert().NoError(err)
	s.Assert().Equal(&pb.Empty{}, res)

	task, err := s.service.GetTask(ctx, &pb.GetTaskRequest{TaskID: 0})
	s.Assert().Nil(task)
	s.Assert().ErrorIs(err, repository.ErrNotFound)
}

func TestSuite(t *testing.T) {
	defer goleak.VerifyNone(t)

	testSuite := new(Suite)
	testSuite.service = service.New(repository.New())
	suite.Run(t, testSuite)
}
