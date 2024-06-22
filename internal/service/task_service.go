package service

import (
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/Stern-Ritter/go_task_manager/internal/errors"
	"github.com/Stern-Ritter/go_task_manager/internal/model"
	"github.com/Stern-Ritter/go_task_manager/internal/storage"
	"github.com/Stern-Ritter/go_task_manager/internal/utils"
)

type TaskService struct {
	store  storage.TaskStore
	logger *zap.Logger
}

func NewTaskService(store storage.TaskStore, logger *zap.Logger) *TaskService {
	return &TaskService{store: store, logger: logger}
}

func (s TaskService) GetNextDate(now string, date string, repeat string) (string, error) {
	parsedNow, err := time.Parse("20060102", now)
	if err != nil {
		return "", errors.NewInvalidDateFormat("invalid task now format", err)
	}
	return utils.NextDate(parsedNow, date, repeat)
}

func (s TaskService) AddTask(t model.Task) (int, error) {
	return s.store.Create(t)
}

func (s TaskService) UpdateTask(t model.Task) error {
	return s.store.Update(t)
}

func (s TaskService) CompleteTask(id int) error {
	t, err := s.store.GetByID(id)
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(t.Repeat)) != 0 {
		return s.store.Complete(t)
	} else {
		return s.store.Delete(t.ID)
	}
}

func (s TaskService) DeleteTask(id int) error {
	return s.store.Delete(id)
}

func (s TaskService) GetTasks(search string) ([]model.Task, error) {
	var tasks []model.Task
	isSearch := len(strings.TrimSpace(search)) > 0
	isValidSearchDate, err := utils.ValidateSearchDate(search)
	if err != nil {
		return tasks, err
	}

	var storeErr error
	switch {
	case isSearch && isValidSearchDate:
		date, _ := time.Parse("02.01.2006", search)
		tasks, storeErr = s.store.GetAllByDate(date.Format("20060102"))
	case isSearch:
		tasks, storeErr = s.store.GetAllByTitleOrComment(strings.Join([]string{"%", search, "%"}, ""))
	default:
		tasks, storeErr = s.store.GetAll()
	}

	return tasks, storeErr
}

func (s TaskService) GetTask(id int) (model.Task, error) {
	return s.store.GetByID(id)
}
