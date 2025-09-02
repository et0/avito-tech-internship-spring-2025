package repository

import "github.com/et0/avito-tech-internship-spring-2025/internal/model"

type Database interface {
	FindByEmail(email string) (*model.User, error)
}
