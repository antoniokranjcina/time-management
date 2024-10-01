package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	locHttp "time-management/internal/location/interface/http"
	repHttp "time-management/internal/report/interface/http"
	appMiddleware "time-management/internal/shared/middleware"
	"time-management/internal/shared/util"
	userHttp "time-management/internal/user/interface/http"
	"time-management/internal/user/role"
	adminHttp "time-management/internal/user/role/admin/interface/http"
	empHttp "time-management/internal/user/role/employee/interface/http"
)

func SetupRoutes(
	locationHandler *locHttp.LocationHandler,
	userHandler *userHttp.UserHandler,
	adminHandler *adminHttp.AdminHandler,
	employeeHandler *empHttp.EmployeeHandler,
	reportHandler *repHttp.ReportHandler,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/login", util.HttpHandler(userHandler.LoginUser))
	r.Post("/logout", util.HttpHandler(userHandler.LogoutUser))

	r.With(appMiddleware.AuthMiddleware).Group(func(r chi.Router) {
		r.Route("/locations", func(r chi.Router) {
			r.With(Role(role.Manager)).
				Post("/", util.HttpHandler(locationHandler.CreateLocation))
			r.With(Role(role.Manager, role.Employee)).
				Get("/", util.HttpHandler(locationHandler.GetLocations))
			r.With(Role(role.Manager, role.Employee)).
				Get("/{id}", util.HttpHandler(locationHandler.GetLocation))
			r.With(Role(role.Manager)).
				Put("/{id}", util.HttpHandler(locationHandler.UpdateLocation))
			r.With(Role()).
				Delete("/{id}", util.HttpHandler(locationHandler.DeleteLocation))
		})
		r.Route("/employees", func(r chi.Router) {
			r.With(Role(role.Manager)).
				Post("/", util.HttpHandler(employeeHandler.CreateEmployee))
			r.With(Role(role.Manager)).
				Get("/", util.HttpHandler(employeeHandler.GetEmployees))
			r.With(Role(role.Manager)).
				Get("/{id}", util.HttpHandler(employeeHandler.GetEmployee))
			r.With(Role()).
				Put("/{id}", util.HttpHandler(employeeHandler.UpdateEmployee))
			r.Route("/password", func(r chi.Router) {
				r.With().
					Patch("/{id}", util.HttpHandler(employeeHandler.ChangePassword))
				r.With(Role(role.Employee)).
					Patch("/", util.HttpHandler(employeeHandler.ChangePassword))
			})
			r.Route("/email", func(r chi.Router) {
				r.With().
					Patch("/{id}", util.HttpHandler(employeeHandler.ChangeEmail))
				r.With(Role(role.Employee)).
					Patch("/", util.HttpHandler(employeeHandler.ChangeEmail))
			})
			r.With(Role(role.Manager)).
				Patch("/{id}/status", util.HttpHandler(employeeHandler.ToggleEmployeeStatus))
			r.With(Role()).
				Delete("/{id}", util.HttpHandler(employeeHandler.DeleteEmployee))
		})
		r.Route("/admins", func(r chi.Router) {
			r.With(Role()).
				Post("/", util.HttpHandler(adminHandler.CreateAdmin))
			r.With(Role()).
				Get("/", util.HttpHandler(adminHandler.GetAdmins))
			r.With(Role()).
				Get("/{id}", util.HttpHandler(adminHandler.GetAdminById))
			r.With(Role()).
				Put("/{id}", util.HttpHandler(adminHandler.UpdateAdmin))
			r.With(Role()).
				Delete("/{id}", util.HttpHandler(adminHandler.DeleteAdmin))
		})
		r.Route("/reports", func(r chi.Router) {
			r.With(Role(role.Employee)).
				Post("/", util.HttpHandler(reportHandler.CreateReport))
			r.With(Role(role.Manager)).
				Get("/", util.HttpHandler(reportHandler.GetReports))
			r.With(Role(role.Manager)).
				Get("/{id}", util.HttpHandler(reportHandler.GetReport))
			r.With(Role(role.Employee, role.Manager)).
				Get("/pending", util.HttpHandler(reportHandler.GetPendingReports))
			r.With(Role(role.Employee, role.Manager)).
				Get("/pending/{id}", util.HttpHandler(reportHandler.GetPendingReport))
			r.With(Role(role.Employee, role.Manager)).
				Put("/pending/{id}", util.HttpHandler(reportHandler.UpdatePendingReport))
			r.With(Role(role.Manager)).
				Get("/denied", util.HttpHandler(reportHandler.GetDeniedReports))
			r.With(Role(role.Manager)).
				Get("/denied/{id}", util.HttpHandler(reportHandler.GetDeniedReport))
			r.With(Role(role.Manager)).
				Patch("/{id}/approve", util.HttpHandler(reportHandler.ApproveReport))
			r.With(Role(role.Manager)).
				Patch("/{id}/deny", util.HttpHandler(reportHandler.DenyReport))
			r.With(Role()).
				Delete("/{id}", util.HttpHandler(reportHandler.DeleteReport))
		})
	})

	return r
}

func Role(roles ...role.Role) func(next http.Handler) http.Handler {
	return appMiddleware.RoleMiddleware(role.Strings(roles...)...)
}
