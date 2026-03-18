package handlers

import (
	"encoding/json"
	"net/http"

	"vpn-startup/backend/internal/repository"
)

type WalletHandler struct {
	WalletRepo *repository.WalletRepo
}

type WalletTopupRequest struct {
	Amount float64 `json:"amount"` // in RUB, e.g. 100, 300, 500
}

// Get current wallet balance and approximate days left.
func (h *WalletHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromRequest(r)
	if claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	bal, _, err := h.WalletRepo.GetBalance(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	// Approx days as in ConfigHandler.
	perDay := repository.RubPerWeek / 7.0
	approxDays := 0.0
	if perDay > 0 {
		approxDays = bal / perDay
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"balance":         bal,
		"approx_days_left": approxDays,
	})
}

// Create topup intent. For now it's a placeholder that immediately credits balance.
// In production, integrate YooMoney: create payment, return payment_url, confirm via webhook.
func (h *WalletHandler) Topup(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromRequest(r)
	if claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	var req WalletTopupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		http.Error(w, `{"error":"invalid amount"}`, http.StatusBadRequest)
		return
	}
	// TODO: integrate YooMoney. For now, credit instantly for demo.
	if err := h.WalletRepo.AddTopup(r.Context(), claims.UserID, req.Amount); err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"amount": req.Amount,
	})
}
