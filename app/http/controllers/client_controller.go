package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

type ClientController struct{}

func NewClientController() *ClientController {
	return &ClientController{}
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
	client := models.Client{
		Name:    ctx.Request().Input("name"),
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

	id := ctx.Request().Route("id")
	var client models.Client
	if err := facades.Orm().Query().Find(&client, id); err != nil || client.ID == 0 {
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
	id := ctx.Request().Route("id")
	var client models.Client
	if err := facades.Orm().Query().Find(&client, id); err != nil || client.ID == 0 {
		return ctx.Response().Redirect(302, "/clients")
	}

	client.Name = ctx.Request().Input("name")
	client.Email = ctx.Request().Input("email")
	client.Phone = ctx.Request().Input("phone")
	client.Company = ctx.Request().Input("company")
	client.Status = ctx.Request().Input("status")
	client.Notes = ctx.Request().Input("notes")

	facades.Orm().Query().Save(&client)

	return ctx.Response().Redirect(302, "/clients/"+id)
}

// Destroy POST /clients/:id/delete
func (c *ClientController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var client models.Client
	if err := facades.Orm().Query().Find(&client, id); err == nil && client.ID != 0 {
		facades.Orm().Query().Delete(&client)
	}
	return ctx.Response().Redirect(302, "/clients")
}
