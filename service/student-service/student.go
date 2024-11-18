package studentservice

import (
	"context"

	"github.com/dapthehuman/learning-management-system/database/models"
	model "github.com/dapthehuman/learning-management-system/database/models"
	"github.com/dapthehuman/learning-management-system/dto"
	"github.com/dapthehuman/learning-management-system/service"
	"goyave.dev/goyave/v5/util/errors"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Repository interface {
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, student *model.User) (*model.User, error)

	EnrollCourse(ctx context.Context, studentID uint64, courseID uint64) (*models.Enrollment, error)
	GetEnrollmentsByUserID(ctx context.Context, studentID uint64) ([]*models.Enrollment, error)
	TrackProgress(ctx context.Context, progress *models.ProgressTracking) (*models.ProgressTracking, error)
	GetProgressByStudentAndCurriculum(ctx context.Context, studentID uint64, curriculumID uint64) ([]*models.ProgressTracking, error)

	ListAchievementByUserID(ctx context.Context, userID uint64) ([]*models.Achievement, error)
	GetAchievementByID(ctx context.Context, id uint64) (*models.Achievement, error)
	CreateAchievement(ctx context.Context, achievement *models.Achievement) (*models.Achievement, error)
	UpdateAchievement(ctx context.Context, achievement *models.Achievement) (*models.Achievement, error)
	DeleteAchievement(ctx context.Context, id uint64) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*dto.User, error) {
	student, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if student == nil {
		return nil, errors.New("Student not found")
	}
	return typeutil.MustConvert[*dto.User](student), nil
}

func (s *Service) GetAll(ctx context.Context) ([]*dto.User, error) {
	students, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[[]*dto.User](students), nil
}

func (s *Service) Update(ctx context.Context, id uint64, updateDTO *dto.UpdateStudentRequest) (*dto.User, error) {
	student, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if student == nil {
		return nil, errors.New("Student not found")
	}

	updatedStudent := typeutil.MustConvert[*model.User](updateDTO)
	updatedStudent.ID = id
	updatedStudent, err = s.repository.Update(ctx, updatedStudent)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.User](updatedStudent), nil
}

func (s *Service) EnrollCourse(ctx context.Context, enrollmentDTO *dto.EnrollStudentRequest) (*dto.Enrollment, error) {
	enrollment := typeutil.MustConvert[*models.Enrollment](enrollmentDTO)
	enrollment, err := s.repository.EnrollCourse(ctx, uint64(enrollment.UserID), uint64(enrollment.CourseID))
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.Enrollment](enrollment), nil
}

func (s *Service) GetEnrollmentsByStudentID(ctx context.Context, studentID uint64) ([]*dto.Enrollment, error) {
	enrollments, err := s.repository.GetEnrollmentsByUserID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[[]*dto.Enrollment](enrollments), nil
}

func (s *Service) TrackProgress(ctx context.Context, progressDTO *dto.TrackProgressRequest) (*dto.ProgressTracking, error) {
	progress := typeutil.MustConvert[*models.ProgressTracking](progressDTO)
	progress, err := s.repository.TrackProgress(ctx, progress)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.ProgressTracking](progress), nil
}

func (s *Service) GetProgressByStudentAndCurriculum(ctx context.Context, studentID, curriculumID uint64) ([]*dto.ProgressTracking, error) {
	progress, err := s.repository.GetProgressByStudentAndCurriculum(ctx, studentID, curriculumID)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[[]*dto.ProgressTracking](progress), nil
}

func (s *Service) GetListAchievements(ctx context.Context, studentID uint64) ([]*dto.Achievement, error) {
	achievements, err := s.repository.ListAchievementByUserID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[[]*dto.Achievement](achievements), nil
}

func (s *Service) CreateAchievement(ctx context.Context, achievementDTO *dto.CreateAchievementRequest) (*dto.Achievement, error) {
	achievement := typeutil.MustConvert[*models.Achievement](achievementDTO)

	achievement, err := s.repository.CreateAchievement(ctx, achievement)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.Achievement](achievement), nil
}

func (s *Service) GetAchievementByID(ctx context.Context, id uint64) (*dto.Achievement, error) {
	achievement, err := s.repository.GetAchievementByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.Achievement](achievement), nil
}

func (s *Service) UpdateAchievement(ctx context.Context, achievementDTO *dto.UpdateAchievementRequest) (*dto.Achievement, error) {
	achievement := typeutil.MustConvert[*models.Achievement](achievementDTO)

	achievement, err := s.repository.UpdateAchievement(ctx, achievement)
	if err != nil {
		return nil, err
	}

	return typeutil.MustConvert[*dto.Achievement](achievement), nil
}

func (s *Service) Name() string {
	return service.Student
}
