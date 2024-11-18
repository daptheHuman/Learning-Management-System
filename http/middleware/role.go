package middleware

import (
	"github.com/golang-jwt/jwt"
	"goyave.dev/goyave/v5"
)

type RoleMiddleware struct {
	goyave.Component
	AllowedRoles []string
}

func NewRoleMiddleware(allowedRoles ...string) *RoleMiddleware {
	return &RoleMiddleware{AllowedRoles: allowedRoles}
}

func (rm *RoleMiddleware) Handle(next goyave.Handler) goyave.Handler {
	return func(response *goyave.Response, request *goyave.Request) {
		// Implement role-based middleware
		claims, ok := request.Extra["user"].(jwt.MapClaims)
		if !ok {
			response.Status(401)
			return
		}

		role, ok := claims["role"].(string)
		allowed := false
		for _, allowedRole := range rm.AllowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !ok || !allowed {
			response.Status(403)
			return
		}

		next(response, request)
	}
}
