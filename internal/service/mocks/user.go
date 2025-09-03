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

func (m *MockUserService) Register(email string, password string, role model.UserRole) (*model.User, error) {
	args := m.Called(email, password, role)
	if user := args.Get(0); user != nil {
		return user.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}
