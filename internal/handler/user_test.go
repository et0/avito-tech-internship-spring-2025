package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/et0/avito-tech-internship-spring-2025/internal/handler"
	"github.com/et0/avito-tech-internship-spring-2025/internal/logging"
	"github.com/et0/avito-tech-internship-spring-2025/internal/model"
	"github.com/et0/avito-tech-internship-spring-2025/internal/service/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type UserTestCase struct {
	name           string
	requestBody    interface{}
	setupMock      func(MockUserService *mocks.MockUserService)
	expectedStatus int
	expectedBody   interface{}
	expectError    bool
}

func TestRegister_TableDriven(t *testing.T) {
	testCases := []UserTestCase{
		{
			name:           "missing_email",
			requestBody:    map[string]string{"password": "", "role": "moderator"},
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Email is required"},
		},
		{
			name:           "missing_password",
			requestBody:    map[string]string{"email": "test@test.com", "role": "moderator"},
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Password is required"},
		},
		{
			name:           "missing_role",
			requestBody:    map[string]string{"email": "test@test.com", "password": "test"},
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Role is required"},
		},
		{
			name:           "empty_password",
			requestBody:    map[string]string{"email": "test@test.com", "password": "", "role": "moderator"},
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Password is required"},
		},
		{
			name:           "empty_role",
			requestBody:    map[string]string{"email": "test@test.com", "password": "test", "role": ""},
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Role is required"},
		},
		{
			name:           "invalid_role",
			requestBody:    map[string]string{"email": "test@test.com", "password": "test", "role": "admin"},
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Role must be 'employee' or 'moderator'"},
		},
		{
			name:           "invalid_email",
			requestBody:    map[string]string{"email": "test", "password": "test", "role": "moderator"},
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Email must be correct"},
		},
		{
			name:           "invalid_json",
			requestBody:    "invalid_json_string",
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Invalid request format"},
		},
		{
			name:        "database_error",
			requestBody: map[string]string{"email": "test@test.com", "password": "test", "role": "moderator"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				MockUserService.On("Register", "test@test.com", "test", model.RoleModerator).
					Return(nil, fmt.Errorf("DB connect failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Failed to create user"},
		},
		{
			name:        "email_already_exists",
			requestBody: map[string]string{"email": "test@test.com", "password": "test", "role": "moderator"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				MockUserService.On("Register", "test@test.com", "test", model.RoleModerator).
					Return(nil, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "User with this email already exists"},
		},
		{
			name:        "successful_registration",
			requestBody: map[string]string{"email": "test@test.com", "password": "test", "role": "moderator"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				MockUserService.On("Register", "test@test.com", "test", model.RoleModerator).
					Return(&model.User{Email: "test@test.com", Role: model.RoleModerator}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   map[string]string{"email": "test@test.com", "role": string(model.RoleModerator)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			MockUserService := new(mocks.MockUserService)
			tc.setupMock(MockUserService)

			log := logging.New()

			handler := handler.NewUserHandler(MockUserService, log)

			// Создание HTTP запроса
			var reqBody []byte
			if tc.requestBody != nil {
				if bodyStr, ok := tc.requestBody.(string); ok {
					reqBody = []byte(bodyStr)
				} else {
					reqBody, _ = json.Marshal(tc.requestBody)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			// Execution
			err := handler.Register(c)

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

func TestDummyLogin_TableDriven(t *testing.T) {
	// Подготовка тестовых случаев
	testCases := []UserTestCase{
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

			log := logging.New()

			handler := handler.NewUserHandler(MockUserService, log)

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

func TestLogin_TableDriven(t *testing.T) {
	testCases := []UserTestCase{
		{
			name:           "missing_email",
			requestBody:    map[string]string{"password": ""},
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Email is required"},
		},
		{
			name:           "missing_password",
			requestBody:    map[string]string{"email": "test@test.com"},
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Password is required"},
		},
		{
			name:           "empty_password",
			requestBody:    map[string]string{"email": "test@test.com", "password": ""},
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Password is required"},
		},
		{
			name:           "invalid_email",
			requestBody:    map[string]string{"email": "test", "password": "test", "role": "moderator"},
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Email must be correct"},
		},
		{
			name:           "invalid_json",
			requestBody:    "invalid_json_string",
			setupMock:      func(MockUserService *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Invalid request format"},
		},
		{
			name:        "database_error",
			requestBody: map[string]string{"email": "test@test.com", "password": "test"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				MockUserService.On("Login", "test@test.com", "test").
					Return("", fmt.Errorf("DB connect failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Failed login"},
		},
		{
			name:        "email_not_found",
			requestBody: map[string]string{"email": "test@test.com", "password": "test"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				MockUserService.On("Login", "test@test.com", "test").
					Return("", fmt.Errorf("User not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Failed login"},
		},
		{
			name:        "wrong_password",
			requestBody: map[string]string{"email": "test@test.com", "password": "test_wrong"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				MockUserService.On("Login", "test@test.com", "test_wrong").
					Return("", fmt.Errorf("Invalid credentials"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "Failed login"},
		},
		{
			name:        "successful_login",
			requestBody: map[string]string{"email": "test@test.com", "password": "test"},
			setupMock: func(MockUserService *mocks.MockUserService) {
				MockUserService.On("Login", "test@test.com", "test").
					Return("correct_token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"token": "correct_token"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			MockUserService := new(mocks.MockUserService)
			tc.setupMock(MockUserService)

			log := logging.New()

			handler := handler.NewUserHandler(MockUserService, log)

			// Создание HTTP запроса
			var reqBody []byte
			if tc.requestBody != nil {
				if bodyStr, ok := tc.requestBody.(string); ok {
					reqBody = []byte(bodyStr)
				} else {
					reqBody, _ = json.Marshal(tc.requestBody)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			// Execution
			err := handler.Login(c)

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
