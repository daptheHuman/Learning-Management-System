package middleware

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"goyave.dev/goyave/v5"
)

type UserAuth struct {
	goyave.Component
}

func NewUserAuth() *UserAuth {
	return &UserAuth{}
}

func (m *UserAuth) Handle(next goyave.Handler) goyave.Handler {
	return func(response *goyave.Response, request *goyave.Request) {
		authHeader := request.Header().Get("Authorization")
		if authHeader == "" {
			response.Status(401)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			response.Status(401)
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("APP_SECRET")), nil
		})

		if err != nil || !token.Valid {
			response.Status(401)
			return
		}

		// Convert claims to the correct type
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			request.Extra["user"] = claims
			next(response, request)
		} else {
			response.Status(401)
			return
		}
	}
}
