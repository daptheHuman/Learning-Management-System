package userservice

import (
	"context"
	"errors"

	model "github.com/dapthehuman/learning-management-system/database/models"
	"github.com/dapthehuman/learning-management-system/dto"
	authDto "github.com/dapthehuman/learning-management-system/dto/auth"
	"github.com/dapthehuman/learning-management-system/service"
	"golang.org/x/crypto/bcrypt"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, student *model.User) (*model.User, error)
	Update(ctx context.Context, student *model.User) (*model.User, error)
	Delete(ctx context.Context, id uint64) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Register(ctx context.Context, credsDTO *authDto.RegisterRequest) (*authDto.RegisterResponse, error) {
	user := typeutil.MustConvert[*model.User](credsDTO)

	// Check if user with email already exists
	_, err := s.repository.GetByEmail(ctx, user.Email)
	if err == nil {
		return nil, errors.New("user with email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credsDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = string(hashedPassword)
	user.Role = "student"

	createdUser, err := s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*authDto.RegisterResponse](createdUser), nil
}

func (s *Service) Login(ctx context.Context, credsDTO *authDto.LoginRequest) (*dto.User, error) {
	user, err := s.repository.GetByEmail(ctx, credsDTO.Email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credsDTO.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return typeutil.MustConvert[*dto.User](user), nil
}

func (s *Service) Name() string {
	return service.User
}
