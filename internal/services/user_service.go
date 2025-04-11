package services

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
	"github.com/rs/zerolog"
)

type IUserService interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	ListUsers(ctx context.Context, page, limit int) ([]models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserService struct {
	log  *zerolog.Logger
	repo repository.IUserRepository
}

func NewUserService(log *zerolog.Logger, repo repository.IUserRepository) IUserService {
	return &UserService{
		log:  log,
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	//// Add business logic here
	//user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	//user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return s.repo.Create(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {

	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {

	return s.repo.GetByUsername(ctx, username)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UserService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	return s.repo.GetByPhoneNumber(ctx, phoneNumber)
}

func (s *UserService) ListUsers(ctx context.Context, page, limit int) ([]models.User, error) {
	return s.repo.List(ctx, page, limit)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return s.repo.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
