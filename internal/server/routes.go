package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Employees
	r.Get("/employees", s.employeesHandler)
	r.Get("/employees/{id}", s.employeeHandler)
	r.Post("/employees", s.addEmployeeHandler)
	r.Delete("/employees/{id}", s.deleteEmployeeHandler)

	// Locations
	r.Get("/locations", s.locationsHandler)
	r.Post("/locations", s.addLocationHandler)
	r.Delete("/locations/{id}", s.deleteLocationHandler)

	// Reports
	r.Get("/reports", s.reportsHandler)
	r.Post("/reports", s.addReportForApproval)
	r.Get("/reports/unapproved", s.unapprovedReportsHandler)
	r.Put("/reports/{id}/approval", s.approveReportHandler)

	return r
}

// Employees
func (s *Server) employeesHandler(w http.ResponseWriter, r *http.Request)      {}
func (s *Server) employeeHandler(w http.ResponseWriter, r *http.Request)       {}
func (s *Server) addEmployeeHandler(w http.ResponseWriter, r *http.Request)    {}
func (s *Server) deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {}

// Locations
func (s *Server) locationsHandler(w http.ResponseWriter, r *http.Request)      {}
func (s *Server) addLocationHandler(w http.ResponseWriter, r *http.Request)    {}
func (s *Server) deleteLocationHandler(w http.ResponseWriter, r *http.Request) {}

// Reports
func (s *Server) reportsHandler(w http.ResponseWriter, r *http.Request)           {}
func (s *Server) addReportForApproval(w http.ResponseWriter, r *http.Request)     {}
func (s *Server) unapprovedReportsHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) approveReportHandler(w http.ResponseWriter, r *http.Request)     {}
