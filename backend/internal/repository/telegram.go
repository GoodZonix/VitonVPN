package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TelegramRepo struct {
	pool *pgxpool.Pool
}

func NewTelegramRepo(pool *pgxpool.Pool) *TelegramRepo {
	return &TelegramRepo{pool: pool}
}

func (r *TelegramRepo) Link(ctx context.Context, tgUserID int64, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO telegram_accounts (tg_user_id, user_id) VALUES ($1, $2)
		 ON CONFLICT (tg_user_id) DO UPDATE SET user_id = EXCLUDED.user_id`,
		tgUserID, userID,
	)
	return err
}

func (r *TelegramRepo) GetUserIDByTelegram(ctx context.Context, tgUserID int64) (uuid.UUID, error) {
	var userID uuid.UUID
	err := r.pool.QueryRow(ctx, `SELECT user_id FROM telegram_accounts WHERE tg_user_id = $1`, tgUserID).Scan(&userID)
	return userID, err
}

