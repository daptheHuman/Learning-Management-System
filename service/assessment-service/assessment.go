package assessmentservice

import (
	"context"

	"github.com/dapthehuman/learning-management-system/database/models"
	dto "github.com/dapthehuman/learning-management-system/dto/assesment"
	"github.com/dapthehuman/learning-management-system/service"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Repository interface {
	Create(ctx context.Context, assessment *models.Assessment) (*models.Assessment, error)
	GetAllByCourseID(ctx context.Context, courseID uint64) ([]*models.Assessment, error)
	GetByID(ctx context.Context, assessmentID uint64) (*models.Assessment, error)
	SubmitAnswer(ctx context.Context, submission *models.Submission) (*models.Submission, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateAssessment(ctx context.Context, createDTO *dto.CreateAssessmentRequest) (*dto.Assessment, error) {
	assessment := typeutil.MustConvert[*models.Assessment](createDTO)
	createdAssessment, err := s.repository.Create(ctx, assessment)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.Assessment](createdAssessment), nil
}

func (s *Service) GetAssessmentByCourseID(ctx context.Context, courseID uint64) ([]*dto.Assessment, error) {
	assessments, err := s.repository.GetAllByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[[]*dto.Assessment](assessments), nil
}

func (s *Service) GetAssessmentByID(ctx context.Context, assessmentID uint64) (*dto.Assessment, error) {
	assessment, err := s.repository.GetByID(ctx, assessmentID)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.Assessment](assessment), nil
}

func (s *Service) SubmitAnswer(ctx context.Context, submission *dto.SubmissionRequest) (*dto.SubmissionResponse, error) {
	submissionModel := typeutil.MustConvert[*models.Submission](submission)
	submittedAnswer, err := s.repository.SubmitAnswer(ctx, submissionModel)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.SubmissionResponse](submittedAnswer), nil
}

func (s *Service) Name() string {
	return service.Assessment
}
