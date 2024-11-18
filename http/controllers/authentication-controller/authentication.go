package auth

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/dapthehuman/learning-management-system/dto"
	authDto "github.com/dapthehuman/learning-management-system/dto/auth"
	"github.com/dapthehuman/learning-management-system/service"
	"github.com/golang-jwt/jwt"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Service interface {
	Register(ctx context.Context, credsDTO *authDto.RegisterRequest) (*authDto.RegisterResponse, error)
	Login(ctx context.Context, loginDTO *authDto.LoginRequest) (*dto.User, error)
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
	subrouter := router.Subrouter("/auth")

	subrouter.Post("/register", ctrl.Register)
	subrouter.Post("/login", ctrl.Login)

}

func (ctrl *Controller) Register(response *goyave.Response, request *goyave.Request) {
	registerDTO := typeutil.MustConvert[*authDto.RegisterRequest](request.Data)
	user, err := ctrl.UserService.Register(request.Context(), registerDTO)
	if err != nil {
		response.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	response.JSON(http.StatusCreated, user)
}

func (ctrl *Controller) Login(response *goyave.Response, request *goyave.Request) {
	loginDTO := typeutil.MustConvert[*authDto.LoginRequest](request.Data)
	user, err := ctrl.UserService.Login(request.Context(), loginDTO)
	if err != nil {
		response.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	// Generate Token
	token := ctrl.generateToken(user)
	response.JSON(http.StatusOK, map[string]string{"token": token})

}

func (ctrl *Controller) generateToken(user *dto.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 3 days
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("APP_SECRET")))
	if err != nil {
		return ""
	}
	return tokenString
}
