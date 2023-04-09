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

func (s *service) Register(ctx context.Context, u *domain.User) (*domain.User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	password := string(hashedPass)

	user := domain.User{
		Email:    u.Email,
		Password: password,
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) Login(ctx context.Context, u *domain.User) (string, error) {
	user, err := s.repo.GetByEmail(ctx, u.Email)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return "", err
	}

	token, err := s.tokenGenerator.Generate(user.Id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) VerifyToken(token string) bool {
	return s.tokenGenerator.Verify(token)
}
