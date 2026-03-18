package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"vpn-startup/backend/internal/models"
)

// Simple time-based billing: fixed RUB per week, billed on each config request.
const (
	RubPerWeek      = 100.0
	secondsInWeek   = 7 * 24 * 3600
	pricePerSecond  = RubPerWeek / secondsInWeek
	minBalanceCutoff = 0.0
)

type WalletRepo struct {
	pool *pgxpool.Pool
}

func NewWalletRepo(pool *pgxpool.Pool) *WalletRepo {
	return &WalletRepo{pool: pool}
}

// ChargeForUsage subtracts balance based on time since last_billing_at.
// Returns updated user and whether access is allowed (balance > 0 OR trial active).
func (r *WalletRepo) ChargeForUsage(ctx context.Context, userID uuid.UUID) (*models.User, bool, error) {
	const q = `
UPDATE users
SET
    wallet_balance = GREATEST(wallet_balance - (EXTRACT(EPOCH FROM (NOW() - last_billing_at)) * $2), 0),
    last_billing_at = NOW()
WHERE id = $1
RETURNING id, email, password_hash, wallet_balance, last_billing_at, trial_expires_at, created_at, updated_at;
`
	var u models.User
	err := r.pool.QueryRow(ctx, q, userID, pricePerSecond).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.WalletBalance,
		&u.LastBillingAt,
		&u.TrialExpiresAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, false, err
	}
	trialActive := u.TrialExpiresAt != nil && u.TrialExpiresAt.After(time.Now())
	allowed := u.WalletBalance > minBalanceCutoff || trialActive
	return &u, allowed, nil
}

// GetBalance returns current wallet balance for user.
func (r *WalletRepo) GetBalance(ctx context.Context, userID uuid.UUID) (float64, time.Time, error) {
	const q = `SELECT wallet_balance, last_billing_at FROM users WHERE id = $1`
	var bal float64
	var last time.Time
	err := r.pool.QueryRow(ctx, q, userID).Scan(&bal, &last)
	return bal, last, err
}

// AddTopup increases wallet balance (after successful YooMoney payment).
func (r *WalletRepo) AddTopup(ctx context.Context, userID uuid.UUID, amount float64) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET wallet_balance = wallet_balance + $1 WHERE id = $2`, amount, userID)
	return err
}

