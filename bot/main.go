package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type tgUpdate struct {
	UpdateID int `json:"update_id"`
	Message  *struct {
		MessageID int `json:"message_id"`
		From      *struct {
			ID int64 `json:"id"`
		} `json:"from"`
		Chat *struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

type tgResp struct {
	OK          bool       `json:"ok"`
	Result      []tgUpdate `json:"result"`
	ErrorCode   int        `json:"error_code"`
	Description string     `json:"description"`
}

type replyKeyboard struct {
	Keyboard        [][]map[string]string `json:"keyboard"`
	ResizeKeyboard  bool                  `json:"resize_keyboard"`
	IsPersistent    bool                  `json:"is_persistent"`
	OneTimeKeyboard bool                  `json:"one_time_keyboard"`
}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	apiBase := os.Getenv("API_BASE_URL") // e.g. http://localhost:8080
	secret := os.Getenv("BOT_SECRET")
	if token == "" || apiBase == "" || secret == "" {
		log.Fatal("set TELEGRAM_BOT_TOKEN, API_BASE_URL, BOT_SECRET")
	}

	offset := 0
	for {
		updates, err := getUpdates(token, offset)
		if err != nil {
			log.Println("getUpdates:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		for _, u := range updates {
			offset = u.UpdateID + 1
			if u.Message == nil || u.Message.From == nil || u.Message.Chat == nil {
				continue
			}
			chatID := u.Message.Chat.ID
			tgUserID := u.Message.From.ID
			text := strings.TrimSpace(u.Message.Text)
			if text == "" {
				continue
			}

			switch {
			case strings.HasPrefix(text, "/start"):
				parts := strings.Fields(text)
				if len(parts) < 2 {
					sendMenu(token, chatID, "Рады видеть вас снова!\n\nОткройте Viton VPN → Кошелёк → «Привязать Telegram-бота» и отправьте сюда команду /start <код>.")
					continue
				}
				code := parts[1]
				if err := apiLink(apiBase, secret, code, tgUserID); err != nil {
					sendMenu(token, chatID, "Не удалось привязать аккаунт.\nПроверьте, что код не истёк, и попробуйте снова.")
					continue
				}
				cab, _ := apiCabinet(apiBase, secret, tgUserID)
				msg := "Рады видеть вас снова!\n\nПерейдите в личный кабинет по ссылке:\n"
				if cab != "" {
					msg += "👉👉 " + cab + " 👈👈\n"
				} else {
					msg += "👉👉 (ссылка недоступна) 👈👈\n"
				}
				msg += "\nСпасибо, что остаётесь с нами!"
				sendMenu(token, chatID, msg)

			case strings.HasPrefix(text, "/topup"):
				amount := 100.0
				parts := strings.Fields(text)
				if len(parts) >= 2 {
					if v, err := strconv.ParseFloat(parts[1], 64); err == nil && v > 0 {
						amount = v
					}
				}
				if err := apiTopup(apiBase, secret, tgUserID, amount); err != nil {
					sendMenu(token, chatID, "Пополнение не выполнено.\nСначала привяжите аккаунт через /start <код>.")
					continue
				}
				link, err := apiConfigLink(apiBase, secret, tgUserID)
				if err != nil {
					sendMenu(token, chatID, "Баланс пополнен, но ссылку получить не удалось (возможно нет серверов).")
					continue
				}
				sendMenu(token, chatID, fmt.Sprintf("Баланс пополнен на %.0f ₽.\n\nВаша ссылка доступа:\n%s", amount, link))

			case strings.HasPrefix(text, "/config"):
				link, err := apiConfigLink(apiBase, secret, tgUserID)
				if err != nil {
					sendMenu(token, chatID, "Ссылку получить не удалось. Проверьте баланс или привязку аккаунта.")
					continue
				}
				sendMenu(token, chatID, "Ваша ссылка доступа:\n"+link)

			case strings.EqualFold(text, "Личный кабинет"):
				cab, err := apiCabinet(apiBase, secret, tgUserID)
				if err != nil || cab == "" {
					sendMenu(token, chatID, "Личный кабинет недоступен. Проверьте привязку аккаунта.")
					continue
				}
				sendMenu(token, chatID, "Перейдите в личный кабинет по ссылке:\n👉👉 "+cab+" 👈👈")

			case strings.EqualFold(text, "Кабинет не работает"):
				sendMenu(token, chatID, "Если кабинет не открывается:\n1) Проверьте интернет\n2) Попробуйте позже\n3) Напишите в поддержку (кнопка «Помощь»)")

			case strings.EqualFold(text, "Пригласить"):
				sendMenu(token, chatID, "Пригласите друга в Viton VPN и получите бонус!\n(реферальная система — в разработке)")

			case strings.EqualFold(text, "Помощь"):
				sendMenu(token, chatID, "Поддержка: support@vitonvpn.app\nОпишите проблему и ваш ID из приложения.")

			case strings.EqualFold(text, "VPN не работает?"):
				sendMenu(token, chatID, "Если VPN не работает:\n- Перезапустите приложение\n- Переключите сервер\n- Проверьте баланс (кошелёк)\n- Если не помогло — напишите в поддержку")

			default:
				sendMenu(token, chatID, "Меню доступно кнопками ниже.\nПривязка: /start <код>")
			}
		}
	}
}

func getUpdates(token string, offset int) ([]tgUpdate, error) {
	u := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?timeout=30&offset=%d", token, offset)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out tgResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if !out.OK {
		if out.ErrorCode != 0 || out.Description != "" {
			return nil, fmt.Errorf("telegram api not ok: %d %s", out.ErrorCode, out.Description)
		}
		return nil, fmt.Errorf("telegram api not ok")
	}
	return out.Result, nil
}

func sendMessage(token string, chatID int64, text string, markup any) {
	body := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	if markup != nil {
		body["reply_markup"] = markup
	}
	b, _ := json.Marshal(body)
	_, _ = http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token), "application/json", bytes.NewReader(b))
}

func sendMenu(token string, chatID int64, text string) {
	kb := replyKeyboard{
		ResizeKeyboard: true,
		IsPersistent:   true,
		Keyboard: [][]map[string]string{
			{{"text": "Личный кабинет"}},
			{{"text": "Кабинет не работает"}, {"text": "VPN не работает?"}},
			{{"text": "Пригласить"}, {"text": "Помощь"}},
		},
	}
	sendMessage(token, chatID, text, kb)
}

func apiLink(apiBase, secret, code string, tgUserID int64) error {
	body := map[string]interface{}{
		"code":             code,
		"telegram_user_id": tgUserID,
	}
	return postJSON(apiBase+"/api/bot/link", secret, body, nil)
}

func apiCabinet(apiBase, secret string, tgUserID int64) (string, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/bot/cabinet?telegram_user_id=%d", apiBase, tgUserID), nil)
	req.Header.Set("X-Bot-Secret", secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	var m map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return "", err
	}
	if v, ok := m["url"].(string); ok && v != "" {
		return v, nil
	}
	return "", fmt.Errorf("missing url")
}

func apiTopup(apiBase, secret string, tgUserID int64, amount float64) error {
	body := map[string]interface{}{
		"telegram_user_id": tgUserID,
		"amount":           amount,
	}
	return postJSON(apiBase+"/api/bot/topup", secret, body, nil)
}

func apiConfigLink(apiBase, secret string, tgUserID int64) (string, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/bot/config?telegram_user_id=%d", apiBase, tgUserID), nil)
	req.Header.Set("X-Bot-Secret", secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	var m map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return "", err
	}
	if v, ok := m["vless_url"].(string); ok && v != "" {
		return v, nil
	}
	return "", fmt.Errorf("missing vless_url")
}

func postJSON(url, secret string, body map[string]interface{}, out any) error {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Bot-Secret", secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}
