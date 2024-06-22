package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/Stern-Ritter/go_task_manager/internal/model"
)

func (s *Server) SignInHandler(res http.ResponseWriter, req *http.Request) {
	authReq := model.AuthRequestDto{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&authReq); err != nil {
		s.Logger.Error("Error decoding auth request", zap.Error(err))
		sendAuthError(res, http.StatusBadRequest, err.Error())
		return
	}

	token, err := s.AuthService.SignIn(authReq)
	if err != nil {
		s.Logger.Error("Error signing in", zap.Error(err))
		sendAuthError(res, http.StatusUnauthorized, err.Error())
		return
	}

	authSuccesDto := model.AuthSuccessDto{
		Token: token,
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(authSuccesDto); err != nil {
		s.Logger.Error("Error encoding auth response", zap.Error(err))
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetNextDateHandler(res http.ResponseWriter, req *http.Request) {
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	next, err := s.TaskService.GetNextDate(now, date, repeat)
	if err != nil {
		s.Logger.Error("Error getting next date for task", zap.Error(err))
		sendTaskError(res, http.StatusBadRequest, err.Error())
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	_, err = io.WriteString(res, next)
	if err != nil {
		s.Logger.Error("Error writing get next date for task response", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) AddTaskHandler(res http.ResponseWriter, req *http.Request) {
	task := model.Task{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&task); err != nil {
		s.Logger.Error("Error decoding add task", zap.Error(err))
		sendTaskError(res, http.StatusBadRequest, err.Error())
		return
	}

	id, err := s.TaskService.AddTask(task)
	if err != nil {
		s.Logger.Error("Error adding task", zap.Error(err))
		sendTaskError(res, http.StatusInternalServerError, "Internal server error")
		return
	}

	succesDto := model.CreateTaskSuccessDto{
		ID: id,
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(succesDto); err != nil {
		s.Logger.Error("Error encoding add task response", zap.Error(err))
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateTaskHandler(res http.ResponseWriter, req *http.Request) {
	task := model.Task{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&task); err != nil {
		s.Logger.Error("Error decoding update task", zap.Error(err))
		sendTaskError(res, http.StatusBadRequest, err.Error())
		return
	}

	err := s.TaskService.UpdateTask(task)
	if err != nil {
		s.Logger.Error("Error updating task", zap.Error(err))
		sendTaskError(res, http.StatusInternalServerError, "Internal server error")
		return
	}

	succesDto := struct{}{}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(succesDto); err != nil {
		s.Logger.Error("Error encoding update task response", zap.Error(err))
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) CompleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	idNumber, err := strconv.Atoi(id)
	if err != nil {
		s.Logger.Error("Error parsing complete task id", zap.Error(err))
		sendTaskError(res, http.StatusBadRequest, err.Error())
		return
	}

	err = s.TaskService.CompleteTask(idNumber)
	if err != nil {
		s.Logger.Error("Error completing task", zap.Error(err))
		sendTaskError(res, http.StatusBadRequest, err.Error())
		return
	}

	succesDto := struct{}{}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(succesDto); err != nil {
		s.Logger.Error("Error encoding complete task response", zap.Error(err))
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	idNumber, err := strconv.Atoi(id)
	if err != nil {
		s.Logger.Error("Error parsing delete task id", zap.Error(err))
		sendTaskError(res, http.StatusBadRequest, err.Error())
		return
	}

	err = s.TaskService.DeleteTask(idNumber)
	if err != nil {
		s.Logger.Error("Error deleting task", zap.Error(err))
		sendTaskError(res, http.StatusBadRequest, err.Error())
		return
	}

	succesDto := struct{}{}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(succesDto); err != nil {
		s.Logger.Error("Error encoding delete task response", zap.Error(err))
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetTasksHandler(res http.ResponseWriter, req *http.Request) {
	search := req.FormValue("search")

	tasks, err := s.TaskService.GetTasks(search)
	if err != nil {
		s.Logger.Error("Error getting tasks", zap.Error(err))
		sendTaskError(res, http.StatusBadRequest, err.Error())
		return
	}

	tasksDto := model.TasksDto{
		Tasks: model.TasksToTasksDto(tasks),
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(tasksDto); err != nil {
		s.Logger.Error("Error encoding get tasks response", zap.Error(err))
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	idNumber, err := strconv.Atoi(id)
	if err != nil {
		s.Logger.Error("Error parsing get get task id", zap.Error(err))
		sendTaskError(res, http.StatusBadRequest, err.Error())
		return
	}

	task, err := s.TaskService.GetTask(idNumber)
	if err != nil {
		s.Logger.Error("Error getting task", zap.Error(err))
		sendTaskError(res, http.StatusBadRequest, err.Error())
		return
	}

	taskDto := model.TaskToTaskDto(task)

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(taskDto); err != nil {
		s.Logger.Error("Error encoding get task response", zap.Error(err))
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func sendTaskError(res http.ResponseWriter, statusCode int, msg string) {
	errorDto := model.CreateTaskErrorDto{
		Error: msg,
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(statusCode)
	enc := json.NewEncoder(res)
	if err := enc.Encode(errorDto); err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func sendAuthError(res http.ResponseWriter, statusCode int, msg string) {
	errorDto := model.AuthFailedDto{
		Error: msg,
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(statusCode)
	enc := json.NewEncoder(res)
	if err := enc.Encode(errorDto); err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
