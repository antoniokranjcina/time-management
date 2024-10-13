package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time-management/internal/location/application/command"
	"time-management/internal/location/application/query"
	locDomain "time-management/internal/location/domain"
	"time-management/internal/shared/util"
	"time-management/internal/user/domain"
)

type LocationHandler struct {
	CreateLocationHandler command.CreateLocationHandler
	GetLocationsHandler   query.GetLocationsHandler
	GetLocationHandler    query.GetLocationHandler
	UpdateLocationHandler command.UpdateLocationHandler
	DeleteLocationHandler command.DeleteLocationHandler
}

func NewLocationHandler(repository locDomain.LocationRepository) *LocationHandler {
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

	location, err := h.CreateLocationHandler.Handle(r.Context(), command.CreateLocationCommand{Name: req.Name})
	if err != nil {
		return util.HandleError(w, err, http.StatusBadRequest)
	}

	return util.WriteJson(w, http.StatusCreated, location)
}

func (h *LocationHandler) GetLocations(w http.ResponseWriter, r *http.Request) error {
	locations, err := h.GetLocationsHandler.Handle(r.Context())
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: domain.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusOK, locations)
}

func (h *LocationHandler) GetLocation(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	location, err := h.GetLocationHandler.Handle(r.Context(), query.GetLocationQuery{Id: id})
	if err != nil {
		return util.HandleError(w, err, http.StatusNotFound)
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

	location, err := h.UpdateLocationHandler.Handle(
		r.Context(),
		command.UpdateLocationCommand{Id: id, Name: requestData.Name},
	)
	if err != nil {
		return util.HandleError(w, err, http.StatusBadRequest)
	}

	return util.WriteJson(w, http.StatusOK, location)
}

func (h *LocationHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := h.DeleteLocationHandler.Handle(r.Context(), command.DeleteLocationCommand{Id: id})
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: domain.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusNoContent, nil)
}
