package user

import (
	"context"
	"storage/domain"
)

type service struct {
	repo           domain.UserRepository
	tokenGenerator tokenGenerator
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	user := domain.User{
		Email:    req.Email,
		Password: hashPassword(req.Password),
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &domain.RegisterResponse{
		Id:    user.Id,
		Email: user.Email,
	}, nil
}

func (s *service) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user.Password != hashPassword(req.Password) {
		return nil, domain.BadRequestError("password is wrong")
	}

	token := s.tokenGenerator.generate(user.Id)
	return &domain.LoginResponse{Token: token}, nil
}

type tokenGenerator interface {
	generate(id int) string
}

func hashPassword(pass string) string {
	return pass
}
