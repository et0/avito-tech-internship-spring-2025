package mocks

import (
	"github.com/et0/avito-tech-internship-spring-2025/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateToken(role model.UserRole) (string, error) {
	args := m.Called(role)
	return args.String(0), args.Error(1)
}
