package main

import (
	"log/slog"
	"os"

	"github.com/et0/avito-tech-internship-spring-2025/internal/config"
	"github.com/et0/avito-tech-internship-spring-2025/internal/handler"
	"github.com/et0/avito-tech-internship-spring-2025/internal/logging"
	"github.com/et0/avito-tech-internship-spring-2025/internal/repository/postgres"
	"github.com/labstack/echo/v4"
)

type App struct {
	Cfg    *config.Config
	Logger *slog.Logger
	Echo   *echo.Echo
	DB     *postgres.Postgres
}

func main() {
	// Logger
	log := logging.New()

	// временное решение
	os.Setenv("CONFIG_PATH", "./config/local.yaml")

	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Error("failed config load", "error", err)
		return
	}

	// DB Postgres
	pg, err := postgres.New(&cfg.DB)
	if err != nil {
		log.Error("failed DB create ", "error", err)
		return
	}
	defer pg.Close()

	// Echo Handler
	e := handler.New(log, pg, []byte(cfg.HTTP.JWTSecret))
	if err := e.Start(":" + cfg.HTTP.Port); err != nil {
		log.Error("failed server start ", "error", err)
	}
}
