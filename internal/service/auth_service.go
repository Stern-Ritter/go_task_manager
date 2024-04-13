package service

import (
	"github.com/Stern-Ritter/go_task_manager/internal/errors"
	"github.com/Stern-Ritter/go_task_manager/internal/utils"

	"github.com/Stern-Ritter/go_task_manager/internal/model"
	"go.uber.org/zap"
)

type AuthService struct {
	rootPassword string
	logger       *zap.Logger
}

func NewAuthService(rootPassword string, logger *zap.Logger) *AuthService {
	return &AuthService{rootPassword: rootPassword, logger: logger}
}

func (s AuthService) SignIn(authReq model.AuthRequestDto) (string, error) {
	if s.rootPassword != authReq.Password {
		return "", errors.NewAuthenticationError("Invalid root user password", nil)
	}

	return utils.HashPassword(authReq.Password)
}
