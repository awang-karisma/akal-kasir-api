package handlers

import (
	"kasir-api/internal"
	"kasir-api/services"
	"net/http"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from != "" && !internal.IsDateValid(from) {
		internal.HandleError(w, http.StatusBadRequest, "Invalid from date")
		return
	}

	if to != "" && !internal.IsDateValid(to) {
		internal.HandleError(w, http.StatusBadRequest, "Invalid to date")
		return
	}

	report, err := h.service.GetReportsRange(from, to)
	if err != nil {
		internal.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	internal.HandleResponse(w, http.StatusOK, report)

}

func (h *ReportHandler) GetReportToday(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetReportToday()
	if err != nil {
		internal.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	internal.HandleResponse(w, http.StatusOK, report)
}

func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetReports(w, r)
	default:
		internal.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *ReportHandler) HandleReportToday(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetReportToday(w, r)
	default:
		internal.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
