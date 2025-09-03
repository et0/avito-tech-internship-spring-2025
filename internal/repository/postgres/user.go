package postgres

import (
	"context"
	"errors"
	"log"

	"github.com/et0/avito-tech-internship-spring-2025/internal/model"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) FindByEmail(email string) (*model.User, error) {
	conn, err := p.Pool.Acquire(context.Background())
	if err != nil {
		log.Fatal("DB connect failed:", err)
	}
	defer conn.Release()

	var user model.User

	err = conn.QueryRow(context.Background(), "SELECT id,email,password,role  FROM users WHERE email = $1 LIMIT 1", email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Role)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *Postgres) CreateUser(email, password string, role model.UserRole) (*model.User, error) {
	conn, err := p.Pool.Acquire(context.Background())
	if err != nil {
		log.Fatal("DB connect failed:", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		"INSERT INTO users (email, password, role) VALUES ($1, $2, $3)",
		email, password, role,
	)
	if err != nil {
		return nil, err
	}

	user, err := p.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
