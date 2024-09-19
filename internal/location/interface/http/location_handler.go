package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time-management/internal/location/application/command"
	"time-management/internal/location/application/query"
	"time-management/internal/location/domain"
	"time-management/internal/shared/util"
)

type LocationHandler struct {
	CreateLocationHandler command.CreateLocationHandler
	GetLocationsHandler   query.GetLocationsHandler
	GetLocationHandler    query.GetLocationHandler
	UpdateLocationHandler command.UpdateLocationHandler
	DeleteLocationHandler command.DeleteLocationHandler
}

func NewLocationHandler(repository domain.LocationRepository) *LocationHandler {
	return &LocationHandler{
		CreateLocationHandler: command.CreateLocationHandler{Repo: repository},
		GetLocationsHandler:   query.GetLocationsHandler{Repo: repository},
		GetLocationHandler:    query.GetLocationHandler{Repo: repository},
		UpdateLocationHandler: command.UpdateLocationHandler{Repo: repository},
		DeleteLocationHandler: command.DeleteLocationHandler{Repo: repository},
	}
}

func (h *LocationHandler) CreateLocation(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}

	location, err := h.CreateLocationHandler.Handle(command.CreateLocationCommand{Name: req.Name})
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

	location, err := h.GetLocationHandler.Handle(query.GetLocationQuery{Id: id})
	if err != nil {
		if errors.Is(err, domain.ErrLocationNotFound) {
			return util.WriteJson(w, http.StatusNotFound, util.ApiError{Error: err.Error()})
		}

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
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: "Invalid request body"})
	}

	location, err := h.UpdateLocationHandler.Handle(command.UpdateLocationCommand{Id: id, Name: requestData.Name})
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: err.Error()})
	}

	return util.WriteJson(w, http.StatusOK, location)
}

func (h *LocationHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := h.DeleteLocationHandler.Handle(command.DeleteLocationCommand{Id: id})
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: "Error deleting location"})
	}

	return util.WriteJson(w, http.StatusNoContent, nil)
}
