package middleware

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

// ApiAuth validates JWT Bearer tokens for API routes.
// Equivalent to Laravel's auth:api middleware.
type ApiAuth struct{}

func NewApiAuth() *ApiAuth {
	return &ApiAuth{}
}

func (m *ApiAuth) Handle(ctx http.Context) {
	header := ctx.Request().Header("Authorization", "")
	if !strings.HasPrefix(header, "Bearer ") {
		ctx.Response().Json(401, http.Json{"message": "Unauthorized"})
		return
	}

	token := strings.TrimPrefix(header, "Bearer ")
	payload, err := facades.Auth(ctx).Guard("api").Parse(token)
	if err != nil || payload == nil {
		ctx.Response().Json(401, http.Json{"message": "Unauthorized"})
		return
	}

	var user models.User
	if err := facades.Auth(ctx).Guard("api").User(&user); err != nil || user.ID == 0 {
		ctx.Response().Json(401, http.Json{"message": "Unauthorized"})
		return
	}

	ctx.Request().Next()
}
