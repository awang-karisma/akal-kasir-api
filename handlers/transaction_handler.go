package handlers

import (
	"encoding/json"
	"kasir-api/internal"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) TransactionHandler {
	return TransactionHandler{service: service}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req []models.CheckoutItem
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		internal.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	checkout, err := h.service.Checkout(req)
	if err != nil {
		internal.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	internal.HandleResponse(w, http.StatusOK, checkout)
}

func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		internal.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
