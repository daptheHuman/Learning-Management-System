package courseservice

import (
	"context"

	model "github.com/dapthehuman/learning-management-system/database/models"
	"github.com/dapthehuman/learning-management-system/dto"
	curriculumDto "github.com/dapthehuman/learning-management-system/dto/curriculum"
	"github.com/dapthehuman/learning-management-system/service"
	"goyave.dev/goyave/v5/util/errors"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Repository interface {
	First(ctx context.Context, id uint64) (*model.Course, error)
	GetAll(ctx context.Context) ([]*model.Course, error)
	Create(ctx context.Context, course *model.Course) (*model.Course, error)
	Update(ctx context.Context, course *model.Course) (*model.Course, error)
	Delete(ctx context.Context, id uint64) error

	CreateCurriculum(ctx context.Context, courseID uint64, curriculum *model.Curriculum) (*model.Curriculum, error)
	GetCurriculum(ctx context.Context, courseID uint64) ([]*model.Curriculum, error)
	GetCurriculumByID(ctx context.Context, id uint64) (*model.Curriculum, error)
	UpdateCurriculum(ctx context.Context, curriculum *model.Curriculum) (*model.Curriculum, error)
	DeleteCurriculum(ctx context.Context, id uint64) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*dto.Course, error) {
	course, err := s.repository.First(ctx, id)

	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("Course not found!")
	}
	return typeutil.MustConvert[*dto.Course](course), nil
}

func (s *Service) GetAll(ctx context.Context) ([]*dto.Course, error) {
	courses, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[[]*dto.Course](courses), nil
}

func (s *Service) Create(ctx context.Context, createDTO *dto.CreateCourseRequest) (*dto.Course, error) {
	course := typeutil.MustConvert[*model.Course](createDTO)
	createdCourse, err := s.repository.Create(ctx, course)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.Course](createdCourse), nil
}

func (s *Service) Update(ctx context.Context, id uint64, updateDTO *dto.UpdateCourseRequest) (*dto.Course, error) {
	course, err := s.repository.First(ctx, id)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("Course not found")
	}

	nCourse := typeutil.MustConvert[*model.Course](updateDTO)
	nCourse.ID = id
	updatedCourse, err := s.repository.Update(ctx, nCourse)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.Course](updatedCourse), nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	course, err := s.repository.First(ctx, id)
	if err != nil {
		return err
	}

	if course == nil {
		return errors.New("Course not found")
	}

	return s.repository.Delete(ctx, id)
}

func (s *Service) CreateCurriculum(ctx context.Context, courseID uint64, createDTO *curriculumDto.CreateCurriculumRequest) (*curriculumDto.Curriculum, error) {
	curriculum := typeutil.MustConvert[*model.Curriculum](createDTO)
	createdCurriculum, err := s.repository.CreateCurriculum(ctx, courseID, curriculum)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*curriculumDto.Curriculum](createdCurriculum), nil
}

func (s *Service) GetCurriculumByCourseID(ctx context.Context, courseID uint64) ([]*curriculumDto.Curriculum, error) {
	curriculum, err := s.repository.GetCurriculum(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[[]*curriculumDto.Curriculum](curriculum), nil
}

func (s *Service) GetCurriculumByID(ctx context.Context, id uint64) (*curriculumDto.Curriculum, error) {
	curriculum, err := s.repository.GetCurriculumByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*curriculumDto.Curriculum](curriculum), nil
}

func (s *Service) UpdateCurriculum(ctx context.Context, courseID uint64, updateDTO *curriculumDto.UpdateCurriculumRequest) (*curriculumDto.Curriculum, error) {
	curriculum := typeutil.MustConvert[*model.Curriculum](updateDTO)
	updatedCurriculum, err := s.repository.UpdateCurriculum(ctx, curriculum)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*curriculumDto.Curriculum](updatedCurriculum), nil
}

func (s *Service) DeleteCurriculum(ctx context.Context, id uint64) error {
	return s.repository.DeleteCurriculum(ctx, id)
}

func (s *Service) Name() string {
	return service.Course
}
