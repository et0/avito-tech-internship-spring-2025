package handler

import (
	"github.com/et0/avito-tech-internship-spring-2025/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(sUS *service.UserService) *UserHandler {
	return &UserHandler{
		service: sUS,
	}
}
