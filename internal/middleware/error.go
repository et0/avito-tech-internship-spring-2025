package middleware

import (
	"log/slog"
	"net/http"

	"github.com/et0/avito-tech-internship-spring-2025/api/gen/openapi"
	"github.com/et0/avito-tech-internship-spring-2025/pkg/errors"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(log *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {

		// Обрабатываем разные типы ошибок
		switch e := err.(type) {
		case *errors.AppError:
			c.JSON(e.Code, openapi.Error{
				Message: e.Message,
			})

		case *echo.HTTPError:
			c.JSON(e.Code, openapi.Error{
				Message: e.Message.(string),
			})

		default:
			c.JSON(http.StatusInternalServerError, openapi.Error{
				Message: "Internal server error",
			})
		}
	}
}
