package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"vpn-startup/backend/internal/bot"
	"vpn-startup/backend/internal/repository"
	"vpn-startup/backend/internal/vless"
)

type CabinetHandler struct {
	CabinetTokens *bot.CabinetTokens
	WalletRepo    *repository.WalletRepo
	ServerRepo    *repository.ServerRepo
	VPNKeyRepo    *repository.VPNKeyRepo
}

func (h *CabinetHandler) Page(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.NotFound(w, r)
		return
	}
	userID, ok, err := h.CabinetTokens.Get(r.Context(), token)
	if err != nil || !ok {
		http.NotFound(w, r)
		return
	}
	bal, _, err := h.WalletRepo.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	perDay := repository.RubPerWeek / 7.0
	approxDays := 0.0
	if perDay > 0 {
		approxDays = bal / perDay
	}

	// Link (if access allowed by balance or trial)
	_, allowed, _ := h.WalletRepo.ChargeForUsage(r.Context(), userID)
	vlessURL := ""
	if allowed {
		vlessURL, _ = h.firstVlessURL(r.Context(), userID)
	}

	data := map[string]any{
		"Balance":     fmt.Sprintf("%.2f", bal),
		"ApproxDays":  fmt.Sprintf("%.1f", approxDays),
		"HasAccess":   allowed,
		"VlessURL":    vlessURL,
		"Token":       token,
		"SupportText": "Поддержка: support@vitonvpn.app",
		"TelegramBot": os.Getenv("TELEGRAM_BOT_USERNAME"),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = cabinetTpl.Execute(w, data)
}

// Demo: emulate YooMoney link + manual confirm.
func (h *CabinetHandler) PayLink(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	userID, ok, err := h.CabinetTokens.Get(r.Context(), token)
	if err != nil || !ok {
		http.NotFound(w, r)
		return
	}
	amount := r.URL.Query().Get("amount")
	if amount == "" {
		amount = "100"
	}
	label := uuid.NewString()

	receiver := os.Getenv("YOOMONEY_RECEIVER") // wallet number, optional
	if receiver == "" {
		// No receiver configured, fallback to demo confirm page.
		http.Redirect(w, r, "/token/"+token+"/confirm?amount="+amount+"&label="+label, http.StatusFound)
		return
	}

	// YooMoney QuickPay link (P2P). Confirmation still needs webhook; we provide manual confirm for dev.
	quickPay := fmt.Sprintf("https://yoomoney.ru/quickpay/confirm.xml?receiver=%s&quickpay-form=shop&targets=%s&paymentType=SB&sum=%s&label=%s",
		receiver, template.URLQueryEscaper("Viton VPN"), amount, label,
	)
	_ = userID // reserved for future: create wallet_topups record
	http.Redirect(w, r, quickPay, http.StatusFound)
}

// Dev-only: instantly credits balance, then shows VLESS link.
func (h *CabinetHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	userID, ok, err := h.CabinetTokens.Get(r.Context(), token)
	if err != nil || !ok {
		http.NotFound(w, r)
		return
	}
	amountStr := r.URL.Query().Get("amount")
	if amountStr == "" {
		amountStr = "100"
	}
	amount, _ := strconv.ParseFloat(amountStr, 64)
	if amount <= 0 {
		amount = 100
	}
	_ = h.WalletRepo.AddTopup(r.Context(), userID, amount)

	http.Redirect(w, r, "/token/"+token, http.StatusFound)
}

func (h *CabinetHandler) firstVlessURL(ctx context.Context, userID uuid.UUID) (string, error) {
	servers, err := h.ServerRepo.ListActive(ctx)
	if err != nil || len(servers) == 0 {
		return "", fmt.Errorf("no servers")
	}
	uuidKey, err := h.VPNKeyRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return "", err
	}
	return vless.BuildVLESSURL(servers[0], uuidKey, "Viton VPN"), nil
}

var cabinetTpl = template.Must(template.New("cabinet").Parse(`
<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <title>Viton VPN — Личный кабинет</title>
  <style>
    body{margin:0;font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Arial;background:#0b0f14;color:#fff}
    .wrap{max-width:720px;margin:0 auto;padding:24px}
    .card{background:#101826;border:1px solid rgba(255,255,255,.08);border-radius:16px;padding:18px;margin:12px 0}
    .row{display:flex;gap:12px;flex-wrap:wrap}
    .btn{display:inline-block;padding:12px 14px;border-radius:12px;text-decoration:none;font-weight:600}
    .btn-green{background:#00E676;color:#000}
    .btn-dark{background:rgba(255,255,255,.08);color:#fff;border:1px solid rgba(255,255,255,.12)}
    .muted{color:rgba(255,255,255,.7)}
    code{display:block;white-space:pre-wrap;word-break:break-all;background:rgba(0,0,0,.35);padding:12px;border-radius:12px;border:1px solid rgba(255,255,255,.08)}
    h1{margin:0 0 10px;font-size:22px}
  </style>
</head>
<body>
  <div class="wrap">
    <h1>Viton VPN — Личный кабинет</h1>
    <div class="card">
      <div class="row" style="justify-content:space-between">
        <div>
          <div class="muted">Баланс</div>
          <div style="font-size:28px;font-weight:700">{{.Balance}} ₽</div>
        </div>
        <div>
          <div class="muted">Осталось</div>
          <div style="font-size:28px;font-weight:700">~{{.ApproxDays}} дн.</div>
        </div>
      </div>
      <div class="muted" style="margin-top:8px">Тариф: 100 ₽ ≈ 7 дней VPN</div>
    </div>

    <div class="card">
      <div style="font-weight:700;margin-bottom:8px">Пополнение YooMoney</div>
      <div class="row">
        <a class="btn btn-green" href="/token/{{.Token}}/pay?amount=100">Пополнить 100 ₽</a>
        <a class="btn btn-dark" href="/token/{{.Token}}/pay?amount=300">Пополнить 300 ₽</a>
        <a class="btn btn-dark" href="/token/{{.Token}}/pay?amount=500">Пополнить 500 ₽</a>
      </div>
      <div class="muted" style="margin-top:10px">
        В dev-режиме без настроенного YooMoney кошелька пополнение подтвердится автоматически.
      </div>
    </div>

    <div class="card">
      <div style="font-weight:700;margin-bottom:8px">Ссылка доступа</div>
      {{if .HasAccess}}
        <div class="muted" style="margin-bottom:10px">Скопируйте VLESS ссылку и добавьте в клиент.</div>
        <code>{{.VlessURL}}</code>
      {{else}}
        <div class="muted">Доступ не активен. Пополните баланс или используйте пробный период.</div>
      {{end}}
    </div>

    <div class="card">
      <div class="muted">{{.SupportText}}</div>
    </div>
  </div>
</body>
</html>
`))

