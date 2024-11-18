package studentservice

import (
	"context"

	model "github.com/dapthehuman/learning-management-system/database/models"
	dto "github.com/dapthehuman/learning-management-system/dto/material"
	"github.com/dapthehuman/learning-management-system/service"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Repository interface {
	Create(ctx context.Context, material *model.Material) (*model.Material, error)
	GetByID(ctx context.Context, id uint64) (*model.Material, error)
	GetByCurriculumID(ctx context.Context, curriculumID uint64) ([]*model.Material, error)
	Update(ctx context.Context, material *model.Material) (*model.Material, error)
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

func (s *Service) Create(ctx context.Context, materialDTO *dto.CreateMaterialRequest) (*dto.CreateMaterialResponse, error) {
	material := typeutil.MustConvert[*model.Material](materialDTO)

	createdMaterial, err := s.repository.Create(ctx, material)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.CreateMaterialResponse](createdMaterial), nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*dto.MaterialResponse, error) {
	material, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.MaterialResponse](material), nil
}

func (s *Service) GetByCurriculumID(ctx context.Context, curriculumID uint64) ([]*dto.MaterialResponse, error) {
	materials, err := s.repository.GetByCurriculumID(ctx, curriculumID)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[[]*dto.MaterialResponse](materials), nil
}

func (s *Service) Update(ctx context.Context, id uint64, updateDTO *dto.UpdateMaterialRequest) (*dto.MaterialResponse, error) {
	_, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updatedMaterial := typeutil.MustConvert[*model.Material](updateDTO)
	updatedMaterial, err = s.repository.Update(ctx, updatedMaterial)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.MaterialResponse](updatedMaterial), nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) Name() string {
	return service.Material
}
