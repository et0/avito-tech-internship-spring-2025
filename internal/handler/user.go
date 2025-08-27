package handler

import (
	"net/http"

	"github.com/et0/avito-tech-internship-spring-2025/api/gen/openapi"
	"github.com/et0/avito-tech-internship-spring-2025/internal/model"
	"github.com/et0/avito-tech-internship-spring-2025/internal/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service service.UserService
}

type UserDummyLoginResponse struct {
	Token openapi.Token `json:"token"`
}

func NewUserHandler(sUS service.UserService) *UserHandler {
	return &UserHandler{
		service: sUS,
	}
}

func (uc *UserHandler) DummyLogin(ctx echo.Context) error {
	var request openapi.PostDummyLoginJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: "Invalid request format"})
	}

	if request.Role != openapi.PostDummyLoginJSONBodyRoleEmployee && request.Role != openapi.PostDummyLoginJSONBodyRoleModerator {
		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: "Role must be 'employee' or 'moderator'"})
	}

	token, err := uc.service.CreateToken(model.UserRole(request.Role))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, UserDummyLoginResponse{token})
}
