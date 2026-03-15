package controllers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

type ClientController struct{}

func NewClientController() *ClientController {
	return &ClientController{}
}

// slugify converts a string to a URL-safe slug.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`[\s-]+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// uniqueSlug ensures the slug is unique in the clients table, appending -2, -3 etc. if needed.
// Pass excludeID > 0 to ignore the current record (used on updates).
func uniqueSlug(base string, excludeID uint) string {
	candidate := base
	for i := 2; ; i++ {
		q := facades.Orm().Query().Model(&models.Client{}).Where("slug", candidate)
		if excludeID > 0 {
			q = q.Where("id != ?", excludeID)
		}
		count, _ := q.Count()
		if count == 0 {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

// Index GET /clients
func (c *ClientController) Index(ctx http.Context) http.Response {
	var user models.User
	_ = facades.Auth(ctx).User(&user)

	var clients []models.Client
	facades.Orm().Query().OrderBy("created_at", "desc").Find(&clients)

	return ctx.Response().View().Make("clients/index.tmpl", map[string]any{
		"user":       user,
		"clients":    clients,
		"activePage": "clients",
	})
}

// Create GET /clients/create
func (c *ClientController) Create(ctx http.Context) http.Response {
	var user models.User
	_ = facades.Auth(ctx).User(&user)

	return ctx.Response().View().Make("clients/create.tmpl", map[string]any{
		"user":       user,
		"activePage": "clients",
	})
}

// Store POST /clients
func (c *ClientController) Store(ctx http.Context) http.Response {
	name := ctx.Request().Input("name")

	rawSlug := ctx.Request().Input("slug")
	if rawSlug == "" {
		rawSlug = name
	}
	slug := uniqueSlug(slugify(rawSlug), 0)

	client := models.Client{
		Name:    name,
		Slug:    slug,
		Email:   ctx.Request().Input("email"),
		Phone:   ctx.Request().Input("phone"),
		Company: ctx.Request().Input("company"),
		Status:  ctx.Request().Input("status"),
		Notes:   ctx.Request().Input("notes"),
	}

	if err := facades.Orm().Query().Create(&client); err != nil {
		var user models.User
		_ = facades.Auth(ctx).User(&user)
		return ctx.Response().View().Make("clients/create.tmpl", map[string]any{
			"user":  user,
			"error": "Failed to create client.",
		})
	}

	return ctx.Response().Redirect(302, "/clients")
}

// Show GET /clients/:id
func (c *ClientController) Show(ctx http.Context) http.Response {
	var user models.User
	_ = facades.Auth(ctx).User(&user)

	slug := ctx.Request().Route("id")
	var client models.Client
	if err := facades.Orm().Query().Where("slug", slug).First(&client); err != nil || client.ID == 0 {
		return ctx.Response().Redirect(302, "/clients")
	}

	return ctx.Response().View().Make("clients/show.tmpl", map[string]any{
		"user":       user,
		"client":     client,
		"activePage": "clients",
	})
}

// Update POST /clients/:id/update
func (c *ClientController) Update(ctx http.Context) http.Response {
	slug := ctx.Request().Route("id")
	var client models.Client
	if err := facades.Orm().Query().Where("slug", slug).First(&client); err != nil || client.ID == 0 {
		return ctx.Response().Redirect(302, "/clients")
	}

	client.Name = ctx.Request().Input("name")
	client.Email = ctx.Request().Input("email")
	client.Phone = ctx.Request().Input("phone")
	client.Company = ctx.Request().Input("company")
	client.Status = ctx.Request().Input("status")
	client.Notes = ctx.Request().Input("notes")

	facades.Orm().Query().Save(&client)

	return ctx.Response().Redirect(302, "/clients/"+slug)
}

// Destroy POST /clients/:id/delete
func (c *ClientController) Destroy(ctx http.Context) http.Response {
	slug := ctx.Request().Route("id")
	var client models.Client
	if err := facades.Orm().Query().Where("slug", slug).First(&client); err == nil && client.ID != 0 {
		facades.Orm().Query().Delete(&client)
	}
	return ctx.Response().Redirect(302, "/clients")
}
