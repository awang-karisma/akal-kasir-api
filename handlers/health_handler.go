package handlers

import (
	"database/sql"
	"kasir-api/internal"
	"net/http"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetHealth(w)
	default:
		internal.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *HealthHandler) GetHealth(w http.ResponseWriter) {
	err := h.db.Ping()
	if err != nil {
		internal.HandleError(w, http.StatusInternalServerError, "Database connection error")
		return
	}
	internal.HandleResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Kasir API is running",
	})
}
