package handler

import (
	"log/slog"

	"github.com/et0/avito-tech-internship-spring-2025/internal/middleware"
	"github.com/et0/avito-tech-internship-spring-2025/internal/repository"
	"github.com/et0/avito-tech-internship-spring-2025/internal/service"
	"github.com/labstack/echo/v4"
)

func New(log *slog.Logger, db repository.Database, jwtSecret []byte) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logging(log))

	// Service
	userService := service.NewUserService(db, jwtSecret)

	// Handler
	userHandler := NewUserHandler(userService, log)

	e.POST("/dummyLogin", userHandler.DummyLogin)
	e.POST("/register", userHandler.Register)

	return e
}
