package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/et0/avito-tech-internship-spring-2025/internal/handler"
	"github.com/et0/avito-tech-internship-spring-2025/internal/model"
	"github.com/et0/avito-tech-internship-spring-2025/internal/service/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type dummyLoginTestCase struct {
	name           string
	requestBody    interface{}
	setupMock      func(MockUserService *mocks.MockUserService)
	expectedStatus int
	expectedBody   interface{}
	expectError    bool
}

func TestDummyLogin_TableDriven(t *testing.T) {
	// Подготовка тестовых случаев
	testCases := []dummyLoginTestCase{
		{
			name:        "success_employee_role",
			requestBody: map[string]string{"role": "employee"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				MockUserService.On("CreateToken", model.RoleEmployee).
					Return("employee_token_123", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"token": "employee_token_123"},
		},
		{
			name:        "success_moderator_role",
			requestBody: map[string]string{"role": "moderator"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				MockUserService.On("CreateToken", model.RoleModerator).
					Return("moderator_token_456", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"token": "moderator_token_456"},
		},
		{
			name:        "invalid_role",
			requestBody: map[string]string{"role": "admin"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				// Мок не должен вызываться для невалидной роли
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Role must be 'employee' or 'moderator'"},
		},
		{
			name:        "empty_role",
			requestBody: map[string]string{"role": ""},
			setupMock: func(MockUserService *mocks.MockUserService) {
				// Мок не должен вызываться для пустой роли
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Role is required"},
		},
		{
			name:        "missing_role_field",
			requestBody: map[string]string{"wrong_field": "employee"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				// Мок не должен вызываться при отсутствии поля role
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Role is required"},
		},
		{
			name:        "invalid_json",
			requestBody: "invalid_json_string",
			setupMock: func(MockUserService *mocks.MockUserService) {
				// Мок не должен вызываться при невалидном JSON
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Invalid request format"},
		},
		{
			name:        "service_error",
			requestBody: map[string]string{"role": "employee"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				MockUserService.On("CreateToken", model.RoleEmployee).
					Return("", fmt.Errorf("failed to generate token"))
			},
			// TODO: change status from StatusBadRequest to StatusInternalServerError
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "failed to generate token"},
		},
		{
			name:        "empty_body",
			requestBody: nil,
			setupMock: func(MockUserService *mocks.MockUserService) {
				// Мок не должен вызываться при пустом теле
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Role is required"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			MockUserService := new(mocks.MockUserService)
			tc.setupMock(MockUserService)

			handler := handler.NewUserHandler(MockUserService)

			// Создание HTTP запроса
			var reqBody []byte
			if tc.requestBody != nil {
				if bodyStr, ok := tc.requestBody.(string); ok {
					reqBody = []byte(bodyStr)
				} else {
					reqBody, _ = json.Marshal(tc.requestBody)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			// Execution
			err := handler.DummyLogin(c)

			// Assertion
			if tc.expectError {
				assert.Error(t, err)
				if httpErr, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tc.expectedStatus, httpErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedStatus, rec.Code)

				// Проверка тела ответа
				var actualResponse map[string]interface{}
				if len(rec.Body.Bytes()) > 0 {
					err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
					assert.NoError(t, err)
				}

				// Для ожидаемого тела в формате map
				if expectedMap, ok := tc.expectedBody.(map[string]string); ok {
					for key, expectedValue := range expectedMap {
						if actualValue, exists := actualResponse[key]; exists {
							assert.Equal(t, expectedValue, actualValue)
						} else {
							assert.Fail(t, "Expected key not found in response: "+key)
						}
					}
				}
			}

			// Verify mock expectations
			MockUserService.AssertExpectations(t)
		})
	}
}
