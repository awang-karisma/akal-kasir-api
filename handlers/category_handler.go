package handlers

import (
	"encoding/json"
	"fmt"
	"kasir-api/internal"
	"kasir-api/models"
	"kasir-api/services"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) CategoryHandler {
	return CategoryHandler{service: service}
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err))
		return
	}
	newCategory, err := h.service.CreateCategory(category)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %s", err))
		return
	}
	internal.HandleResponse(w, http.StatusCreated, newCategory)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid uuid")
		return
	}

	category, err := h.service.GetCategoryByID(id.String())
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if category.ID == "" {
		internal.HandleError(w, http.StatusNotFound, "Category not found")
		return
	}

	internal.HandleResponse(w, http.StatusOK, category)
}

func (h *CategoryHandler) UpdateCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid uuid")
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category, err = h.service.UpdateCategoryByID(id.String(), category)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if category.ID == "" {
		internal.HandleError(w, http.StatusNotFound, "Category not found")
		return
	}

	internal.HandleResponse(w, http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid uuid")
		return
	}

	deletedCategory, err := h.service.DeleteCategoryByID(id.String())
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if deletedCategory.ID == "" {
		internal.HandleError(w, http.StatusNotFound, "Category not found")
		return
	}

	internal.HandleResponse(w, http.StatusOK, deletedCategory)
}

func (h *CategoryHandler) HandleCategory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetCategories(w, r)
	case http.MethodPost:
		h.CreateCategory(w, r)
	default:
		internal.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetCategoryByID(w, r)
	case http.MethodPut:
		h.UpdateCategoryByID(w, r)
	case http.MethodDelete:
		h.DeleteCategoryByID(w, r)
	default:
		internal.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
