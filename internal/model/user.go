package model

import "time"

type UserRole string

const (
	RoleEmployee  UserRole = "employee"
	RoleModerator UserRole = "moderator"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
