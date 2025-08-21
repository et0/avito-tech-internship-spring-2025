package handler

import (
	"log/slog"

	"github.com/et0/avito-tech-internship-spring-2025/internal/middleware"
	"github.com/et0/avito-tech-internship-spring-2025/internal/repository/postgres"
	"github.com/et0/avito-tech-internship-spring-2025/internal/service"
	"github.com/labstack/echo/v4"
)

func New(log *slog.Logger, db *postgres.Postgres, jwtSecret []byte) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logging(log))

	// Service
	userService := service.NewUserService(db, jwtSecret)

	// Handler
	userHandler := NewUserHandler(userService)

	e.POST("/dummyLogin", userHandler.DummyLogin)

	return e
}
