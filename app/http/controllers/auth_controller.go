package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// ShowLoginForm GET /login
func (c *AuthController) ShowLoginForm(ctx http.Context) http.Response {
	return ctx.Response().View().Make("auth/login.tmpl", map[string]any{})
}

// ShowRegisterForm GET /register
func (c *AuthController) ShowRegisterForm(ctx http.Context) http.Response {
	return ctx.Response().View().Make("auth/register.tmpl", map[string]any{})
}

// Login POST /login
func (c *AuthController) Login(ctx http.Context) http.Response {
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")

	var user models.User
	if err := facades.Orm().Query().Where("email = ?", email).First(&user); err != nil {
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"error": "No account found with that email.",
		})
	}

	ok := facades.Hash().Check(password, user.Password)
	if !ok {
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"error": "Incorrect password.",
		})
	}

	if _, err := facades.Auth(ctx).Login(&user); err != nil {
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"error": "Login failed. Please try again.",
		})
	}

	return ctx.Response().Redirect(302, "/dashboard")
}

// Register POST /register
func (c *AuthController) Register(ctx http.Context) http.Response {
	name := ctx.Request().Input("name")
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")

	// Check email not taken
	var existing models.User
	err := facades.Orm().Query().Where("email = ?", email).First(&existing)
	if err == nil && existing.ID != 0 {
		return ctx.Response().View().Make("auth/register.tmpl", map[string]any{
			"error": "An account with that email already exists.",
		})
	}

	hashed, err := facades.Hash().Make(password)
	if err != nil {
		return ctx.Response().View().Make("auth/register.tmpl", map[string]any{
			"error": "Registration failed.",
		})
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: hashed,
		Role:     models.RoleAdmin, // first user gets admin — change as needed
	}

	if err := facades.Orm().Query().Create(&user); err != nil {
		return ctx.Response().View().Make("auth/register.tmpl", map[string]any{
			"error": "Could not create account.",
		})
	}

	if _, err := facades.Auth(ctx).Login(&user); err != nil {
		return ctx.Response().Redirect(302, "/login")
	}

	return ctx.Response().Redirect(302, "/dashboard")
}

// Logout POST /logout
func (c *AuthController) Logout(ctx http.Context) http.Response {
	_ = facades.Auth(ctx).Logout()
	return ctx.Response().Redirect(302, "/login")
}
