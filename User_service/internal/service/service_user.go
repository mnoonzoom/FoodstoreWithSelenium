package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"user/internal/dao"
	"user/internal/model"
)

type UserService struct {
	repo *dao.UserRepository
}

func NewUserService(repo *dao.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, user model.User) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)

	return s.repo.CreateUser(ctx, user)
}
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}
