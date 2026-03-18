package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"vpn-startup/backend/internal/bot"
	"vpn-startup/backend/internal/api/handlers"
	"vpn-startup/backend/internal/api/middleware"
	"vpn-startup/backend/internal/auth"
	"vpn-startup/backend/internal/config"
	"vpn-startup/backend/internal/repository"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatal("db:", err)
	}
	defer pool.Close()

	ropt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatal("redis:", err)
	}
	rdb := redis.NewClient(ropt)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("redis ping:", err)
	}
	defer rdb.Close()

	jwt := auth.NewJWT(cfg.JWTSecret)
	userRepo := repository.NewUserRepo(pool)
	deviceRepo := repository.NewDeviceRepo(pool)
	serverRepo := repository.NewServerRepo(pool)
	vpnKeyRepo := repository.NewVPNKeyRepo(pool)
	walletRepo := repository.NewWalletRepo(pool)
	tgRepo := repository.NewTelegramRepo(pool)
	linkCodes := bot.NewLinkCodes(rdb)
	cabinetTokens := bot.NewCabinetTokens(rdb)

	authH := &handlers.AuthHandler{UserRepo: userRepo, DeviceRepo: deviceRepo, JWT: jwt, Cfg: cfg}
	serverH := &handlers.ServerHandler{ServerRepo: serverRepo}
	configH := &handlers.ConfigHandler{
		ServerRepo: serverRepo, VPNKeyRepo: vpnKeyRepo,
		DeviceRepo: deviceRepo, WalletRepo: walletRepo, MaxDevices: cfg.MaxDevices,
	}
	walletH := &handlers.WalletHandler{WalletRepo: walletRepo}
	tgLinkH := &handlers.TelegramLinkHandler{LinkCodes: linkCodes}
	botH := &handlers.BotHandler{
		LinkCodes: linkCodes, TelegramRepo: tgRepo, WalletRepo: walletRepo,
		ServerRepo: serverRepo, VPNKeyRepo: vpnKeyRepo,
		CabinetTokens: cabinetTokens,
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(middleware.RateLimit(rdb, "rl", 100, time.Minute))

	r.Post("/api/login", authH.Login)
	r.Post("/api/register", authH.Register)
	r.Get("/api/servers", serverH.List)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWT(jwt))
		r.Get("/api/config", configH.Get)
		r.Get("/api/wallet", walletH.Get)
		r.Post("/api/wallet/topup", walletH.Topup)
		r.Post("/api/telegram/link", tgLinkH.CreateCode)
	})

	// Telegram bot callbacks (protected by X-Bot-Secret)
	r.Post("/api/bot/link", botH.LinkTelegram)
	r.Post("/api/bot/topup", botH.Topup)
	r.Get("/api/bot/config", botH.ConfigLink)
	r.Get("/api/bot/cabinet", botH.Cabinet)

	// Tokenized cabinet page (hitvpn-like)
	cabinetH := &handlers.CabinetHandler{
		CabinetTokens: cabinetTokens,
		WalletRepo:    walletRepo,
		ServerRepo:    serverRepo,
		VPNKeyRepo:    vpnKeyRepo,
	}
	r.Get("/token/{token}", cabinetH.Page)
	r.Get("/token/{token}/pay", cabinetH.PayLink)
	r.Get("/token/{token}/confirm", cabinetH.Confirm)

	port := cfg.HTTPPort
	if p := os.Getenv("PORT"); p != "" {
		if pn, e := strconv.Atoi(p); e == nil {
			port = pn
		}
	}
	log.Printf("API listening on :%d", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), r); err != nil {
		log.Fatal(err)
	}
}
