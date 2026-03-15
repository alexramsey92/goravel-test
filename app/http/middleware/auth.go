package middleware

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

// Auth redirects unauthenticated users to /login (web guard).
// This is the Goravel equivalent of Laravel's 'auth' middleware.
type Auth struct{}

func NewAuth() *Auth {
	return &Auth{}
}

func (m *Auth) Handle(ctx http.Context) {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || user.ID == 0 {
		ctx.Response().Redirect(302, "/login")
		return
	}
	ctx.Request().Next()
}
