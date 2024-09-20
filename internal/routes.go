package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	empHttp "time-management/internal/employee/interface/http"
	locHttp "time-management/internal/location/interface/http"
	"time-management/internal/shared/util"
)

func SetupRoutes(
	locationHandler *locHttp.LocationHandler,
	employeeHandler *empHttp.EmployeeHandler,
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

	return r
}
