package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"vpn-startup/backend/internal/bot"
	"vpn-startup/backend/internal/repository"
	"vpn-startup/backend/internal/vless"
)

// Authenticated endpoint for the app: issues a one-time code to link Telegram.
type TelegramLinkHandler struct {
	LinkCodes *bot.LinkCodes
}

func (h *TelegramLinkHandler) CreateCode(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromRequest(r)
	if claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	code, err := h.LinkCodes.Create(r.Context(), claims.UserID, 10*time.Minute)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	botUser := os.Getenv("TELEGRAM_BOT_USERNAME") // e.g. vitonvpn_bot
	deepLink := ""
	if botUser != "" {
		deepLink = "https://t.me/" + botUser + "?start=" + code
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":      code,
		"deep_link": deepLink,
	})
}

// Bot endpoints protected by X-Bot-Secret. Bot links Telegram user to app user.
type BotHandler struct {
	LinkCodes   *bot.LinkCodes
	TelegramRepo *repository.TelegramRepo
	WalletRepo  *repository.WalletRepo
	ServerRepo  *repository.ServerRepo
	VPNKeyRepo  *repository.VPNKeyRepo
	CabinetTokens *bot.CabinetTokens
}

type botLinkRequest struct {
	Code         string `json:"code"`
	TelegramUser int64  `json:"telegram_user_id"`
}

func (h *BotHandler) LinkTelegram(w http.ResponseWriter, r *http.Request) {
	if !botSecretOK(r) {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}
	var req botLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Code == "" || req.TelegramUser == 0 {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	userID, ok, err := h.LinkCodes.Consume(r.Context(), req.Code)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, `{"error":"invalid_or_expired_code"}`, http.StatusBadRequest)
		return
	}
	if err := h.TelegramRepo.Link(r.Context(), req.TelegramUser, userID); err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok", "user_id": userID.String()})
}

type botTopupRequest struct {
	TelegramUser int64   `json:"telegram_user_id"`
	Amount       float64 `json:"amount"`
}

// For now: credits wallet instantly (demo). In production: verify payment before credit.
func (h *BotHandler) Topup(w http.ResponseWriter, r *http.Request) {
	if !botSecretOK(r) {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}
	var req botTopupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TelegramUser == 0 || req.Amount <= 0 {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	userID, err := h.TelegramRepo.GetUserIDByTelegram(r.Context(), req.TelegramUser)
	if err != nil {
		http.Error(w, `{"error":"not_linked"}`, http.StatusBadRequest)
		return
	}
	if err := h.WalletRepo.AddTopup(r.Context(), userID, req.Amount); err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
}

// Returns a ready VLESS link for the linked user (first server).
func (h *BotHandler) ConfigLink(w http.ResponseWriter, r *http.Request) {
	if !botSecretOK(r) {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}
	tg := r.URL.Query().Get("telegram_user_id")
	if tg == "" {
		http.Error(w, `{"error":"telegram_user_id required"}`, http.StatusBadRequest)
		return
	}
	tgID, err := parseInt64(tg)
	if err != nil || tgID == 0 {
		http.Error(w, `{"error":"invalid telegram_user_id"}`, http.StatusBadRequest)
		return
	}
	userID, err := h.TelegramRepo.GetUserIDByTelegram(r.Context(), tgID)
	if err != nil {
		http.Error(w, `{"error":"not_linked"}`, http.StatusBadRequest)
		return
	}

	// Enforce billing / trial to decide if link should be issued.
	_, allowed, err := h.WalletRepo.ChargeForUsage(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, `{"error":"no_access"}`, http.StatusPaymentRequired)
		return
	}

	servers, err := h.ServerRepo.ListActive(r.Context())
	if err != nil || len(servers) == 0 {
		http.Error(w, `{"error":"no_servers"}`, http.StatusServiceUnavailable)
		return
	}
	userUUID, err := h.VPNKeyRepo.GetOrCreate(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	link := vless.BuildVLESSURL(servers[0], userUUID, "Viton VPN")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"vless_url": link})
}

// Returns a "personal cabinet" URL (tokenized). Actual web UI can be added later.
func (h *BotHandler) Cabinet(w http.ResponseWriter, r *http.Request) {
	if !botSecretOK(r) {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}
	tg := r.URL.Query().Get("telegram_user_id")
	if tg == "" {
		http.Error(w, `{"error":"telegram_user_id required"}`, http.StatusBadRequest)
		return
	}
	tgID, err := parseInt64(tg)
	if err != nil || tgID == 0 {
		http.Error(w, `{"error":"invalid telegram_user_id"}`, http.StatusBadRequest)
		return
	}
	userID, err := h.TelegramRepo.GetUserIDByTelegram(r.Context(), tgID)
	if err != nil {
		http.Error(w, `{"error":"not_linked"}`, http.StatusBadRequest)
		return
	}
	token, err := h.CabinetTokens.Create(r.Context(), userID, 24*time.Hour)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	base := os.Getenv("CABINET_BASE_URL") // e.g. https://vitonvpn.app/token
	if base == "" {
		base = "https://vitonvpn.app/token"
	}
	url := base + "/" + token
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"url": url})
}

func botSecretOK(r *http.Request) bool {
	secret := os.Getenv("BOT_SECRET")
	if secret == "" {
		return false
	}
	return r.Header.Get("X-Bot-Secret") == secret
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

