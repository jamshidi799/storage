package user

import (
	"context"
	"gorm.io/gorm"
	"log"
	"storage/domain"
)

type user struct {
	ID       int
	Email    string `gorm:"unique, index"`
	Password string
}

type postgresRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
	if err := db.AutoMigrate(user{}); err != nil {
		log.Println(err)
	}

	return &postgresRepo{db: db}
}

func (p *postgresRepo) Create(ctx context.Context, user *domain.User) error {
	u := convertToModel(user)
	return p.db.WithContext(ctx).Create(u).Error
}

func (p *postgresRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u user
	err := p.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u.toUser(), err
}

func convertToModel(u *domain.User) *user {
	return &user{
		ID:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (u *user) toUser() *domain.User {
	return &domain.User{
		Id:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}
}
