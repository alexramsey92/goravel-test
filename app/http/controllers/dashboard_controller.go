package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

// Index GET /dashboard
func (c *DashboardController) Index(ctx http.Context) http.Response {
	var user models.User
	_ = facades.Auth(ctx).User(&user)

	clientCount, _ := facades.Orm().Query().Model(&models.Client{}).Count()
	pageCount, _ := facades.Orm().Query().Model(&models.Page{}).Count()
	leadCount, _ := facades.Orm().Query().Model(&models.Client{}).Where("status = ?", models.ClientStatusLead).Count()

	return ctx.Response().View().Make("dashboard.tmpl", map[string]any{
		"user":        user,
		"clientCount": clientCount,
		"pageCount":   pageCount,
		"leadCount":   leadCount,
		"activePage":  "dashboard",
	})
}
