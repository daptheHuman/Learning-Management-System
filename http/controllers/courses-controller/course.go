package courses

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dapthehuman/learning-management-system/dto"
	curriculumDto "github.com/dapthehuman/learning-management-system/dto/curriculum"
	"github.com/dapthehuman/learning-management-system/http/middleware"
	"github.com/dapthehuman/learning-management-system/service"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Service interface {
	GetByID(ctx context.Context, id uint64) (*dto.Course, error)
	GetAll(ctx context.Context) ([]*dto.Course, error)
	Create(ctx context.Context, createDTO *dto.CreateCourseRequest) (*dto.Course, error)
	Update(ctx context.Context, id uint64, updateDTO *dto.UpdateCourseRequest) (*dto.Course, error)
	Delete(ctx context.Context, id uint64) error

	GetCurriculumByCourseID(ctx context.Context, courseID uint64) ([]*curriculumDto.Curriculum, error)
	GetCurriculumByID(ctx context.Context, id uint64) (*curriculumDto.Curriculum, error)
	CreateCurriculum(ctx context.Context, courseID uint64, createDTO *curriculumDto.CreateCurriculumRequest) (*curriculumDto.Curriculum, error)
	UpdateCurriculum(ctx context.Context, courseID uint64, updateDTO *curriculumDto.UpdateCurriculumRequest) (*curriculumDto.Curriculum, error)
	DeleteCurriculum(ctx context.Context, courseID uint64) error
}

type Controller struct {
	goyave.Component
	CourseService Service
}

func NewController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) Init(server *goyave.Server) {
	ctrl.CourseService = server.Service(service.Course).(Service)
	ctrl.Component.Init(server)
}

func (ctrl *Controller) RegisterRoutes(router *goyave.Router) {
	subrouter := router.Subrouter("/courses")

	authMiddleware := middleware.NewUserAuth()
	subrouter.Middleware(authMiddleware)
	subrouter.Get("/", ctrl.Index)
	subrouter.Get("/{id}", ctrl.Show)

	// CRUD routes
	instructorOnly := middleware.NewRoleMiddleware("admin", "instructor")
	instructorRouter := subrouter.Group()
	instructorRouter.Middleware(instructorOnly)
	instructorRouter.Post("/", ctrl.Create)
	instructorRouter.Put("/{id}", ctrl.Update)
	instructorRouter.Delete("/{id}", ctrl.Delete)

	// Curriculum-related routes nested under a course
	subrouter.Get("/{id}/curriculums", ctrl.ListCurriculum)                             // List curriculum for a course
	subrouter.Get("/{id}/curriculums/{curriculum_id}", ctrl.ShowCurriculum)             // List curriculum for a course
	instructorRouter.Post("/{id}/curriculums", ctrl.AddCuriculum)                       // Add a curriculum section
	instructorRouter.Put("/{id}/curriculums/{curriculum_id}", ctrl.UpdateCurriculum)    // Update a curriculum section
	instructorRouter.Delete("/{id}/curriculums/{curriculum_id}", ctrl.DeleteCurriculum) // Delete a curriculum section
}

func (ctrl *Controller) Index(response *goyave.Response, request *goyave.Request) {
	courses, err := ctrl.CourseService.GetAll(request.Context())
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	response.JSON(http.StatusOK, courses)
}

func (ctrl *Controller) Show(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}
	course, err := ctrl.CourseService.GetByID(request.Context(), id)
	if err != nil {
		response.JSON(http.StatusNotFound, map[string]string{"error": "Course not found"})
		return
	}
	response.JSON(http.StatusOK, course)
}

func (ctrl *Controller) Create(response *goyave.Response, request *goyave.Request) {
	createDTO := typeutil.MustConvert[*dto.CreateCourseRequest](request.Data)
	course, err := ctrl.CourseService.Create(request.Context(), createDTO)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusCreated, course)
}

func (ctrl *Controller) Update(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}

	updateDTO := typeutil.MustConvert[*dto.UpdateCourseRequest](request.Data)

	course, err := ctrl.CourseService.Update(request.Context(), id, updateDTO)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, course)
}

func (ctrl *Controller) Delete(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}

	if err := ctrl.CourseService.Delete(request.Context(), id); err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK,
		map[string]string{"message": "Course deleted successfully"})
}

func (ctrl *Controller) AddCuriculum(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["id"], 10, 64)

	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}

	// check if course ID is valid and exists
	_, err = ctrl.CourseService.GetByID(request.Context(), id)
	if err != nil {
		response.JSON(http.StatusNotFound, map[string]string{"error": "Course not found"})
		return
	}

	createDTO := typeutil.MustConvert[*curriculumDto.CreateCurriculumRequest](request.Data)
	curriculum, err := ctrl.CourseService.CreateCurriculum(request.Context(), id, createDTO)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusCreated, curriculum)
}

func (ctrl *Controller) ListCurriculum(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}
	curriculum, err := ctrl.CourseService.GetCurriculumByCourseID(request.Context(), id)
	if err != nil {
		response.JSON(http.StatusNotFound, map[string]string{"error": "Curriculum not found"})
		return
	}
	response.JSON(http.StatusOK, curriculum)
}

func (ctrl *Controller) ShowCurriculum(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["curriculum_id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid curriculum ID"})
		return
	}
	curriculum, err := ctrl.CourseService.GetCurriculumByID(request.Context(), id)
	if err != nil {
		response.JSON(http.StatusNotFound, map[string]string{"error": "Curriculum not found"})
		return
	}
	response.JSON(http.StatusOK, curriculum)
}

func (ctrl *Controller) UpdateCurriculum(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["curriculum_id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid curriculum ID"})
		return
	}
	updateDTO := typeutil.MustConvert[*curriculumDto.UpdateCurriculumRequest](request.Data)

	curriculum, err := ctrl.CourseService.UpdateCurriculum(request.Context(), id, updateDTO)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, curriculum)
}

func (ctrl *Controller) DeleteCurriculum(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["curriculum_id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid curriculum ID"})
		return
	}
	if err := ctrl.CourseService.DeleteCurriculum(request.Context(), id); err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, map[string]string{"message": "Curriculum deleted successfully"})
}
