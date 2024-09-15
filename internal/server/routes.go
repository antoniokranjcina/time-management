package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/employees", func(r chi.Router) {
		r.Get("/", s.employeesHandler)
		r.Get("/{id}", s.employeeHandler)
		r.Post("/", s.addEmployeeHandler)
		r.Delete("/{id}", s.deleteEmployeeHandler)
	})

	r.Route("/locations", func(r chi.Router) {
		r.Get("/", MakeHttpHandleFunc(s.locationsHandler))
		r.Get("/{id}", MakeHttpHandleFunc(s.locationHandler))
		r.Post("/", MakeHttpHandleFunc(s.createLocationHandler))
		r.Delete("/{id}", MakeHttpHandleFunc(s.deleteLocationHandler))
		r.Put("/{id}", MakeHttpHandleFunc(s.updateLocationHandler))
	})

	r.Route("/reports", func(r chi.Router) {
		r.Get("/", s.reportsHandler)
		r.Post("/", s.addReportForApproval)
		r.Post("/unapproved", s.unapprovedReportsHandler)
		r.Put("/{id}/approval", s.approveReportHandler)
	})

	return r
}

// Employees
func (s *Server) employeesHandler(w http.ResponseWriter, r *http.Request)      {}
func (s *Server) employeeHandler(w http.ResponseWriter, r *http.Request)       {}
func (s *Server) addEmployeeHandler(w http.ResponseWriter, r *http.Request)    {}
func (s *Server) deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {}

// Reports
func (s *Server) reportsHandler(w http.ResponseWriter, r *http.Request)           {}
func (s *Server) addReportForApproval(w http.ResponseWriter, r *http.Request)     {}
func (s *Server) unapprovedReportsHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) approveReportHandler(w http.ResponseWriter, r *http.Request)     {}
