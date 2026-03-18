package handlers

import (
	"encoding/json"
	"net/http"

	"vpn-startup/backend/internal/repository"
)

type ServerHandler struct {
	ServerRepo *repository.ServerRepo
}

type ServerListItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Type     string `json:"type"`
	IsActive bool   `json:"is_active"`
}

func (h *ServerHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.ServerRepo.ListActive(r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	out := make([]ServerListItem, 0, len(list))
	for _, s := range list {
		out = append(out, ServerListItem{
			ID:       s.ID.String(),
			Name:     s.Name,
			Region:   s.Region,
			Host:     s.Host,
			Port:     s.Port,
			Type:     s.Type,
			IsActive: s.IsActive,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"servers": out})
}
