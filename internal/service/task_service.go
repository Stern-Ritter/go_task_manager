package service

import (
	"strconv"
	"strings"
	"time"

	"github.com/Stern-Ritter/go_task_manager/internal/errors"
	"github.com/Stern-Ritter/go_task_manager/internal/model"
	"github.com/Stern-Ritter/go_task_manager/internal/storage"
	"github.com/Stern-Ritter/go_task_manager/internal/utils"
	"go.uber.org/zap"
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

func (s TaskService) CompleteTask(id string) error {
	idNumber, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	t, err := s.store.GetByID(idNumber)
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(t.Repeat)) != 0 {
		return s.store.Complete(t)
	} else {
		return s.store.Delete(t.ID)
	}
}

func (s TaskService) DeleteTask(id string) error {
	idNumber, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.store.Delete(idNumber)
}

func (s TaskService) GetTasks(search string) ([]model.Task, error) {
	var tasks []model.Task
	isSearch := len(strings.TrimSpace(search)) > 0
	isValidSerachDate, err := utils.ValidateSearchDate(search)
	if err != nil {
		return tasks, err
	}

	var storeErr error
	switch {
	case isSearch && isValidSerachDate:
		date, _ := time.Parse("02.01.2006", search)
		tasks, storeErr = s.store.GetAllByDate(date.Format("20060102"))
	case isSearch:
		tasks, storeErr = s.store.GetAllByTitleOrComment(strings.Join([]string{"%", search, "%"}, ""))
	default:
		tasks, storeErr = s.store.GetAll()
	}

	return tasks, storeErr
}

func (s TaskService) GetTask(id string) (model.Task, error) {
	idNumber, err := strconv.Atoi(id)
	if err != nil {
		return model.Task{}, err
	}
	return s.store.GetByID(idNumber)
}
