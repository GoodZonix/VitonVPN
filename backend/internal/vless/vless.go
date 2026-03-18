package vless

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"vpn-startup/backend/internal/models"
)

// BuildVLESSURL builds a VLESS URL for a server and user UUID.
// Format: vless://UUID@host:port?type=tcp&security=reality&pbk=...&sid=...&sni=...&fp=chrome#name
func BuildVLESSURL(server models.Server, userUUID uuid.UUID, remark string) string {
	u := url.URL{
		Scheme: "vless",
		User:   url.User(userUUID.String()),
		Host:   fmt.Sprintf("%s:%d", server.Host, server.Port),
		Path:   "",
	}
	q := u.Query()
	q.Set("type", "tcp")
	q.Set("flow", "xtls-rprx-vision")
	switch server.Type {
	case "reality":
		q.Set("security", "reality")
		if server.RealityPubKey != "" {
			q.Set("pbk", server.RealityPubKey)
		}
		if server.RealityShortID != "" {
			q.Set("sid", server.RealityShortID)
		}
		if server.RealitySNI != "" {
			q.Set("sni", server.RealitySNI)
		}
		q.Set("fp", "chrome")
	case "tls":
		q.Set("security", "tls")
		if server.RealitySNI != "" {
			q.Set("sni", server.RealitySNI)
		}
		q.Set("alpn", "h2,http/1.1")
	default:
		q.Set("security", "reality")
		q.Set("fp", "chrome")
	}
	u.RawQuery = q.Encode()
	if remark != "" {
		u.Fragment = remark
	}
	return u.String()
}
