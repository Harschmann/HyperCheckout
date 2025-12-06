package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Harschmann/hyper-checkout/internal/repository"
)

// PurchaseRequest defines what the user sends us
type PurchaseRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Handler struct {
	Repo *repository.Store
}

func NewHandler(repo *repository.Store) *Handler {
	return &Handler{Repo: repo}
}

// HandlePurchase is the controller for POST /purchase
func (h *Handler) HandlePurchase(w http.ResponseWriter, r *http.Request) {
	var req PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 2. Call the Repository
	// We pass r.Context() so if the user disconnects, the DB query stops.
	err := h.Repo.PurchaseProduct(r.Context(), req.UserID, req.ProductID, req.Quantity)

	// 3. Handle Errors
	if err != nil {
		if err == repository.ErrOutofStock {
			http.Error(w, "Product is out of stock", http.StatusConflict) // 409 Conflict
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError) // 500 Error
		}
		return
	}

	// 4. Success Response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Order placed successfully"}`))
}
