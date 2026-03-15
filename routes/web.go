package routes

import (
	"github.com/goravel/framework/contracts/route"
	sessionmiddleware "github.com/goravel/framework/session/middleware"

	"goravel/app/facades"
	"goravel/app/http/controllers"
	"goravel/app/http/middleware"
)

func Web() {
	facades.Route().Static("public", "./public")

	auth := controllers.NewAuthController()
	dashboard := controllers.NewDashboardController()
	clients := controllers.NewClientController()
	pages := controllers.NewPageController()

	// All web routes use StartSession so facades.Auth(ctx) has a session driver.
	// Equivalent to Laravel's web middleware group.
	facades.Route().Middleware(sessionmiddleware.StartSession()).Group(func(router route.Router) {

		// ── Guest routes ──────────────────────────────────────────────────────
		router.Get("/", auth.ShowLoginForm)
		router.Get("/login", auth.ShowLoginForm)
		router.Post("/login", auth.Login)
		router.Get("/register", auth.ShowRegisterForm)
		router.Post("/register", auth.Register)
		router.Post("/logout", auth.Logout)

		// ── Auth-protected routes ──────────────────────────────────────────────
		authMiddleware := middleware.NewAuth()
		router.Middleware(authMiddleware.Handle).Group(func(r route.Router) {
			// Dashboard
			r.Get("/dashboard", dashboard.Index)

			// Clients — CRM
			r.Get("/clients", clients.Index)
			r.Get("/clients/create", clients.Create)
			r.Post("/clients", clients.Store)
			r.Get("/clients/:id", clients.Show)
			r.Post("/clients/:id/update", clients.Update)
			r.Post("/clients/:id/delete", clients.Destroy)

			// Pages — CMS
			r.Get("/pages", pages.Index)
			r.Get("/pages/create", pages.Create)
			r.Post("/pages", pages.Store)
			r.Post("/pages/:id/delete", pages.Destroy)
		})
	})
}
