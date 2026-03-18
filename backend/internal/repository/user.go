package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"vpn-startup/backend/internal/models"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, email, passwordHash string) (*models.User, error) {
	u := &models.User{ID: uuid.New()}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (id, email, password_hash, trial_expires_at) VALUES ($1, $2, $3, NOW() + interval '2 days')
		 RETURNING id, email, wallet_balance, last_billing_at, trial_expires_at, created_at, updated_at`,
		u.ID, email, passwordHash,
	).Scan(&u.ID, &u.Email, &u.WalletBalance, &u.LastBillingAt, &u.TrialExpiresAt, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, wallet_balance, last_billing_at, trial_expires_at, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.WalletBalance, &u.LastBillingAt, &u.TrialExpiresAt, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	u := &models.User{ID: id}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, wallet_balance, last_billing_at, trial_expires_at, created_at, updated_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.WalletBalance, &u.LastBillingAt, &u.TrialExpiresAt, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
