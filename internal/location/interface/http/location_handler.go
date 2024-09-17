package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time-management/internal/location/application/commands"
	"time-management/internal/location/application/queries"
	"time-management/internal/location/domain"
	"time-management/internal/shared/util"
)

type LocationHandler struct {
	CreateLocationHandler commands.CreateLocationHandler
	GetLocationsHandler   queries.GetLocationsHandler
	GetLocationHandler    queries.GetLocationHandler
	UpdateLocationHandler commands.UpdateLocationHandler
	DeleteLocationHandler commands.DeleteLocationHandler
}

func NewLocationHandler(repository domain.LocationRepository) *LocationHandler {
	return &LocationHandler{
		CreateLocationHandler: commands.CreateLocationHandler{Repo: repository},
		GetLocationsHandler:   queries.GetLocationsHandler{Repo: repository},
		GetLocationHandler:    queries.GetLocationHandler{Repo: repository},
		UpdateLocationHandler: commands.UpdateLocationHandler{Repo: repository},
		DeleteLocationHandler: commands.DeleteLocationHandler{Repo: repository},
	}
}

func (h *LocationHandler) CreateLocation(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}

	location, err := h.CreateLocationHandler.Handle(commands.CreateLocationCommand{Name: req.Name})
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: err.Error()})
	}

	return util.WriteJson(w, http.StatusCreated, location)
}

func (h *LocationHandler) GetLocations(w http.ResponseWriter, r *http.Request) error {
	locations, err := h.GetLocationsHandler.Handle()
	if err != nil {
		log.Println(err)
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: err.Error()})
	}

	return util.WriteJson(w, http.StatusOK, locations)
}

func (h *LocationHandler) GetLocation(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	location, err := h.GetLocationHandler.Handle(queries.GetLocationQuery{Id: id})
	if err != nil {
		log.Println(err)
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: err.Error()})
	}

	return util.WriteJson(w, http.StatusOK, location)
}

func (h *LocationHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var requestData struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.Name == "" {
		log.Println("Invalid request body:", err)
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: "Invalid request body"})
	}

	location, err := h.UpdateLocationHandler.Handle(commands.UpdateLocationCommand{Id: id, Name: requestData.Name})
	if err != nil {
		log.Println("Error updating location:", err)
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: err.Error()})
	}

	return util.WriteJson(w, http.StatusOK, location)
}

func (h *LocationHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := h.DeleteLocationHandler.Handle(commands.DeleteLocationCommand{ID: id})
	if err != nil {
		log.Println(err)
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: "Error deleting location"})
	}
	return util.WriteJson(w, http.StatusNoContent, nil)
}
