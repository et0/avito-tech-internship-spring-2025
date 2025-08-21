package service

import (
	"github.com/et0/avito-tech-internship-spring-2025/internal/repository/postgres"
)

type UserService struct {
	db        *postgres.Postgres
	jwtSecret []byte
}

func NewUserService(db *postgres.Postgres, jwtSecret []byte) *UserService {
	return &UserService{db, jwtSecret}
}
