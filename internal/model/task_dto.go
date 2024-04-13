package model

import "strconv"

type TaskDto struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksDto struct {
	Tasks []TaskDto `json:"tasks"`
}

type CreateTaskSuccessDto struct {
	ID int `json:"id"`
}

type CreateTaskErrorDto struct {
	Error string `json:"error"`
}

func TaskToTaskDto(task Task) TaskDto {
	return TaskDto{
		ID:      strconv.Itoa(task.ID),
		Date:    task.Date.Format("20060102"),
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}
func TasksToTasksDto(tasks []Task) []TaskDto {
	dto := make([]TaskDto, len(tasks))
	for idx, task := range tasks {
		dto[idx] = TaskToTaskDto(task)
	}
	return dto
}
