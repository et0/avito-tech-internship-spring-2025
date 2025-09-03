package handler

import (
	"log/slog"
	"net/http"

	"github.com/et0/avito-tech-internship-spring-2025/api/gen/openapi"
	"github.com/et0/avito-tech-internship-spring-2025/internal/model"
	"github.com/et0/avito-tech-internship-spring-2025/internal/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service service.UserService
	log     *slog.Logger
}

type UserDummyLoginResponse struct {
	Token openapi.Token `json:"token"`
}

type UserRegisterResponse struct {
	Email string           `json:"email"`
	Role  openapi.UserRole `json:"role"`
}

func NewUserHandler(sUS service.UserService, log *slog.Logger) *UserHandler {
	return &UserHandler{
		service: sUS,
		log:     log,
	}
}

func (uc *UserHandler) DummyLogin(ctx echo.Context) error {
	var request openapi.PostDummyLoginJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: "Invalid request format"})
	}

	if request.Role == "" {
		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: "Role is required"})
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

func (u *UserHandler) Register(ctx echo.Context) error {
	var request openapi.PostRegisterJSONRequestBody

	if err := ctx.Bind(&request); err != nil {
		u.log.Error("handler.user.register", "error", err)

		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: "Invalid request format"})
	}

	if request.Email == "" {
		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: "Email is required"})
	}

	if request.Password == "" {
		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: "Password is required"})
	}

	if request.Role == "" {
		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: "Role is required"})
	}

	if request.Role != openapi.Employee && request.Role != openapi.Moderator {
		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: "Role must be 'employee' or 'moderator'"})
	}

	user, err := u.service.Register(string(request.Email), request.Password, model.UserRole(request.Role))
	if err != nil {
		u.log.Error("register", "error", err)

		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: "Failed to create user"})
	}
	if user == nil {
		return ctx.JSON(http.StatusBadRequest, openapi.Error{Message: "User with this email already exists"})
	}

	return ctx.JSON(http.StatusCreated, UserRegisterResponse{Email: user.Email, Role: openapi.UserRole(user.Role)})
}
