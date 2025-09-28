package handler

import "log/slog"

type PvzHandler struct {
	log *slog.Logger
}

func NewPvzHandler(log *slog.Logger) *PvzHandler {
	return &PvzHandler{
		log: log,
	}
}
