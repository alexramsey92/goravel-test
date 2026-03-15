package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"

	"goravel/app/facades"
	"goravel/app/http/middleware"
	"goravel/app/models"
)

// Api registers JSON API routes — useful for testing with Postman / curl.
// All /api routes return JSON; auth uses JWT via Authorization: Bearer <token>.
func Api() {
	facades.Route().Prefix("/api").Group(func(router route.Router) {

		// ── Auth ─────────────────────────────────────────────────────────────
		router.Post("/auth/register", func(ctx http.Context) http.Response {
			name := ctx.Request().Input("name")
			email := ctx.Request().Input("email")
			password := ctx.Request().Input("password")

			var existing models.User
			err := facades.Orm().Query().Where("email = ?", email).First(&existing)
			if err == nil && existing.ID != 0 {
				return ctx.Response().Json(422, http.Json{"message": "Email already taken"})
			}

			hashed, err := facades.Hash().Make(password)
			if err != nil {
				return ctx.Response().Json(500, http.Json{"message": "Internal error"})
			}

			user := models.User{Name: name, Email: email, Password: hashed, Role: models.RoleUser}
			if err := facades.Orm().Query().Create(&user); err != nil {
				return ctx.Response().Json(500, http.Json{"message": "Could not create user"})
			}

			token, _ := facades.Auth(ctx).Guard("api").Login(&user)
			return ctx.Response().Success().Json(http.Json{
				"token": token,
				"user":  map[string]any{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role},
			})
		})

		router.Post("/auth/login", func(ctx http.Context) http.Response {
			email := ctx.Request().Input("email")
			password := ctx.Request().Input("password")

			var user models.User
			if err := facades.Orm().Query().Where("email = ?", email).First(&user); err != nil {
				return ctx.Response().Json(401, http.Json{"message": "Invalid credentials"})
			}

			ok := facades.Hash().Check(password, user.Password)
			if !ok {
				return ctx.Response().Json(401, http.Json{"message": "Invalid credentials"})
			}

			token, err := facades.Auth(ctx).Guard("api").Login(&user)
			if err != nil {
				return ctx.Response().Json(500, http.Json{"message": "Login failed"})
			}

			return ctx.Response().Success().Json(http.Json{
				"token": token,
				"user":  map[string]any{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role},
			})
		})

		// ── Protected API routes (JWT Bearer) ─────────────────────────────────
		apiAuth := middleware.NewApiAuth()
		router.Prefix("/v1").Middleware(apiAuth.Handle).Group(func(r route.Router) {

			r.Get("/me", func(ctx http.Context) http.Response {
				var user models.User
				_ = facades.Auth(ctx).Guard("api").User(&user)
				return ctx.Response().Success().Json(http.Json{
					"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role,
				})
			})

			// Clients
			r.Get("/clients", func(ctx http.Context) http.Response {
				var clients []models.Client
				facades.Orm().Query().OrderBy("created_at", "desc").Find(&clients)
				return ctx.Response().Success().Json(http.Json{"data": clients})
			})

			r.Post("/clients", func(ctx http.Context) http.Response {
				client := models.Client{
					Name:    ctx.Request().Input("name"),
					Email:   ctx.Request().Input("email"),
					Phone:   ctx.Request().Input("phone"),
					Company: ctx.Request().Input("company"),
					Status:  ctx.Request().Input("status"),
					Notes:   ctx.Request().Input("notes"),
				}
				if err := facades.Orm().Query().Create(&client); err != nil {
					return ctx.Response().Json(500, http.Json{"message": "Failed to create client"})
				}
				return ctx.Response().Json(201, http.Json{"data": client})
			})

			r.Get("/clients/:id", func(ctx http.Context) http.Response {
				var client models.Client
				if err := facades.Orm().Query().Find(&client, ctx.Request().Route("id")); err != nil || client.ID == 0 {
					return ctx.Response().Json(404, http.Json{"message": "Not found"})
				}
				return ctx.Response().Success().Json(http.Json{"data": client})
			})

			r.Delete("/clients/:id", func(ctx http.Context) http.Response {
				var client models.Client
				if err := facades.Orm().Query().Find(&client, ctx.Request().Route("id")); err != nil || client.ID == 0 {
					return ctx.Response().Json(404, http.Json{"message": "Not found"})
				}
				facades.Orm().Query().Delete(&client)
				return ctx.Response().Json(204, nil)
			})

			// Pages
			r.Get("/pages", func(ctx http.Context) http.Response {
				var pages []models.Page
				facades.Orm().Query().OrderBy("created_at", "desc").Find(&pages)
				return ctx.Response().Success().Json(http.Json{"data": pages})
			})

			r.Post("/pages", func(ctx http.Context) http.Response {
				page := models.Page{
					Title:     ctx.Request().Input("title"),
					Slug:      ctx.Request().Input("slug"),
					Content:   ctx.Request().Input("content"),
					Published: ctx.Request().Input("published") == "true",
				}
				if err := facades.Orm().Query().Create(&page); err != nil {
					return ctx.Response().Json(500, http.Json{"message": "Failed to create page"})
				}
				return ctx.Response().Json(201, http.Json{"data": page})
			})
		})
	})
}
