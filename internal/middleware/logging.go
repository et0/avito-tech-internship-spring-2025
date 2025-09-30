package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"time"

	"github.com/et0/avito-tech-internship-spring-2025/pkg/errors"
	"github.com/labstack/echo/v4"
)

func Logging(log *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestBody := getRequestBody(c)

			err := next(c)

			logRequest(c, log, requestBody, err)

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

func logRequest(c echo.Context, log *slog.Logger, requestBody []byte, err error) {
	status := c.Response().Status

	// Если есть ошибка, а код ответа по прежнему 200
	if err != nil && status == 200 {
		switch e := err.(type) {
		case *errors.AppError:
			status = e.Code
		case *echo.HTTPError:
			status = e.Code
		default:
			status = 500
		}
	}

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

	if err != nil {
		logArgs = append(logArgs, "error", err.Error())
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
