package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	locHttp "time-management/internal/location/interface/http"
	repHttp "time-management/internal/report/interface/http"
	"time-management/internal/shared/util"
	adminHttp "time-management/internal/user/role/admin/interface/http"
	empHttp "time-management/internal/user/role/employee/interface/http"
)

func SetupRoutes(
	locationHandler *locHttp.LocationHandler,
	adminHandler *adminHttp.AdminHandler,
	employeeHandler *empHttp.EmployeeHandler,
	reportHandler *repHttp.ReportHandler,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/locations", func(r chi.Router) {
		r.Post("/", util.HttpHandler(locationHandler.CreateLocation))
		r.Get("/", util.HttpHandler(locationHandler.GetLocations))
		r.Get("/{id}", util.HttpHandler(locationHandler.GetLocation))
		r.Put("/{id}", util.HttpHandler(locationHandler.UpdateLocation))
		r.Delete("/{id}", util.HttpHandler(locationHandler.DeleteLocation))
	})

	r.Route("/employees", func(r chi.Router) {
		r.Post("/", util.HttpHandler(employeeHandler.CreateEmployee))
		r.Get("/", util.HttpHandler(employeeHandler.GetEmployees))
		r.Get("/{id}", util.HttpHandler(employeeHandler.GetEmployee))
		r.Put("/{id}", util.HttpHandler(employeeHandler.UpdateEmployee))
		r.Patch("/{id}/password", util.HttpHandler(employeeHandler.ChangePassword))
		r.Patch("/{id}/email", util.HttpHandler(employeeHandler.ChangeEmail))
		r.Patch("/{id}/status", util.HttpHandler(employeeHandler.ToggleEmployeeStatus))
		r.Delete("/{id}", util.HttpHandler(employeeHandler.DeleteEmployee))
	})

	r.Route("/admins", func(r chi.Router) {
		r.Post("/", util.HttpHandler(adminHandler.CreateAdmin))
		r.Get("/", util.HttpHandler(adminHandler.GetAdmins))
		r.Get("/{id}", util.HttpHandler(adminHandler.GetAdminById))
		r.Put("/{id}", util.HttpHandler(adminHandler.UpdateAdmin))
		r.Delete("/{id}", util.HttpHandler(adminHandler.DeleteAdmin))
	})

	r.Route("/reports", func(r chi.Router) {
		r.Post("/", util.HttpHandler(reportHandler.CreateReport))
		r.Get("/", util.HttpHandler(reportHandler.GetReports))
		r.Get("/{id}", util.HttpHandler(reportHandler.GetReport))
		r.Get("/pending", util.HttpHandler(reportHandler.GetPendingReports))
		r.Get("/pending/{id}", util.HttpHandler(reportHandler.GetPendingReport))
		r.Put("/pending/{id}", util.HttpHandler(reportHandler.UpdatePendingReport))
		r.Get("/denied", util.HttpHandler(reportHandler.GetDeniedReports))
		r.Get("/denied/{id}", util.HttpHandler(reportHandler.GetDeniedReport))
		r.Patch("/{id}/approve", util.HttpHandler(reportHandler.ApproveReport))
		r.Patch("/{id}/deny", util.HttpHandler(reportHandler.DenyReport))
		r.Delete("/{id}", util.HttpHandler(reportHandler.DeleteReport))
	})

	return r
}
