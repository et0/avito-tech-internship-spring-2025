package service

import (
	"fmt"
	"time"

	"github.com/et0/avito-tech-internship-spring-2025/internal/model"
	"github.com/et0/avito-tech-internship-spring-2025/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateToken(role model.UserRole) (string, error)
	Register(email string, password string, role model.UserRole) (*model.User, error)
}

type userService struct {
	db        repository.Database
	jwtSecret []byte
}

func NewUserService(db repository.Database, jwtSecret []byte) *userService {
	return &userService{db, jwtSecret}
}

func (uS *userService) CreateToken(role model.UserRole) (string, error) {
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

func (uS *userService) Register(email string, password string, role model.UserRole) (*model.User, error) {
	existingUser, err := uS.db.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	user, err := uS.db.CreateUser(email, string(hashedPassword), role)
	if err != nil {
		return nil, err
	}

	return user, nil
}
