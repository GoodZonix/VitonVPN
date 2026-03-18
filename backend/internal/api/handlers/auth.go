package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"vpn-startup/backend/internal/auth"
	"vpn-startup/backend/internal/config"
	"vpn-startup/backend/internal/repository"
)

type AuthHandler struct {
	UserRepo   *repository.UserRepo
	DeviceRepo *repository.DeviceRepo
	JWT        *auth.JWT
	Cfg        *config.Config
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
	DeviceName string `json:"device_name"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
	DeviceName string `json:"device_name"`
}

type AuthUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type AuthResponse struct {
	Token     string   `json:"token"`
	ExpiresIn int      `json:"expires_in"`
	User      AuthUser `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" || req.DeviceID == "" {
		http.Error(w, `{"error":"email, password, device_id required"}`, http.StatusBadRequest)
		return
	}
	user, err := h.UserRepo.GetByEmail(r.Context(), req.Email)
	if err != nil || user == nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	count, _ := h.DeviceRepo.CountByUser(r.Context(), user.ID)
	if count >= h.Cfg.MaxDevices {
		_, _ = h.DeviceRepo.Upsert(r.Context(), user.ID, req.DeviceID, req.DeviceName)
	}
	_, _ = h.DeviceRepo.Upsert(r.Context(), user.ID, req.DeviceID, req.DeviceName)
	token, err := h.JWT.Sign(user.ID, user.Email, req.DeviceID, 7*24*time.Hour)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token:     token,
		ExpiresIn: int((7 * 24 * time.Hour).Seconds()),
		User:      AuthUser{ID: user.ID.String(), Email: user.Email},
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" || req.DeviceID == "" {
		http.Error(w, `{"error":"email, password, device_id required"}`, http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		http.Error(w, `{"error":"password too short"}`, http.StatusBadRequest)
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	user, err := h.UserRepo.Create(r.Context(), req.Email, hash)
	if err != nil {
		http.Error(w, `{"error":"email already exists"}`, http.StatusConflict)
		return
	}
	_, _ = h.DeviceRepo.Upsert(r.Context(), user.ID, req.DeviceID, req.DeviceName)
	token, _ := h.JWT.Sign(user.ID, user.Email, req.DeviceID, 7*24*time.Hour)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token:     token,
		ExpiresIn: int((7 * 24 * time.Hour).Seconds()),
		User:      AuthUser{ID: user.ID.String(), Email: user.Email},
	})
}
