package middleware

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

// Admin restricts access to users with the "admin" role.
// Equivalent to smbgen's CompanyAdministrator middleware.
type Admin struct{}

func NewAdmin() *Admin {
	return &Admin{}
}

func (m *Admin) Handle(ctx http.Context) {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || !user.IsAdmin() {
		ctx.Response().Json(403, http.Json{
			"message": "Forbidden",
		})
		return
	}
	ctx.Request().Next()
}
