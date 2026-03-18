package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	WalletBalance float64   `json:"wallet_balance"`
	LastBillingAt time.Time `json:"last_billing_at"`
	TrialExpiresAt *time.Time `json:"trial_expires_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Device struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	DeviceID  string    `json:"device_id"`
	Name      string    `json:"name"`
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
}

type Server struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Region          string    `json:"region"`
	Host            string    `json:"host"`
	Port            int       `json:"port"`
	Type            string    `json:"type"`
	RealityPubKey   string    `json:"-"`
	RealityShortID  string    `json:"-"`
	RealitySNI      string    `json:"-"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Subscription struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Plan       string    `json:"plan"`
	StartedAt time.Time `json:"started_at"`
	ExpiresAt time.Time `json:"expires_at"`
	AutoRenew  bool      `json:"auto_renew"`
	ExternalID string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserVPNKey struct {
	UserID    uuid.UUID `json:"user_id"`
	VlessUUID uuid.UUID `json:"vless_uuid"`
	CreatedAt time.Time `json:"created_at"`
}
