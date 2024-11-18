package assessments

import (
	"context"
	"net/http"
	"strconv"

	dto "github.com/dapthehuman/learning-management-system/dto/assesment"
	"github.com/dapthehuman/learning-management-system/http/middleware"
	"github.com/dapthehuman/learning-management-system/service"
	"github.com/golang-jwt/jwt"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Service interface {
	CreateAssessment(ctx context.Context, createDTO *dto.CreateAssessmentRequest) (*dto.Assessment, error)
	GetAssessmentByCourseID(ctx context.Context, courseID uint64) ([]*dto.Assessment, error)
	GetAssessmentByID(ctx context.Context, assessmentID uint64) (*dto.Assessment, error)
	SubmitAnswer(ctx context.Context, submission *dto.SubmissionRequest) (*dto.SubmissionResponse, error)
}

type Controller struct {
	goyave.Component
	assessmentService Service
}

func NewController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) Init(server *goyave.Server) {
	ctrl.assessmentService = server.Service(service.Assessment).(Service)
	ctrl.Component.Init(server)
}

func (ctrl *Controller) RegisterRoutes(router *goyave.Router) {
	subrouter := router.Subrouter("/assessment")

	authMiddleware := middleware.NewUserAuth()
	subrouter.Middleware(authMiddleware)

	roleMiddleware := middleware.NewRoleMiddleware("admin", "instructor")
	subrouter.Post("/", ctrl.Create).Middleware(roleMiddleware)
	subrouter.Get("/course/{courseID}", ctrl.GetByCourseID)
	subrouter.Post("/submit", ctrl.SubmitAnswer)
}

func (ctrl *Controller) Create(response *goyave.Response, request *goyave.Request) {
	createDTO := typeutil.MustConvert[*dto.CreateAssessmentRequest](request.Data)
	assessment, err := ctrl.assessmentService.CreateAssessment(request.Context(), createDTO)
	if err != nil {
		response.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusCreated, assessment)
}

func (ctrl *Controller) GetByCourseID(response *goyave.Response, request *goyave.Request) {
	courseID, err := strconv.ParseUint(request.RouteParams["courseID"], 10, 64)
	if err != nil {
		response.JSON(400, map[string]string{"error": "Invalid course ID"})
		return
	}

	assessments, err := ctrl.assessmentService.GetAssessmentByCourseID(request.Context(), courseID)
	if err != nil {
		response.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, assessments)
}

func (ctrl *Controller) SubmitAnswer(response *goyave.Response, request *goyave.Request) {
	submission := typeutil.MustConvert[*dto.SubmissionRequest](request.Data)

	user := request.Extra["user"].(jwt.MapClaims)
	userID := uint64(user["user_id"].(float64))
	submission.UserID = int(userID)

	submittedAnswer, err := ctrl.assessmentService.SubmitAnswer(request.Context(), submission)
	if err != nil {
		response.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusCreated, submittedAnswer)
}
