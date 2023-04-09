package user

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"storage/domain"
)

type service struct {
	repo           domain.UserRepository
	tokenGenerator domain.TokenGenerator
}

func NewUserService(repo domain.UserRepository, tg domain.TokenGenerator) domain.UserService {
	return &service{
		repo:           repo,
		tokenGenerator: tg,
	}
}

func (s *service) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	password := string(hashedPass)

	user := domain.User{
		Email:    req.Email,
		Password: password,
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	token, err := s.tokenGenerator.Generate(user.Id)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{Token: token}, nil
}

func (s *service) VerifyToken(token string) bool {
	return s.tokenGenerator.Verify(token)
}
