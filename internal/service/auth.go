package service

import "context"

type AuthService interface {
	register(ctx context.Context)
	login(ctx context.Context)
}
