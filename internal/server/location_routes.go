package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func (s *Server) locationsHandler(w http.ResponseWriter, r *http.Request) error {
	locations, err := s.db.GetLocations()
	if err != nil {
		log.Println("Error fetching locations:", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Error fetching locations"})
	}

	return WriteJson(w, http.StatusOK, locations)
}

func (s *Server) locationHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	location, err := s.db.GetLocationById(id)
	if err != nil {
		log.Println("Error fetching location:", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Error fetching location"})
	}
	return WriteJson(w, http.StatusOK, location)
}

func (s *Server) createLocationHandler(w http.ResponseWriter, r *http.Request) error {
	var requestData struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.Name == "" {
		log.Println("Invalid request body:", err)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request body"})
	}

	loc, err := s.db.CreateLocation(requestData.Name)
	if err != nil {
		log.Println("Error creating location:", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Error creating location"})
	}

	return WriteJson(w, http.StatusCreated, loc)
}

func (s *Server) deleteLocationHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := s.db.DeleteLocationById(id)
	if err != nil {
		log.Println("Error deleting location:", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Error deleting location"})
	}
	return WriteJson(w, http.StatusNoContent, nil)
}

func (s *Server) updateLocationHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var requestData struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.Name == "" {
		log.Println("Invalid request body:", err)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request body"})
	}

	loc, err := s.db.UpdateLocation(id, requestData.Name)
	if err != nil {
		log.Println("Error updating location:", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Error updating location"})
	}

	return WriteJson(w, http.StatusOK, loc)
}
