package students

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dapthehuman/learning-management-system/dto"
	"github.com/dapthehuman/learning-management-system/http/middleware"
	"github.com/dapthehuman/learning-management-system/service"
	"github.com/golang-jwt/jwt"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Service interface {
	GetByID(ctx context.Context, id uint64) (*dto.User, error)
	GetAll(ctx context.Context) ([]*dto.User, error)
	Update(ctx context.Context, id uint64, updateDTO *dto.UpdateStudentRequest) (*dto.User, error)

	EnrollCourse(ctx context.Context, enrollmentDTO *dto.EnrollStudentRequest) (*dto.Enrollment, error)
	GetEnrollmentsByStudentID(ctx context.Context, studentID uint64) ([]*dto.Enrollment, error)
	TrackProgress(ctx context.Context, progressDTO *dto.TrackProgressRequest) (*dto.ProgressTracking, error)
	GetProgressByStudentAndCurriculum(ctx context.Context, studentID, curriculumID uint64) ([]*dto.ProgressTracking, error)

	GetListAchievements(ctx context.Context, studentID uint64) ([]*dto.Achievement, error)
	CreateAchievement(ctx context.Context, achievementDTO *dto.CreateAchievementRequest) (*dto.Achievement, error)
}

type Controller struct {
	goyave.Component
	StudentService Service
}

func NewController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) Init(server *goyave.Server) {
	ctrl.StudentService = server.Service(service.Student).(Service)
	ctrl.Component.Init(server)
}

func (ctrl *Controller) RegisterRoutes(router *goyave.Router) {
	// Student routes
	authMiddleware := middleware.NewUserAuth()
	studentRouter := router.Group()
	studentRouter.Middleware(authMiddleware)

	studentSubrouter := studentRouter.Subrouter("/me")
	studentSubrouter.Get("/", ctrl.ShowCurrentUser)
	studentSubrouter.Put("/", ctrl.Update)

	studentSubrouter.Post("/enroll", ctrl.EnrollCourse)
	studentSubrouter.Post("/progress", ctrl.TrackProgressCurrentUser)
	studentSubrouter.Get("/progress/{curriculum_id}", ctrl.GetProgressCurrentUser)

	studentSubrouter.Get("/achievements", ctrl.GetAchievementsByUserID)
	studentSubrouter.Post("/achievements", ctrl.CreateAchievement)

	// Instructor routes
	instructorSubrouter := studentRouter.Subrouter("/students")
	instructorOnly := middleware.NewRoleMiddleware("admin", "instructor")
	instructorSubrouter.Middleware(instructorOnly)
	instructorSubrouter.Get("/", ctrl.Index)
	instructorSubrouter.Get("/{studentID}", ctrl.Show)
	instructorSubrouter.Put("/{studentID}", ctrl.Update)

	instructorSubrouter.Get("/{studentID}/enrollments", ctrl.GetEnrollmentsByStudentID)

	// Progress tracking
	instructorSubrouter.Post("/progress", ctrl.TrackProgress)
	instructorSubrouter.Get("/{studentID}/progress/{curriculumID}", ctrl.GetProgressByStudentAndCurriculum)
}

func (ctrl *Controller) ShowCurrentUser(response *goyave.Response, request *goyave.Request) {
	user := request.Extra["user"].(jwt.MapClaims)
	userID := uint64(user["user_id"].(float64))

	student, err := ctrl.StudentService.GetByID(request.Context(), userID)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	response.JSON(http.StatusOK, student)
}

func (ctrl *Controller) Index(response *goyave.Response, request *goyave.Request) {
	students, err := ctrl.StudentService.GetAll(request.Context())

	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	response.JSON(http.StatusOK, students)
}

func (ctrl *Controller) Show(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["studentID"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
		return
	}
	student, err := ctrl.StudentService.GetByID(request.Context(), id)
	if err != nil {
		response.JSON(http.StatusNotFound, map[string]string{"error": "Student not found"})
		return
	}
	response.JSON(http.StatusOK, student)
}

func (ctrl *Controller) Update(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["studentID"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
		return
	}

	updateDTO := typeutil.MustConvert[*dto.UpdateStudentRequest](request.Data)

	student, err := ctrl.StudentService.Update(request.Context(), id, updateDTO)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, student)
}

func (ctrl *Controller) EnrollCourse(response *goyave.Response, request *goyave.Request) {
	enrollmentDTO := typeutil.MustConvert[*dto.EnrollStudentRequest](request.Data)

	user := request.Extra["user"].(jwt.MapClaims)
	userID := uint64(user["user_id"].(float64))
	enrollmentDTO.UserID = userID

	enrollment, err := ctrl.StudentService.EnrollCourse(request.Context(), enrollmentDTO)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusCreated, enrollment)
}

func (ctrl *Controller) GetEnrollmentsByStudentID(response *goyave.Response, request *goyave.Request) {
	studentID, err := strconv.ParseUint(request.RouteParams["studentID"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
		return
	}

	enrollments, err := ctrl.StudentService.GetEnrollmentsByStudentID(request.Context(), studentID)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, enrollments)
}

func (ctrl *Controller) TrackProgress(response *goyave.Response, request *goyave.Request) {
	progressDTO := typeutil.MustConvert[*dto.TrackProgressRequest](request.Data)

	progress, err := ctrl.StudentService.TrackProgress(request.Context(), progressDTO)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusCreated, progress)
}

func (ctrl *Controller) TrackProgressCurrentUser(response *goyave.Response, request *goyave.Request) {

	user := request.Extra["user"].(jwt.MapClaims)
	userID := uint64(user["user_id"].(float64))

	fmt.Println(userID)

	progressDTO := typeutil.MustConvert[*dto.TrackProgressRequest](request.Data)
	progressDTO.UserID = int(userID)

	progress, err := ctrl.StudentService.TrackProgress(request.Context(), progressDTO)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusCreated, progress)
}

func (ctrl *Controller) GetProgressCurrentUser(response *goyave.Response, request *goyave.Request) {
	curriculumID, err := strconv.ParseUint(request.RouteParams["curriculum_id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid curriculum ID"})
		return
	}

	user := request.Extra["user"].(jwt.MapClaims)
	userID := uint64(user["user_id"].(float64))

	progress, err := ctrl.StudentService.GetProgressByStudentAndCurriculum(request.Context(), userID, curriculumID)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, progress)
}

func (ctrl *Controller) GetProgressByStudentAndCurriculum(response *goyave.Response, request *goyave.Request) {
	studentID, err := strconv.ParseUint(request.RouteParams["studentID"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
		return
	}

	curriculumID, err := strconv.ParseUint(request.RouteParams["curriculumID"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid curriculum ID"})
		return
	}

	progress, err := ctrl.StudentService.GetProgressByStudentAndCurriculum(request.Context(), studentID, curriculumID)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, progress)
}

func (ctrl *Controller) GetAchievementsByUserID(response *goyave.Response, request *goyave.Request) {
	user := request.Extra["user"].(jwt.MapClaims)
	userID := uint64(user["user_id"].(float64))

	achievements, err := ctrl.StudentService.GetListAchievements(request.Context(), userID)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, achievements)
}

func (ctrl *Controller) CreateAchievement(response *goyave.Response, request *goyave.Request) {
	achievementDTO := typeutil.MustConvert[*dto.CreateAchievementRequest](request.Data)

	user := request.Extra["user"].(jwt.MapClaims)
	userID := uint64(user["user_id"].(float64))
	achievementDTO.UserID = int(userID)

	achievement, err := ctrl.StudentService.CreateAchievement(request.Context(), achievementDTO)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusCreated, achievement)
}
