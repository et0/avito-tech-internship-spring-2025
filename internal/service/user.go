package service

import (
	"fmt"
	"time"

	"github.com/et0/avito-tech-internship-spring-2025/internal/model"
	"github.com/et0/avito-tech-internship-spring-2025/internal/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	CreateToken(role model.UserRole) (string, error)
}

type userService struct {
	db        *postgres.Postgres
	jwtSecret []byte
}

func NewUserService(db *postgres.Postgres, jwtSecret []byte) *userService {
	return &userService{db, jwtSecret}
}

func (uS *userService) CreateToken(role model.UserRole) (string, error) {
	if role != model.RoleEmployee && role != model.RoleModerator {
		return "", fmt.Errorf("Role must be 'employee' or 'moderator'")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(uS.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	return tokenString, nil
}
