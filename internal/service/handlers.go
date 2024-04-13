package service

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Stern-Ritter/go_task_manager/internal/model"
)

func (s *Server) SignInHandler(res http.ResponseWriter, req *http.Request) {
	authReq := model.AuthRequestDto{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&authReq); err != nil {
		sendAuthError(res, http.StatusBadRequest, err)
		return
	}

	token, err := s.AuthService.SignIn(authReq)
	if err != nil {
		sendAuthError(res, http.StatusUnauthorized, err)
		return
	}

	authSuccesDto := model.AuthSuccessDto{
		Token: token,
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(authSuccesDto); err != nil {
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
		sendTaskError(res, http.StatusBadRequest, err)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	_, err = io.WriteString(res, next)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) AddTaskHandler(res http.ResponseWriter, req *http.Request) {
	task := model.Task{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&task); err != nil {
		sendTaskError(res, http.StatusBadRequest, err)
		return
	}

	id, err := s.TaskService.AddTask(task)
	if err != nil {
		sendTaskError(res, http.StatusBadRequest, err)
		return
	}

	succesDto := model.CreateTaskSuccessDto{
		ID: id,
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(succesDto); err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateTaskHandler(res http.ResponseWriter, req *http.Request) {
	task := model.Task{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&task); err != nil {
		sendTaskError(res, http.StatusBadRequest, err)
		return
	}

	err := s.TaskService.UpdateTask(task)
	if err != nil {
		sendTaskError(res, http.StatusBadRequest, err)
		return
	}

	succesDto := struct{}{}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(succesDto); err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) CompleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	err := s.TaskService.CompleteTask(id)
	if err != nil {
		sendTaskError(res, http.StatusBadRequest, err)
		return
	}

	succesDto := struct{}{}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(succesDto); err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	err := s.TaskService.DeleteTask(id)
	if err != nil {
		sendTaskError(res, http.StatusBadRequest, err)
		return
	}

	succesDto := struct{}{}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(succesDto); err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetTasksHandler(res http.ResponseWriter, req *http.Request) {
	search := req.FormValue("search")

	tasks, err := s.TaskService.GetTasks(search)
	if err != nil {
		sendTaskError(res, http.StatusBadRequest, err)
		return
	}

	tasksDto := model.TasksDto{
		Tasks: model.TasksToTasksDto(tasks),
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(tasksDto); err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	task, err := s.TaskService.GetTask(id)
	if err != nil {
		sendTaskError(res, http.StatusBadRequest, err)
		return
	}

	taskDto := model.TaskToTaskDto(task)

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(res)
	if err := enc.Encode(taskDto); err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func sendTaskError(res http.ResponseWriter, statusCode int, err error) {
	errorDto := model.CreateTaskErrorDto{
		Error: err.Error(),
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(statusCode)
	enc := json.NewEncoder(res)
	if err := enc.Encode(errorDto); err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func sendAuthError(res http.ResponseWriter, statusCode int, err error) {
	errorDto := model.AuthFailedDto{
		Error: err.Error(),
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(statusCode)
	enc := json.NewEncoder(res)
	if err := enc.Encode(errorDto); err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
