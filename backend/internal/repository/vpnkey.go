package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"vpn-startup/backend/internal/models"
)

type VPNKeyRepo struct {
	pool *pgxpool.Pool
}

func NewVPNKeyRepo(pool *pgxpool.Pool) *VPNKeyRepo {
	return &VPNKeyRepo{pool: pool}
}

func (r *VPNKeyRepo) GetOrCreate(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var vlessUUID uuid.UUID
	err := r.pool.QueryRow(ctx, `SELECT vless_uuid FROM user_vpn_keys WHERE user_id = $1`, userID).Scan(&vlessUUID)
	if err == nil {
		return vlessUUID, nil
	}
	vlessUUID = uuid.New()
	_, err = r.pool.Exec(ctx, `INSERT INTO user_vpn_keys (user_id, vless_uuid) VALUES ($1, $2) ON CONFLICT (user_id) DO NOTHING`, userID, vlessUUID)
	if err != nil {
		return uuid.Nil, err
	}
	// If ON CONFLICT DO NOTHING, we might have lost the race; select again
	_ = r.pool.QueryRow(ctx, `SELECT vless_uuid FROM user_vpn_keys WHERE user_id = $1`, userID).Scan(&vlessUUID)
	return vlessUUID, nil
}

func (r *VPNKeyRepo) Get(ctx context.Context, userID uuid.UUID) (*models.UserVPNKey, error) {
	k := &models.UserVPNKey{UserID: userID}
	err := r.pool.QueryRow(ctx, `SELECT vless_uuid, created_at FROM user_vpn_keys WHERE user_id = $1`, userID).Scan(&k.VlessUUID, &k.CreatedAt)
	if err != nil {
		return nil, err
	}
	return k, nil
}
