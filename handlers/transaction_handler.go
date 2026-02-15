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

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if from == "" || to == "" {
		transactions, err := h.service.GetTransactions()
		if err != nil {
			internal.HandleError(w, http.StatusInternalServerError, err.Error())
			return
		}
		internal.HandleResponse(w, http.StatusOK, transactions)
		return
	}
	if from != "" && !internal.IsDateValid(from) {
		internal.HandleError(w, http.StatusBadRequest, "Invalid from date")
		return
	}
	if to != "" && !internal.IsDateValid(to) {
		internal.HandleError(w, http.StatusBadRequest, "Invalid to date")
		return
	}

	transactions, err := h.service.GetTransactionsRange(from, to)
	if err != nil {
		internal.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	internal.HandleResponse(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) HandleTransactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTransactions(w, r)
	default:
		internal.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
