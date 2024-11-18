package admins

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dapthehuman/learning-management-system/dto"
	"github.com/dapthehuman/learning-management-system/http/middleware"
	"github.com/dapthehuman/learning-management-system/service"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Service interface {
	GetAllUsers(ctx context.Context) ([]*dto.User, error)
	GetUserByID(ctx context.Context, userID uint64) (*dto.User, error)
	UpdateUser(ctx context.Context, userID uint64, updateDTO *dto.UpdateUserRequest) (*dto.User, error)
	DeleteUser(ctx context.Context, id uint64) error
}

type Controller struct {
	goyave.Component
	UserService Service
}

func NewController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) Init(server *goyave.Server) {
	ctrl.UserService = server.Service(service.User).(Service)
	ctrl.Component.Init(server)
}

func (ctrl *Controller) RegisterRoutes(router *goyave.Router) {
	subrouter := router.Subrouter("/admin")

	authMiddleware := middleware.NewUserAuth()
	subrouter.Middleware(authMiddleware)

	// Enrollment-related routes
	adminOnly := middleware.NewRoleMiddleware("admin", "instructor")
	adminSubrouter := subrouter.Group().Middleware(adminOnly)
	adminSubrouter.Get("/users", ctrl.GetUsers)
	adminSubrouter.Get("/users/{id}", ctrl.GetUserByID)
	adminSubrouter.Patch("/users/{id}", ctrl.UpdateUser)
	adminSubrouter.Delete("/users/{id}", ctrl.DeleteUser)

}

func (ctrl *Controller) GetUsers(response *goyave.Response, request *goyave.Request) {
	users, err := ctrl.UserService.GetAllUsers(request.Context())
	if err != nil {
		response.Error(err)
		return
	}
	response.JSON(http.StatusOK, users)
}

func (ctrl *Controller) GetUserByID(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["id"], 10, 64)
	user, err := ctrl.UserService.GetUserByID(request.Context(), id)
	if err != nil {
		response.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}
	response.JSON(http.StatusOK, user)
}

func (ctrl *Controller) UpdateUser(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	updateDTO := typeutil.MustConvert[*dto.UpdateUserRequest](request.Data)
	user, err := ctrl.UserService.UpdateUser(request.Context(), id, updateDTO)
	if err != nil {
		response.Error(err)
		return
	}
	response.JSON(http.StatusOK, user)
}

func (ctrl *Controller) DeleteUser(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseUint(request.RouteParams["id"], 10, 64)
	if err != nil {
		response.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	if err := ctrl.UserService.DeleteUser(request.Context(), id); err != nil {
		response.Error(err)
		return
	}
	response.Status(http.StatusNoContent)
}
