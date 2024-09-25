package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time-management/internal/employee/application/command"
	"time-management/internal/employee/application/query"
	"time-management/internal/employee/domain"
	"time-management/internal/shared/util"
)

type EmployeeHandler struct {
	CreateEmployeeHandler command.CreateEmployeeHandler
	DeleteEmployeeHandler command.DeleteEmployeeHandler
	UpdateEmailHandler    command.UpdateEmailHandler
	UpdateEmployeeHandler command.UpdateEmployeeHandler
	UpdatePasswordHandler command.UpdatePasswordHandler
	ToggleStatusHandler   command.ToggleStatusHandler
	GetEmployeesHandler   query.GetEmployeesHandler
	GetEmployeeHandler    query.GetEmployeeHandler
}

func NewEmployeeHandler(repository domain.EmployeeRepository) *EmployeeHandler {
	return &EmployeeHandler{
		CreateEmployeeHandler: command.CreateEmployeeHandler{Repo: repository},
		DeleteEmployeeHandler: command.DeleteEmployeeHandler{Repo: repository},
		UpdateEmailHandler:    command.UpdateEmailHandler{Repo: repository},
		UpdateEmployeeHandler: command.UpdateEmployeeHandler{Repo: repository},
		UpdatePasswordHandler: command.UpdatePasswordHandler{Repo: repository},
		ToggleStatusHandler:   command.ToggleStatusHandler{Repo: repository},
		GetEmployeesHandler:   query.GetEmployeesHandler{Repo: repository},
		GetEmployeeHandler:    query.GetEmployeeHandler{Repo: repository},
	}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}

	employee, err := h.CreateEmployeeHandler.Handle(command.CreateEmployeeCommand{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		return util.HandleError(w, err, http.StatusBadRequest)
	}

	return util.WriteJson(w, http.StatusCreated, employee)
}

func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) error {
	employees, err := h.GetEmployeesHandler.Handle()
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: util.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusOK, employees)
}

func (h *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	employee, err := h.GetEmployeeHandler.Handle(query.GetEmployeeQuery{Id: id})
	if err != nil {
		return util.HandleError(w, err, http.StatusNotFound)
	}

	return util.WriteJson(w, http.StatusOK, employee)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}

	updatedEmployee, err := h.UpdateEmployeeHandler.Handle(command.UpdateEmployeeCommand{
		Id:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		return util.HandleError(w, err, http.StatusBadRequest)
	}

	return util.WriteJson(w, http.StatusOK, updatedEmployee)
}

func (h *EmployeeHandler) ChangePassword(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}

	err := h.UpdatePasswordHandler.Handle(command.UpdatePasswordCommand{
		Id:       id,
		Password: req.Password,
	})
	if err != nil {
		return util.HandleError(w, err, http.StatusBadRequest)
	}

	return util.WriteJson(w, http.StatusOK, nil)
}

func (h *EmployeeHandler) ChangeEmail(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}

	err := h.UpdateEmailHandler.Handle(command.UpdateEmailCommand{
		Id:    id,
		Email: req.Email,
	})
	if err != nil {
		return util.HandleError(w, err, http.StatusBadRequest)
	}

	return util.WriteJson(w, http.StatusOK, nil)
}

func (h *EmployeeHandler) ToggleEmployeeStatus(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var req struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}

	status, err := h.ToggleStatusHandler.Handle(command.ToggleStatusCommand{
		Id:     id,
		Active: req.Active,
	})
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: util.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusOK, status)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := h.DeleteEmployeeHandler.Handle(command.DeleteEmployeeCommand{Id: id})
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: util.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusNoContent, nil)
}
