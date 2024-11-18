package adminservice

import (
	"context"
	"errors"

	model "github.com/dapthehuman/learning-management-system/database/models"
	"github.com/dapthehuman/learning-management-system/dto"
	"github.com/dapthehuman/learning-management-system/service"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Repository interface {
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetUserByID(ctx context.Context, id uint64) (*model.User, error)
	UpdateUser(ctx context.Context, userID uint64, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id uint64) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAllUsers(ctx context.Context) ([]*dto.User, error) {
	users, err := s.repository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[[]*dto.User](users), nil
}

func (s *Service) GetUserByID(ctx context.Context, userID uint64) (*dto.User, error) {
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.User](user), nil
}

func (s *Service) UpdateUser(ctx context.Context, userID uint64, updateDTO *dto.UpdateUserRequest) (*dto.User, error) {
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	user = typeutil.MustConvert[*model.User](updateDTO)
	updatedUser, err := s.repository.UpdateUser(ctx, userID, user)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.User](updatedUser), nil
}

func (s *Service) DeleteUser(ctx context.Context, id uint64) error {
	return s.repository.DeleteUser(ctx, id)
}

func (s *Service) Name() string {
	return service.Admin
}
