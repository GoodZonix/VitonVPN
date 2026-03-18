package handlers

import (
	"encoding/json"
	"math"
	"net/http"

	"vpn-startup/backend/internal/repository"
	"vpn-startup/backend/internal/vless"
)

type ConfigHandler struct {
	ServerRepo *repository.ServerRepo
	VPNKeyRepo *repository.VPNKeyRepo
	DeviceRepo *repository.DeviceRepo
	WalletRepo *repository.WalletRepo
	MaxDevices int
}

type ConfigResponse struct {
	VlessURLs []string         `json:"vless_urls"`
	Servers   []ServerListItem `json:"servers"`
	Wallet    *struct {
		Balance        float64 `json:"balance"`
		ApproxDaysLeft float64 `json:"approx_days_left"`
		Active         bool    `json:"active"`
	} `json:"wallet,omitempty"`
}

func (h *ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromRequest(r)
	if claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	userID := claims.UserID

	// Charge for usage since last_billing_at and get updated balance.
	user, allowed, err := h.WalletRepo.ChargeForUsage(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}

	servers, err := h.ServerRepo.ListActive(r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}

	out := make([]ServerListItem, 0, len(servers))
	for _, s := range servers {
		out = append(out, ServerListItem{
			ID: s.ID.String(), Name: s.Name, Region: s.Region, Host: s.Host, Port: s.Port, Type: s.Type, IsActive: s.IsActive,
		})
	}

	resp := ConfigResponse{
		VlessURLs: []string{},
		Servers:   out,
	}

	// Approximate days left: balance / (RUB per week / 7)
	const rubPerWeek = repository.RubPerWeek
	approxDays := 0.0
	if rubPerWeek > 0 {
		perDay := rubPerWeek / 7.0
		if perDay > 0 {
			approxDays = user.WalletBalance / perDay
		}
	}
	approxDays = math.Round(approxDays*10) / 10 // 1 знак после запятой

	resp.Wallet = &struct {
		Balance        float64 `json:"balance"`
		ApproxDaysLeft float64 `json:"approx_days_left"`
		Active         bool    `json:"active"`
	}{
		Balance:        user.WalletBalance,
		ApproxDaysLeft: approxDays,
		Active:         allowed,
	}

	if allowed {
		userUUID, err := h.VPNKeyRepo.GetOrCreate(r.Context(), userID)
		if err != nil {
			http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
			return
		}
		urls := make([]string, 0, len(servers))
		for _, s := range servers {
			urls = append(urls, vless.BuildVLESSURL(s, userUUID, s.Name))
		}
		resp.VlessURLs = urls
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
