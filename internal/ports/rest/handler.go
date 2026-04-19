package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/R-accoo-n/opog-lab3/internal"
)

type ProductHandler struct {
	service internal.Products
}

func NewProductHandler(service internal.Products) ProductHandler {
	return ProductHandler{service: service}
}

func (h ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProduct(w, r)
	case http.MethodPost:
		h.CreateProduct(w, r)
	default:
		http.Error(w, fmt.Sprintf("method %s is not supported", r.Method), http.StatusMethodNotAllowed)
	}
}

func (h ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateProduct(r.Context(), internal.CreateProductPayload{
		Name: req.Name,
		Category: internal.Category{
			Name: req.Category.Name,
			Tax:  req.Category.Tax,
		},
		Price: req.Price,
	})
	if err != nil {
		if errors.Is(err, internal.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id.String()})
}

func (h ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	rawID := r.URL.Query().Get("id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	product, finalPrice, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, internal.ErrNoResource) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, internal.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ProductResponse{
		ID:   product.ID.String(),
		Name: product.Name,
		Category: CategoryPayload{
			Name: product.Category.Name,
			Tax:  product.Category.Tax,
		},
		Price:      product.Price,
		FinalPrice: finalPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
