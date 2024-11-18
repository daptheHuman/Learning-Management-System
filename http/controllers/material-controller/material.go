package materials

import (
	"context"
	"net/http"
	"strconv"

	dto "github.com/dapthehuman/learning-management-system/dto/material"
	"github.com/dapthehuman/learning-management-system/http/middleware"
	"github.com/dapthehuman/learning-management-system/service"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Service interface {
	Create(ctx context.Context, materialDTO *dto.CreateMaterialRequest) (*dto.CreateMaterialResponse, error)
	GetByID(ctx context.Context, id uint64) (*dto.MaterialResponse, error)
	GetByCurriculumID(ctx context.Context, curriculumID uint64) ([]*dto.MaterialResponse, error)
	Update(ctx context.Context, id uint64, updateDTO *dto.UpdateMaterialRequest) (*dto.MaterialResponse, error)
	Delete(ctx context.Context, id uint64) error
}

type Controller struct {
	goyave.Component
	MaterialService Service
}

func NewController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) Init(server *goyave.Server) {
	ctrl.MaterialService = server.Service(service.Material).(Service)
	ctrl.Component.Init(server)
}

func (ctrl *Controller) RegisterRoutes(router *goyave.Router) {
	subrouter := router.Subrouter("/curriculums/{curriculum_id}/materials")

	authMiddleware := middleware.NewUserAuth()
	subrouter.Middleware(authMiddleware)

	subrouter.Get("/", ctrl.Index)    // Get all materials for a curriculum
	subrouter.Get("/{id}", ctrl.Show) // Get a material by ID

	instructorOnly := middleware.NewRoleMiddleware("admin", "instructor")
	instructorRouter := subrouter.Group().Middleware(instructorOnly)
	instructorRouter.Post("/", ctrl.Create)       // Create a material for a curriculum
	instructorRouter.Put("/{id}", ctrl.Update)    // Update a material by ID
	instructorRouter.Delete("/{id}", ctrl.Delete) // Delete a material by ID
}

func (ctrl *Controller) Create(response *goyave.Response, request *goyave.Request) {
	materialDTO := typeutil.MustConvert[*dto.CreateMaterialRequest](request.Data)

	curriculumID, err := strconv.ParseUint(request.RouteParams["curriculum_id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid curriculum ID"})
		return
	}

	materialDTO.CurriculumID = curriculumID

	createdMaterial, err := ctrl.MaterialService.Create(request.Context(), materialDTO)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusCreated, createdMaterial)
}

func (ctrl *Controller) Index(response *goyave.Response, request *goyave.Request) {
	curriculumID, err := strconv.ParseUint(request.RouteParams["curriculum_id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid curriculum ID"})
		return
	}

	materials, err := ctrl.MaterialService.GetByCurriculumID(request.Context(), curriculumID)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, materials)
}

func (ctrl *Controller) Show(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid material ID"})
		return
	}

	material, err := ctrl.MaterialService.GetByID(request.Context(), id)
	if err != nil {
		response.JSON(http.StatusNotFound, map[string]string{"error": "Material not found"})
		return
	}

	response.JSON(http.StatusOK, material)
}

func (ctrl *Controller) Update(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid material ID"})
		return
	}

	updateDTO := typeutil.MustConvert[*dto.UpdateMaterialRequest](request.Data)

	material, err := ctrl.MaterialService.Update(request.Context(), id, updateDTO)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, material)
}

func (ctrl *Controller) Delete(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid material ID"})
		return
	}

	err = ctrl.MaterialService.Delete(request.Context(), id)
	if err != nil {
		response.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusOK, map[string]string{"message": "Material deleted successfully"})
}
