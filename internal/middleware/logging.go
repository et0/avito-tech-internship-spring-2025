package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func Logging(log *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestBody := getRequestBody(c)

			err := next(c)

			logRequest(c, log, requestBody)

			return err
		}
	}
}

func getRequestBody(c echo.Context) []byte {
	var requestBody []byte
	var originalBody io.ReadCloser

	// Сохраняем оригинальное тело
	originalBody = c.Request().Body

	// Читаем тело
	requestBody, _ = io.ReadAll(io.LimitReader(originalBody, 1024))

	// Восстанавливаем тело для последующих обработчиков
	c.Request().Body = io.NopCloser(bytes.NewBuffer(requestBody))

	// Не забываем закрыть оригинальное тело
	defer originalBody.Close()

	return requestBody
}

func logRequest(c echo.Context, log *slog.Logger, requestBody []byte) {
	status := c.Response().Status

	logArgs := []any{
		"status", status,
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
		"query", c.Request().URL.RawQuery,
		"request_body", string(requestBody),
		"ip", c.RealIP(),
		"duration", time.Since(time.Now()).Milliseconds(),
		"duration_human", time.Since(time.Now()).String(),
	}

	switch {
	case status >= 500:
		log.Error("request completed", logArgs...)
	case status >= 400:
		log.Warn("request completed", logArgs...)
	default:
		log.Info("request completed", logArgs...)
	}
}
