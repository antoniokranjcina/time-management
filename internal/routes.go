package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	http2 "time-management/internal/locations/interfaces/http"
	"time-management/internal/shared/util"
)

func SetupRoutes(locationHandler *http2.LocationHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/locations", func(r chi.Router) {
		r.Get("/", util.HttpHandler(locationHandler.GetLocations))
		r.Get("/{id}", util.HttpHandler(locationHandler.GetLocation))
		r.Post("/", util.HttpHandler(locationHandler.CreateLocation))
		r.Put("/{id}", util.HttpHandler(locationHandler.UpdateLocation))
		r.Delete("/{id}", util.HttpHandler(locationHandler.DeleteLocation))
	})

	return r
}
