package handler

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

type PvzHandler struct {
	log *slog.Logger
}

func NewPvzHandler(log *slog.Logger) *PvzHandler {
	return &PvzHandler{
		log: log,
	}
}

func (ph *PvzHandler) Create(ctx echo.Context) error {
	return nil
}
