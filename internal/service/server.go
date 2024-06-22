package service

import (
	"go.uber.org/zap"

	"github.com/Stern-Ritter/go_task_manager/internal/config"
)

type Server struct {
	AuthService *AuthService
	TaskService *TaskService
	Config      *config.ServerConfig
	Logger      *zap.Logger
}

func NewServer(authService *AuthService, taskService *TaskService, config *config.ServerConfig,
	logger *zap.Logger) *Server {
	return &Server{AuthService: authService, TaskService: taskService, Config: config, Logger: logger}
}
