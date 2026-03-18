package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"vpn-startup/backend/internal/models"
)

type SubscriptionRepo struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepo(pool *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{pool: pool}
}

func (r *SubscriptionRepo) ActiveByUser(ctx context.Context, userID uuid.UUID) (*models.Subscription, error) {
	s := &models.Subscription{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, plan, started_at, expires_at, auto_renew, created_at, updated_at
		 FROM subscriptions WHERE user_id = $1 AND expires_at > NOW() ORDER BY expires_at DESC LIMIT 1`,
		userID,
	).Scan(&s.ID, &s.UserID, &s.Plan, &s.StartedAt, &s.ExpiresAt, &s.AutoRenew, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SubscriptionRepo) Create(ctx context.Context, userID uuid.UUID, plan string, duration time.Duration, externalID string) (*models.Subscription, error) {
	now := time.Now()
	exp := now.Add(duration)
	s := &models.Subscription{UserID: userID, Plan: plan, StartedAt: now, ExpiresAt: exp, ExternalID: externalID}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO subscriptions (user_id, plan, started_at, expires_at, external_id) VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, plan, started_at, expires_at, auto_renew, created_at, updated_at`,
		userID, plan, now, exp, externalID,
	).Scan(&s.ID, &s.UserID, &s.Plan, &s.StartedAt, &s.ExpiresAt, &s.AutoRenew, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *SubscriptionRepo) Extend(ctx context.Context, subscriptionID uuid.UUID, extendBy time.Duration) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE subscriptions SET expires_at = expires_at + $1, updated_at = NOW() WHERE id = $2`,
		extendBy, subscriptionID,
	)
	return err
}
