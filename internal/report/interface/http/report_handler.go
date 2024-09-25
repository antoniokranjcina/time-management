package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time-management/internal/report/application/command"
	"time-management/internal/report/application/query"
	"time-management/internal/report/infrastructure/repository"
	"time-management/internal/shared/util"
)

type ReportHandler struct {
	CreateReportHandler        command.CreateReportHandler
	GetReportsHandler          query.GetReportsHandler
	GetPendingReportsHandler   query.GetPendingReportsHandler
	GetDeniedReportsHandler    query.GetDeniedReportsHandler
	GetReportHandler           query.GetReportHandler
	GetPendingReportHandler    query.GetPendingReportHandler
	GetDeniedReportHandler     query.GetDeniedReportHandler
	UpdatePendingReportHandler command.UpdatePendingReportHandler
	ApproveReportHandler       command.ApproveReportHandler
	DenyReportHandler          command.DenyReportHandler
	DeleteReportHandler        command.DeleteReportHandler
}

func NewReportHandler(repository *repository.PgReportRepository) *ReportHandler {
	return &ReportHandler{
		CreateReportHandler:        command.CreateReportHandler{Repo: repository},
		GetPendingReportsHandler:   query.GetPendingReportsHandler{Repo: repository},
		GetDeniedReportsHandler:    query.GetDeniedReportsHandler{Repo: repository},
		GetReportsHandler:          query.GetReportsHandler{Repo: repository},
		GetReportHandler:           query.GetReportHandler{Repo: repository},
		GetPendingReportHandler:    query.GetPendingReportHandler{Repo: repository},
		GetDeniedReportHandler:     query.GetDeniedReportHandler{Repo: repository},
		UpdatePendingReportHandler: command.UpdatePendingReportHandler{Repo: repository},
		ApproveReportHandler:       command.ApproveReportHandler{Repo: repository},
		DenyReportHandler:          command.DenyReportHandler{Repo: repository},
		DeleteReportHandler:        command.DeleteReportHandler{Repo: repository},
	}
}

func (h *ReportHandler) CreateReport(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		EmployeeId       string `json:"employee_id"`
		LocationId       string `json:"location_id"`
		WorkingHours     int64  `json:"working_hours"`
		MaintenanceHours int64  `json:"maintenance_hours"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}
	if req.WorkingHours < 0 || req.MaintenanceHours < 0 {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: "hours cannot be negative"})
	}

	report, err := h.CreateReportHandler.Handle(command.CreateReportCommand{
		EmployeeId:       req.EmployeeId,
		LocationId:       req.LocationId,
		WorkingHours:     uint64(req.WorkingHours),
		MaintenanceHours: uint64(req.MaintenanceHours),
	})
	if err != nil {
		return util.HandleError(w, err, http.StatusBadRequest)
	}

	return util.WriteJson(w, http.StatusCreated, report)
}

func (h *ReportHandler) GetReports(w http.ResponseWriter, r *http.Request) error {
	reports, err := h.GetReportsHandler.Handle()
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: util.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusOK, reports)
}

func (h *ReportHandler) GetPendingReports(w http.ResponseWriter, r *http.Request) error {
	reports, err := h.GetPendingReportsHandler.Handle()
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: util.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusOK, reports)
}

func (h *ReportHandler) GetDeniedReports(w http.ResponseWriter, r *http.Request) error {
	reports, err := h.GetDeniedReportsHandler.Handle()
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: util.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusOK, reports)
}

func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	report, err := h.GetReportHandler.Handle(query.GetReportQuery{Id: id})
	if err != nil {
		return util.HandleError(w, err, http.StatusNotFound)
	}

	return util.WriteJson(w, http.StatusOK, report)
}

func (h *ReportHandler) GetPendingReport(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	report, err := h.GetPendingReportHandler.Handle(query.GetPendingReportQuery{Id: id})
	if err != nil {
		return util.HandleError(w, err, http.StatusNotFound)
	}

	return util.WriteJson(w, http.StatusOK, report)
}

func (h *ReportHandler) GetDeniedReport(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	report, err := h.GetDeniedReportHandler.Handle(query.GetDeniedReportQuery{Id: id})
	if err != nil {
		return util.HandleError(w, err, http.StatusNotFound)
	}

	return util.WriteJson(w, http.StatusOK, report)
}

func (h *ReportHandler) UpdatePendingReport(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var req struct {
		LocationId       string `json:"location_id"`
		WorkingHours     int64  `json:"working_hours"`
		MaintenanceHours int64  `json:"maintenance_hours"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: err.Error()})
	}
	if req.WorkingHours < 0 || req.MaintenanceHours < 0 {
		return util.WriteJson(w, http.StatusBadRequest, util.ApiError{Error: "hours cannot be negative"})
	}

	updatedReport, err := h.UpdatePendingReportHandler.Handle(command.UpdatePendingReportCommand{
		Id:               id,
		LocationId:       req.LocationId,
		WorkingHours:     uint64(req.WorkingHours),
		MaintenanceHours: uint64(req.MaintenanceHours),
	})
	if err != nil {
		return util.HandleError(w, err, http.StatusBadRequest)
	}

	return util.WriteJson(w, http.StatusOK, updatedReport)
}

func (h *ReportHandler) ApproveReport(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := h.ApproveReportHandler.Handle(command.ApproveReportCommand{Id: id})
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: util.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusOK, nil)
}

func (h *ReportHandler) DenyReport(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := h.DenyReportHandler.Handle(command.DenyReportCommand{Id: id})
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: util.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusOK, nil)
}

func (h *ReportHandler) DeleteReport(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := h.DeleteReportHandler.Handle(command.DeleteReport{Id: id})
	if err != nil {
		return util.WriteJson(w, http.StatusInternalServerError, util.ApiError{Error: util.ErrInternalServer.Error()})
	}

	return util.WriteJson(w, http.StatusOK, nil)
}
