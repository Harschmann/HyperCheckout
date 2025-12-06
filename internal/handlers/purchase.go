package handlers

import (
	"encoding/json"
	"log"
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

	err := h.Repo.PurchaseProduct(r.Context(), req.UserID, req.ProductID, req.Quantity)

	if err != nil {
		if err == repository.ErrOutofStock {
			http.Error(w, "Product is out of stock", http.StatusConflict)
		} else {
			// ---------------------------------------------------------
			// NEW: Print the actual error to the terminal!
			// ---------------------------------------------------------
			log.Printf("❌ Transaction Failed: %v", err)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Order placed successfully"}`))
}
