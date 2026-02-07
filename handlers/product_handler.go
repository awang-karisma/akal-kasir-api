package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"kasir-api/internal"
	"kasir-api/models"
	"kasir-api/services"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) ProductHandler {
	return ProductHandler{service: service}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	products, err := h.service.GetProducts(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err))
		return
	}
	newProduct, err := h.service.CreateProduct(product)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %s", err))
		return
	}
	internal.HandleResponse(w, http.StatusCreated, newProduct)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid uuid")
		return
	}

	product, err := h.service.GetProductByID(id.String())
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if product.ID == "" {
		internal.HandleError(w, http.StatusNotFound, "Product not found")
		return
	}

	internal.HandleResponse(w, http.StatusOK, product)
}

func (h *ProductHandler) UpdateProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid uuid")
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err = h.service.UpdateProductByID(id.String(), product)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if product.ID == "" {
		internal.HandleError(w, http.StatusNotFound, "Product not found")
		return
	}

	internal.HandleResponse(w, http.StatusOK, product)
}

func (h *ProductHandler) DeleteProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid uuid")
		return
	}

	deletedProduct, err := h.service.DeleteProductByID(id.String())
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if deletedProduct.ID == "" {
		internal.HandleError(w, http.StatusNotFound, "Product not found")
		return
	}

	internal.HandleResponse(w, http.StatusOK, deletedProduct)
}

func (h *ProductHandler) HandleProduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProducts(w, r)
	case http.MethodPost:
		h.CreateProduct(w, r)
	default:
		internal.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProductByID(w, r)
	case http.MethodPut:
		h.UpdateProductByID(w, r)
	case http.MethodDelete:
		h.DeleteProductByID(w, r)
	default:
		internal.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *ProductHandler) AddCategoryToProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid product uuid")
		return
	}

	var req struct {
		CategoryID string `json:"category_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	_, err = uuid.Parse(req.CategoryID)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid category uuid")
		return
	}

	err = h.service.AddCategoryToProduct(id.String(), req.CategoryID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			internal.HandleError(w, http.StatusNotFound, "Product or category not found")
			return
		}
		internal.HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	internal.HandleResponse(w, http.StatusCreated, map[string]string{"message": "Category added to product successfully"})
}

func (h *ProductHandler) RemoveCategoryFromProduct(w http.ResponseWriter, r *http.Request) {
	log.Print(mux.Vars(r))
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid product uuid")
		return
	}

	var req struct {
		CategoryID string `json:"category_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	_, err = uuid.Parse(req.CategoryID)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid category uuid")
		return
	}

	err = h.service.RemoveCategoryFromProduct(id.String(), req.CategoryID)
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	internal.HandleResponse(w, http.StatusOK, map[string]string{"message": "Category removed from product successfully"})
}

func (h *ProductHandler) GetProductCategories(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusBadRequest, "Invalid product uuid")
		return
	}

	// Check if product exists
	product, err := h.service.GetProductByID(id.String())
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if product.ID == "" {
		internal.HandleError(w, http.StatusNotFound, "Product not found")
		return
	}

	categories, err := h.service.GetCategoriesByProductID(id.String())
	if err != nil {
		log.Println(err)
		internal.HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	internal.HandleResponse(w, http.StatusOK, categories)
}

func (h *ProductHandler) HandleProductCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProductCategories(w, r)
	case http.MethodPost:
		h.AddCategoryToProduct(w, r)
	case http.MethodDelete:
		h.RemoveCategoryFromProduct(w, r)
	default:
		internal.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
