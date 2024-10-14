package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time-management/internal/shared/util"
	"time-management/internal/user/domain"
	"time-management/internal/user/role/admin/application/command"
	"time-management/internal/user/role/admin/application/query"
)

type AdminHandler struct {
	CreateAdminHandler command.CreateAdminHandler
	GetAdminsHandler   query.GetAdminsHandler
	GetAdminHandler    query.GetAdminHandler
	UpdateAdminHandler command.UpdateAdminHandler
	DeleteAdminHandler command.DeleteAdminHandler
}

func NewAdminHandler(repository domain.UserRepository) *AdminHandler {
	return &AdminHandler{
		CreateAdminHandler: command.CreateAdminHandler{Repo: repository},
		GetAdminsHandler:   query.GetAdminsHandler{Repo: repository},
		GetAdminHandler:    query.GetAdminHandler{Repo: repository},
		UpdateAdminHandler: command.UpdateAdminHandler{Repo: repository},
		DeleteAdminHandler: command.DeleteAdminHandler{Repo: repository},
	}
}

func (h *AdminHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}

	cmd := command.CreateAdminCommand{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}
	admin, err := h.CreateAdminHandler.Handle(r.Context(), cmd)
	if err != nil {
		return util.HandleError(w, err, http.StatusBadRequest)
	}

	return util.WriteJson(w, http.StatusCreated, admin)
}

func (h *AdminHandler) GetAdmins(w http.ResponseWriter, r *http.Request) error {
	admins, err := h.GetAdminsHandler.Handle(r.Context())
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: domain.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusOK, admins)
}

func (h *AdminHandler) GetAdminById(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	admin, err := h.GetAdminHandler.Handle(r.Context(), query.GetAdminQuery{Id: id})
	if err != nil {
		return util.HandleError(w, err, http.StatusNotFound)
	}

	return util.WriteJson(w, http.StatusOK, admin)
}

func (h *AdminHandler) UpdateAdmin(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}

	cmd := command.UpdateAdminCommand{
		Id:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	updatedAdmin, err := h.UpdateAdminHandler.Handle(r.Context(), cmd)
	if err != nil {
		return util.HandleError(w, err, http.StatusBadRequest)
	}

	return util.WriteJson(w, http.StatusOK, updatedAdmin)
}

func (h *AdminHandler) DeleteAdmin(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := h.DeleteAdminHandler.Handle(r.Context(), command.DeleteAdminCommand{Id: id})
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: domain.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusNoContent, nil)
}
