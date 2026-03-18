package handlers

import (
	"net/http"

	"vpn-startup/backend/internal/auth"
)

func ClaimsFromRequest(r *http.Request) *auth.Claims {
	ctx := r.Context()
	claims, _ := ctx.Value("claims").(*auth.Claims)
	return claims
}
