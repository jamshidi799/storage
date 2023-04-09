package domain

import "context"

type User struct {
	Id       int
	Email    string
	Password string
}

type UserService interface {
	Register(ctx context.Context, req *User) (*User, error)
	Login(ctx context.Context, req *User) (string, error)
	VerifyToken(token string) bool
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type TokenGenerator interface {
	Generate(id int) (string, error)
	Verify(token string) bool
}
