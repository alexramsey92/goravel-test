package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

type PageController struct{}

func NewPageController() *PageController {
	return &PageController{}
}

// Index GET /pages
func (c *PageController) Index(ctx http.Context) http.Response {
	var user models.User
	_ = facades.Auth(ctx).User(&user)

	var pages []models.Page
	facades.Orm().Query().OrderBy("created_at", "desc").Find(&pages)

	return ctx.Response().View().Make("pages/index.tmpl", map[string]any{
		"user":       user,
		"pages":      pages,
		"activePage": "pages",
	})
}

// Create GET /pages/create
func (c *PageController) Create(ctx http.Context) http.Response {
	var user models.User
	_ = facades.Auth(ctx).User(&user)

	return ctx.Response().View().Make("pages/create.tmpl", map[string]any{
		"user":       user,
		"activePage": "pages",
	})
}

// Store POST /pages
func (c *PageController) Store(ctx http.Context) http.Response {
	published := ctx.Request().Input("published") == "on"

	page := models.Page{
		Title:     ctx.Request().Input("title"),
		Slug:      ctx.Request().Input("slug"),
		Content:   ctx.Request().Input("content"),
		Published: published,
	}

	if err := facades.Orm().Query().Create(&page); err != nil {
		var user models.User
		_ = facades.Auth(ctx).User(&user)
		return ctx.Response().View().Make("pages/create.tmpl", map[string]any{
			"user":  user,
			"error": "Failed to create page. Slug may already be taken.",
		})
	}

	return ctx.Response().Redirect(302, "/pages")
}

// Destroy POST /pages/:id/delete
func (c *PageController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var page models.Page
	if err := facades.Orm().Query().Find(&page, id); err == nil && page.ID != 0 {
		facades.Orm().Query().Delete(&page)
	}
	return ctx.Response().Redirect(302, "/pages")
}
